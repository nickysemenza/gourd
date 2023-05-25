package db

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/volatiletech/null/v8"
)

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
	ID       string    `db:"id"`
	BlurHash string    `db:"blur_hash"`
	Source   string    `db:"source"`
	TakenAt  null.Time `db:"taken_at"`
}

func (c *Client) UpsertGPhotos(ctx context.Context, photos ...GPhoto) error {
	q := c.psql.Insert("gphotos_photos").Columns("id", "album_id", "creation_time", "image_id")
	for _, photo := range photos {
		q = q.Values(photo.PhotoID, photo.AlbumID, photo.Created, photo.ImageID)
	}
	q = q.Suffix("ON CONFLICT (id) DO UPDATE SET last_seen = ?, image_id = excluded.image_id", time.Now())
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
	q = c.psql.Select("id", "blur_hash", "source", "taken_at").From("images").Where(sq.Eq{"id": ids})
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
