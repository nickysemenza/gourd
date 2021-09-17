package rs_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type parseMethod string

const (
	Ingredient   parseMethod = "parse"
	Scrape       parseMethod = "scrape"
	Amount       parseMethod = "parse_amount"
	RecipeDecode parseMethod = "decode_recipe"
)

func Call(ctx context.Context, text string, kind parseMethod, target interface{}) error {
	ctx, span := otel.Tracer("rs_client").Start(ctx, "Parse")
	defer span.End()
	url := fmt.Sprintf("http://localhost:8080/%s?text=%s", kind, url.QueryEscape(text))

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("rs Call to %s failed: %w", url, err)
	}
	log.WithField("kind", kind).Debugf("rs: parsed %s", text)

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(target)

}

func Convert(ctx context.Context, body, target interface{}) error {
	ctx, span := otel.Tracer("rs_client").Start(ctx, "Convert")
	defer span.End()
	url := fmt.Sprint("http://localhost:8080/convert")

	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(body)
	if err != nil {
		return err
	}

	span.AddEvent("sending", trace.WithAttributes(attribute.String("recipe", spew.Sdump(body))))

	req, _ := http.NewRequestWithContext(ctx, "POST", url, payloadBuf)
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	log.Debugf("rs: converted %v", body)

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(target)

}
