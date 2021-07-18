package main

import (
	"errors"
	"fmt"
	"os"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"

	jaegerp "go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
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

		tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))

		otel.SetTracerProvider(tp)
	} else if endpoint != "" {
		_, err := jaeger.InstallNewPipeline(
			jaeger.WithCollectorEndpoint(endpoint),
			jaeger.WithProcess(jaeger.Process{
				ServiceName: "gourd",
				Tags: []attribute.KeyValue{
					attribute.String("exporter", "jaeger"),
				},
			}),
			jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, jaegerp.Jaeger{}))
		return err
	}

	return nil
}

func setupEnv() {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5555)
	viper.SetDefault("DB_USER", "gourd")
	viper.SetDefault("DB_PASSWORD", "gourd")
	viper.SetDefault("DB_DBNAME", "food")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 20)
	viper.SetDefault("PORT", 4242)
	viper.SetDefault("HTTP_TIMEOUT", "30s")
	viper.SetDefault("SENTRY_DSN", "https://6a67b0ba96a744d2877fc1b21984aa18@o83311.ingest.sentry.io/5778936")

	viper.SetDefault("JAEGER_ENDPOINT", "http://localhost:14268/api/traces")
	viper.SetDefault("BYPASS_AUTH", false)

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		var configErr *viper.ConfigFileNotFoundError
		if errors.As(err, &configErr) {
			panic(fmt.Errorf("Fatal error config file: %s \n", configErr))
		}
	}

}

func setupMisc() {
	// env vars
	setupEnv()
	viper.AutomaticEnv()
	// if err := viper.WriteConfig(); err != nil {
	// 	panic(err)
	// }

	// tracing
	if err := initTracer(); err != nil {
		err := fmt.Errorf("failed to init tracer: %w", err)
		log.Fatal(err)
	}
	log.Infof("tracer initialized")

	if err := sentry.Init(sentry.ClientOptions{
		Dsn: viper.GetString("SENTRY_DSN"),
	}); err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	log.Infof("sentry initialized")
}
