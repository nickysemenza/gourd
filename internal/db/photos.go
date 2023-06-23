package db

import (
	"context"
	"database/sql"

	"github.com/nickysemenza/gourd/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (c *Client) SaveImage(ctx context.Context, tx *sql.Tx, items ...*models.Image) error {
	ctx, span := c.tracer.Start(ctx, "saveImage")
	defer span.End()
	for _, r := range items {
		if err := r.Insert(ctx, tx, boil.Infer()); err != nil {
			return err
		}
	}
	return nil
}

// TODO?  func that Returns photos related to notion; google photos are tied to meals not recipes

func (c *Client) GetPhotosForMeal(ctx context.Context, meal string) (models.GphotosPhotoSlice, error) {
	ctx, span := c.tracer.Start(ctx, "GetPhotosForMeal")
	defer span.End()

	return models.GphotosPhotos(
		qm.InnerJoin("meal_gphoto on meal_gphoto.gphotos_id = gphotos_photos.id"),
		qm.Where("meal_id = ?", meal),
		qm.Load(
			models.GphotosPhotoRels.Image,
		),
	).All(ctx, c.db)
}

func (c *Client) GetNotionPhotosForMeal(ctx context.Context, meal string) (models.NotionImageSlice, error) {
	ctx, span := c.tracer.Start(ctx, "GetNotionPhotosForMeal")
	defer span.End()
	images, err := models.NotionImages(
		qm.InnerJoin("notion_recipe on notion_recipe.page_id = notion_image.page_id"),
		qm.InnerJoin("notion_meal on notion_meal.notion_id = notion_recipe.notion_id"),
		qm.Where("meal_id = ?", meal),
		qm.Load(
			qm.Rels(
				models.NotionImageRels.Image,
			),
		),
		qm.Load(
			qm.Rels(
				models.NotionImageRels.Page,
				models.NotionRecipeRels.Meals,
			),
		),
	).All(ctx, c.db)
	return images, err

}
