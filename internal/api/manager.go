package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mitsakis.org/workerpool"

	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/internal/clients/notion"
	"github.com/nickysemenza/gourd/internal/clients/rs_client"
	"github.com/nickysemenza/gourd/internal/common"
	"github.com/nickysemenza/gourd/internal/db"
	"github.com/nickysemenza/gourd/internal/db/models"
	"github.com/nickysemenza/gourd/internal/image"
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

type notionSyncResult struct {
	nRecipes []models.NotionRecipe
	nImages  []db.NotionImage
	images   []db.Image
}

func (a *API) processNotionRecipe(ctx context.Context, nRecipe notion.Recipe) (res notionSyncResult, err error) {
	ctx, span := a.tracer.Start(ctx, "processNotionRecipe")
	defer span.End()
	dbnr, err := a.notionRecipeToDB(ctx, nRecipe)
	if err != nil {
		return notionSyncResult{}, err
	}
	res.nRecipes = append(res.nRecipes, *dbnr)
	for _, nPhoto := range nRecipe.Photos {

		l := log.WithField("block_id", nPhoto.BlockID)
		// nPhoto.BlockID
		exists, err := a.db.DoesNotionImageExist(ctx, nPhoto.BlockID)
		if err != nil {
			return notionSyncResult{}, err
		}
		if exists {
			l.Println("already exists")
			continue
		}

		bh, image, err := image.GetBlurHash(ctx, nPhoto.URL)
		if err != nil {
			return notionSyncResult{}, err
		}
		id := common.ID("notion_image")
		err = a.ImageStore.SaveImage(ctx, id, image)
		if err != nil {
			return notionSyncResult{}, err
		}

		res.images = append(res.images, db.Image{
			ID:       id,
			BlurHash: bh,
			Source:   "notion",
		})
		res.nImages = append(res.nImages, db.NotionImage{
			PageID:  nRecipe.PageID,
			BlockID: nPhoto.BlockID,
			ImageID: id,
		})
	}
	return
}

func (a *API) syncRecipeFromNotion(ctx context.Context, lookbackDays int) error {
	ctx, span := a.tracer.Start(ctx, "syncRecipeFromNotion")
	defer span.End()
	nRecipes, err := a.Notion.GetAll(ctx, lookbackDays, "")
	if err != nil {
		return err
	}
	log.Infof("got %d notion recipes", len(nRecipes))

	p, _ := workerpool.NewPoolWithResults(8, func(job workerpool.Job[notion.Recipe], workerID int) (notionSyncResult, error) {
		return a.processNotionRecipe(ctx, job.Payload)
	})
	go func() {
		for _, nRecipe := range nRecipes {
			p.Submit(nRecipe)
		}
		p.StopAndWait()
	}()

	summary := notionSyncResult{}

	for result := range p.Results {
		res := result.Value
		summary.nRecipes = append(summary.nRecipes, res.nRecipes...)
		summary.nImages = append(summary.nImages, res.nImages...)
		summary.images = append(summary.images, res.images...)
	}
	err = a.db.SaveImage(ctx, summary.images)
	if err != nil {
		return err
	}
	err = a.db.SaveNotionRecipes(ctx, summary.nRecipes)
	if err != nil {
		return err
	}
	err = a.db.UpsertNotionImages(ctx, summary.nImages)
	if err != nil {
		return err
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
