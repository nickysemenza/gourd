package api

import (
	"context"
	"fmt"
	"net/http"

	sq "github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/db"
)

func (a *API) notionPhotosFromDBPhoto(ctx context.Context, photos []db.NotionImage) ([]GooglePhoto, error) {
	items := []GooglePhoto{}
	for _, aa := range photos {
		bh := aa.Image.BlurHash
		// url := aa.Image.ID
		url := a.Manager.ImageStore.GetImageURL(ctx, aa.Image.ID)
		// bh := aa.Image.BlurHash
		// url, err := a.Manager.Notion.ImageFromBlock(ctx, notionapi.BlockID(aa.BlockID))
		// if err != nil {
		// 	return nil, err
		// }
		items = append(items, GooglePhoto{
			Id:       aa.BlockID,
			Created:  aa.LastSeen,
			BlurHash: &bh,
			Width:    300,
			Height:   400,
			BaseUrl:  url,
			Source:   GooglePhotoSourceNotion,
		})
	}
	return items, nil
}

func (a *API) googlePhotosFromDBPhoto(ctx context.Context, photos []db.GPhoto, getURLs bool) ([]GooglePhoto, []string, error) {
	ctx, span := a.tracer.Start(ctx, "fromDBPhoto")
	defer span.End()
	// items := []GooglePhoto{}
	var ids []string
	// for _, p := range photos {
	// 	gp := GooglePhoto{Id: p.PhotoID, Created: p.Created, Source: GooglePhotoSourceGoogle}
	// 	// if p.BlurHash.Valid {
	// 	// 	s := p.BlurHash.String
	// 	// 	gp.BlurHash = &s
	// 	// }
	// 	gp.BlurHash = &p.Image.BlurHash
	// 	items = append(items, gp)
	// 	ids = append(ids, p.PhotoID)
	// }

	// if getURLs {
	// 	results, err := a.Manager.Photos.GetMediaItems(ctx, ids)
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}
	// 	for x, item := range items {
	// 		val, ok := results[item.Id]
	// 		if !ok {
	// 			continue
	// 		}
	// 		items[x].BaseUrl = val.BaseUrl
	// 		items[x].Width = val.MediaMetadata.Width
	// 		items[x].Height = val.MediaMetadata.Height
	// 	}
	// }
	items := []GooglePhoto{}
	for _, aa := range photos {
		ids = append(ids, aa.PhotoID)
		bh := aa.Image.BlurHash
		// url := aa.Image.ID
		url := a.Manager.ImageStore.GetImageURL(ctx, aa.Image.ID)
		// bh := aa.Image.BlurHash
		// url, err := a.Manager.Notion.ImageFromBlock(ctx, notionapi.BlockID(aa.BlockID))
		// if err != nil {
		// 	return nil, err
		// }
		items = append(items, GooglePhoto{
			Id:       aa.PhotoID,
			Created:  aa.Seen,
			BlurHash: &bh,
			Width:    300,
			Height:   400,
			BaseUrl:  url,
			Source:   GooglePhotoSourceGoogle,
		})
	}
	// return items, nil
	return items, ids, nil
}
func (a *API) ListPhotos(c echo.Context, params ListPhotosParams) error {
	ctx := c.Request().Context()
	photos, err := a.Manager.DB().GetPhotos(ctx)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	items, _, err := a.googlePhotosFromDBPhoto(ctx, photos, true)
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

			test := a.transformRecipes(ctx, recipeDetailsById[mr.RecipeID], true)

			mrs = append(mrs, MealRecipe{Multiplier: mr.Multiplier, Recipe: test[0]})
		}
		meal.Recipes = &mrs

		googlePhotosDB, err := a.DB().GetPhotosForMeal(ctx, m.ID)
		if err != nil {
			return nil, err
		}

		googlePhotos, gIDs, err := a.googlePhotosFromDBPhoto(ctx, googlePhotosDB, false)
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
		gphotoIDs = append(gphotoIDs, gIDs...)
		items = append(items, meal)
	}
	// urls, err := a.Manager.Photos.GetMediaItems(ctx, gphotoIDs)
	// if err != nil {
	// 	return nil, err
	// }
	// for x, item := range items {
	// 	// if meals[x].Notion != nil {
	// 	// 	images, _, err := a.Manager.Notion.ImagesFromPage(ctx, notionapi.ObjectID(*meals[x].Notion))
	// 	// 	if err != nil {
	// 	// 		return nil, err
	// 	// 	}
	// 	// 	for _, url := range images {
	// 	// 		items[x].Photos = append(items[x].Photos, GooglePhoto{BaseUrl: url.URL})
	// 	// 	}
	// 	// }
	// 	for y, photo := range item.Photos {
	// 		val, ok := urls[photo.Id]
	// 		if !ok {
	// 			continue
	// 		}
	// 		items[x].Photos[y].BaseUrl = val.BaseUrl
	// 		items[x].Photos[y].Width = val.MediaMetadata.Width
	// 		items[x].Photos[y].Height = val.MediaMetadata.Height
	// 	}
	// }
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
