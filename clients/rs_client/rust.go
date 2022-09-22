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
	// ParseIngredient parseMethod = "parse"
	Scrape       parseMethod = "scrape"
	ParseAmount  parseMethod = "parse_amount"
	RecipeDecode parseMethod = "decode_recipe"
)

type Client struct {
	baseurl string
}

func New(baseurl string) *Client {
	return &Client{baseurl}
}
func (c *Client) Call(ctx context.Context, text string, kind parseMethod, target interface{}) error {
	ctx, span := otel.Tracer("rs_client").Start(ctx, "Parse")
	defer span.End()
	url := fmt.Sprintf("%s%s?text=%s", c.baseurl, kind, url.QueryEscape(text))

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Add("User-Agent", "gourd")

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	res, err := client.Do(req)
	log.WithError(err).WithField("kind", kind).Debugf("rs: parsed %s", text)

	if err != nil {
		return fmt.Errorf("rs Call to %s failed: %w", url, err)
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusBadRequest || res.StatusCode == http.StatusInternalServerError {
		return nil
	}
	return json.NewDecoder(res.Body).Decode(target)

}

func (c *Client) ConvertUnit(ctx context.Context, body, target interface{}) error {
	ctx, span := otel.Tracer("rs_client").Start(ctx, "Convert")
	defer span.End()

	return c.Post(ctx, "convert", body, target)
}

func (c *Client) Post(ctx context.Context, route string, body, target interface{}) error {

	route = c.baseurl + route

	ctx, span := otel.Tracer("rs_client").Start(ctx, "post")
	defer span.End()
	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(body)
	if err != nil {
		return err
	}

	span.AddEvent("sending", trace.WithAttributes(attribute.String("body", spew.Sdump(body))))
	defer func() {
		span.AddEvent("got", trace.WithAttributes(attribute.String("target", spew.Sdump(target))))
	}()

	req, _ := http.NewRequestWithContext(ctx, "POST", route, payloadBuf)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "gourd")

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(target)

}