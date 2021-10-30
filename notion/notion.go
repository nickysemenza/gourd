package notion

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jomei/notionapi"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
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

type NotionMeal struct {
	Title  string     `json:"title,omitempty"`
	Time   *time.Time `json:"time,omitempty"`
	Tags   []string   `json:"tags,omitempty"`
	Photos []string   `json:"photos,omitempty"`
	PageID string     `json:"page_id,omitempty"`
	URL    string     `json:"url,omitempty"`
}

func (c *Client) Dump(ctx context.Context) ([]NotionMeal, error) {
	ctx, span := otel.Tracer("notion").Start(ctx, "Dump")
	defer span.End()

	var cursor notionapi.Cursor
	meals := []NotionMeal{}
	for {
		resp, err := c.client.Database.Query(ctx, c.database,
			&notionapi.DatabaseQueryRequest{
				// PropertyFilter: &notionapi.PropertyFilter{
				// 	Property:    "Tags",
				// 	MultiSelect: &notionapi.MultiSelectFilterCondition{Contains: "done"},
				// },
				PageSize:    100,
				StartCursor: cursor,
			})
		if err != nil {
			return nil, err
		}
		for _, page := range resp.Results {
			if page.Object != notionapi.ObjectTypePage {
				continue
			}
			meal := NotionMeal{
				Title:  page.Properties["Name"].(*notionapi.TitleProperty).Title[0].Text.Content,
				PageID: page.ID.String(),
				URL:    page.URL,
			}
			date := page.Properties["Date"].(*notionapi.DateProperty).Date.Start
			if date != nil {
				meal.Time = zero.TimeFrom(time.Time(*date).Add(time.Hour * 9)).Ptr()
			}
			for _, ms := range page.Properties["Tags"].(*notionapi.MultiSelectProperty).MultiSelect {
				meal.Tags = append(meal.Tags, ms.Name)
			}

			// on each page, get all the blocks that are images
			meal.Photos, err = c.ImagesFromPage(ctx, page.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get images for page %s: %w", page.ID, err)
			}

			meals = append(meals, meal)

		}
		spew.Dump(resp.HasMore, resp.NextCursor, resp.Object)
		cursor = resp.NextCursor
		if !resp.HasMore {
			return meals, nil
		}
	}

}

func (c *Client) ImagesFromPage(ctx context.Context, pageID notionapi.ObjectID) ([]string, error) {
	ctx, span := otel.Tracer("notion").Start(ctx, "imagesFromPage")
	defer span.End()

	var images []string
	var childCursor notionapi.Cursor
	for {
		children, err := c.client.Block.GetChildren(ctx, notionapi.BlockID(pageID), &notionapi.Pagination{PageSize: 100, StartCursor: childCursor})
		if err != nil {
			return nil, err
		}

		for _, block := range children.Results {
			log.Info(block.GetType())
			if block.GetType() != notionapi.BlockTypeImage {
				continue
			}
			i := block.(*notionapi.ImageBlock)
			images = append(images, i.Image.File.URL)
		}

		childCursor = notionapi.Cursor(children.NextCursor)
		if !children.HasMore {
			break
		}
	}
	return images, nil

}
