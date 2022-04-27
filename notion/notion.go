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
	dbId     notionapi.DatabaseID
	block    notionapi.BlockService
	db       notionapi.DatabaseService
	page     notionapi.PageService
	testOnly bool
}

func New(token, database string, testOnly bool) *Client {
	hClient := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	client := notionapi.NewClient(notionapi.Token(token), notionapi.WithHTTPClient(&hClient))
	return &Client{
		// client:   client,
		dbId:     notionapi.DatabaseID(database),
		db:       client.Database,
		block:    client.Block,
		page:     client.Page,
		testOnly: testOnly,
	}
}

type NotionPhoto struct {
	BlockID string `json:"block_id,omitempty"`
	URL     string `json:"url,omitempty"`
}
type NotionRecipe struct {
	Title     string         `json:"title,omitempty"`
	Time      *time.Time     `json:"time,omitempty"`
	Tags      []string       `json:"tags,omitempty"`
	Photos    []NotionPhoto  `json:"photos,omitempty"`
	PageID    string         `json:"page_id,omitempty"`
	NotionURL string         `json:"notion_url,omitempty"`
	SourceURL string         `json:"source_url,omitempty"`
	Raw       string         `json:"raw,omitempty"`
	Children  []NotionRecipe `json:"children,omitempty"`
}

func removeDuplicateValues(intSlice []NotionPhoto) []NotionPhoto {
	keys := make(map[string]bool)
	list := []NotionPhoto{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry.BlockID]; !value {
			keys[entry.BlockID] = true
			list = append(list, entry)
		}
	}
	return list
}

// nolint: exhaustive
func (c *Client) processPage(ctx context.Context, page notionapi.Page) (recipe *NotionRecipe, err error) {
	switch page.Object {
	case "column_list", notionapi.ObjectTypePage:
		var name string
		var titlePropName string
		switch page.Parent.Type {
		case "database_id":
			titlePropName = "Name"
		case "page_id":
			// sub page
			titlePropName = "title"
		default:
			return nil, fmt.Errorf("unknown parent type %s", page.Parent.Type)
		}

		nameProp, nameOk := page.Properties[titlePropName]
		if nameOk && len(nameProp.(*notionapi.TitleProperty).Title) != 1 {
			err = fmt.Errorf("page %s has no title", page.ID)
			log.Error(err)
			return nil, nil
		}
		name = nameProp.(*notionapi.TitleProperty).Title[0].Text.Content
		meal := NotionRecipe{
			Title:     name,
			PageID:    page.ID.String(),
			NotionURL: page.URL,
		}
		log.WithField("page_id", meal.PageID).Info(meal.Title)

		if dateProp, ok := page.Properties["Date"]; ok {
			date := dateProp.(*notionapi.DateProperty).Date.Start
			if date != nil {
				utcTime := time.Time(*date)
				dinnerTime := utcTime.Add(time.Hour * (3 + 24))
				meal.Time = zero.TimeFrom(dinnerTime).Ptr()
			}
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
		return c.detailsFromPage(ctx, page.ID, meal)

	default:
		return nil, fmt.Errorf("unknown page type %s", page.Object)
	}
}
func (c *Client) PageById(ctx context.Context, id notionapi.PageID) (*NotionRecipe, error) {
	ctx, span := otel.Tracer("notion").Start(ctx, "PageById")
	defer span.End()

	page, err := c.page.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return c.processPage(ctx, *page)
	// return page, nil
}
func (c *Client) GetAll(ctx context.Context, lookbackDays int) ([]NotionRecipe, error) {
	ctx, span := otel.Tracer("notion").Start(ctx, "Dump")
	defer span.End()

	var cursor notionapi.Cursor
	meals := []NotionRecipe{}

	filter := notionapi.CompoundFilter{
		notionapi.FilterOperatorAND: {
			notionapi.PropertyFilter{
				Property:    "Tags",
				MultiSelect: &notionapi.MultiSelectFilterCondition{DoesNotContain: "dining"},
			},
		},
	}

	daysAgo := notionapi.Date(time.Now().AddDate(0, 0, -lookbackDays))

	testFilter := notionapi.PropertyFilter{
		Property: "Date",
		Date:     &notionapi.DateFilterCondition{OnOrAfter: &daysAgo},
	}

	filter["and"] = append(filter["and"], testFilter)
	for {
		resp, err := c.db.Query(ctx, c.dbId, &notionapi.DatabaseQueryRequest{
			CompoundFilter: &filter,
			PageSize:       100,
			StartCursor:    cursor,
		})
		if err != nil {
			return nil, err
		}

		for _, page := range resp.Results {
			span.AddEvent("page", trace.WithAttributes(attribute.String("page", spew.Sdump(page))))
			// if page.ID.String() != "e19dd5c9-5894-4bbc-b048-76e4293e32c0" {
			// 	continue
			// }
			meal, err := c.processPage(ctx, page)
			if err != nil {
				return nil, err
			}
			if meal != nil {
				meals = append(meals, *meal)
			}
		}
		cursor = resp.NextCursor
		if !resp.HasMore {
			return meals, nil
		}
	}

}

func (c *Client) detailsFromPage(ctx context.Context, pageID notionapi.ObjectID, meal NotionRecipe) (*NotionRecipe, error) {
	ctx, span := otel.Tracer("notion").Start(ctx, "detailsFromPage")
	defer span.End()

	var childCursor notionapi.Cursor
	for {
		children, err := c.block.GetChildren(ctx, notionapi.BlockID(pageID), &notionapi.Pagination{PageSize: 100, StartCursor: childCursor})
		if err != nil {
			return nil, err
		}

		for _, block := range children.Results {
			blockID := block.GetID().String()
			span.AddEvent("block", trace.WithAttributes(attribute.String("block", spew.Sdump(block))))
			l := log.WithField("page_id", pageID).
				WithField("block_id", blockID)
			l.
				Infof("\tfound notion %s", block.GetType())
			// nolint: exhaustive
			switch block.GetType() {
			case notionapi.BlockTypeImage:
				i := block.(*notionapi.ImageBlock)
				l.Infof("appending via image")
				meal.Photos = append(meal.Photos, NotionPhoto{URL: i.Image.File.URL, BlockID: blockID})
			case notionapi.BlockTypeColumnList, notionapi.BlockTypeColumn:
				foo, err := c.detailsFromPage(ctx, notionapi.ObjectID(blockID), meal)
				if err != nil {
					return nil, err
				}
				l.Infof("appending via col")
				meal.Photos = append(meal.Photos, foo.Photos...)
			case notionapi.BlockTypeCode:
				i := block.(*notionapi.CodeBlock)
				if len(i.Code.Text) > 0 {
					if text := i.Code.Text[0].Text.Content; strings.HasPrefix(text, "name:") {
						meal.Raw = text
					}
				} else {
					l.Errorf("found code block with no text")
					spew.Dump(i)
				}
			case notionapi.BlockTypeChildPage:
				// treat as top level page?

				i := block.(*notionapi.ChildPageBlock)
				foo, err := c.PageById(ctx, notionapi.PageID(i.ID))
				if err != nil {
					return nil, err
				}
				if foo != nil {
					meal.Children = append(meal.Children, *foo)
				}
			}
		}

		childCursor = notionapi.Cursor(children.NextCursor)
		if !children.HasMore {
			break
		}
	}

	meal.Photos = removeDuplicateValues(meal.Photos)
	return &meal, nil

}
