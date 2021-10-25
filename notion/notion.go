package notion

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/jomei/notionapi"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	client   *notionapi.Client
	database notionapi.DatabaseID
}

func New(token, database string) *Client {
	hClient := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	return &Client{
		client:   notionapi.NewClient(notionapi.Token(token), notionapi.WithHTTPClient(&hClient)),
		database: notionapi.DatabaseID(database),
	}
}

func (c *Client) Dump(ctx context.Context) error {
	resp, err := c.client.Database.Query(ctx, c.database, &notionapi.DatabaseQueryRequest{PropertyFilter: &notionapi.PropertyFilter{Property: "Tags", MultiSelect: &notionapi.MultiSelectFilterCondition{Contains: "done"}}})
	if err != nil {
		return err
	}
	for _, res := range resp.Results {
		if res.Object == "page" {
			children, err := c.client.Block.GetChildren(context.Background(), notionapi.BlockID(res.ID), &notionapi.Pagination{PageSize: 100})
			if err != nil {
				return err
			}
			title := res.Properties["Name"].(*notionapi.TitleProperty).Title[0].Text.Content
			date := res.Properties["Date"].(*notionapi.DateProperty).Date.Start
			time := time.Time(*date)
			tags := []string{}
			for _, ms := range res.Properties["Tags"].(*notionapi.MultiSelectProperty).MultiSelect {
				tags = append(tags, ms.Name)
			}

			log.WithField("time", time).WithField("title", title).WithField("tags", strings.Join(tags, ",")).Info("checking")
			for _, block := range children.Results {
				log.Info(block.GetType())
				if block.GetType() == "image" {
					i := block.(*notionapi.ImageBlock)
					log.Infof("url: %s, page: %v", i.Image.File.URL, title)
				}
			}

		}

	}

	return nil

}
