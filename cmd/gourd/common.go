package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/credentials"

	jaegerp "go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func initTracer(jaegerURL, honeycombKey, name, env string) error {
	// Create the Jaeger exporter

	var exp tracesdk.SpanExporter
	var err error
	if jaegerURL != "" {
		exp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerURL)))

	} else if honeycombKey != "" {
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint("api.honeycomb.io:443"),
			otlptracegrpc.WithHeaders(map[string]string{
				"x-honeycomb-team": honeycombKey,
			}),
			otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
		}
		exp, err = otlptrace.New(context.Background(), otlptracegrpc.NewClient(opts...))
	}
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// tracesdk.WithSpanProcessor(tracesdk.NewBatchSpanProcessor(exp)),
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

	viper.Debug()
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
	if err := initTracer(viper.GetString("JAEGER_ENDPOINT"), viper.GetString("HONEYCOMB_KEY"), "gourd", "dev"); err != nil {
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
