package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *API) indexRecipeDetails(ctx context.Context, w *RecipeWrapper) {
	ctx, span := a.tracer.Start(ctx, "indexRecipeDetails")
	defer span.End()

	toIndex := []RecipeDetail{w.Detail}
	if w.OtherVersions != nil {
		toIndex = append(toIndex, *w.OtherVersions...)
	}

	if err := a.R.Send(ctx, "index_recipe_detail", toIndex, nil); err != nil {
		l(ctx).Error(err)
	}
}

func (a *API) Misc(c echo.Context) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "Misc")
	defer span.End()

	// items, err := a.imagesFromRecipeDetailId(ctx, "rd_08c6db27")
	// items, err := a.Notion.PageById(ctx, "f6a5d0759d4a4becb95adf696b1cccb0")
	items, err := a.ingredientUsage(ctx, EntitySummary{Id: "rd_2dfbb24c", Kind: IngredientKindRecipe})

	if err != nil {
		return handleErr(c, err)
	}
	// s := spew.Sdump(recipes)
	// // s = strings.ReplaceAll(s, "\n", "<br/>")
	// s = fmt.Sprintf("<html>%s</html>", s)
	return c.JSON(http.StatusOK, items)
}
