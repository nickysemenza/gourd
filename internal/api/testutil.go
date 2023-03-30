package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// ExtractNames pulls names out of ingredient detail
func ExtractNames(inp []IngredientDetail) (res []string) {
	for _, x := range inp {
		res = append(res, x.Ingredient.Name)
	}
	return
}

// MustInsert inserts a CompactRecipe
func MustInsert(t *testing.T, a *API, cr CompactRecipe) string {
	t.Helper()
	ctx := context.Background()

	r, err := a.RecipeFromCompact(ctx, cr)
	require.NoError(t, err)

	ids, err := a.CreateRecipeDetails(ctx, r.Detail)
	require.NoError(t, err)
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
	ing, err := apiManager.DB().IngredientByName(context.TODO(), name)
	require.NoError(t, err)
	return IngredientID(ing.Id)
}
