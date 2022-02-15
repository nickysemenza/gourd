package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/nickysemenza/gourd/common"
	"github.com/nickysemenza/gourd/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gopkg.in/guregu/null.v4/zero"
)

func (c *Client) GetKV(ctx context.Context, key string) (string, error) {
	var json string
	q := c.psql.Select("value").From("kv").Where(sq.Eq{"key": key})
	err := c.getContext(ctx, q, &json)
	return json, err
}
func (c *Client) SetKV(ctx context.Context, key string, json string) error {
	q := c.psql.Insert("kv").Columns("key", "value").Values(key, json).Suffix("ON CONFLICT (key) DO UPDATE SET value = ?", json)
	_, err := c.execContext(ctx, q)
	return err
}

type GAlbum struct {
	ID      string `db:"id"`
	Usecase string `db:"usecase"`
}

func (c *Client) GetAlbums(ctx context.Context) ([]GAlbum, error) {
	var albums []GAlbum
	q := c.psql.Select("id", "usecase").From("gphotos_albums")
	err := c.selectContext(ctx, q, &albums)
	return albums, err
}

type GPhoto struct {
	AlbumID string    `db:"album_id"`
	PhotoID string    `db:"id"`
	Created time.Time `db:"creation_time"`
	Seen    time.Time `db:"last_seen"`
	ImageID string    `db:"image_id"`
	// MetadataJSON types.JSONText `db:"media_metadata"`
	Image Image
}

type Image struct {
	ID       string `db:"id"`
	BlurHash string `db:"blur_hash"`
	Source   string `db:"source"`
}

type NotionRecipe struct {
	PageID    string `db:"page_id"`
	PageTitle string `db:"page_title"`
	// Meta      string      `db:"meta"`
	LastSeen time.Time   `db:"last_seen"`
	Recipe   zero.String `db:"recipe_id"`
	AteAt    zero.Time   `db:"ate_at"`
}

func (c *Client) UpsertPhotos(ctx context.Context, photos []GPhoto) error {
	q := c.psql.Insert("gphotos_photos").Columns("id", "album_id", "creation_time", "image_id")
	for _, photo := range photos {
		q = q.Values(photo.PhotoID, photo.AlbumID, photo.Created, photo.ImageID)
	}
	q = q.Suffix("ON CONFLICT (id) DO UPDATE SET last_seen = ?, image_id = excluded.image_id", time.Now())
	_, err := c.execContext(ctx, q)
	return err
}

type NotionImage struct {
	BlockID  string    `db:"block_id"`
	PageID   string    `db:"page_id"`
	ImageID  string    `db:"image_id"`
	LastSeen time.Time `db:"last_seen"`
	Image    Image
}

func (c *Client) UpsertNotionImages(ctx context.Context, photos []NotionImage) error {
	q := c.psql.Insert("notion_image").Columns("block_id", "page_id", "image_id")
	for _, photo := range photos {
		q = q.Values(photo.BlockID, photo.PageID, photo.ImageID)
	}
	q = q.Suffix("ON CONFLICT (block_id,page_id) DO UPDATE SET last_seen = ?, image_id = excluded.image_id", time.Now())
	_, err := c.execContext(ctx, q)
	return err
}

