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
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer() error {
	endpoint := viper.GetString("JAEGER_ENDPOINT")
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	// Create and install Jaeger export pipeline
	if projectID != "" {
		// google cloud trace
		// env: GOOGLE_CLOUD_PROJECT=xx GOOGLE_APPLICATION_CREDENTIALS=x.json
		// exporter, err := texporter.NewExporter(texporter.WithProjectID(projectID))
		// if err != nil {
		// 	return fmt.Errorf("texporter.NewExporter: %w", err)
		// }
		exporter, err := texporter.New(texporter.WithProjectID(projectID))
		if err != nil {
			return fmt.Errorf("texporter.NewExporter: %w", err)
		}
		tp := tracesdk.NewTracerProvider(
			tracesdk.WithBatcher(exporter),
		)

		// tp := tracesdk.NewTracerProvider(tracesdk.WithSyncer(exporter))

		otel.SetTracerProvider(tp)
	} else if endpoint != "" {
		tp, err := tracerProvider(endpoint)
		if err != nil {
			return err
		}
		otel.SetTracerProvider(tp)

		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
			jaegerp.Jaeger{},
		))
		return err
	}

	return nil
}
func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("gourd"),
			attribute.String("environment", "dev"),
		)),
	)
	return tp, nil
}

func setupEnv() error {
	viper.SetDefault("LOG_LEVEL", "info")
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
	viper.SetDefault("RS_URI", "http://localhost:8080/")
	viper.SetDefault("PDFLATEX_BINARY", "pdflatex")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if errors.As(err, &viper.ConfigFileNotFoundError{}) {
		log.Errorf("config file err: %s ", err)
		return nil
	}

	return err

}

func setupMisc() error {
	// env vars
	err := setupEnv()
	if err != nil {
		return err
	}

	viper.AutomaticEnv()

	level, err := log.ParseLevel(viper.GetString("LOG_LEVEL"))
	if err != nil {
		return err
	}
	log.SetLevel(level)

	// tracing
	if err := initTracer(); err != nil {
		err := fmt.Errorf("failed to init tracer: %w", err)
		return err
	}
	log.Infof("tracer initialized")

	if err := sentry.Init(sentry.ClientOptions{
		Dsn: viper.GetString("SENTRY_DSN"),
	}); err != nil {
		return fmt.Errorf("sentry.Init: %w", err)
	}
	log.Infof("sentry initialized")
	return nil
}
