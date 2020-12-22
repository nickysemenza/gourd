package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/deepmap/oapi-codegen/pkg/testutil"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/nickysemenza/gourd/common"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/manager"
	"github.com/stretchr/testify/require"
)

func TestAPI(t *testing.T) {
	require := require.New(t)

	tdb := db.NewDB(t)
	m := manager.New(tdb, nil, nil)

	apiManager := NewAPI(m)
	e := echo.New()
	e.Use(echo_middleware.Logger())
	RegisterHandlers(e, apiManager)

	rName := fmt.Sprintf("recipe-%s", common.UUID())
	iName := fmt.Sprintf("ing-%s", common.UUID())

	newIngredient := Ingredient{Name: iName}

	{
		result := testutil.NewRequest().Post("/ingredients").WithJsonBody(newIngredient).Go(t, e)
		require.Equal(http.StatusCreated, result.Code(), result.Recorder.Body)
		err := result.UnmarshalBodyToObject(&newIngredient)
		require.NoError(err)
	}

	{
		var results PaginatedIngredients
		result := testutil.NewRequest().Get("/ingredients?limit=100").Go(t, e)
		require.Equal(http.StatusOK, result.Code(), result.Recorder.Body)
		err := result.UnmarshalBodyToObject(&results)
		require.NoError(err)

		found := false
		for _, e := range *results.Ingredients {
			if e.Ingredient.Name == newIngredient.Name {
				found = true
				require.NotEmpty(e.Ingredient.Id)
			}
		}
		require.True(found)
	}

	makeRecipe := func(newRecipe RecipeWrapper) RecipeWrapper {
		result := testutil.NewRequest().Post("/recipes").WithJsonBody(newRecipe).Go(t, e)
		require.Equal(http.StatusCreated, result.Code(), result.Recorder.Body)

		var resultRecipe RecipeWrapper
		err := result.UnmarshalBodyToObject(&resultRecipe)
		require.NoError(err)
		return resultRecipe
	}
	id := ""
	{
		w := 12.5
		newRecipe := RecipeWrapper{
			Detail: Recipe{Name: rName},
			Sections: []RecipeSection{{Minutes: 3,
				Instructions: []SectionInstruction{{Instruction: "mix"}},
				Ingredients:  []SectionIngredient{{Grams: w, Ingredient: &newIngredient, Kind: "ingredient"}},
			}},
		}
		resultRecipe := makeRecipe(newRecipe)

		require.Equal(resultRecipe.Detail.Name, newRecipe.Detail.Name)
		id = resultRecipe.Detail.Id

		newRecipe.Detail.Name += "sub"
		newRecipe.Sections[0].Ingredients = append(newRecipe.Sections[0].Ingredients, SectionIngredient{
			Grams:  w,
			Recipe: &Recipe{Id: resultRecipe.Id},
			Kind:   "recipe"})
		makeRecipe(newRecipe)

	}

	{
		result := testutil.NewRequest().Get("/recipes?offset=0&limit=10").Go(t, e)
		require.Equal(http.StatusOK, result.Code())
		var results PaginatedRecipes
		err := result.UnmarshalBodyToObject(&results)
		require.NoError(err)
		// require.Contains(results, name)
		// require.Equal(resultRecipe.Detail.Name, newRecipe.Detail.Name)
	}
	{
		result := testutil.NewRequest().Get("/recipes/"+id).Go(t, e)
		require.Equal(http.StatusOK, result.Code())
		var results RecipeWrapper
		err := result.UnmarshalBodyToObject(&results)
		require.NoError(err)
		// require.Contains(results, name)
		require.Equal(results.Detail.Name, rName)
	}
}