func (c *Client) getPhotos(ctx context.Context, addons func(q sq.SelectBuilder) sq.SelectBuilder) ([]GPhoto, error) {
	ctx, span := c.tracer.Start(ctx, "db.getPhotos")
	defer span.End()
	q := c.psql.Select("id", "album_id", "creation_time", "last_seen", "image_id").From("gphotos_photos").OrderBy("creation_time DESC")
	q = addons(q)
	var results []GPhoto
	err := c.selectContext(ctx, q, &results)
	// return results, err

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
func (c *Client) GetPhotos(ctx context.Context) ([]GPhoto, error) {
	return c.getPhotos(ctx, func(q sq.SelectBuilder) sq.SelectBuilder { return q })
}
func (c *Client) GetAllPhotos(ctx context.Context) (map[string]GPhoto, error) {
	ctx, span := c.tracer.Start(ctx, "db.GetAllPhotos")
	defer span.End()
	photos, err := c.getPhotos(ctx, func(q sq.SelectBuilder) sq.SelectBuilder { return q })
	if err != nil {
		return nil, err
	}
	byId := make(map[string]GPhoto)
	for _, p := range photos {
		byId[p.PhotoID] = p
	}
	return byId, nil
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

func (c *Client) SaveImage(ctx context.Context, items []Image) (err error) {
	ctx, span := c.tracer.Start(ctx, "db.SaveImage")
	defer span.End()

	q := c.psql.Insert("images").Columns("id", "blur_hash", "source")
	for _, r := range items {
		q = q.Values(r.ID, r.BlurHash, r.Source)
	}
	_, err = c.execContext(ctx, q)

	return
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

	for _, r := range items {
		r.LastSeen = time.Now()
		err := r.Upsert(ctx, c.db, true, []string{"page_id"}, boil.Whitelist("page_title", "recipe_id", "last_seen", "meta"), boil.Infer())
		if err != nil {
			return err
		}
	}

	return
}

func (c *Client) SyncMealsFromGPhotos(ctx context.Context) error {
	q := c.psql.Select("id", "album_id", "creation_time").From("gphotos_photos").
		LeftJoin("meal_gphoto on gphotos_photos.id = meal_gphoto.gphotos_id").Where(sq.Eq{"meal_id": nil})
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

func (c *Client) SyncNotionMealFromNotionRecipe(ctx context.Context) error {
	q := c.psql.Select("page_id", "ate_at", "recipe_id").From("notion_recipe").
		LeftJoin("notion_meal on notion_recipe.page_id = notion_meal.notion_recipe").Where(sq.Eq{"meal_id": nil})
	var missingMeals []NotionRecipe
	err := c.selectContext(ctx, q, &missingMeals)
	if err != nil {
		return err
	}

	for _, m := range missingMeals {
		if !m.AteAt.Valid {
			continue
		}

		mealID, err := c.MealIDInRange(ctx, m.AteAt.Time,
			// fmt.Sprintf("meal on %s", m.AteAt.Time),
			m.AteAt.Time.Format("Mon Jan 2"),
		)
		if err != nil {
			return err
		}

		if m.Recipe.Valid {
			err = c.AddRecipeToMeal(ctx, mealID, m.Recipe.ValueOrZero(), 1)
			if err != nil {
				return err
			}
		}

		q := c.psql.Insert("notion_meal").Columns("meal_id", "notion_recipe").Values(mealID, m.PageID)
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

// Returns photos related to notion; google photos are tied to meals not recipes
func (c *Client) GetPhotosWithRecipe(ctx context.Context, recipeID ...string) (images map[string][]Image, err error) {
	ctx, span := c.tracer.Start(ctx, "GetPhotosWithRecipe")
	defer span.End()

	type res struct {
		Recipe string `db:"recipe_id"`
		Image
	}
	res2 := []res{}
	q := c.psql.Select("notion_recipe.recipe_id as recipe_id", "id", "blur_hash", "source").From("images").
		LeftJoin("notion_image on notion_image.image = images.id").
		LeftJoin("notion_recipe on notion_recipe.page_id = notion_image.page_id").
		Where(sq.Eq{"notion_recipe.recipe_id": recipeID})

	err = c.selectContext(ctx, q, &res2)
	images = make(map[string][]Image)
	for _, rec := range res2 {
		images[rec.Recipe] = append(images[rec.Recipe], rec.Image)
	}
	return
}

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

func (c *Client) AddRecipeToMeal(ctx context.Context, mealId, recipeId string, multiplier float64) error {
	ctx, span := c.tracer.Start(ctx, "AddRecipeToMeal")
	defer span.End()

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

func (c *Client) GetPhotosForMeal(ctx context.Context, meal string) ([]GPhoto, error) {
	ctx, span := c.tracer.Start(ctx, "GetPhotosForMeal")
	defer span.End()
	return c.getPhotos(ctx, func(q sq.SelectBuilder) sq.SelectBuilder {
		return q.LeftJoin("meal_gphoto on meal_gphoto.gphotos_id = gphotos_photos.id").
			Where(sq.Eq{"meal_id": meal})
	})
}

func (c *Client) GetNotionPhotosForMeal(ctx context.Context, meal string) ([]NotionImage, error) {
	ctx, span := c.tracer.Start(ctx, "GetNotionPhotosForMeal")
	defer span.End()
	return c.getNotionPhotos(ctx, func(q sq.SelectBuilder) sq.SelectBuilder {
		return q.LeftJoin("notion_recipe on notion_recipe.page_id = notion_image.page_id").
			LeftJoin("notion_meal on notion_meal.notion_recipe = notion_recipe.page_id").
			Where(sq.Eq{"meal_id": meal})
	})
}
