package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/nickysemenza/gourd/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/exp/slices"
)

type NotionImage struct {
	BlockID  string    `db:"block_id"`
	PageID   string    `db:"page_id"`
	ImageID  string    `db:"image_id"`
	LastSeen time.Time `db:"last_seen"`
	Image    Image
}

func (c *Client) UpsertNotionImages(ctx context.Context, photos []models.NotionImage) error {
	ctx, span := c.tracer.Start(ctx, "UpsertNotionImages")
	defer span.End()

	for _, photo := range photos {
		err := photo.Upsert(ctx, c.db, true, []string{"block_id", "page_id"}, boil.Infer(), boil.Infer())
		if err != nil {
			return err
		}

	}
	return nil
}

func (c *Client) getNotionPhotos(ctx context.Context, addons func(q sq.SelectBuilder) sq.SelectBuilder) ([]NotionImage, error) {
	ctx, span := c.tracer.Start(ctx, "db.getPhotos")
	defer span.End()
	q := c.psql.Select("block_id", "notion_image.page_id", "notion_image.last_seen", "image_id").From("notion_image").OrderBy("last_seen DESC")
	q = addons(q)
	var results []NotionImage
	err := c.selectContext(ctx, q, &results)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, r := range results {
		ids = append(ids, r.ImageID)
	}

	images := []Image{}
	q = c.psql.Select("id", "blur_hash", "source").From("images").Where(sq.Eq{"id": ids})
	err = c.selectContext(ctx, q, &images)
	if err != nil {
		return nil, err
	}
	for i, r := range results {
		for _, img := range images {
			if img.ID == r.ImageID {
				results[i].Image = img
			}
		}
	}
	return results, err
}

func (c *Client) DoesNotionImageExist(ctx context.Context, blockID string) (exists bool, err error) {
	res, err := models.NotionImages(qm.Where("block_id = ?", blockID)).One(ctx, c.db)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return res != nil, err
}

func (c *Client) SaveNotionRecipes(ctx context.Context, items []models.NotionRecipe) (err error) {
	ctx, span := c.tracer.Start(ctx, "db.SaveNotionRecipes")
	defer span.End()

	if len(items) == 0 {
		return nil
	}

	for _, r := range items {
		r.LastSeen = time.Now()
		err := r.Upsert(ctx, c.db, true, []string{"page_id"}, boil.Whitelist("page_title", "recipe_id", "last_seen", "meta", "scale"), boil.Infer())
		if err != nil {
			return err
		}
	}

	return
}

func (c *Client) SyncNotionMealFromNotionRecipe(ctx context.Context) error {
	ctx = boil.WithDebug(ctx, true)

	ctx, span := c.tracer.Start(ctx, "SyncNotionMealFromNotionRecipe")
	defer span.End()
	missingMeals, err := models.NotionRecipes(
		qm.Load(models.NotionRecipeRels.Meals),
		qm.Load(models.NotionRecipeRels.Recipe),
		// Load(models.)
		// qm.SQL("select * from notion_recipe left join notion_meal on notion_meal.notion_id = notion_recipe.notion_id where notion_meal.meal_id is null"),
	).All(ctx, c.db)

	if err != nil {
		return err
	}

	for _, m := range missingMeals {
		if !m.AteAt.Valid {
			continue
		}
		if len(m.R.Meals) > 0 {
			continue
		}
		var meta NotionRecipeMeta
		err := m.Meta.Unmarshal(&meta)
		if err != nil {
			return err
		}
		suffix := "meal"
		for _, t := range meta.Tags {
			if slices.Contains([]string{"dinner", "lunch", "breakfast"}, t) {
				suffix = t
			}
		}

		mealID, err := c.MealIDInRange(ctx, m.AteAt.Time,
			fmt.Sprintf("%s %s", m.AteAt.Time.Add(time.Hour*-24).Format("Mon Jan 2"), suffix),
		)
		if err != nil {
			return fmt.Errorf("failed to find meal in range: %w", err)
		}

		if m.R.Recipe != nil {
			var mult *float64
			if !m.Scale.IsZero() {
				val, _ := m.Scale.Float64()
				mult = &val
			}
			err = c.AddRecipeToMeal(ctx, mealID, m.R.Recipe.ID, mult)
			// m.R.Recipe.AddMealRecipes(ctx, c.db, false, &models.MealRecipe{MealID: mealID})
			if err != nil {
				return fmt.Errorf("failed to add recipe to meal: %w", err)
			}
		}

		err = m.AddMeals(ctx, c.db, false, &models.Meal{ID: mealID})
		if err != nil {
			return fmt.Errorf("failed to add meals: %w", err)
		}

	}

	return nil

}
