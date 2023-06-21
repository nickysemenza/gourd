//go:build integration
// +build integration

package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestUsage(t *testing.T) {
	_, apiManager := makeHandler(t)
	ctx := context.Background()
	ctx = boil.WithDebug(ctx, true)
	rdSub := MustInsert(ctx, t, apiManager, NewCompact("subA", []string{
		"1 tsp salt",
		"1 gram salt",
		"1 gram sugar",
		"1 tsp pepper",
		"1 gram pepper",
		"1 tbsp water",
	}, []string{}))

	// rdMain is a recipe that has 1 recipe of subA
	// should end up with:
	// salt:
	//   from main: 1tsp = 3 gram + 1 gram
	//   from subA: n/a
	//   TOTAL: 4 gram
	// sugar:
	//   from main: n/a
	//   from subA: 1 gram
	//   TOTAL: 1 gram
	// pepper:
	//   from main: 1 tsp + 1 gram
	//   from subA: 1 tsp + 1 gram
	//   TOTAL: 2 tsp + 2 gram pepper
	// water:
	//   from main: 1 gram
	//   from subA: 1 tbsp = 3 tsp = 3 gram
	//   TOTAL: 4 gram
	rdMain := MustInsert(ctx, t, apiManager, NewCompact("main", []string{
		"1 tsp pepper",
		"1 gram pepper",
		"1 recipe subA",
		"1 gram water"},
		[]string{}))

	rdSMallMain := MustInsert(ctx, t, apiManager, NewCompact("smallmain", []string{
		"0.5 recipe subA",
		"1 tsp water",
	}, []string{}))

	require.NoError(t, apiManager.insertIngredientMappings(ctx, []IngredientMapping{
		{Name: "salt", UnitMappings: []UnitMapping{{
			A: Amount{Value: 1, Unit: "tsp"},
			B: Amount{Value: 3, Unit: "g"},
		}}},
		{Name: "water", UnitMappings: []UnitMapping{{
			A: Amount{Value: 1, Unit: "tsp"},
			B: Amount{Value: 1, Unit: "g"},
		}}},
	}))

	// check subA first
	{
		res, err := apiManager.ingredientUsage(ctx, EntitySummary{Id: rdSub, Multiplier: 1, Kind: IngredientKindRecipe})
		require.NoError(t, err)

		amt := firstAmount(res[IngIDFromName(t, apiManager, "water")].Sum, true)
		require.NotNil(t, amt)
		require.Equal(t, 3.0, amt.Value, amt)
	}

	for _, mult := range []float64{1, 0.5, 2.0} {
		t.Run(fmt.Sprint(mult), func(t *testing.T) {
			require := require.New(t)
			{
				res, err := apiManager.ingredientUsage(ctx, EntitySummary{Id: rdMain, Multiplier: mult, Kind: IngredientKindRecipe})
				require.NoError(err)

				spew.Dump(res)

				// should only have grams
				amt := firstAmount(res[IngIDFromName(t, apiManager, "water")].Sum, true)
				require.Equal(4*mult, amt.Value, amt)
				// require.Nil(firstAmount(res[ingID(t, apiManager, "water")].Sum, false))

				// should have both grams and non grams
				amt = firstAmount(res[IngIDFromName(t, apiManager, "pepper")].Sum, true)
				require.Equal(2*mult, amt.Value, amt)
				amt = firstAmount(res[IngIDFromName(t, apiManager, "pepper")].Sum, false)
				require.Equal(2*mult, amt.Value, amt)

				// should only have grams
				require.Equal(1*mult, firstAmount(res[IngIDFromName(t, apiManager, "sugar")].Sum, true).Value)
				require.Nil(firstAmount(res[IngIDFromName(t, apiManager, "sugar")].Sum, false))

				ingSalt := IngIDFromName(t, apiManager, "salt")
				// should only have grams

				//TODO: this is broken
				// require.Equal(4*mult, firstAmount(res[ingSalt].Sum, true).Value)

				// require.Nil(firstAmount(res[ingSalt].Sum, false))
				require.Equal("salt", res[ingSalt].Meta.Name)
				require.Equal(string(ingSalt), res[ingSalt].Meta.Id)
			}
			{
				res, err := apiManager.ingredientUsage(ctx, EntitySummary{Id: rdSMallMain, Multiplier: mult, Kind: IngredientKindRecipe})
				require.NoError(err)

				require.Equal(mult*2.5, firstAmount(res[IngIDFromName(t, apiManager, "water")].Sum, true).Value)
				require.Equal(mult*0.5, firstAmount(res[IngIDFromName(t, apiManager, "sugar")].Sum, true).Value)

				// should have both grams and non grams
				require.Equal(.5*mult, firstAmount(res[IngIDFromName(t, apiManager, "pepper")].Sum, true).Value)
				require.Equal(.5*mult, firstAmount(res[IngIDFromName(t, apiManager, "pepper")].Sum, false).Value)

			}
		})
	}
}
