package notion

import (
	"context"
	"fmt"

	"github.com/kjk/notionapi"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
)

// GetPage fetches a snippet of a notion page.
func GetPage(ctx context.Context, pageID string) (*Page, error) {
	tr := global.Tracer("notion")
	_, span := tr.Start(ctx, "notion: GetPage")
	span.SetAttributes(core.Key("page_id").String(pageID))
	defer span.End()
	client := &notionapi.Client{}
	page, err := client.DownloadPage(pageID)
	if err != nil {
		return nil, fmt.Errorf("failed to download notion page: %w", err)
	}
	imageUrls := []string{}
	page.ForEachBlock(func(b *notionapi.Block) {
		if b.Type == "image" {
			imageUrls = append(imageUrls, b.ImageURL)
		}
	})
	return &Page{Title: page.Root().Title, ImageURLs: imageUrls}, nil
}

// Page is a snippet.
type Page struct {
	Title     string
	ImageURLs []string
}
