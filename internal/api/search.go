package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (a *API) Search(c echo.Context, params SearchParams) error {
	ctx := c.Request().Context()

	ctx, span := a.tracer.Start(ctx, "Search")
	defer span.End()

	listMeta := parsePagination(params.Offset, params.Limit)

	recipes, recipesCount, err := a.searchRecipes(ctx, listMeta, string(params.Name))
	if err != nil {
		return handleErr(c, err)
	}
	ingredients, ingredientsCount, err := a.searchIngredients(ctx, listMeta, string(params.Name))
	if err != nil {
		return handleErr(c, err)
	}

	// todo: use this
	listMeta.setTotalCount(+uint64(recipesCount + ingredientsCount))

	return c.JSON(http.StatusOK, SearchResult{Recipes: &recipes, Ingredients: &ingredients})
}

func (a *API) searchRecipes(ctx context.Context, pagination Items, name string) ([]RecipeWrapper, int64, error) {
	return a.RecipeListV2(ctx, pagination, qm.InnerJoin("recipe_details rd on rd.recipe_id = recipes.id"),
		qm.Where("lower(name) like lower(?)", fmt.Sprintf("%%%s%%", name)))

}
func (a *API) searchIngredients(ctx context.Context, pagination Items, name string) ([]IngredientWrapper, int64, error) {
	return a.IngredientListV2(ctx, pagination,
		qm.Where("lower(name) like lower(?)", fmt.Sprintf("%%%s%%", name)))

}
