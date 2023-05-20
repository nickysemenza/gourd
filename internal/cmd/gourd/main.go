package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	defer cleanupTracer()

	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		log.Error(err)
	}
}

func cleanupTracer() {
	log.Debug("cleaning up tracer")

	err := otel.GetTracerProvider().(*tracesdk.TracerProvider).ForceFlush(context.Background())
	if err != nil {
		log.Error(err)
	}
}
