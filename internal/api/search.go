package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *API) Search(c echo.Context, params SearchParams) error {
	ctx := c.Request().Context()

	ctx, span := a.tracer.Start(ctx, "Search")
	defer span.End()

	_, listMeta := parsePagination(params.Offset, params.Limit)

	recipes, recipesCount, err := a.DB().GetRecipesDetails(ctx, string(params.Name))
	if err != nil {
		return handleErr(c, err)
	}
	ingredients, ingredientsCount, err := a.DB().GetIngredients(ctx, string(params.Name), nil)
	if err != nil {
		return handleErr(c, err)
	}

	listMeta.setTotalCount(recipesCount + ingredientsCount)

	var resRecipes []RecipeWrapper
	var resIngredients []Ingredient

	for _, x := range recipes {
		r, err := a.transformRecipe(ctx, x, true)
		if err != nil {
			return handleErr(c, err)
		}
		resRecipes = append(resRecipes, RecipeWrapper{Detail: r.Detail, Id: x.RecipeId})
	}
	for _, x := range ingredients {
		i := transformIngredient(x)
		resIngredients = append(resIngredients, i)
	}

	return c.JSON(http.StatusOK, SearchResult{Recipes: &resRecipes, Ingredients: &resIngredients})
}
