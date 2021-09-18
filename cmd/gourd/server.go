package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/luna-duclos/instrumentedsql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/nickysemenza/gourd/api"
	"github.com/nickysemenza/gourd/auth"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/google"
	"github.com/nickysemenza/gourd/manager"
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
	logger := log.New()
	logger.SetLevel(log.DebugLevel)
	// setupMisc()

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
	if err := autoMigrate(dbConn); err != nil {
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

	r := rs_client.New(viper.GetString("RS_URI"))
	m := manager.New(dbClient, gClient, auth, r)
	apiManager := api.NewAPI(m)

	// server
	return &server.Server{
		Manager:     m,
		HTTPPort:    viper.GetUint("PORT"),
		DB:          dbClient,
		HTTPTimeout: viper.GetDuration("HTTP_TIMEOUT"),
		HTTPHost:    viper.GetString("HTTP_HOST"),
		APIManager:  apiManager,
		BypassAuth:  viper.GetBool("BYPASS_AUTH"),
		Logger:      logger,
	}, nil
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

func autoMigrate(dbConn *sql.DB) error {
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to migrate: %w", err)
		}
	}
	log.Info("db: migrated")
	return nil
}
