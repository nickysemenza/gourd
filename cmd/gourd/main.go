package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	if err := setupMisc(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Debug("cleaning up tracer")
		err := otel.GetTracerProvider().(*tracesdk.TracerProvider).ForceFlush(context.Background())
		if err != nil {
			log.Error(err)
		}
	}()
	ctx, span := otel.Tracer("client").Start(context.Background(), "gourd main")
	defer span.End()

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		log.Error(err)
	}
}
