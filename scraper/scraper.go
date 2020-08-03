// Package scraper retrieves recipes from websites containing it in json+ld format.
// Compatable website include cooking.nytimes.com, seriouseats.com
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
	"github.com/nickysemenza/gourd/parser"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/plugin/httptrace"
)

// GetIngredients returns ingredients before/after they have been run through the parser.
func GetIngredients(ctx context.Context, url string) (interface{}, error) {
	ctx, span := global.Tracer("scraper").Start(ctx, "scraper.GetIngredients")
	defer span.End()
	html, err := getHTML(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape %s: %w", url, err)
	}
	recipe, err := extractRecipeJSONLD(ctx, html)
	if err != nil {
		return nil, fmt.Errorf("failed to extract ld+json from %s: %w", url, err)
	}
	spew.Dump(recipe.RecipeIngredient, err)
	output := []string{}
	for _, item := range recipe.RecipeIngredient {
		i, err := parser.Parse(ctx, item)
		if err != nil {
			msg := fmt.Sprintf("failed to parse: %s", err)
			log.Errorf(msg)
			output = append(output, msg)
			continue
		}
		output = append(output, i.ToString())
		fmt.Println(i.ToString())
	}
	return map[string]interface{}{"input": recipe.RecipeIngredient, "output": output}, nil
}
func extractRecipeJSONLD(ctx context.Context, html string) (*Recipe, error) {
	_, span := global.Tracer("scraper").Start(ctx, "scraper.extractRecipeJSONLD")
	defer span.End()
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
	ctx, span := global.Tracer("scraper").Start(ctx, "scraper.getHTML")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	ctx, req = httptrace.W3C(ctx, req)
	httptrace.Inject(ctx, req)

	// nolint:gosec
	res, err := http.DefaultClient.Do(req)
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
