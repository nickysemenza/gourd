package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

var DefaultPagination = parsePagination(nil, nil)

// ExtractNames pulls names out of ingredient detail
func ExtractNames(inp []IngredientDetail) (res []string) {
	for _, x := range inp {
		res = append(res, x.Ingredient.Name)
	}
	return
}

// MustInsert inserts a CompactRecipe
func MustInsert(ctx context.Context, t *testing.T, a *API, cr CompactRecipe) string {
	t.Helper()
	t.Logf("inserting compact %s %s", cr.Name, cr.Id)

	r, err := a.RecipeFromCompact(ctx, cr)
	require.NoError(t, err)

	ids, err := a.CreateRecipeDetails(ctx, r.Detail)
	require.NoError(t, err)
	require.NotEmpty(t, ids[0])
	return ids[0]
}

// NewCompact builds a compact recipe
func NewCompact(name string, ingredients, instructions []string) CompactRecipe {
	return CompactRecipe{
		Name: name,
		Sections: []CompactRecipeSection{{
			Ingredients:  ingredients,
			Instructions: instructions,
		}},
	}
}

// IngIDFromName turns name to id
func IngIDFromName(t *testing.T, apiManager *API, name string) IngredientID {
	t.Helper()
	ing, err := apiManager.ingredientByName(context.Background(), name)
	require.NoError(t, err)
	return IngredientID(ing.ID)
}
