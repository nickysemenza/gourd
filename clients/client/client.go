package client

import (
	"net/http"

	"github.com/nickysemenza/gourd/api"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	*api.ClientWithResponses
}

func New(url string) (*Client, error) {
	c, err := api.NewClientWithResponses(url,
		api.WithHTTPClient(&http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}),
	)
	if err != nil {
		return nil, err
	}
	return &Client{c}, nil
}
