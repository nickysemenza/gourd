package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/clients/notion"
	"github.com/nickysemenza/gourd/clients/rs_client"
	"github.com/nickysemenza/gourd/common"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/image"
	"github.com/nickysemenza/gourd/models"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"
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

func (a *API) notionRecipeToDB(ctx context.Context, nRecipe notion.Recipe) (*models.NotionRecipe, error) {
	ctx, span := a.tracer.Start(ctx, "notionRecipeToDB")
	defer span.End()

	output := RecipeDetailInput{}
	if nRecipe.Raw != "" {
		err := a.R.Call(ctx, nRecipe.Raw, rs_client.RecipeDecode, &output)
		if err != nil {
			return nil, fmt.Errorf("failed to decode recipe: %w", err)
		}
		// output.Sources = &[]api.RecipeSource{{Title: }}
	} else if nRecipe.SourceURL != "" {
		r, err := a.FetchAndTransform(ctx, nRecipe.SourceURL, a.IngredientIdByName)
		if err != nil {
			return nil, err
		}
		output = r.Detail
	}
	output.Name = nRecipe.Title
	output.Date = nRecipe.Time
	nRecipe.Tags = append(nRecipe.Tags, "notion")
	output.Tags = nRecipe.Tags

	r, err := a.CreateRecipe(ctx, &RecipeWrapperInput{Detail: output})
	if err != nil {
		return nil, err
	}

	for _, child := range nRecipe.Children {
		//todo: child images aren't saved
		_, err := a.notionRecipeToDB(ctx, child)
		if err != nil {
			return nil, err
		}
	}
	nr := models.NotionRecipe{
		PageID:    nRecipe.PageID,
		PageTitle: nRecipe.Title,
		AteAt:     null.TimeFromPtr(nRecipe.Time),
		RecipeID:  null.StringFrom(r.Id),
	}
	m := db.NotionRecipeMeta{Tags: nRecipe.Tags}
	err = nr.Meta.Marshal(m)
	return &nr, err
}
func (a *API) syncRecipeFromNotion(ctx context.Context, lookbackDays int) error {
	ctx, span := a.tracer.Start(ctx, "syncRecipeFromNotion")
	defer span.End()
	nRecipes, err := a.Notion.GetAll(ctx, lookbackDays, "")
	if err != nil {
		return err
	}
	log.Infof("got %d notion recipes", len(nRecipes))

	var notionRecipes []models.NotionRecipe
	var notionImages []db.NotionImage
	var images []db.Image
	for _, nRecipe := range nRecipes {
		dbnr, err := a.notionRecipeToDB(ctx, nRecipe)
		if err != nil {
			return err
		}
		notionRecipes = append(notionRecipes, *dbnr)
		for _, nPhoto := range nRecipe.Photos {

			l := log.WithField("block_id", nPhoto.BlockID)
			// nPhoto.BlockID
			exists, err := a.db.DoesNotionImageExist(ctx, nPhoto.BlockID)
			if err != nil {
				return err
			}
			if exists {
				l.Println("already exists")
				continue
			}

			bh, image, err := image.GetBlurHash(ctx, nPhoto.URL)
			if err != nil {
				return err
			}
			id := common.ID("notion_image")
			err = a.ImageStore.SaveImage(ctx, id, image)
			if err != nil {
				return err
			}

			images = append(images, db.Image{
				ID:       id,
				BlurHash: bh,
				Source:   "notion",
			})
			notionImages = append(notionImages, db.NotionImage{
				PageID:  nRecipe.PageID,
				BlockID: nPhoto.BlockID,
				ImageID: id,
			})
		}
	}

	if len(images) > 0 {
		err = a.db.SaveImage(ctx, images)
		if err != nil {
			return err
		}
	}

	if len(notionRecipes) > 0 {
		err = a.db.SaveNotionRecipes(ctx, notionRecipes)
		if err != nil {
			return err
		}
	}
	if len(notionImages) > 0 {
		err = a.db.UpsertNotionImages(ctx, notionImages)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *API) DoSync(c echo.Context, params DoSyncParams) error {
	ctx := c.Request().Context()
	err := a.Sync(ctx, params.LookbackDays)
	if err != nil {
		return handleErr(c, err)
	}
	return c.JSON(http.StatusOK, nil)
}
func (a *API) Sync(ctx context.Context, lookbackDays int) error {
	now := time.Now()
	ctx, span := a.tracer.Start(ctx, "DoSync")
	defer span.End()

	if err := a.syncRecipeFromNotion(ctx, lookbackDays); err != nil {
		return fmt.Errorf("notion: %w", err)
	}
	log.Infof("synced recipes from notion")
	if err := a.DB().SyncNotionMealFromNotionRecipe(ctx); err != nil {
		return fmt.Errorf("notion meal: %w", err)
	}
	if a.GPhotos != nil {
		if err := a.GPhotos.SyncAlbums(ctx); err != nil {
			return fmt.Errorf("gphotos: %w", err)
		}
	}
	if err := a.DB().SyncMealsFromGPhotos(ctx); err != nil {
		return fmt.Errorf("gphotos meal: %w", err)
	}
	log.Infof("sync complete in %s", time.Since(now))
	return nil
}
