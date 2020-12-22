package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3/zero"
)

//nolint: funlen
func TestInsertGet(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)
	db := NewDB(t)

	all1, _, err := db.GetRecipes(ctx, "")
	require.NoError(err)

	insertedDetail, err := db.InsertRecipe(ctx, &RecipeDetail{
		Name:     fmt.Sprintf("r-%d", time.Now().Unix()),
		Sections: []Section{{Minutes: zero.IntFrom(33)}},
	})

	require.NoError(err)
	all2, _, err := db.GetRecipes(ctx, "")
	require.NoError(err)

	ingEgg, err := db.IngredientByName(ctx, "egg")
	require.NoError(err)
	ingFlour, err := db.IngredientByName(ctx, "flour")
	require.NoError(err)
	ingWater, err := db.IngredientByName(ctx, "water")
	require.NoError(err)

	require.Equal(1, len(all2)-len(all1), "inserting 1 recipe should increase length of getAll by 1")
	r, err := db.GetRecipeDetailByUUIDFull(ctx, insertedDetail.UUID)
	require.NoError(err)
	r.TotalMinutes = zero.IntFrom(3)
	r.Unit = zero.StringFrom("items")
	r.Sections = []Section{{
		Minutes:      zero.IntFrom(88),
		Instructions: []SectionInstruction{{Instruction: "add flour"}},
		Ingredients: []SectionIngredient{{
			Grams:          zero.FloatFrom(52),
			IngredientUUID: zero.StringFrom(ingFlour.UUID),
		}},
	}, {
		Minutes:      zero.IntFrom(1),
		Instructions: []SectionInstruction{{Instruction: "add more flour"}, {Instruction: "mix"}},
		Ingredients: []SectionIngredient{{
			Grams:          zero.FloatFrom(1),
			IngredientUUID: zero.StringFrom(ingFlour.UUID),
		}, {
			Grams:          zero.FloatFrom(178),
			IngredientUUID: zero.StringFrom(ingWater.UUID),
			Amount:         zero.FloatFrom(.7),
			Unit:           zero.StringFrom("c"),
		}, {

			Grams:          zero.FloatFrom(60),
			IngredientUUID: zero.StringFrom(ingEgg.UUID),
			Amount:         zero.FloatFrom(1),
			Unit:           zero.StringFrom("large egg"),
		}},
	}}

	r2, err := db.InsertRecipe(ctx, r)
	require.NoError(err)
	require.EqualValues(3, r2.TotalMinutes.Int64)
	require.EqualValues("items", r2.Unit.String)
	require.EqualValues("add flour", r2.Sections[0].Instructions[0].Instruction)
	require.EqualValues(.7, r2.Sections[1].Ingredients[1].Amount.Float64)

	_, err = db.InsertRecipe(ctx, &RecipeDetail{
		Name: fmt.Sprintf("r2-%d", time.Now().Unix()),
		Sections: []Section{{
			Minutes: zero.IntFrom(33),
			Ingredients: []SectionIngredient{{
				Grams:      zero.FloatFrom(52),
				RecipeUUID: zero.StringFrom(r2.RecipeUUID),
			}}}},
	})
	require.NoError(err)
}
