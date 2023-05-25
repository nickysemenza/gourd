package db

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (c *Client) SaveImage(ctx context.Context, items ...Image) (err error) {
	ctx, span := c.tracer.Start(ctx, "db.SaveImage")
	defer span.End()

	if len(items) == 0 {
		return nil
	}

	q := c.psql.Insert("images").Columns("id", "blur_hash", "source", "taken_at")
	for _, r := range items {
		q = q.Values(r.ID, r.BlurHash, r.Source, r.TakenAt)
	}
	_, err = c.execContext(ctx, q)

	return
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
			LeftJoin("notion_meal on notion_meal.notion_id = notion_recipe.notion_id").
			Where(sq.Eq{"meal_id": meal})
	})
}
