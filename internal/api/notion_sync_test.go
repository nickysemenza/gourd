//go:build integration
// +build integration

package api

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestSync(t *testing.T) {
	require := require.New(t)
	apiManager := makeAPI(t)
	ctx := context.Background()
	err := apiManager.Sync(ctx, 14)
	require.NoError(err)

	items, err := apiManager.RecipeListV2(ctx, 10, 0)
	require.NoError(err)
	require.Len(items, 3)
	rd := items[0].Detail.Id
	res, err := apiManager.recipeById(ctx, rd)
	require.NoError(err)

	require.Len(res.Detail.Sections, 1)
	require.Equal("bread", res.Detail.Sections[0].Ingredients[0].Ingredient.Ingredient.Name)
	require.Equal("eat", res.Detail.Sections[0].Instructions[0].Instruction)

	meals, err := apiManager.listMeals(ctx)
	require.NoError(err)
	require.Len(meals, 1)
	spew.Dump(meals)

	// l, err := apiManager.Latex(ctx, rd)
	// require.NoError(err)
	// require.Greater(len(l), 10000)

}
