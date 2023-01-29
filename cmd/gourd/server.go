package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/luna-duclos/instrumentedsql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/nickysemenza/gourd/internal/api"
	"github.com/nickysemenza/gourd/internal/auth"
	"github.com/nickysemenza/gourd/internal/clients/google"
	"github.com/nickysemenza/gourd/internal/clients/notion"
	"github.com/nickysemenza/gourd/internal/clients/rs_client"
	"github.com/nickysemenza/gourd/internal/db"
	"github.com/nickysemenza/gourd/internal/image"
	"github.com/nickysemenza/gourd/internal/server"
)

func getDBConn(dsn string, kind db.Kind) (*sql.DB, error) {
	driver := "instrumented-postgres-" + string(kind)
	sql.Register(driver, instrumentedsql.WrapDriver(&pq.Driver{},
		instrumentedsql.WithLogger(instrumentedsql.LoggerFunc(func(ctx context.Context, msg string, keyvals ...interface{}) {
			log.WithContext(ctx).WithFields(log.Fields{"client": "db", "kind": string(kind)}).Debugf("%s %v", msg, keyvals)
		})),
		instrumentedsql.WithTracer(NewTracer(true)),
		instrumentedsql.WithOpsExcluded(instrumentedsql.OpSQLRowsNext),
	))

	if !strings.Contains(dsn, "sslmode=disable") {
		dsn += "?sslmode=disable"
	}
	log.Info("connecting to db: ", dsn)
	dbConn, err := sql.Open(driver, dsn)
	return dbConn, err
}
func makeServer(ctx context.Context) (*server.Server, error) {

	// postgres database
	dbConn, err := getDBConn(viper.GetString("DATABASE_URL"), db.Gourd)
	if err != nil {
		err := fmt.Errorf("failed to init db conn: %w", err)
		log.Fatal(err)
	}

	dbConnUSDA, err := getDBConn(viper.GetString("DATABASE_URL_USDA"), db.USDA)
	if err != nil {
		err := fmt.Errorf("failed to init db conn: %w", err)
		log.Fatal(err)
	}

	dbClient, err := db.New(dbConn, db.Gourd)
	if err != nil {
		err := fmt.Errorf("failed to init db: %w", err)
		log.Fatal(err)
	}

	dbClientUSDA, err := db.New(dbConnUSDA, db.USDA)
	if err != nil {
		err := fmt.Errorf("failed to init db: %w", err)
		log.Fatal(err)
	}

	// postgres internal/db/migrations
	if err := db.AutoMigrate(ctx, dbConn, true); err != nil {
		err := fmt.Errorf("failed to migrate db: %w", err)
		log.Fatal(err)
	}

	gClient := google.New(dbClient,
		viper.GetString("GOOGLE_CLIENT_ID"),
		viper.GetString("GOOGLE_CLIENT_SECRET"),
		viper.GetString("GOOGLE_REDIRECT_URL"),
	)
	auth, err := auth.New(viper.GetString("JWT_KEY"))
	if err != nil {
		return nil, err
	}

	s := &server.Server{
		HTTPPort:    viper.GetUint("PORT"),
		DB:          dbClient,
		HTTPTimeout: viper.GetDuration("HTTP_TIMEOUT"),
		HTTPHost:    viper.GetString("HTTP_HOST"),
		BypassAuth:  viper.GetBool("BYPASS_AUTH"),
		Logger:      log.New(),
	}

	r := rs_client.New(viper.GetString("RS_URI"))
	n := notion.New(viper.GetString("notion_secret"), viper.GetString("notion_db"))
	//
	var i image.Store
	if viper.GetBool("image_s3") {
		i, err = image.NewS3Store(viper.GetString("s3_endpoint"), viper.GetString("s3_region"), viper.GetString("s3_bucket"), viper.GetString("s3_keyid"), viper.GetString("s3_appkey"))
	} else {
		i, err = image.NewLocalImageStore(s.GetBaseURL())
	}
	if err != nil {
		return nil, err
	}
	m := api.New(dbClient, dbClientUSDA, gClient, auth, r, n, i)
	s.APIManager = m

	// server
	return s, nil
}
func runServer(ctx context.Context) error {
	s, err := makeServer(ctx)
	if err != nil {
		err := fmt.Errorf("failed to make server: %w", err)
		log.Fatal(err)
	}
	return s.Run(ctx)
}
