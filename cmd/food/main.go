package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nickysemenza/food/db"
	"github.com/nickysemenza/food/manager"
	"github.com/nickysemenza/food/server"
	"github.com/spf13/viper"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"

	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() {
	// Create and install Jaeger export pipeline
	_, _, err := jaeger.NewExportPipeline(
		jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "food",
			Tags: []core.KeyValue{
				key.String("exporter", "jaeger"),
				key.Float64("float", 312.23),
			},
		}),
		jaeger.RegisterAsGlobal(),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5555)
	viper.SetDefault("DB_USER", "food")
	viper.SetDefault("DB_PASSWORD", "food")
	viper.SetDefault("DB_DBNAME", "food")

	viper.AutomaticEnv()

	initTracer()

	dbConn, err := sqlx.Open("postgres", db.ConnnectionString(
		viper.GetString("DB_HOST"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_DBNAME"),
		viper.GetInt64("DB_PORT")))

	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	if err = dbConn.Ping(); err != nil {
		panic(err)
	}
	dbClient := db.New(dbConn)
	m := manager.New(dbClient)
	s := server.Server{Manager: m, HTTPPort: 4242, DB: dbClient}

	spew.Dump(s.Run())

}
