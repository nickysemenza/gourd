package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint: unused
func extractNames(inp []IngredientDetail) (res []string) {
	for _, x := range inp {
		res = append(res, x.Ingredient.Name)
	}
	return
}

//nolint: unused
func mustInsert(t *testing.T, a *API, cr CompactRecipe) string {
	t.Helper()
	ctx := context.Background()

	r, err := a.RecipeFromCompact(ctx, cr)
	require.NoError(t, err)

	ids, err := a.CreateRecipeDetails(ctx, r.Detail)
	require.NoError(t, err)
	return ids[0]
}

//nolint: unused
func newCompact(name string, ingredients, instructions []string) CompactRecipe {
	return CompactRecipe{
		Meta: CompactRecipeMeta{Name: name},
		Sections: []CompactRecipeSection{{
			Ingredients:  ingredients,
			Instructions: instructions,
		}},
	}
}

//nolint: unused
func ingID(t *testing.T, apiManager *API, name string) IngredientID {
	t.Helper()
	ing, err := apiManager.DB().IngredientByName(context.TODO(), name)
	require.NoError(t, err)
	return IngredientID(ing.Id)
}
