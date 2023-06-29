package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/nickysemenza/gourd/internal/common"
	"github.com/nickysemenza/gourd/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (c *Client) GetAllMeals(ctx context.Context) (Meals, error) {
	ctx, span := c.tracer.Start(ctx, "GetAllMeals")
	defer span.End()
	q := c.psql.Select("id", "name", "ate_at").From("meals").OrderBy("ate_at DESC")
	var results Meals
	err := c.selectContext(ctx, q, &results)
	return results, err
}

func (c *Client) GetMealById(ctx context.Context, id string) (*Meal, error) {
	ctx, span := c.tracer.Start(ctx, "GetMealById")
	defer span.End()
	q := c.psql.Select("id", "name", "ate_at").From("meals")
	var result Meal
	err := c.getContext(ctx, q, &result)
	return &result, err
}

func (c *Client) AddRecipeToMeal(ctx context.Context, mealId, recipeId string, m *float64, tx *sql.Tx) error {
	ctx, span := c.tracer.Start(ctx, "AddRecipeToMeal")
	defer span.End()

	mr := models.MealRecipe{
		MealID:     mealId,
		RecipeID:   recipeId,
		Multiplier: common.NullDecimalFromFloat(m),
	}
	err := mr.Upsert(ctx, tx, true,
		[]string{models.MealRecipeColumns.RecipeID, models.MealRecipeColumns.MealID},
		boil.Whitelist(models.MealRecipeColumns.Multiplier),
		boil.Infer())
	return err
}

type MealRecipe struct {
	MealID     string  `db:"meal_id"`
	RecipeID   string  `db:"recipe_id"`
	Multiplier float64 `db:"multiplier"`
}
type MealRecipes []MealRecipe

func (r MealRecipes) ByMealID() map[string][]MealRecipe {
	m := make(map[string][]MealRecipe)
	for _, x := range r {
		m[x.MealID] = append(m[x.MealID], x)
	}
	return m
}
func (r MealRecipes) RecipeIDs() []string {
	m := []string{}
	for _, x := range r {
		m = append(m, x.RecipeID)
	}
	return m
}

func (c *Client) GetMealRecipes(ctx context.Context, mealID ...string) (MealRecipes, error) {
	ctx, span := c.tracer.Start(ctx, "GetMealRecipes")
	defer span.End()
	q := c.psql.Select("*").From("meal_recipe").Where(sq.Eq{"meal_id": mealID})
	var results MealRecipes
	err := c.selectContext(ctx, q, &results)
	return results, err
}

func (c *Client) MealIDInRange(ctx context.Context, t time.Time, name string, tx *sql.Tx) (mealID string, err error) {
	ctx, span := c.tracer.Start(ctx, "db.MealIDInRange")
	defer span.End()

	err = c.db.GetContext(ctx, &mealID, `select id from meals
WHERE ate_at > $1::timestamp - INTERVAL '1 hour'
AND ate_at < $1::timestamp + INTERVAL '1 hour' limit 1`, pq.FormatTimestamp(t))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return
		}
		// err no rows, need to insert
		mealID = common.ID("m")

		newMeal := models.Meal{
			ID:    mealID,
			AteAt: t,
			Name:  name,
		}
		err = newMeal.Insert(ctx, tx, boil.Infer())
		if err != nil {
			err = fmt.Errorf("failed to insert meal %v: %w", newMeal, err)
		}
	}
	return
}

func (c *Client) SyncMealsFromGPhotos(ctx context.Context) error {

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	missingMeals, err := models.GphotosPhotos(
		qm.InnerJoin("meal_gphoto on gphotos_photos.id = meal_gphoto.gphotos_id"),
		qm.InnerJoin("gphotos_albums on gphotos_photos.album_id = gphotos_albums.id"),
		qm.Where("meal_id IS NULL"),
		qm.Where("usecase = ?", "food"),
		qm.Load(
			models.GphotosPhotoRels.Image,
		),
	).All(ctx, tx)

	if err != nil {
		return err
	}

	for _, m := range missingMeals {

		mealID, err := c.MealIDInRange(ctx, m.CreationTime, fmt.Sprintf("meal on %s", m.CreationTime), tx)
		if err != nil {
			return err
		}

		photo := models.MealGphoto{
			MealID:    mealID,
			GphotosID: m.ID,
		}
		err = photo.Insert(ctx, tx, boil.Infer())

		if err != nil {
			return err
		}
	}

	return tx.Commit()

}

type Meal struct {
	ID    string    `db:"id"`
	Name  string    `db:"name"`
	AteAt time.Time `db:"ate_at"`
}
type Meals []Meal

func (r Meals) MealIDs() []string {
	m := []string{}
	for _, x := range r {
		m = append(m, x.ID)
	}
	return m
}

func (c *Client) GetMealsWithRecipe(ctx context.Context, recipeID string) (Meals, error) {
	ctx, span := c.tracer.Start(ctx, "GetMealsWithRecipe")
	defer span.End()
	q := c.psql.Select("id", "name", "ate_at").
		From("meals").
		LeftJoin("meal_recipe on meals.id = meal_recipe.meal_id").
		// LeftJoin("recipe_details on recipe_details.recipe = meal_recipe.recipe_id").
		Where(sq.Eq{"meal_recipe.recipe_id": recipeID}).
		OrderBy("ate_at DESC")
	var results Meals
	err := c.selectContext(ctx, q, &results)
	return results, err
}
