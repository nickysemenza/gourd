// Package scraper retrieves recipes from websites containing it in json+ld format.
// Compatable website include cooking.nytimes.com, seriouseats.com
package scraper

import (
	"context"

	"github.com/nickysemenza/gourd/api"
	"github.com/nickysemenza/gourd/rs_client"
	"go.opentelemetry.io/otel"
)

// FetchAndTransform returns a recipe.
func FetchAndTransform(ctx context.Context, addr string, ingredientToId func(ctx context.Context, name string, kind string) (string, error)) (*api.RecipeWrapper, error) {
	ctx, span := otel.Tracer("scraper").Start(ctx, "scraper.GetIngredients")
	defer span.End()

	r := api.RecipeWrapper{}
	err := rs_client.Call(ctx, addr, rs_client.Scrape, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
