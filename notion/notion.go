package notion

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jomei/notionapi"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/guregu/null.v4/zero"
)

type Client struct {
	// client   *notionapi.Client
	dbId  notionapi.DatabaseID
	block notionapi.BlockService
	db    notionapi.DatabaseService
}

func New(token, database string) *Client {
	hClient := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	client := notionapi.NewClient(notionapi.Token(token), notionapi.WithHTTPClient(&hClient))
	return &Client{
		// client:   client,
		dbId:  notionapi.DatabaseID(database),
		db:    client.Database,
		block: client.Block,
	}
}

type NotionPhoto struct {
	BlockID string `json:"block_id,omitempty"`
	URL     string `json:"url,omitempty"`
}
type NotionRecipe struct {
	Title     string        `json:"title,omitempty"`
	Time      *time.Time    `json:"time,omitempty"`
	Tags      []string      `json:"tags,omitempty"`
	Photos    []NotionPhoto `json:"photos,omitempty"`
	PageID    string        `json:"page_id,omitempty"`
	NotionURL string        `json:"notion_url,omitempty"`
	SourceURL string        `json:"source_url,omitempty"`
	Raw       string        `json:"raw,omitempty"`
}

func (c *Client) Dump(ctx context.Context) ([]NotionRecipe, error) {
	ctx, span := otel.Tracer("notion").Start(ctx, "Dump")
	defer span.End()

	var cursor notionapi.Cursor
	meals := []NotionRecipe{}

	filter := &notionapi.PropertyFilter{
		Property:    "Tags",
		MultiSelect: &notionapi.MultiSelectFilterCondition{DoesNotContain: "dining"},
	}
	for {
		resp, err := c.db.Query(ctx, c.dbId, &notionapi.DatabaseQueryRequest{
			PropertyFilter: filter,
			PageSize:       100,
			StartCursor:    cursor,
		})

		if err != nil {
			return nil, err
		}
		for _, page := range resp.Results {
			span.AddEvent("page", trace.WithAttributes(attribute.String("page", spew.Sdump(page))))
			switch page.Object {
			case "column_list", notionapi.ObjectTypePage:
				if len(page.Properties["Name"].(*notionapi.TitleProperty).Title) != 1 {
					return nil, fmt.Errorf("page %s has no title", page.ID)
				}
				meal := NotionRecipe{
					Title:     page.Properties["Name"].(*notionapi.TitleProperty).Title[0].Text.Content,
					PageID:    page.ID.String(),
					NotionURL: page.URL,
				}
				log.WithField("page_id", meal.PageID).Info(meal.Title)
				date := page.Properties["Date"].(*notionapi.DateProperty).Date.Start
				if date != nil {
					utcTime := time.Time(*date) //todo: this is slightly wrong
					meal.Time = zero.TimeFrom(utcTime.Add(time.Hour * 9)).Ptr()
				}
				if tags, ok := page.Properties["Tags"]; ok {
					for _, ms := range tags.(*notionapi.MultiSelectProperty).MultiSelect {
						meal.Tags = append(meal.Tags, ms.Name)
					}
				}
				if url, ok := page.Properties["source"]; ok {
					meal.SourceURL = url.(*notionapi.URLProperty).URL
				}

				// on each page, get all the blocks that are images
				meal.Photos, meal.Raw, err = c.imagesFromPage(ctx, page.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to get images for page %s: %w", page.ID, err)
				}
				meals = append(meals, meal)

			}
		}
		cursor = resp.NextCursor
		if !resp.HasMore {
			return meals, nil
		}
	}

}

func (c *Client) imagesFromPage(ctx context.Context, pageID notionapi.ObjectID) (images []NotionPhoto, raw string, err error) {
	ctx, span := otel.Tracer("notion").Start(ctx, "imagesFromPage")
	defer span.End()

	var childCursor notionapi.Cursor
	for {
		children, err := c.block.GetChildren(ctx, notionapi.BlockID(pageID), &notionapi.Pagination{PageSize: 100, StartCursor: childCursor})
		if err != nil {
			return nil, "", err
		}

		for _, block := range children.Results {
			span.AddEvent("block", trace.WithAttributes(attribute.String("block", spew.Sdump(block))))
			log.WithField("page_id", pageID).Infof("\tfound notion %s", block.GetType())
			switch block.GetType() {
			case notionapi.BlockTypeImage:
				i := block.(*notionapi.ImageBlock)
				images = append(images, NotionPhoto{URL: i.Image.File.URL, BlockID: i.ID.String()})
			case "column", "column_list":
				i := block.(*notionapi.UnsupportedBlock)
				images2, _, err := c.imagesFromPage(ctx, notionapi.ObjectID(i.ID))
				if err != nil {
					return nil, "", err
				}
				images = append(images, images2...)
			case notionapi.BlockTypeCode:
				i := block.(*notionapi.CodeBlock)
				if text := i.Code.Text[0].Text.Content; strings.HasPrefix(text, "name:") {
					raw = text
				}

			}
		}

		childCursor = notionapi.Cursor(children.NextCursor)
		if !children.HasMore {
			break
		}
	}
	return

}
