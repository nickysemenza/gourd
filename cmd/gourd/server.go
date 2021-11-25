package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/luna-duclos/instrumentedsql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/nickysemenza/gourd/api"
	"github.com/nickysemenza/gourd/auth"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/google"
	"github.com/nickysemenza/gourd/image"
	"github.com/nickysemenza/gourd/notion"
	"github.com/nickysemenza/gourd/rs_client"
	"github.com/nickysemenza/gourd/server"
)

func getDBConn() (*sql.DB, error) {
	sql.Register("instrumented-postgres", instrumentedsql.WrapDriver(&pq.Driver{},
		instrumentedsql.WithLogger(instrumentedsql.LoggerFunc(func(ctx context.Context, msg string, keyvals ...interface{}) {
			log.WithContext(ctx).WithFields(log.Fields{"client": "db"}).Debugf("%s %v", msg, keyvals)
		})),
		instrumentedsql.WithTracer(NewTracer(true)),
		instrumentedsql.WithOpsExcluded(instrumentedsql.OpSQLRowsNext),
	))
	dbConn, err := sql.Open("instrumented-postgres", db.ConnnectionString(
		viper.GetString("DB_HOST"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_DBNAME"),
		viper.GetInt64("DB_PORT")))
	return dbConn, err
}
func makeServer() (*server.Server, error) {

	// postgres database
	dbConn, err := getDBConn()
	if err != nil {
		err := fmt.Errorf("failed to init db conn: %w", err)
		log.Fatal(err)
	}
	dbConn.SetMaxOpenConns(viper.GetInt("DB_MAX_OPEN_CONNS"))
	// defer dbConn.Close()

	dbClient, err := db.New(dbConn)
	if err != nil {
		err := fmt.Errorf("failed to init db: %w", err)
		log.Fatal(err)
	}

	// postgres db/migrations
	if err := db.AutoMigrate(dbConn, true); err != nil {
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
	n := notion.New(viper.GetString("notion_secret"), viper.GetString("notion_db"), viper.GetBool("notion_test_only"))
	i, err := image.NewLocalImageStore(s.GetBaseURL())
	if err != nil {
		return nil, err
	}
	m := api.New(dbClient, gClient, auth, r, n, i)
	s.APIManager = m

	// server
	return s, nil
}
func runServer() {
	ctx := context.Background()
	s, err := makeServer()
	if err != nil {
		err := fmt.Errorf("failed to make server: %w", err)
		log.Fatal(err)
	}
	log.Fatal(s.Run(ctx))
}
