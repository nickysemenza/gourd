//go:build integration
// +build integration

package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUsage(t *testing.T) {
	require := require.New(t)
	_, apiManager := makeHandler(t)
	ctx := context.Background()

	MustInsert(t, apiManager, NewCompact("sub", []string{
		"1 tsp salt",
		"1 gram salt",
		"1 gram sugar",
		"1 tsp pepper",
		"1 gram pepper",
		"1 tbsp water",
	}, []string{}))

	// rdMain is a recipe that has 1 recipe of sub
	// should end up with:
	// salt:
	//   from main: 1tsp = 3 gram + 1 gram
	//   from sub: n/a
	//   TOTAL: 4 gram
	// sugar:
	//   from main: n/a
	//   from sub: 1 gram
	//   TOTAL: 1 gram
	// pepper:
	//   from main: 1 tsp + 1 gram
	//   from sub: 1 tsp + 1 gram
	//   TOTAL: 2 tsp + 2 gram pepper
	// water:
	//   from main: 1 gram
	//   from sub: 1 tbsp = 3 tsp = 3 gram
	//   TOTAL: 4 gram
	rdMain := MustInsert(t, apiManager, NewCompact("main", []string{
		"1 tsp pepper",
		"1 gram pepper",
		"1 recipe sub",
		"1 gram water"},
		[]string{}))

	rdSMallMain := MustInsert(t, apiManager, NewCompact("smallmain", []string{
		"0.5 recipe sub",
		"1 tsp water",
	}, []string{}))

	require.NoError(apiManager.insertIngredientMappings(ctx, []IngredientMapping{
		{Name: "salt", UnitMappings: []UnitMapping{{
			A: Amount{Value: 1, Unit: "tsp"},
			B: Amount{Value: 3, Unit: "g"},
		}}},
		{Name: "water", UnitMappings: []UnitMapping{{
			A: Amount{Value: 1, Unit: "tsp"},
			B: Amount{Value: 1, Unit: "g"},
		}}},
	}))

	for _, mult := range []float64{1, 0.5, 2.0} {
		{
			res, err := apiManager.IngredientUsage(ctx, EntitySummary{Id: rdMain, Multiplier: mult, Kind: IngredientKindRecipe})
			require.NoError(err)

			// should only have grams
			require.Equal(1*mult, firstAmount(res[IngIDFromName(t, apiManager, "sugar")].Sum, true).Value)
			require.Nil(firstAmount(res[IngIDFromName(t, apiManager, "sugar")].Sum, false))

			// should only have grams
			require.Equal(4*mult, firstAmount(res[IngIDFromName(t, apiManager, "water")].Sum, true).Value)
			// require.Nil(firstAmount(res[ingID(t, apiManager, "water")].Sum, false))

			// should have both grams and non grams
			require.Equal(2*mult, firstAmount(res[IngIDFromName(t, apiManager, "pepper")].Sum, true).Value)
			require.Equal(2*mult, firstAmount(res[IngIDFromName(t, apiManager, "pepper")].Sum, false).Value)

			ingSalt := IngIDFromName(t, apiManager, "salt")
			// should only have grams
			require.Equal(4*mult, firstAmount(res[ingSalt].Sum, true).Value)
			// require.Nil(firstAmount(res[ingSalt].Sum, false))
			require.Equal("salt", res[ingSalt].Meta.Name)
			require.Equal(string(ingSalt), res[ingSalt].Meta.Id)
		}
		{
			res, err := apiManager.IngredientUsage(ctx, EntitySummary{Id: rdSMallMain, Multiplier: mult, Kind: IngredientKindRecipe})
			require.NoError(err)

			require.Equal(mult*2.5, firstAmount(res[IngIDFromName(t, apiManager, "water")].Sum, true).Value)
			require.Equal(mult*0.5, firstAmount(res[IngIDFromName(t, apiManager, "sugar")].Sum, true).Value)

			// should have both grams and non grams
			require.Equal(.5*mult, firstAmount(res[IngIDFromName(t, apiManager, "pepper")].Sum, true).Value)
			require.Equal(.5*mult, firstAmount(res[IngIDFromName(t, apiManager, "pepper")].Sum, false).Value)

		}

	}
}
