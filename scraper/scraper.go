package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"
	"github.com/nickysemenza/food/parser"
	log "github.com/sirupsen/logrus"
)

func ExampleScrape(ctx context.Context, url string) error {
	html, err := getHTML(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to scrape %s: %w", url, err)
	}
	recipe, err := extractRecipeJSONLD(ctx, html)
	if err != nil {
		return fmt.Errorf("failed to extract ld+json from %s: %w", url, err)
	}
	spew.Dump(recipe.RecipeIngredient, err)
	for _, item := range recipe.RecipeIngredient {
		i, err := parser.Parse(ctx, item)
		if err != nil {
			log.Errorf("failed to parse: %s", err)
			continue
		}
		fmt.Println(i.ToString())
	}
	return nil
}
func extractRecipeJSONLD(ctx context.Context, html string) (*Recipe, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("goquery parse: %w", err)
	}

	recipe := Recipe{}
	// Find the review items
	doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
		var r Recipe
		err = json.Unmarshal([]byte(s.Text()), &r)
		if err == nil && r.Type == "Recipe" {
			recipe = r
		} else if err != nil {
			err = fmt.Errorf("failed to parse ld+json: %w", err)
		}
	})
	return &recipe, err
}
func getHTML(ctx context.Context, url string) (string, error) {
	// nolint:gosec
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("got status code: %d %s while downloading %s", res.StatusCode, res.Status, url)
	}
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(res.Body); err != nil {
		return "", err
	}
	return buf.String(), nil
}
