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

// Client is a notio nclient
type Client struct {
	dbID  notionapi.DatabaseID
	block notionapi.BlockService
	db    notionapi.DatabaseService
	page  notionapi.PageService
}

// New makes a notion client
func New(token, database string) *Client {
	hClient := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	client := notionapi.NewClient(notionapi.Token(token), notionapi.WithHTTPClient(&hClient))
	return &Client{
		dbID:  notionapi.DatabaseID(database),
		db:    client.Database,
		block: client.Block,
		page:  client.Page,
	}
}

// Photo holds a photo from notion
type Photo struct {
	BlockID string `json:"block_id,omitempty"`
	URL     string `json:"url,omitempty"`
}

// Recipe holds notion recipe page
type Recipe struct {
	Title     string     `json:"title,omitempty"`
	Time      *time.Time `json:"time,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
	Photos    []Photo    `json:"photos,omitempty"`
	PageID    string     `json:"page_id,omitempty"`
	NotionURL string     `json:"notion_url,omitempty"`
	SourceURL string     `json:"source_url,omitempty"`
	Raw       string     `json:"raw,omitempty"`
	Children  []Recipe   `json:"children,omitempty"`
	Debug     []string   `json:"debug,omitempty"`
	Scale     *float64   `json:"scale,omitempty"`
}

func (r *Recipe) addDebug(l *log.Entry, format string, args ...interface{}) {
	l.Infof(format, args...)
	r.Debug = append(r.Debug, fmt.Sprintf(format, args...))
}

func removeDuplicateValues(intSlice []Photo) []Photo {
	keys := make(map[string]bool)
	list := []Photo{}

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
func (c *Client) processPage(ctx context.Context, page notionapi.Page) (recipe *Recipe, err error) {

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
		r := Recipe{
			Title:     name,
			PageID:    page.ID.String(),
			NotionURL: page.URL,
		}
		l := log.WithField("page_id", r.PageID)
		r.addDebug(l, "processing page %s (%s)", r.PageID, page.Object)
		l.Info(r.Title)

		if dateProp, ok := page.Properties["Date"]; ok {
			date := dateProp.(*notionapi.DateProperty).Date.Start
			if date != nil {
				utcTime := time.Time(*date)
				dinnerTime := utcTime.Add(time.Hour * (3 + 24))
				r.Time = zero.TimeFrom(dinnerTime).Ptr()
			}
		}
		if tags, ok := page.Properties["Tags"]; ok {
			for _, ms := range tags.(*notionapi.MultiSelectProperty).MultiSelect {
				r.Tags = append(r.Tags, ms.Name)
			}
		}
		if url, ok := page.Properties["source"]; ok {
			r.SourceURL = url.(*notionapi.URLProperty).URL
		}
		if num, ok := page.Properties["scale"]; ok && num != nil {
			raw := num.(*notionapi.NumberProperty).Number
			if raw != 0 {
				r.Scale = &raw
			}
		}

		// on each page, get all the blocks that are images
		return c.detailsFromPage(ctx, page.ID, r)

	default:
		return nil, fmt.Errorf("unknown page type %s", page.Object)
	}
}

// PageByID gets page by id
func (c *Client) PageByID(ctx context.Context, id notionapi.PageID) (*Recipe, error) {
	ctx, span := otel.Tracer("notion").Start(ctx, "PageById")
	defer span.End()

	page, err := c.page.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return c.processPage(ctx, *page)
	// return page, nil
}

// GetAll looks back the specified number of days. shortcut to filter by a single page id
func (c *Client) GetAll(ctx context.Context, lookback time.Duration, pageID string) ([]Recipe, error) {
	ctx, span := otel.Tracer("notion").Start(ctx, "GetAll")
	defer span.End()

	var cursor notionapi.Cursor
	recipes := []Recipe{}

	daysAgo := notionapi.Date(time.Now().Add(-lookback))

	filter := notionapi.AndCompoundFilter{
		notionapi.PropertyFilter{
			Property:    "Tags",
			MultiSelect: &notionapi.MultiSelectFilterCondition{DoesNotContain: "dining"},
		},
		notionapi.PropertyFilter{
			Property: "Date",
			Date:     &notionapi.DateFilterCondition{OnOrAfter: &daysAgo},
		},
	}

	for {
		req := &notionapi.DatabaseQueryRequest{
			Filter:      &filter,
			PageSize:    100,
			StartCursor: cursor,
		}
		resp, err := c.db.Query(ctx, c.dbID, req)
		if err != nil {
			return nil, err
		}
		span.AddEvent("data", trace.WithAttributes(
			attribute.String("data", spew.Sdump(resp)),
			attribute.String("req", spew.Sdump(req)),
		))

		for _, page := range resp.Results {
			span.AddEvent("page", trace.WithAttributes(attribute.String("page", spew.Sdump(page))))
			if pageID != "" && page.ID.String() != pageID {
				continue
			}
			recipe, err := c.processPage(ctx, page)
			if err != nil {
				return nil, err
			}
			if recipe != nil {
				recipes = append(recipes, *recipe)
			}
		}
		cursor = resp.NextCursor
		if !resp.HasMore {
			return recipes, nil
		}
	}

}

func (c *Client) detailsFromPage(ctx context.Context, pageID notionapi.ObjectID, r Recipe) (*Recipe, error) {
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
			r.addDebug(l, "\tfound notion %s", block.GetType())
			switch block.GetType() {
			case notionapi.BlockTypeImage:
				i := block.(*notionapi.ImageBlock)
				r.addDebug(l, "found an image")
				r.Photos = append(r.Photos, Photo{URL: i.Image.File.URL, BlockID: blockID})
			case notionapi.BlockTypeColumnList, notionapi.BlockTypeColumn:
				r.addDebug(l, "will get details from %s", block.GetType())
				foo, err := c.detailsFromPage(ctx, notionapi.ObjectID(blockID), r)
				if err != nil {
					return nil, err
				}
				r.Photos = append(r.Photos, foo.Photos...)
			case notionapi.BlockTypeCode:
				i := block.(*notionapi.CodeBlock)
				if len(i.Code.RichText) > 0 {
					if text := i.Code.RichText[0].Text.Content; strings.HasPrefix(text, "name:") {
						r.addDebug(l, "found code block")
						r.Raw = text
					}
				} else {
					r.addDebug(l, "found code block with no text")
				}
			case notionapi.BlockTypeChildPage:
				// treat as top level page?
				r.addDebug(l, "found child page")
				i := block.(*notionapi.ChildPageBlock)
				foo, err := c.PageByID(ctx, notionapi.PageID(i.ID))
				if err != nil {
					return nil, err
				}
				if foo != nil {
					r.Children = append(r.Children, *foo)
				}
			case notionapi.BlockTypeLinkToPage:
				i := block.(*notionapi.LinkToPageBlock)
				linkedTo, err := c.PageByID(ctx, i.LinkToPage.PageID)
				if err != nil {
					return nil, err
				}
				if linkedTo != nil {
					r.addDebug(l, "linked to page %s, using that as raw/source", linkedTo.PageID)
					r.Raw = linkedTo.Raw
					r.Children = append(r.Children, linkedTo.Children...)
					r.SourceURL = linkedTo.SourceURL
				}
			case notionapi.BlockCallout, notionapi.BlockQuote, notionapi.BlockTypeBookmark,
				notionapi.BlockTypeBreadcrumb, notionapi.BlockTypeBulletedListItem, notionapi.BlockTypeChildDatabase,
				notionapi.BlockTypeDivider, notionapi.BlockTypeEmbed, notionapi.BlockTypeEquation, notionapi.BlockTypeFile,
				notionapi.BlockTypeHeading1, notionapi.BlockTypeHeading2, notionapi.BlockTypeHeading3, notionapi.BlockTypeLinkPreview,
				notionapi.BlockTypeNumberedListItem, notionapi.BlockTypeParagraph, notionapi.BlockTypePdf,
				notionapi.BlockTypeSyncedBlock, notionapi.BlockTypeTableBlock, notionapi.BlockTypeTableOfContents,
				notionapi.BlockTypeTableRowBlock, notionapi.BlockTypeTemplate, notionapi.BlockTypeToDo,
				notionapi.BlockTypeToggle, notionapi.BlockTypeUnsupported, notionapi.BlockTypeVideo:
				// not supported
			}
		}

		childCursor = notionapi.Cursor(children.NextCursor)
		if !children.HasMore {
			break
		}
	}

	r.Photos = removeDuplicateValues(r.Photos)
	return &r, nil

}
