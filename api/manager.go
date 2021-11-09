package api

import (
	"context"
	"fmt"

	"github.com/nickysemenza/gourd/common"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/image"
	"github.com/nickysemenza/gourd/rs_client"
	"gopkg.in/guregu/null.v4/zero"
)

func (a *API) DB() *db.Client {
	return a.db
}

func (a *API) ProcessGoogleAuth(ctx context.Context, code string) (jwt string, rawUser map[string]interface{}, err error) {
	err = a.Google.Finish(ctx, code)
	if err != nil {
		return
	}
	user, err := a.Google.GetUserInfo(ctx)
	if err != nil {
		return
	}
	rawUser = map[string]interface{}{"raw": user}

	jwt, err = a.Auth.GetJWT(user)
	return
}

func (a *API) syncRecipeFromNotion(ctx context.Context) error {
	ctx, span := a.tracer.Start(ctx, "syncRecipeFromNotion")
	defer span.End()
	nRecipes, err := a.Notion.Dump(ctx)
	if err != nil {
		return err
	}

	var foo []db.NotionRecipe
	var bar []db.NotionImage
	var img []db.Image
	for _, nRecipe := range nRecipes {

		output := RecipeDetailInput{}
		if nRecipe.Raw != "" {
			err = a.R.Call(ctx, nRecipe.Raw, rs_client.RecipeDecode, &output)
			if err != nil {
				return fmt.Errorf("failed to decode recipe: %w", err)
			}
			// output.Sources = &[]api.RecipeSource{{Title: }}
		} else if nRecipe.SourceURL != "" {
			r, err := a.FetchAndTransform(ctx, nRecipe.SourceURL, a.IngredientIdByName)
			if err != nil {
				return err
			}
			output = r.Detail
		}
		output.Name = nRecipe.Title

		r, err := a.CreateRecipe(ctx, &RecipeWrapperInput{Detail: output})
		if err != nil {
			return err
		}
		foo = append(foo, db.NotionRecipe{
			PageID:    nRecipe.PageID,
			PageTitle: nRecipe.Title,
			AteAt:     zero.TimeFromPtr(nRecipe.Time),
			Recipe:    zero.StringFrom(r.Id),
		})
		for _, nPhoto := range nRecipe.Photos {
			bh, image, err := image.GetBlurHash(ctx, nPhoto.URL)
			if err != nil {
				return err
			}
			id := common.ID("notion_image")
			err = a.ImageStore.SaveImage(ctx, id, image)
			if err != nil {
				return err
			}

			i := db.Image{
				ID:       id,
				BlurHash: bh,
				Source:   "notion",
			}
			img = append(img, i)
			bar = append(bar, db.NotionImage{
				PageID:  nRecipe.PageID,
				BlockID: nPhoto.BlockID,
				ImageID: id,
			})
		}
	}

	err = a.db.SaveImage(ctx, img)
	if err != nil {
		return err
	}

	err = a.db.SaveNotionRecipes(ctx, foo)
	if err != nil {
		return err
	}
	err = a.db.UpsertNotionImages(ctx, bar)
	if err != nil {
		return err
	}
	return nil
}

func (a *API) DoSync(ctx context.Context) error {
	ctx, span := a.tracer.Start(ctx, "DoSync")
	defer span.End()

	err := a.syncRecipeFromNotion(ctx)
	if err != nil {
		return err
	}
	err = a.DB().SyncNotionMealFromNotionRecipe(ctx)
	if err != nil {
		return err
	}
	err = a.GPhotos.SyncAlbums(ctx)
	if err != nil {
		return err
	}
	err = a.DB().SyncMealsFromPhotos(ctx)
	if err != nil {
		return err
	}
	return nil
}
