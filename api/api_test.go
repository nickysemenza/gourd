//go:build integration
// +build integration

package api

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/deepmap/oapi-codegen/pkg/testutil"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/nickysemenza/gourd/common"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/rs_client"
	"github.com/stretchr/testify/require"
)

func makeHandler(t *testing.T) (*echo.Echo, *API) {
	t.Helper()
	tdb := db.NewDB(t)
	apiManager := New(tdb, nil, nil, rs_client.New("http://localhost:8080/"), nil, nil)

	e := echo.New()
	e.Use(echo_middleware.Logger())
	RegisterHandlers(e, apiManager)
	return e, apiManager
}
func TestAPI(t *testing.T) {
	require := require.New(t)

	e, _ := makeHandler(t)
	rName := fmt.Sprintf("recipe-%s", common.ID(""))
	iName := fmt.Sprintf("ing-%s", common.ID(""))

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
			Detail: RecipeDetail{Name: rName,
				Sections: []RecipeSection{{Duration: &TimeRange{Min: 3},
					Instructions: []SectionInstruction{{Instruction: "mix"}},
					Ingredients:  []SectionIngredient{{Amounts: []Amount{{Unit: "grams", Value: w}}, Ingredient: &IngredientDetail{Ingredient: newIngredient}, Kind: "ingredient"}},
				}}},
		}
		resultRecipe := makeRecipe(newRecipe)

		require.Equal(resultRecipe.Detail.Name, newRecipe.Detail.Name)
		id = resultRecipe.Detail.Id

		newRecipe.Detail.Name += "sub"
		newRecipe.Detail.Sections[0].Ingredients = append(newRecipe.Detail.Sections[0].Ingredients, SectionIngredient{
			Amounts: []Amount{{Unit: "grams", Value: w}},
			Recipe:  &RecipeDetail{Id: resultRecipe.Id},
			Kind:    "recipe"})
		makeRecipe(newRecipe)

	}

	{
		result := testutil.NewRequest().Get("/recipes?offset=0&limit=10").Go(t, e)
		require.Equal(http.StatusOK, result.Code())
		var results PaginatedRecipeWrappers
		require.NoError(result.UnmarshalBodyToObject(&results))
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

func TestRecipeReferencingRecipe(t *testing.T) {
	require := require.New(t)
	_, apiManager := makeHandler(t)
	ctx := context.Background()
	r, err := apiManager.RecipeFromFile(ctx, "../testdata/dep_1.yaml")
	require.NoError(err)
	err = apiManager.CreateRecipeDetails(ctx, r...)
	require.NoError(err)
}

func TestSearches(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()
	e, api := makeHandler(t)

	rName := fmt.Sprintf("recipe-%s", common.ID(""))
	iName := fmt.Sprintf("ing-%s", common.ID(""))
	err := api.CreateRecipeDetails(ctx, RecipeDetailInput{
		Name:     rName,
		Sections: []RecipeSectionInput{{Ingredients: []SectionIngredientInput{{Kind: "ingredient", Name: &iName}}}},
	})
	require.NoError(err)

	ingId := SearchByKind(t, e, iName, "ingredient")
	require.NotEmpty(ingId)

	result := testutil.NewRequest().Post("/ingredients/"+ingId+"/convert_to_recipe").Go(t, e)
	require.Equal(http.StatusCreated, result.Code())
	var results RecipeDetail
	require.NoError(result.UnmarshalBodyToObject(&results))

	recId := SearchByKind(t, e, iName, "recipe")
	require.NotEmpty(recId)

	// require.Equal(results.Id, recId)

}

func SearchByKind(t *testing.T, e *echo.Echo, name string, kind string) string {
	require := require.New(t)
	result := testutil.NewRequest().Get("/search?name="+name).Go(t, e)
	require.Equal(http.StatusOK, result.Code())
	var results SearchResult
	require.NoError(result.UnmarshalBodyToObject(&results))
	id := ""
	switch kind {
	case "ingredient":
		for _, x := range *results.Ingredients {
			if x.Name == name {
				id = x.Id
			}
		}
	case "recipe":
		for _, x := range *results.Recipes {
			if x.Detail.Name == name {
				id = x.Id
			}
		}
	default:
		t.Fatalf("bad kind: %s", kind)
	}

	return id
}
