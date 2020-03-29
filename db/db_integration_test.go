package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bxcodec/faker/v3"
)

func TestInsertGet(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)
	db := newDB(t)
	uuid := faker.UUIDHyphenated()
	err := db.InsertRecipe(ctx, &Recipe{
		Name:     faker.Username(),
		UUID:     uuid,
		Sections: []Section{{Minutes: sql.NullInt64{Valid: true, Int64: 33}}},
	})

	require.NoError(err)
	r, err := db.GetRecipeByUUID(ctx, uuid)
	require.NoError(err)
	r.TotalMinutes = sql.NullInt64{Valid: true, Int64: 3}
	r.Unit = sql.NullString{Valid: true, String: "items"}
	r.Sections = []Section{{
		Minutes:      sql.NullInt64{Valid: true, Int64: 88},
		Instructions: []SectionInstruction{{Instruction: "add flour"}},
		Ingredients: []SectionIngredient{{
			Grams: sql.NullFloat64{Valid: true, Float64: 52},
			Name:  "flour",
		}},
	}, {
		Minutes:      sql.NullInt64{Valid: true, Int64: 1},
		Instructions: []SectionInstruction{{Instruction: "add more flour"}, {Instruction: "mix"}},
		Ingredients: []SectionIngredient{{
			Grams: sql.NullFloat64{Valid: true, Float64: 1},
			Name:  "flour",
		}},
	}}

	err = db.UpdateRecipe(ctx, r)
	require.NoError(err)
	r2, err := db.GetRecipeByUUID(ctx, uuid)
	require.NoError(err)
	require.EqualValues(3, r2.TotalMinutes.Int64)
	require.EqualValues("items", r2.Unit.String)
	require.EqualValues("add flour", r2.Sections[0].Instructions[0].Instruction)

}
