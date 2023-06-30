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

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ctx, span := c.tracer.Start(ctx, "SyncNotionMealFromNotionRecipe")
	defer span.End()
	all, err := models.NotionRecipes(
		qm.Load(qm.Rels(models.NotionRecipeRels.Meals, models.MealRels.MealRecipes)),

		qm.Load(models.NotionRecipeRels.Recipe),
	).All(ctx, tx)

	if err != nil {
		return err
	}

	for _, notionRecipe := range all {
		if !notionRecipe.AteAt.Valid {
			continue
		}
		notionAteAtTime := notionRecipe.AteAt.Time
		if len(notionRecipe.R.Meals) > 0 {

			needsAdjustment := false
			// checking meals for drift

			for _, linkedMeal := range notionRecipe.R.Meals {
				drift := linkedMeal.AteAt.Sub(notionAteAtTime)
				if drift > 0 {
					needsAdjustment = true
					fmt.Printf("time difference between %s and %s (%s) is %s\n", linkedMeal.Name, notionRecipe.PageTitle, notionRecipe.NotionID, drift)

					_, err := models.MealRecipes(qm.Where("meal_id = ? and recipe_id = ?", linkedMeal.ID, notionRecipe.RecipeID)).DeleteAll(ctx, tx)
					if err != nil {
						return err
					}

					err = notionRecipe.RemoveMeals(ctx, tx, linkedMeal)
					if err != nil {
						return err
					}

				}

				// notionAteAtTime.Sub(x.AteAt)
			}
			if !needsAdjustment {
				continue
			}
		}
		var meta NotionRecipeMeta
		err := notionRecipe.Meta.Unmarshal(&meta)
		if err != nil {
			return err
		}
		suffix := "meal"
		for _, t := range meta.Tags {
			if slices.Contains([]string{"dinner", "lunch", "breakfast"}, t) {
				suffix = t
			}
		}

		mealID, err := c.MealIDInRange(ctx, notionRecipe.AteAt.Time,
			fmt.Sprintf("%s %s", notionRecipe.AteAt.Time.Add(time.Hour*-24).Format("Mon Jan 2"), suffix),
			tx,
		)
		if err != nil {
			return fmt.Errorf("failed to find meal in range: %w", err)
		}

		if notionRecipe.R.Recipe != nil {
			var mult *float64
			if !notionRecipe.Scale.IsZero() {
				val, _ := notionRecipe.Scale.Float64()
				mult = &val
			}
			err = c.AddRecipeToMeal(ctx, mealID, notionRecipe.R.Recipe.ID, mult, tx)
			// m.R.Recipe.AddMealRecipes(ctx, c.db, false, &models.MealRecipe{MealID: mealID})
			if err != nil {
				return fmt.Errorf("failed to add recipe to meal: %w", err)
			}
		}

		err = notionRecipe.AddMeals(ctx, tx, false, &models.Meal{ID: mealID})
		if err != nil {
			return fmt.Errorf("failed to add meals: %w", err)
		}

	}

	return tx.Commit()

}
