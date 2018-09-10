package model

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestRecipe_CreateOrUpdate(t *testing.T) {

	require := require.New(t)
	db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	require.NoError(err)
	defer db.Close()
	db = DBReset(db)
	db = DBMigrate(db)
	require.NotNil(db)

	//getting nonexistant should fail
	_, err = GetRecipeFromSlug(db, "foo")
	require.Error(err)

	r := Recipe{
		Slug: "test-1",
		Sections: []Section{
			{
				Minutes: 4,
				Ingredients: []SectionIngredient{
					{Item: Ingredient{Name: "flour"}, Grams: 100},
					{Item: Ingredient{Name: "water"}, Grams: 200},
				},
			},
			{
				Minutes: 2,
				Ingredients: []SectionIngredient{
					{Item: Ingredient{Name: "water"}, Grams: 10},
					{Item: Ingredient{Name: "flour"}, Grams: 600},
				},
			},
		},
	}
	//basil insert-and-retrieve smoke test
	err = r.CreateOrUpdate(db, false)
	require.NoError(err)
	r2, err := GetRecipeFromSlug(db, "test-1")
	require.NoError(err)
	require.Equal("test-1", r2.Slug, "slug should match")
	flourID := r2.Sections[0].Ingredients[0].Item.ID
	require.Equal("flour", r2.Sections[0].Ingredients[0].Item.Name)
	require.Equal("flour", r2.Sections[1].Ingredients[1].Item.Name)
	require.Equal(
		r2.Sections[0].Ingredients[0].Item.ID,
		r2.Sections[1].Ingredients[1].Item.ID,
		"both flours should have same ID",
	)
	require.NotEqual(
		r2.Sections[0].Ingredients[1].Item.ID,
		r2.Sections[1].Ingredients[1].Item.ID,
		"flour and sugar should have different IDs",
	)

	r2.Sections[0].Ingredients[0].Item.Name = "bread flour"
	//test something
	err = r2.CreateOrUpdate(db, false)
	require.NoError(err)
	r2, err = GetRecipeFromSlug(db, "test-1")

	require.NotEqual(flourID, r2.Sections[0].Ingredients[0].Item.ID)
	require.Equal("bread flour", r2.Sections[0].Ingredients[0].Item.Name)

	require.Equal(flourID, r2.Sections[1].Ingredients[1].Item.ID)
	require.Equal("flour", r2.Sections[1].Ingredients[1].Item.Name)

}
