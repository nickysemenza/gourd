// Package scraper retrieves recipes from websites containing it in json+ld format.
// Compatable website include cooking.nytimes.com, seriouseats.com
package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"
	"github.com/nickysemenza/gourd/api"
	"github.com/nickysemenza/gourd/parser"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel"
)

// FetchAndTransform returns a recipe.
func FetchAndTransform(ctx context.Context, addr string, ingredientToUUID func(ctx context.Context, name string, kind string) (string, error)) (*api.RecipeWrapper, error) {
	ctx, span := otel.Tracer("scraper").Start(ctx, "scraper.GetIngredients")
	defer span.End()
	html, err := getHTML(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape %s: %w", addr, err)
	}
	recipe, err := extractRecipeJSONLD(ctx, html)
	if err != nil {
		return nil, fmt.Errorf("failed to extract ld+json from %s: %w", addr, err)
	}
	spew.Dump(recipe.RecipeIngredient, err)
	debugMsgs := []string{}

	section := api.RecipeSection{}

	for _, item := range recipe.RecipeIngredient {
		i, err := parser.Parse(ctx, item)
		if err != nil {
			msg := fmt.Sprintf("failed to parse: %s", err)
			log.Errorf(msg)
			debugMsgs = append(debugMsgs, msg)
			continue
		}
		debugMsgs = append(debugMsgs, i.ToString())
		fmt.Println(i.ToString())
		uuid, err := ingredientToUUID(ctx, i.Name, "ingredient")
		if err != nil {
			return nil, fmt.Errorf("failed to map ingredient %s to uuid: %w", i.Name, err)
		}
		spew.Dump(uuid)
		section.Ingredients = append(section.Ingredients, api.SectionIngredient{
			Ingredient: &api.Ingredient{Id: uuid},
			Kind:       "ingredient",
			Amount:     &i.Volume.Value,
			Unit:       i.Volume.Unit,
			Adjective:  i.Modifier,
			Grams:      i.Grams(),
		})
	}
	for _, item := range recipe.RecipeInstructions {
		section.Instructions = append(section.Instructions, api.SectionInstruction{Instruction: item.Text})
	}
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	source := fmt.Sprintf("todo: %s %s", u.Host, addr)
	r := api.RecipeWrapper{
		Detail: api.Recipe{
			Name:   recipe.Name,
			Source: &source,
			// Source: &api.Source{
			// 	Name: u.Host,
			// 	Meta: addr,
			// },
		},
		Sections: []api.RecipeSection{section},
	}
	spew.Dump(r, debugMsgs)

	return &r, nil
}
func extractRecipeJSONLD(ctx context.Context, html string) (*Recipe, error) {
	_, span := otel.Tracer("scraper").Start(ctx, "scraper.extractRecipeJSONLD")
	defer span.End()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("goquery parse: %w", err)
	}

	recipe := Recipe{}
	// Find the recipe items
	doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
		var r Recipe
		err = json.Unmarshal([]byte(
			strings.ReplaceAll(s.Text(), "\n", ""),
		), &r)
		if err == nil && r.Type == "Recipe" {
			recipe = r
		} else if err != nil {
			err = fmt.Errorf("failed to parse ld+json (%s): %w", s.Text(), err)
		}
	})
	return &recipe, err
}
func getHTML(ctx context.Context, url string) (string, error) {
	ctx, span := otel.Tracer("scraper").Start(ctx, "scraper.getHTML")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	ctx, req = otelhttptrace.W3C(ctx, req)
	otelhttptrace.Inject(ctx, req)

	// nolint:gosec
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err := fmt.Errorf("failed to get html: %w", err)
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got status code: %d %s while downloading %s", res.StatusCode, res.Status, url)
	}
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(res.Body); err != nil {
		return "", err
	}
	return buf.String(), nil
}
