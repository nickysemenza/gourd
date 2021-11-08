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

	// var filter *notionapi.PropertyFilter
	// if false {
	filter := &notionapi.PropertyFilter{
		Property:    "Tags",
		MultiSelect: &notionapi.MultiSelectFilterCondition{DoesNotContain: "dining"},
	}
	// }
	for {
		resp, err := c.client.Database.Query(ctx, c.database, &notionapi.DatabaseQueryRequest{
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
				meal := NotionRecipe{
					Title:     page.Properties["Name"].(*notionapi.TitleProperty).Title[0].Text.Content,
					PageID:    page.ID.String(),
					NotionURL: page.URL,
				}
				date := page.Properties["Date"].(*notionapi.DateProperty).Date.Start
				if date != nil {
					meal.Time = zero.TimeFrom(time.Time(*date).Add(time.Hour * 9)).Ptr()
				}
				for _, ms := range page.Properties["Tags"].(*notionapi.MultiSelectProperty).MultiSelect {
					meal.Tags = append(meal.Tags, ms.Name)
				}
				meal.SourceURL = page.Properties["source"].(*notionapi.URLProperty).URL

				// on each page, get all the blocks that are images
				meal.Photos, meal.Raw, err = c.ImagesFromPage(ctx, page.ID)
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

func (c *Client) ImageFromBlock(ctx context.Context, blockID notionapi.BlockID) (NotionPhoto, error) {
	ctx, span := otel.Tracer("notion").Start(ctx, "ImageFromBlock")
	defer span.End()

	block, err := c.client.Block.Get(ctx, blockID)
	if err != nil {
		return NotionPhoto{}, err
	}
	if block.GetType() != notionapi.BlockTypeImage {
		return NotionPhoto{}, fmt.Errorf("block %s is not an image", blockID)
	}
	span.AddEvent("block", trace.WithAttributes(attribute.String("block", spew.Sdump(block))))
	i := block.(*notionapi.ImageBlock)
	return NotionPhoto{
		BlockID: blockID.String(),
		URL:     i.Image.File.URL,
	}, nil
}

func (c *Client) ImagesFromPage(ctx context.Context, pageID notionapi.ObjectID) (images []NotionPhoto, raw string, err error) {
	ctx, span := otel.Tracer("notion").Start(ctx, "imagesFromPage")
	defer span.End()

	var childCursor notionapi.Cursor
	for {
		children, err := c.client.Block.GetChildren(ctx, notionapi.BlockID(pageID), &notionapi.Pagination{PageSize: 100, StartCursor: childCursor})
		if err != nil {
			return nil, "", err
		}

		for _, block := range children.Results {
			span.AddEvent("block", trace.WithAttributes(attribute.String("block", spew.Sdump(block))))
			log.Info(block.GetType())
			switch block.GetType() {
			case notionapi.BlockTypeImage:
				i := block.(*notionapi.ImageBlock)
				images = append(images, NotionPhoto{URL: i.Image.File.URL, BlockID: i.ID.String()})
			case "column", "column_list":
				i := block.(*notionapi.UnsupportedBlock)
				images2, _, err := c.ImagesFromPage(ctx, notionapi.ObjectID(i.ID))
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
