package main

import (
	"context"
	"database/sql"
	"os"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/nickysemenza/food/db"
	"github.com/nickysemenza/food/manager"
	"github.com/nickysemenza/food/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
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

	if false {
		// google cloud trace
		// env: GOOGLE_CLOUD_PROJECT=xx GOOGLE_APPLICATION_CREDENTIALS=x.json

		projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
		exporter, err := texporter.NewExporter(texporter.WithProjectID(projectID))
		if err != nil {
			log.Fatalf("texporter.NewExporter: %v", err)
		}

		tp, err := sdktrace.NewProvider(sdktrace.WithSyncer(exporter))
		if err != nil {
			log.Fatal(err)
		}
		global.SetTraceProvider(tp)
	}

	return err
}

func main() {
	ctx := context.Background()
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
	if err = dbx.PingContext(ctx); err != nil {
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
