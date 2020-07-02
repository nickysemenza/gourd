package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/getsentry/sentry-go"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/luna-duclos/instrumentedsql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/nickysemenza/food/db"
	"github.com/nickysemenza/food/manager"
	"github.com/nickysemenza/food/server"
)

func initTracer() error {
	endpoint := viper.GetString("JAEGER_ENDPOINT")
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	// Create and install Jaeger export pipeline
	if projectID != "" {
		// google cloud trace
		// env: GOOGLE_CLOUD_PROJECT=xx GOOGLE_APPLICATION_CREDENTIALS=x.json
		exporter, err := texporter.NewExporter(texporter.WithProjectID(projectID))
		if err != nil {
			return fmt.Errorf("texporter.NewExporter: %w", err)
		}

		tp, err := sdktrace.NewProvider(sdktrace.WithSyncer(exporter))
		if err != nil {
			return err
		}
		global.SetTraceProvider(tp)
	} else if endpoint != "" {
		_, _, err := jaeger.NewExportPipeline(
			jaeger.WithCollectorEndpoint(endpoint),
			jaeger.WithProcess(jaeger.Process{
				ServiceName: "food",
				Tags: []core.KeyValue{
					key.String("exporter", "jaeger"),
				},
			}),
			jaeger.RegisterAsGlobal(),
			jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		)
		return err
	}

	return nil
}

func setupEnv() {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5555)
	viper.SetDefault("DB_USER", "food")
	viper.SetDefault("DB_PASSWORD", "food")
	viper.SetDefault("DB_DBNAME", "food")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 20)
	viper.SetDefault("PORT", 4242)
	viper.SetDefault("HTTP_TIMEOUT", "30s")
	viper.SetDefault("SENTRY_DSN", "https://8220ab8a2b3d4c3c9cf7f636ec183c7a@o83311.ingest.sentry.io/5298706")

	viper.SetDefault("JAEGER_ENDPOINT", "http://localhost:14268/api/traces")
}
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
func main() {
	ctx := context.Background()

	// env vars
	setupEnv()
	viper.AutomaticEnv()

	// tracing
	if err := initTracer(); err != nil {
		log.Fatal(err)
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn: viper.GetString("SENTRY_DSN"),
	}); err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

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
