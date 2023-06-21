package api

import (
	"context"

	"github.com/nickysemenza/gourd/internal/db/models"
	// . "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// func (a *API) imagesFromRecipeDetailId(ctx context.Context, id string) ([]Photo, error) {
// 	rd, err := models.RecipeDetails(
// 		Where("recipe_details.id = ?", id),
// 		Load(Rels(
// 			models.RecipeDetailRels.Recipe,
// 			models.RecipeRels.NotionRecipes,
// 			models.NotionRecipeRels.PageNotionImages,
// 			models.NotionImageRels.Image,
// 		)),
// 		Load(Rels(models.RecipeDetailRels.Recipe,
// 			models.RecipeRels.MealRecipes,
// 			models.MealRecipeRels.Meal,
// 			models.MealRels.MealGphotos,
// 			models.MealGphotoRels.Gphoto,
// 			models.GphotosPhotoRels.Image,
// 		)),
// 	).
// 		One(ctx, a.db.DB())
// 	if err != nil {
// 		return nil, err
// 	}
// 	gp := models.GphotosPhotoSlice{}
// 	for _, m := range rd.R.Recipe.R.MealRecipes {
// 		for _, x := range m.R.Meal.R.MealGphotos {
// 			gp = append(gp, x.R.Gphoto)
// 		}

// 	}

// 	return a.imagesFromModel(ctx, rd.R.Recipe.R.NotionRecipes, gp)

// }

func (a *API) imageFromModel(ctx context.Context, p *models.Image) (*Photo, error) {
	if p == nil {
		return nil, nil
	}
	url, err := a.ImageStore.GetImageURL(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	return &Photo{
		Id:       p.ID,
		TakenAt:  p.TakenAt.Ptr(),
		BlurHash: &p.BlurHash,
		Width:    300,
		Height:   400,
		BaseUrl:  url,
		Source:   PhotoSource(p.Source),
	}, nil

}

func (a *API) imagesFromModel(ctx context.Context, nr models.NotionRecipeSlice, gp models.GphotosPhotoSlice) ([]Photo, error) {
	items := []Photo{}
	for _, notionRecipe := range nr {
		for _, notionImage := range notionRecipe.R.PageNotionImages {
			photo, err := a.imageFromModel(ctx, notionImage.R.Image)
			if err != nil {
				return nil, err
			}
			items = append(items, *photo)
		}
	}

	for _, gPhoto := range gp {

		image, err := a.imageFromModel(ctx, gPhoto.R.Image)
		if err != nil {
			return nil, err
		}
		items = append(items, *image)
	}

	return items, nil
}
