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
	"github.com/nickysemenza/gourd/graph/model"
	"github.com/nickysemenza/gourd/parser"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/plugin/httptrace"
)

// FetchAndTransform returns a recipe.
func FetchAndTransform(ctx context.Context, url string, ingredientToUUID func(ctx context.Context, name string, kind model.SectionIngredientKind) (string, error)) (*model.RecipeInput, error) {
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

	section := &model.SectionInput{}

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
		uuid, err := ingredientToUUID(ctx, i.Name, model.SectionIngredientKindIngredient)
		if err != nil {
			return nil, fmt.Errorf("failed to map ingredient %s to uuid: %w", i.Name, err)
		}
		spew.Dump(uuid)
		section.Ingredients = append(section.Ingredients, &model.SectionIngredientInput{
			InfoUUID:  uuid,
			Kind:      model.SectionIngredientKindIngredient,
			Amount:    &i.Volume.Value,
			Unit:      &i.Volume.Unit,
			Adjective: &i.Modifier,
			Grams:     i.Grams(),
		})
	}
	for _, item := range recipe.RecipeInstructions {
		section.Instructions = append(section.Instructions, &model.SectionInstructionInput{Instruction: item.Text})
	}

	r := model.RecipeInput{
		Name:     recipe.Name,
		Sections: []*model.SectionInput{section},
	}
	spew.Dump(r)

	return &r, nil
}
func extractRecipeJSONLD(ctx context.Context, html string) (*Recipe, error) {
	_, span := global.Tracer("scraper").Start(ctx, "scraper.extractRecipeJSONLD")
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
			strings.Replace(s.Text(), "\n", "", -1),
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
