//go:build integration
// +build integration

package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// for now, one degree of seperation away which should be sufficient
// todo: rename child to related since can go either drection (up or down the tree)
func TestChildrenParentIngredients(t *testing.T) {
	require := require.New(t)
	_, apiManager := makeHandler(t)
	ctx := context.Background()
	ctx = boil.WithDebug(ctx, true)

	require.NoError(apiManager.insertIngredientMappings(ctx,
		[]IngredientMapping{
			{Name: "aparent", Aliases: []string{"bchild"}},
			// {Name: "bchild", Aliases: []string{"c"}},
		},
	))

	ingA := IngIDFromName(t, apiManager, "aparent")
	ingB := IngIDFromName(t, apiManager, "bchild")
	// ingC := ingID(t, apiManager, "c")

	resA, err := apiManager.ingredientById(ctx, ingA)
	require.NoError(err)
	resB, err := apiManager.ingredientById(ctx, ingB)
	require.NoError(err)
	// resC, err := apiManager.ingredientById(ctx, ingC)
	// require.NoError(err)

	require.Equal([]string{"bchild"}, ExtractNames(*resA.Children))
	require.Equal([]string{"aparent"}, ExtractNames(*resB.Children))
	// require.Equal([]string{"bchild"}, extractNames(*resC.Children))

	rd := MustInsert(ctx, t, apiManager, NewCompact("main", []string{
		"1 gram aparent",
		"1 gram bchild",
	},
		[]string{}))

	wrapper, err := apiManager.recipeById(ctx, rd)
	require.NoError(err)
	ings := wrapper.Detail.Sections[0].Ingredients

	require.Len(ings, 2)
	require.Equal("aparent", ings[0].Ingredient.Ingredient.Name)
	require.Equal("bchild", ings[1].Ingredient.Ingredient.Name)
	// require.Equal("c", ings[2].Ingredient.Ingredient.Name)

	require.Equal([]string{"bchild"}, ExtractNames(*ings[0].Ingredient.Children))
	require.Equal([]string{"aparent"}, ExtractNames(*ings[1].Ingredient.Children))
	// require.Equal([]string{"bchild"}, extractNames(*ings[2].Ingredient.Children))
}

func TestConvertIngToRecipe(t *testing.T) {
	_, apiManager := makeHandler(t)
	ctx := context.Background()
	ctx = boil.WithDebug(ctx, true)
	MustInsert(ctx, t, apiManager, NewCompact("testabc", []string{
		"1 tsp salt",
		"1 gram salt",
		"1 gram sugar",
	}, []string{}))

	ing := IngIDFromName(t, apiManager, "salt")
	_, err := apiManager.convertIngredientToRecipe(ctx, string(ing))
	require.NoError(t, err)
	// todo: more assertions here
}
