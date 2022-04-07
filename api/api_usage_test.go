//go:build integration
// +build integration

package api

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestUsage(t *testing.T) {
	require := require.New(t)
	_, apiManager := makeHandler(t)
	ctx := context.Background()

	r, err := apiManager.RecipeFromText(ctx, `
name: sub
---
1 tsp salt
1 gram salt
1 gram sugar
1 tsp common`)
	require.NoError(err)
	r2, err := apiManager.RecipeFromText(ctx, `
name: main
---
1 recipe sub
1 gram common`)
	require.NoError(err)
	r3, err := apiManager.RecipeFromText(ctx, `
name: smallmain
---
0.5 recipe sub
1 gram common`)
	require.NoError(err)
	ids, err := apiManager.CreateRecipeDetails(ctx, *r, *r2, *r3)
	require.NoError(err)
	require.Len(ids, 3)
	rdMain := ids[1]
	rdSMallMain := ids[2]

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
		require.Equal("salt", res[ingSalt].Ing.Name)
		require.Equal(string(ingSalt), res[ingSalt].Ing.Id)
	}
	{
		res, err := apiManager.IngredientUsage(ctx, []EntitySummary{{Id: rdSMallMain, Kind: IngredientKindRecipe}})
		require.NoError(err)

		require.Equal(0.5, firstAmount(res[ingID(t, apiManager, "sugar")].Sum, true).Value, "sugar should be double the original")
		require.Equal(1.5, firstAmount(res[ingID(t, apiManager, "common")].Sum, true).Value, "common should be double the original")

	}
}

func ingID(t *testing.T, apiManager *API, name string) IngredientID {
	t.Helper()
	ing, err := apiManager.DB().IngredientByName(context.TODO(), name)
	require.NoError(t, err)
	return IngredientID(ing.Id)
}
