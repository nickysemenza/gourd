//go:build integration
// +build integration

package api

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func mustInsert(t *testing.T, a *API, cr CompactRecipe) string {
	t.Helper()
	ctx := context.Background()

	r, err := a.RecipeFromCompact(ctx, cr)
	require.NoError(t, err)

	ids, err := a.CreateRecipeDetails(ctx, r.Detail)
	require.NoError(t, err)
	return ids[0]
}
func newCompact(name string, ingredients, instructions []string) CompactRecipe {
	return CompactRecipe{
		Meta: CompactRecipeMeta{Name: name},
		Sections: []CompactRecipeSection{{
			Ingredients:  ingredients,
			Instructions: instructions,
		}},
	}
}
func TestUsage(t *testing.T) {
	require := require.New(t)
	_, apiManager := makeHandler(t)
	ctx := context.Background()

	mustInsert(t, apiManager, newCompact("sub", []string{
		"1 tsp salt",
		"1 gram salt",
		"1 gram sugar",
		"1 tsp pepper",
		"1 gram pepper",
		"1 tsp common",
	}, []string{}))

	rdMain := mustInsert(t, apiManager, newCompact("main", []string{
		"1 tsp pepper",
		"1 gram pepper",
		"1 recipe sub",
		"1 gram common"},
		[]string{}))
	rdSMallMain := mustInsert(t, apiManager, newCompact("smallmain", []string{
		"0.5 recipe sub",
		"1 gram common",
	}, []string{}))

	require.NoError(apiManager.loadIngredientMappings(ctx, []IngredientMapping{
		{Name: "salt", UnitMappings: []UnitMapping{{
			A: Amount{Value: 1, Unit: "tsp"},
			B: Amount{Value: 3, Unit: "g"},
		}}},
		{Name: "common", UnitMappings: []UnitMapping{{
			A: Amount{Value: 1, Unit: "tsp"},
			B: Amount{Value: 1, Unit: "g"},
		}}},
	}))

	{
		res, err := apiManager.IngredientUsage(ctx, []EntitySummary{{Id: rdMain, Multiplier: 2, Kind: IngredientKindRecipe}})
		require.NoError(err)

		require.Equal(2.0, firstAmount(res[ingID(t, apiManager, "sugar")].Sum, true).Value, "sugar should be double the original")
		require.Equal(4.0, firstAmount(res[ingID(t, apiManager, "common")].Sum, true).Value, "common should be double the original")

		ingSalt := ingID(t, apiManager, "salt")
		// require.Len(res[ingSalt].Sum, 1)
		saltGrams := firstAmount(res[ingSalt].Sum, true)
		require.Equal(8.0, saltGrams.Value, spew.Sdump(saltGrams))
		require.Equal("salt", res[ingSalt].Meta.Name)
		require.Equal(string(ingSalt), res[ingSalt].Meta.Id)
	}
	{
		res, err := apiManager.IngredientUsage(ctx, []EntitySummary{{Id: rdSMallMain, Kind: IngredientKindRecipe}})
		require.NoError(err)

		require.Equal(0.5, firstAmount(res[ingID(t, apiManager, "sugar")].Sum, true).Value, "sugar should be double the original")
		require.Equal(1.5, firstAmount(res[ingID(t, apiManager, "common")].Sum, true).Value, "common should be double the original")
		spew.Dump(res)

	}
}

func ingID(t *testing.T, apiManager *API, name string) IngredientID {
	t.Helper()
	ing, err := apiManager.DB().IngredientByName(context.TODO(), name)
	require.NoError(t, err)
	return IngredientID(ing.Id)
}
