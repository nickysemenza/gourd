//go:build integration
// +build integration

package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestUsage(t *testing.T) {
	require := require.New(t)
	_, apiManager := makeHandler(t)
	ctx := context.Background()

	r, err := apiManager.RecipeFromText(ctx, `
name: a
---
1 tsp b
1 gram b
1 tsp c`)
	require.NoError(err)
	r2, err := apiManager.RecipeFromText(ctx, `
name: pep
---
1 recipe a
1 gram c`)
	require.NoError(err)
	ids, err := apiManager.CreateRecipeDetails(ctx, *r, *r2)
	require.NoError(err)
	require.Len(ids, 2)
	rd := ids[1]
	fmt.Println(rd)

	err = apiManager.loadIngredientMappings(ctx, []IngredientMapping{{Name: "b", UnitMappings: []UnitMapping{{
		A: Amount{Value: 1, Unit: "tsp"},
		B: Amount{Value: 3, Unit: "g"},
	}}}})
	require.NoError(err)

	res, err := apiManager.IngredientUsage(ctx, []EntitySummary{{Id: rd, Multiplier: 2, Kind: IngredientKindRecipe}})
	require.NoError(err)
	spew.Dump(res)

	ing, err := apiManager.DB().IngredientByName(ctx, "b")
	require.NoError(err)
	id := IngredientID(ing.Id)
	require.Len(res[id].Sum, 1)
	require.Equal(4.0, res[id].Sum[0].Value)
	require.Equal("b", res[id].Ing.Name)
	require.Equal(string(id), res[id].Ing.Id)
}
