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
	ctx, span := otel.Tracer("rs_client").Start(ctx, "Call")
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

func (c *Client) Send(ctx context.Context, route string, body, target interface{}) error {

	route = c.baseurl + route

	ctx, span := otel.Tracer("rs_client").Start(ctx, "send")
	defer span.End()

	span.AddEvent("sending", trace.WithAttributes(attribute.String("body", spew.Sdump(body))))
	defer func() {
		span.AddEvent("got", trace.WithAttributes(attribute.String("target", spew.Sdump(target))))
	}()
	var req *http.Request
	var err error
	if body == nil {
		req, err = http.NewRequestWithContext(ctx, "GET", route, nil)
	} else {
		payloadBuf := new(bytes.Buffer)
		err = json.NewEncoder(payloadBuf).Encode(body)
		if err != nil {
			return err
		}
		req, err = http.NewRequestWithContext(ctx, "POST", route, payloadBuf)
	}
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "gourd")

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil
	}

	if target == nil {
		return nil
	}
	return json.NewDecoder(res.Body).Decode(target)

}
