package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/nickysemenza/gourd/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/exp/slices"
)

func (c *Client) UpsertNotionImages(ctx context.Context, tx *sql.Tx, photos models.NotionImageSlice) error {
	ctx, span := c.tracer.Start(ctx, "UpsertNotionImages")
	defer span.End()

	for _, photo := range photos {
		err := photo.Upsert(ctx, tx, true, []string{"block_id", "page_id"}, boil.Infer(), boil.Infer())
		if err != nil {
			return err
		}

	}
	return nil
}

func (c *Client) DoesNotionImageExist(ctx context.Context, blockID string) (exists bool, err error) {
	res, err := models.NotionImages(qm.Where("block_id = ?", blockID)).One(ctx, c.db)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return res != nil, err
}

func (c *Client) SyncNotionMealFromNotionRecipe(ctx context.Context) error {
	ctx = boil.WithDebug(ctx, true)

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ctx, span := c.tracer.Start(ctx, "SyncNotionMealFromNotionRecipe")
	defer span.End()
	missingMeals, err := models.NotionRecipes(
		qm.Load(models.NotionRecipeRels.Meals),
		qm.Load(models.NotionRecipeRels.Recipe),
	).All(ctx, tx)

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
			tx,
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
			err = c.AddRecipeToMeal(ctx, mealID, m.R.Recipe.ID, mult, tx)
			// m.R.Recipe.AddMealRecipes(ctx, c.db, false, &models.MealRecipe{MealID: mealID})
			if err != nil {
				return fmt.Errorf("failed to add recipe to meal: %w", err)
			}
		}

		err = m.AddMeals(ctx, tx, false, &models.Meal{ID: mealID})
		if err != nil {
			return fmt.Errorf("failed to add meals: %w", err)
		}

	}

	return tx.Commit()

}
