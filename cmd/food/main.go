package main

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/nickysemenza/food/db"
	"github.com/nickysemenza/food/manager"
	"github.com/nickysemenza/food/server"
	"github.com/spf13/viper"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"

	"github.com/luna-duclos/instrumentedsql"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer(endpoint string) error {
	// Create and install Jaeger export pipeline
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

func main() {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5555)
	viper.SetDefault("DB_USER", "food")
	viper.SetDefault("DB_PASSWORD", "food")
	viper.SetDefault("DB_DBNAME", "food")

	viper.SetDefault("JAEGER_ENDPOINT", "http://localhost:14268/api/traces")

	viper.AutomaticEnv()

	err := initTracer(viper.GetString("JAEGER_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}

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

	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	dbx := sqlx.NewDb(dbConn, "postgres")
	if err = dbx.Ping(); err != nil {
		log.Fatal(err)
	}
	dbClient := db.New(dbx)

	s := server.Server{
		Manager:  manager.New(dbClient),
		HTTPPort: 4242,
		DB:       dbClient,
	}

	log.Fatal(s.Run())

}
