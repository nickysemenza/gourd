//go:build integration
// +build integration

package api

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// for now, one degree of seperation away which should be sufficient
// todo: rename child to related since can go either drection (up or down the tree)
func TestInsertGet(t *testing.T) {
	require := require.New(t)
	_, apiManager := makeHandler(t)
	ctx := context.Background()
	ctx = boil.WithDebug(ctx, true)

	all1, _, err := apiManager.RecipeListV2(ctx, DefaultPagination)
	require.NoError(err)

	inserted, err := apiManager.CreateRecipe(ctx, &RecipeWrapperInput{
		Detail: RecipeDetailInput{
			Name: fmt.Sprintf("r-%d", time.Now().Unix()),
			Sections: []RecipeSectionInput{
				{
					Duration:     &Amount{Value: 7, UpperValue: null.Float64From(20).Ptr()},
					Instructions: []SectionInstructionInput{{Instruction: "add salt"}},
					Ingredients: []SectionIngredientInput{{
						Kind:    IngredientKindIngredient,
						Amounts: []Amount{{Unit: "gram", Value: 1}},
						Name:    null.StringFrom("salt").Ptr(),
					}},
				},
			},
		},
	})
	require.NoError(err)
	all2, _, err := apiManager.RecipeListV2(ctx, DefaultPagination)
	require.NoError(err)

	ingEgg, err := apiManager.ingredientByName(ctx, "egg")
	require.NoError(err)
	ingFlour, err := apiManager.ingredientByName(ctx, "flour")
	require.NoError(err)
	ingWater, err := apiManager.ingredientByName(ctx, "water")
	require.NoError(err)

	require.Equal(1, len(all2)-len(all1), "inserting 1 recipe should increase length of getAll by 1")
	rw, err :=
		apiManager.recipeByDetailID(ctx, inserted.Detail.Id)
	require.NoError(err)
	r := RecipeDetailInput{
		Name: rw.Detail.Name,
	}
	r.ServingInfo.Unit = "items"
	r.Sections = []RecipeSectionInput{{
		Duration:     &Amount{Value: 7, UpperValue: null.Float64From(20).Ptr()},
		Instructions: []SectionInstructionInput{{Instruction: "add flour"}},
		Ingredients: []SectionIngredientInput{{
			Kind:    IngredientKindIngredient,
			Amounts: []Amount{{Unit: "grams", Value: 1}},
			Name:    &ingFlour.Name,
		}},
	}, {
		Duration:     &Amount{Value: 7, UpperValue: null.Float64From(20).Ptr()},
		Instructions: []SectionInstructionInput{{Instruction: "add more flour"}, {Instruction: "mix"}},
		Ingredients: []SectionIngredientInput{{
			Kind:    IngredientKindIngredient,
			Name:    &ingFlour.Name,
			Amounts: []Amount{{Unit: "grams", Value: 1}},
		}, {
			Kind:    IngredientKindIngredient,
			Name:    &ingWater.Name,
			Amounts: []Amount{{Unit: "grams", Value: 1}, {Unit: "c", Value: .7}},
		}, {
			Kind:    IngredientKindIngredient,
			Name:    &ingEgg.Name,
			Amounts: []Amount{{Unit: "grams", Value: 60}, {Unit: "large egg", Value: 1}},
		}},
	}}

	dbVersion, err := apiManager.CreateRecipe(ctx, &RecipeWrapperInput{Id: &rw.Id, Detail: r})

	r2w, err := apiManager.recipeById(ctx, dbVersion.Detail.Id)
	require.NoError(err)
	r2 := r2w.Detail
	require.EqualValues("items", r2.ServingInfo.Unit)
	require.EqualValues("add flour", r2.Sections[0].Instructions[0].Instruction)
	require.EqualValues(.7, r2.Sections[1].Ingredients[1].Amounts[1].Value)

	// _, err = db.InsertRecipe(ctx, &RecipeDetail{
	// 	Name: fmt.Sprintf("r2-%d", time.Now().Unix()),
	// 	Sections: []Section{{
	// 		TimeRange: `{"upper_value": 69, "value": 7}`,
	// 		Ingredients: []SectionIngredient{{
	// 			Amounts:  []Amount{{Unit: "grams", Value: 52}},
	// 			RecipeId: zero.StringFrom(r2.RecipeId),
	// 		}}}},
	// })
	// require.NoError(err)
}

func TestUnits(t *testing.T) {
	_, apiManager := makeHandler(t)
	ctx := context.Background()
	ctx = boil.WithDebug(ctx, true)
	detail1 := MustInsert(ctx, t, apiManager, NewCompact("testabc", []string{
		"1 tsp foo",
	}, []string{}))

	require.NoError(t, apiManager.insertIngredientMappings(ctx,
		[]IngredientMapping{
			{
				Name:    "foo",
				Aliases: []string{"bar"},
				UnitMappings: []UnitMapping{
					{
						A: Amount{Unit: "tsp", Value: 1},
						B: Amount{Unit: "g", Value: 2},
					},
				},
			},
		},
	))

	r2w, err := apiManager.recipeById(ctx, detail1)
	require.NoError(t, err)
	require.NotNil(t, r2w)

	require.Len(t, r2w.Detail.Sections, 1)
	require.Len(t, r2w.Detail.Sections[0].Ingredients, 1)
	amounts := r2w.Detail.Sections[0].Ingredients[0].Amounts
	require.Len(t, amounts, 2)
	require.Contains(t, amounts, Amount{Unit: "g", Value: 0.5, Source: null.StringFrom("calculated").Ptr()})
	require.Contains(t, amounts, Amount{Unit: "tsp", Value: 1, Source: null.StringFrom("db").Ptr()})
}
