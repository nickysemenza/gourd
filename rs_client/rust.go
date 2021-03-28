package rs_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type parseMethod string

const (
	Ingredient parseMethod = "parse"
	Amount     parseMethod = "parse_amount"
)

func Parse(ctx context.Context, text string, kind parseMethod, target interface{}) error {
	ctx, span := otel.Tracer("rs_client").Start(ctx, "Parse")
	defer span.End()
	url := fmt.Sprintf("http://localhost:8080/%s?text=%s", kind, url.QueryEscape(text))

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(target)

}
