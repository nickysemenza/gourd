//go:build integration
// +build integration

package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4/zero"
)

//nolint: funlen
func TestInsertGet(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)
	db := NewTestDB(t)

	all1, _, err := db.GetRecipesDetails(ctx, "")
	require.NoError(err)

	insertedDetail, err := db.InsertRecipe(ctx, &RecipeDetail{
		Name:     fmt.Sprintf("r-%d", time.Now().Unix()),
		Sections: []Section{{TimeRange: `{"max": 69, "min": 7}`}},
	})

	require.NoError(err)
	all2, _, err := db.GetRecipesDetails(ctx, "")
	require.NoError(err)

	ingEgg, err := db.IngredientByName(ctx, "egg")
	require.NoError(err)
	ingFlour, err := db.IngredientByName(ctx, "flour")
	require.NoError(err)
	ingWater, err := db.IngredientByName(ctx, "water")
	require.NoError(err)

	require.Equal(1, len(all2)-len(all1), "inserting 1 recipe should increase length of getAll by 1")
	r, err := db.GetRecipeDetailByIdFull(ctx, insertedDetail.Id)
	require.NoError(err)
	r.Unit = zero.StringFrom("items")
	r.Sections = []Section{{
		TimeRange:    `{"max": 69, "min": 7}`,
		Instructions: []SectionInstruction{{Instruction: "add flour"}},
		Ingredients: []SectionIngredient{{
			Amounts:      []Amount{{Unit: "grams", Value: 1}},
			IngredientId: zero.StringFrom(ingFlour.Id),
		}},
	}, {
		TimeRange:    `{"max": 69, "min": 7}`,
		Instructions: []SectionInstruction{{Instruction: "add more flour"}, {Instruction: "mix"}},
		Ingredients: []SectionIngredient{{
			IngredientId: zero.StringFrom(ingFlour.Id),
			Amounts:      []Amount{{Unit: "grams", Value: 1}},
		}, {
			IngredientId: zero.StringFrom(ingWater.Id),
			Amounts:      []Amount{{Unit: "grams", Value: 1}, {Unit: "c", Value: .7}},
		}, {

			IngredientId: zero.StringFrom(ingEgg.Id),
			Amounts:      []Amount{{Unit: "grams", Value: 60}, {Unit: "large egg", Value: 1}},
		}},
	}}

	r2, err := db.InsertRecipe(ctx, r)
	require.NoError(err)
	require.EqualValues("items", r2.Unit.String)
	require.EqualValues("add flour", r2.Sections[0].Instructions[0].Instruction)
	require.EqualValues(.7, r2.Sections[1].Ingredients[1].Amounts[1].Value)

	_, err = db.InsertRecipe(ctx, &RecipeDetail{
		Name: fmt.Sprintf("r2-%d", time.Now().Unix()),
		Sections: []Section{{
			TimeRange: `{"max": 69, "min": 7}`,
			Ingredients: []SectionIngredient{{
				Amounts:  []Amount{{Unit: "grams", Value: 52}},
				RecipeId: zero.StringFrom(r2.RecipeId),
			}}}},
	})
	require.NoError(err)
}
