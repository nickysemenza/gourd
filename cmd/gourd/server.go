package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/luna-duclos/instrumentedsql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/manager"
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
func runServer() {
	ctx := context.Background()

	// setupMisc()

	// postgres database
	dbConn, err := getDBConn()
	if err != nil {
		log.Fatal(err)
	}
	dbConn.SetMaxOpenConns(viper.GetInt("DB_MAX_OPEN_CONNS"))
	defer dbConn.Close()

	dbClient, err := db.New(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	// postgres migrations
	if err := autoMigrate(dbConn); err != nil {
		log.Fatal(err)
	}

	// server
	s := server.Server{
		Manager:     manager.New(dbClient),
		HTTPPort:    viper.GetUint("PORT"),
		DB:          dbClient,
		HTTPTimeout: viper.GetDuration("HTTP_TIMEOUT"),
		HTTPHost:    viper.GetString("HTTP_HOST"),
	}
	log.Fatal(s.Run(ctx))
}

func autoMigrate(dbConn *sql.DB) error {
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return fmt.Errorf("failed to migrate: %w", err)
		}
	}
	log.Info("db: migrated")
	return nil
}
