//go:build integration
// +build integration

package api

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

// for now, one degree of seperation away which should be sufficient
// todo: rename child to related since can go either drection (up or down the tree)
func TestChildrenParentIngredients(t *testing.T) {
	require := require.New(t)
	_, apiManager := makeHandler(t)
	ctx := context.Background()

	require.NoError(apiManager.loadIngredientMappings(ctx,
		[]IngredientMapping{
			{Name: "a", Aliases: []string{"b"}},
			{Name: "b", Aliases: []string{"c"}},
		},
	))

	ingA := ingID(t, apiManager, "a")
	ingB := ingID(t, apiManager, "b")
	ingC := ingID(t, apiManager, "c")

	resA, err := apiManager.ingredientById(ctx, ingA)
	require.NoError(err)
	resB, err := apiManager.ingredientById(ctx, ingB)
	require.NoError(err)
	resC, err := apiManager.ingredientById(ctx, ingC)
	require.NoError(err)

	require.Equal([]string{"b"}, extractNames(*resA.Children))
	require.Equal([]string{"a", "c"}, extractNames(*resB.Children))
	require.Equal([]string{"b"}, extractNames(*resC.Children))

	rd := mustInsert(t, apiManager, newCompact("main", []string{
		"1 gram a",
		"1 gram b",
		"1 gram c"},
		[]string{}))

	wrapper, err := apiManager.recipeById(ctx, rd)
	require.NoError(err)
	spew.Dump(wrapper)
	ings := wrapper.Detail.Sections[0].Ingredients
	require.Len(ings, 3)
	require.Equal("a", ings[0].Ingredient.Ingredient.Name)
	require.Equal("b", ings[1].Ingredient.Ingredient.Name)
	require.Equal("c", ings[2].Ingredient.Ingredient.Name)

	require.Equal([]string{"b"}, extractNames(*ings[0].Ingredient.Children))
	require.Equal([]string{"a", "c"}, extractNames(*ings[1].Ingredient.Children))
	require.Equal([]string{"b"}, extractNames(*ings[2].Ingredient.Children))
}
