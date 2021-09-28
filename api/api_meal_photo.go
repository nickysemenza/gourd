package api

import (
	"context"
	"fmt"
	"net/http"

	sq "github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/db"
)

func (a *API) fromDBPhoto(ctx context.Context, photos []db.Photo, getURLs bool) ([]GooglePhoto, []string, error) {
	ctx, span := a.tracer.Start(ctx, "fromDBPhoto")
	defer span.End()
	items := []GooglePhoto{}
	var ids []string
	for _, p := range photos {
		gp := GooglePhoto{Id: p.PhotoID, Created: p.Created}
		if p.BlurHash.Valid {
			s := p.BlurHash.String
			gp.BlurHash = &s
		}
		items = append(items, gp)
		ids = append(ids, p.PhotoID)
	}

	if getURLs {
		results, err := a.Manager.Photos.GetMediaItems(ctx, ids)
		if err != nil {
			return nil, nil, err
		}
		for x, item := range items {
			val, ok := results[item.Id]
			if !ok {
				continue
			}
			items[x].BaseUrl = val.BaseUrl
			items[x].Width = val.MediaMetadata.Width
			items[x].Height = val.MediaMetadata.Height
		}
	}
	return items, ids, nil
}
func (a *API) ListPhotos(c echo.Context, params ListPhotosParams) error {
	ctx := c.Request().Context()
	photos, err := a.Manager.DB().GetPhotos(ctx)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	items, _, err := a.fromDBPhoto(ctx, photos, true)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}

	resp := PaginatedPhotos{
		Photos: &items,
	}
	return c.JSON(http.StatusOK, resp)
}

func (a *API) GetMealInfo(ctx context.Context, meals db.Meals) ([]Meal, error) {
	items := []Meal{}
	mealIds := meals.MealIDs()
	mealRecipes, err := a.DB().GetMealRecipes(ctx, mealIds...)
	if err != nil {
		return nil, err
	}

	recipesDetails, err := a.DB().GetRecipeDetailWhere(ctx, sq.Eq{"recipe": mealRecipes.RecipeIDs()})
	if err != nil {
		return nil, err
	}
	recipeDetailsById := recipesDetails.ByRecipeId()

	var gphotoIDs []string
	for _, m := range meals {
		meal := Meal{Id: m.ID,
			Name:  m.Name,
			AteAt: m.AteAt}

		mrs := []MealRecipe{}
		for _, mr := range mealRecipes.ByMealID()[m.ID] {

			test := a.transformRecipes(ctx, recipeDetailsById[mr.RecipeID])

			mrs = append(mrs, MealRecipe{Multiplier: mr.Multiplier, Recipe: test[0]})
		}
		meal.Recipes = &mrs

		photos, err := a.DB().GetPhotosForMeal(ctx, m.ID)
		if err != nil {
			return nil, err
		}

		photos2, gIDs, err := a.fromDBPhoto(ctx, photos, false)
		if err != nil {
			return nil, err
		}
		meal.Photos = photos2
		gphotoIDs = append(gphotoIDs, gIDs...)

		items = append(items, meal)
	}
	urls, err := a.Manager.Photos.GetMediaItems(ctx, gphotoIDs)
	if err != nil {
		return nil, err
	}
	for x, item := range items {
		for y, photo := range item.Photos {
			val, ok := urls[photo.Id]
			if !ok {
				continue
			}
			items[x].Photos[y].BaseUrl = val.BaseUrl
			items[x].Photos[y].Width = val.MediaMetadata.Width
			items[x].Photos[y].Height = val.MediaMetadata.Height
		}
	}
	return items, nil
}
func (a *API) ListMeals(c echo.Context, params ListMealsParams) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "ListMeals")
	defer span.End()

	meals, err := a.DB().GetAllMeals(ctx)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	items, err := a.GetMealInfo(ctx, meals)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}

	resp := PaginatedMeals{
		Meals: &items,
	}
	return c.JSON(http.StatusOK, resp)
}

func (a *API) GetMealById(c echo.Context, mealId string) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "GetMealById")
	defer span.End()

	meal, err := a.DB().GetMealById(ctx, mealId)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	items, err := a.GetMealInfo(ctx, []db.Meal{*meal})
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, items[0])
}
func (a *API) UpdateRecipesForMeal(c echo.Context, mealId string) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "UpdateRecipesForMeal")
	defer span.End()

	var r MealRecipeUpdate
	if err := c.Bind(&r); err != nil {
		err = fmt.Errorf("invalid format for input: %w", err)
		return sendErr(c, http.StatusBadRequest, err)
	}
	switch r.Action {
	case MealRecipeUpdateActionAdd:
		err := a.DB().AddRecipeToMeal(ctx, mealId, r.RecipeId, r.Multiplier)
		if err != nil {
			return sendErr(c, http.StatusInternalServerError, err)
		}
		return a.GetMealById(c, mealId)
	case MealRecipeUpdateActionRemove:
		return sendErr(c, http.StatusBadRequest, fmt.Errorf("unsupported %s", r.Action))
	default:
		return sendErr(c, http.StatusBadRequest, fmt.Errorf("unknown action %s", r.Action))
	}

}
