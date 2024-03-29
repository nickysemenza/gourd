package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/internal/db"
	"github.com/nickysemenza/gourd/internal/db/models"
)

func (a *API) notionPhotosFromDBPhoto(ctx context.Context, photos models.NotionImageSlice) ([]Photo, error) {
	items := []Photo{}
	for _, aa := range photos {
		item, err := a.photoFromModel(ctx, aa.R.Image, Notion)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	return items, nil
}
func (a *API) photoFromModel(ctx context.Context, p *models.Image, source PhotoSource) (*Photo, error) {
	bh := p.BlurHash
	url, err := a.ImageStore.GetImageURL(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	return &Photo{
		Id:       p.ID,
		TakenAt:  p.TakenAt.Ptr(),
		BlurHash: &bh,
		Width:    300,
		Height:   400,
		BaseUrl:  url,
		Source:   source,
	}, nil
}
func (a *API) googlePhotosFromDBPhoto(ctx context.Context, photos models.GphotosPhotoSlice) ([]Photo, error) {
	ctx, span := a.tracer.Start(ctx, "fromDBPhoto")
	defer span.End()

	items := []Photo{}
	for _, aa := range photos {
		item, err := a.photoFromModel(ctx, aa.R.Image, Google)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	return items, nil
}
func (a *API) ListPhotos(c echo.Context, params ListPhotosParams) error {
	ctx := c.Request().Context()
	photos, err := a.DB().GetPhotos(ctx)
	if err != nil {
		return handleErr(c, err)
	}
	items, err := a.googlePhotosFromDBPhoto(ctx, photos)
	if err != nil {
		return handleErr(c, err)
	}

	resp := PaginatedPhotos{
		Photos: &items,
	}
	return c.JSON(http.StatusOK, resp)
}

func (a *API) GetMealInfo(ctx context.Context, meals db.Meals) ([]Meal, error) {
	ctx, span := a.tracer.Start(ctx, "GetMealInfo")
	defer span.End()

	items := []Meal{}
	mealIds := meals.MealIDs()
	mealRecipes, err := a.DB().GetMealRecipes(ctx, mealIds...)
	if err != nil {
		return nil, err
	}

	for _, m := range meals {
		meal := Meal{Id: m.ID,
			Name:  m.Name,
			AteAt: m.AteAt}

		mrs := []MealRecipe{}
		for _, mr := range mealRecipes.ByMealID()[m.ID] {

			wrapper, err := a.recipeByWrapperID(ctx, mr.RecipeID)

			if err != nil {
				return nil, err
			}

			mrs = append(mrs, MealRecipe{Multiplier: mr.Multiplier, Recipe: wrapper.Detail})
		}
		meal.Recipes = &mrs

		googlePhotosDB, err := a.DB().GetPhotosForMeal(ctx, m.ID)
		if err != nil {
			return nil, err
		}

		googlePhotos, err := a.googlePhotosFromDBPhoto(ctx, googlePhotosDB)
		if err != nil {
			return nil, err
		}

		notionPhotosDB, err := a.DB().GetNotionPhotosForMeal(ctx, m.ID)
		if err != nil {
			return nil, err
		}
		notionPhotos, err := a.notionPhotosFromDBPhoto(ctx, notionPhotosDB)
		if err != nil {
			return nil, err
		}

		googlePhotos = append(googlePhotos, notionPhotos...)
		meal.Photos = googlePhotos
		items = append(items, meal)
	}
	return items, nil
}

func (a *API) listMeals(ctx context.Context) ([]Meal, error) {
	ctx, span := a.tracer.Start(ctx, "ListMeals")
	defer span.End()

	meals, err := a.DB().GetAllMeals(ctx)
	if err != nil {
		return nil, err
	}
	items, err := a.GetMealInfo(ctx, meals)
	if err != nil {
		return nil, err
	}
	return items, nil
}
func (a *API) ListMeals(c echo.Context, _ ListMealsParams) error {
	ctx := c.Request().Context()

	items, err := a.listMeals(ctx)
	if err != nil {
		return handleErr(c, err)
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
		return handleErr(c, err)
	}
	items, err := a.GetMealInfo(ctx, []db.Meal{*meal})
	if err != nil {
		return handleErr(c, err)
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
	case Add:
		tx := a.tx(ctx)
		err := a.DB().AddRecipeToMeal(ctx, mealId, r.RecipeId, &r.Multiplier, tx)
		if err != nil {
			return handleErr(c, err)
		}
		if err := tx.Commit(); err != nil {
			return handleErr(c, err)
		}
		return a.GetMealById(c, mealId)
	case Remove:
		return sendErr(c, http.StatusBadRequest, fmt.Errorf("unsupported %s", r.Action))
	default:
		return sendErr(c, http.StatusBadRequest, fmt.Errorf("unknown action %s", r.Action))
	}

}
