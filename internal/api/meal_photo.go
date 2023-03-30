package api

import (
	"context"
	"fmt"
	"net/http"

	sq "github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/internal/db"
)

func (a *API) notionPhotosFromDBPhoto(ctx context.Context, photos []db.NotionImage) ([]Photo, error) {
	items := []Photo{}
	for _, aa := range photos {
		bh := aa.Image.BlurHash
		url, err := a.ImageStore.GetImageURL(ctx, aa.Image.ID)
		if err != nil {
			return nil, err
		}
		items = append(items, Photo{
			Id:       aa.BlockID,
			TakenAt:  aa.Image.TakenAt.Ptr(),
			BlurHash: &bh,
			Width:    300,
			Height:   400,
			BaseUrl:  url,
			Source:   PhotoSourceNotion,
		})
	}
	return items, nil
}

func (a *API) googlePhotosFromDBPhoto(ctx context.Context, photos []db.GPhoto) ([]Photo, error) {
	ctx, span := a.tracer.Start(ctx, "fromDBPhoto")
	defer span.End()

	items := []Photo{}
	for _, aa := range photos {
		bh := aa.Image.BlurHash
		url, err := a.ImageStore.GetImageURL(ctx, aa.Image.ID)
		if err != nil {
			return nil, err
		}
		items = append(items, Photo{
			Id:       aa.PhotoID,
			TakenAt:  aa.Image.TakenAt.Ptr(),
			BlurHash: &bh,
			Width:    300,
			Height:   400,
			BaseUrl:  url,
			Source:   PhotoSourceGoogle,
		})
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

	recipesDetails, err := a.DB().GetRecipeDetailWhere(ctx, sq.Eq{"recipe_id": mealRecipes.RecipeIDs()})
	if err != nil {
		return nil, err
	}
	recipeDetailsById := recipesDetails.ByRecipeId()

	for _, m := range meals {
		meal := Meal{Id: m.ID,
			Name:  m.Name,
			AteAt: m.AteAt}

		mrs := []MealRecipe{}
		for _, mr := range mealRecipes.ByMealID()[m.ID] {

			test, err := a.transformRecipes(ctx, recipeDetailsById[mr.RecipeID], true)
			if err != nil {
				return nil, err
			}

			mrs = append(mrs, MealRecipe{Multiplier: mr.Multiplier, Recipe: test[0].Detail})
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
	case MealRecipeUpdateActionAdd:
		err := a.DB().AddRecipeToMeal(ctx, mealId, r.RecipeId, &r.Multiplier)
		if err != nil {
			return handleErr(c, err)
		}
		return a.GetMealById(c, mealId)
	case MealRecipeUpdateActionRemove:
		return sendErr(c, http.StatusBadRequest, fmt.Errorf("unsupported %s", r.Action))
	default:
		return sendErr(c, http.StatusBadRequest, fmt.Errorf("unknown action %s", r.Action))
	}

}
