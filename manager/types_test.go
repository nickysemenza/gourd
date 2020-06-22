package manager

import (
	"testing"

	"github.com/nickysemenza/food/db"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3/zero"
)

func TestFromRecipe(t *testing.T) {
	dbr := db.Recipe{Name: "name", UUID: "123"}
	dbr.TotalMinutes = zero.IntFrom(3)
	dbr.Unit = zero.StringFrom("items")

	dbr.Sections = []db.Section{{
		Minutes:      zero.IntFrom(88),
		Instructions: []db.SectionInstruction{{Instruction: "add flour"}},
		Ingredients: []db.SectionIngredient{{
			Grams: zero.FloatFrom(52),
		}},
	}, {
		Minutes:      zero.IntFrom(1),
		Instructions: []db.SectionInstruction{{Instruction: "add more flour"}, {Instruction: "mix"}},
		Ingredients: []db.SectionIngredient{{
			Grams:  zero.FloatFrom(120),
			Amount: zero.FloatFrom(2),
			Unit:   zero.StringFrom("cup"),
		}},
	}}

	dbr2 := FromRecipe(&dbr).toDB()
	require.EqualValues(t, dbr, *dbr2)
}
