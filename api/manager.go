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

func (m *API) DB() *db.Client {
	return m.db
}

func (m *API) ProcessGoogleAuth(ctx context.Context, code string) (jwt string, rawUser map[string]interface{}, err error) {
	err = m.Google.Finish(ctx, code)
	if err != nil {
		return
	}
	user, err := m.Google.GetUserInfo(ctx)
	if err != nil {
		return
	}
	rawUser = map[string]interface{}{"raw": user}

	jwt, err = m.Auth.GetJWT(user)
	return
}

func (m *API) SyncNotionToMeals(ctx context.Context) error {
	nRecipes, err := m.Notion.Dump(ctx)
	if err != nil {
		return err
	}

	var foo []db.NotionRecipe
	var bar []db.NotionImage
	var img []db.Image
	for _, nRecipe := range nRecipes {

		output := RecipeDetailInput{}
		if nRecipe.Raw != "" {
			err = m.R.Call(ctx, nRecipe.Raw, rs_client.RecipeDecode, &output)
			if err != nil {
				return fmt.Errorf("failed to decode recipe: %w", err)
			}
			// output.Sources = &[]api.RecipeSource{{Title: }}
		} else if nRecipe.SourceURL != "" {
			r, err := m.FetchAndTransform(ctx, nRecipe.SourceURL, m.IngredientIdByName)
			if err != nil {
				return err
			}
			output = r.Detail
		}
		output.Name = nRecipe.Title

		r, err := m.CreateRecipe(ctx, &RecipeWrapperInput{Detail: output})
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
			err = m.ImageStore.SaveImage(ctx, id, image)
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

	err = m.db.SaveImage(ctx, img)
	if err != nil {
		return err
	}

	err = m.db.SaveNotionRecipes(ctx, foo)
	if err != nil {
		return err
	}
	err = m.db.UpsertNotionImages(ctx, bar)
	if err != nil {
		return err
	}
	return nil
}
