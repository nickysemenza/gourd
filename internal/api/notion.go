package api

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/nickysemenza/gourd/internal/clients/notion"
	"github.com/nickysemenza/gourd/internal/clients/rs_client"
	"github.com/nickysemenza/gourd/internal/common"
	"github.com/nickysemenza/gourd/internal/db"
	"github.com/nickysemenza/gourd/internal/db/models"
	"github.com/nickysemenza/gourd/internal/image"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.mitsakis.org/workerpool"
)

func (a *API) notionRecipeToInput(ctx context.Context, nRecipe notion.Recipe) (*RecipeDetailInput, error) {
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

	l(ctx).Debugf("notionRecipeToInput: %s (%v)", output.Name, output.Tags)
	return &output, nil
}

// todo: make this transactional
func (a *API) notionRecipeToDB(ctx context.Context, nRecipe notion.Recipe) (*models.NotionRecipe, error) {
	ctx, span := a.tracer.Start(ctx, "notionRecipeToDB")
	defer span.End()

	output, err := a.notionRecipeToInput(ctx, nRecipe)
	if err != nil {
		return nil, err
	}

	r, err := a.CreateRecipe(ctx, &RecipeWrapperInput{Detail: *output})
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
		NotionID:  nRecipe.UID,
		PageID:    nRecipe.PageID,
		PageTitle: nRecipe.Title,
		AteAt:     null.TimeFromPtr(nRecipe.Time),
		RecipeID:  null.StringFrom(r.Id),
	}

	nr.Scale = common.NullDecimalFromFloat(nRecipe.Scale)
	m := db.NotionRecipeMeta{Tags: nRecipe.Tags}
	err = nr.Meta.Marshal(m)
	return &nr, err
}

type notionSyncResult struct {
	nRecipes models.NotionRecipeSlice
	nImages  models.NotionImageSlice
	images   models.ImageSlice
}

func (a *API) processNotionRecipe(ctx context.Context, nRecipe notion.Recipe) (res notionSyncResult, err error) {
	ctx, span := a.tracer.Start(ctx, "processNotionRecipe")
	defer span.End()
	dbnr, err := a.notionRecipeToDB(ctx, nRecipe)
	if err != nil {
		return notionSyncResult{}, err
	}
	res.nRecipes = append(res.nRecipes, dbnr)
	for _, nPhoto := range nRecipe.Photos {

		l := l(ctx).WithField("block_id", nPhoto.BlockID)
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

		res.images = append(res.images, &models.Image{
			ID:       id,
			BlurHash: bh,
			Source:   "notion",
			TakenAt:  null.TimeFromPtr(nRecipe.Time),
		})
		res.nImages = append(res.nImages, &models.NotionImage{
			PageID:   nRecipe.PageID,
			BlockID:  nPhoto.BlockID,
			ImageID:  id,
			LastSeen: time.Now(),
		})
	}
	return
}

func (a *API) syncRecipeFromNotion(ctx context.Context, lookback time.Duration) error {
	ctx, span := a.tracer.Start(ctx, "syncRecipeFromNotion")
	defer span.End()
	nRecipes, err := a.Notion.GetAll(ctx, lookback, "")
	if err != nil {
		return err
	}
	l(ctx).Infof("got %d notion recipes", len(nRecipes))

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
		if result.Error != nil {
			l(ctx).Error(result.Error)
		} else {
			res := result.Value
			summary.nRecipes = append(summary.nRecipes, res.nRecipes...)
			summary.nImages = append(summary.nImages, res.nImages...)
			summary.images = append(summary.images, res.images...)
		}
	}

	tx := a.tx(ctx)
	err = a.db.SaveImage(ctx, tx, summary.images...)
	if err != nil {
		return err
	}
	err = a.saveNotionRecipes(ctx, tx, summary.nRecipes)
	if err != nil {
		return err
	}
	err = a.db.UpsertNotionImages(ctx, tx, summary.nImages)
	if err != nil {
		return err
	}

	res, err := models.NotionRecipes(
		qm.Where("ate_at > ?", time.Now().Add(-lookback)),
		qm.Where("last_seen < ?", time.Now().Add(-time.Hour*12)),
	).DeleteAll(ctx, tx, false)
	if err != nil {
		return err
	}
	l(ctx).Infof("updated %d recipes", len(summary.nRecipes))
	l(ctx).Warnf("deleted %d stale recipes", res)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (a *API) saveNotionRecipes(ctx context.Context, tx *sql.Tx, items models.NotionRecipeSlice) error {
	ctx, span := a.tracer.Start(ctx, "saveNotionRecipes")
	defer span.End()
	for _, r := range items {
		err := r.Upsert(ctx, tx, true,
			[]string{
				models.NotionRecipeColumns.PageID,
			},
			boil.Infer(),
			boil.Infer(),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
