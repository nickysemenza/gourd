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

func (c *Client) AddRecipeToMeal(ctx context.Context, mealId, recipeId string, m *float64) error {
	ctx, span := c.tracer.Start(ctx, "AddRecipeToMeal")
	defer span.End()

	multiplier := 1.0
	if m != nil {
		multiplier = *m
	}

	c.psql.Insert("meals")
	q := c.psql.Insert("meal_recipe").Columns("meal_id", "recipe_id", "multiplier").
		Values(mealId, recipeId, multiplier).
		Suffix("ON CONFLICT (recipe_id,meal_id) DO UPDATE SET multiplier = ?", multiplier)
	_, err := c.execContext(ctx, q)
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

func (c *Client) MealIDInRange(ctx context.Context, t time.Time, name string) (mealID string, err error) {
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
		iq := c.psql.Insert("meals").Columns("id", "ate_at", "name").Values(mealID, t, name)
		_, err = c.execContext(ctx, iq)
	}
	return
}

func (c *Client) SyncMealsFromGPhotos(ctx context.Context) error {
	q := c.psql.Select("gphotos_photos.id as id", "album_id", "creation_time").From("gphotos_photos").
		LeftJoin("meal_gphoto on gphotos_photos.id = meal_gphoto.gphotos_id").
		LeftJoin("gphotos_albums on gphotos_photos.album_id = gphotos_albums.id").
		Where(sq.Eq{"meal_id": nil}).Where(sq.Eq{"usecase": "food"})
	var missingMeals []GPhoto
	err := c.selectContext(ctx, q, &missingMeals)
	if err != nil {
		return err
	}

	for _, m := range missingMeals {

		mealID, err := c.MealIDInRange(ctx, m.Created, fmt.Sprintf("meal on %s", m.Created))
		if err != nil {
			return err
		}

		q := c.psql.Insert("meal_gphoto").Columns("meal_id", "gphotos_id").Values(mealID, m.PhotoID)
		_, err = c.execContext(ctx, q)
		if err != nil {
			return err
		}
	}

	return nil

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
