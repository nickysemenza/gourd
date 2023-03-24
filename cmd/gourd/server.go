package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
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
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func getDBConn(dsn string, kind db.Kind) (*sql.DB, error) {

	if !strings.Contains(dsn, "sslmode=disable") {
		dsn += "?sslmode=disable"
	}
	log.Info("connecting to db: ", dsn)
	dbConn, err := otelsql.Open("postgres", dsn,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithDBName(fmt.Sprintf("%s_db", kind)))

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
		i, err = image.NewS3Store(ctx, viper.GetString("s3_endpoint"), viper.GetString("s3_region"), viper.GetString("s3_bucket"), viper.GetString("s3_keyid"), viper.GetString("s3_appkey"), viper.GetString("s3_prefix"))
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
