package main

import (
	"errors"
	"fmt"

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

func initTracer(url, name, env string) error {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
			attribute.String("environment", env),
		)),
	)
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		jaegerp.Jaeger{},
	))
	return nil
}

func setupEnv() error {
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("DATABASE_URL", "postgres://gourd:gourd@localhost:5555/food")
	viper.SetDefault("DATABASE_URL_USDA", "postgres://gourd:gourd@localhost:5556/usda")
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
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if errors.As(err, &viper.ConfigFileNotFoundError{}) {
		log.Errorf("config file err: %s ", err)
		return nil
	}

	return err

}

func setupMisc() error {
	if err := setupEnv(); err != nil {
		return err
	}

	level, err := log.ParseLevel(viper.GetString("LOG_LEVEL"))
	if err != nil {
		return err
	}
	log.SetLevel(level)
	log.SetReportCaller(true)

	// tracing
	if err := initTracer(viper.GetString("JAEGER_ENDPOINT"), "gourd", "dev"); err != nil {
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
