package db

import (
	"context"

	"github.com/nickysemenza/gourd/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (c *Client) GetAlbums(ctx context.Context) (models.GphotosAlbumSlice, error) {
	return models.GphotosAlbums(qm.Load(models.GphotosAlbumRels.AlbumGphotosPhotos)).All(ctx, c.db)

}

func (c *Client) GetPhotos(ctx context.Context) (models.GphotosPhotoSlice, error) {
	return models.GphotosPhotos(qm.Load(models.GphotosPhotoRels.Image)).All(ctx, c.db)
}
