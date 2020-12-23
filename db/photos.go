package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/nickysemenza/gourd/common"
	"gopkg.in/guregu/null.v3/zero"
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

type Photo struct {
	AlbumID  string      `db:"album_id"`
	PhotoID  string      `db:"id"`
	Created  time.Time   `db:"creation_time"`
	Seen     time.Time   `db:"last_seen"`
	BlurHash zero.String `db:"blur_hash"`
	// MetadataJSON types.JSONText `db:"media_metadata"`
}

func (c *Client) UpsertPhotos(ctx context.Context, photos []Photo) error {
	q := c.psql.Insert("gphotos_photos").Columns("id", "album_id", "creation_time", "blur_hash")
	for _, photo := range photos {
		q = q.Values(photo.PhotoID, photo.AlbumID, photo.Created, photo.BlurHash)
	}
	q = q.Suffix("ON CONFLICT (id) DO UPDATE SET last_seen = ?, blur_hash = excluded.blur_hash", time.Now())
	_, err := c.execContext(ctx, q)
	return err
}

func (c *Client) getPhotos(ctx context.Context, addons func(q sq.SelectBuilder) sq.SelectBuilder) ([]Photo, error) {
	ctx, span := c.tracer.Start(ctx, "db.getPhotos")
	defer span.End()
	q := c.psql.Select("id", "album_id", "creation_time", "last_seen", "blur_hash").From("gphotos_photos").OrderBy("creation_time DESC")
	q = addons(q)
	var results []Photo
	err := c.selectContext(ctx, q, &results)
	return results, err
}
func (c *Client) GetPhotos(ctx context.Context) ([]Photo, error) {
	return c.getPhotos(ctx, func(q sq.SelectBuilder) sq.SelectBuilder { return q })
}
func (c *Client) GetAllPhotos(ctx context.Context) (map[string]Photo, error) {
	ctx, span := c.tracer.Start(ctx, "db.GetAllPhotos")
	defer span.End()
	photos, err := c.getPhotos(ctx, func(q sq.SelectBuilder) sq.SelectBuilder { return q })
	if err != nil {
		return nil, err
	}
	byId := make(map[string]Photo)
	for _, p := range photos {
		byId[p.PhotoID] = p
	}
	return byId, nil
}

func (c *Client) SyncMealsFromPhotos(ctx context.Context) error {
	q := c.psql.Select("id", "album_id", "creation_time").From("gphotos_photos").
		LeftJoin("meal_photo on gphotos_photos.id = meal_photo.gphotos_id").Where(sq.Eq{"meal": nil})
	var missingMeals []Photo
	err := c.selectContext(ctx, q, &missingMeals)
	if err != nil {
		return err
	}

	for _, m := range missingMeals {
		target := pq.FormatTimestamp(m.Created)
		// select mealID from meals WHERE ate_at > now() - INTERVAL '1 hour' AND ate_at < now() + INTERVAL '1 hour' limit 1
		// q = c.psql.Select("mealID").From("meals").
		// 	Where(sq.GtOrEq{"ate_at": fmt.Sprintf("timestamp '%s' - INTERVAL '1 hour'", target)}).
		// 	Limit(1)
		var mealID string
		// err := c.getContext(ctx, q, id)
		err := c.db.GetContext(ctx, &mealID, `select id from meals
WHERE ate_at > $1::timestamp - INTERVAL '1 hour'
AND ate_at < $1::timestamp + INTERVAL '1 hour' limit 1`, target)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			// insert
			mealID = common.UUID()
			iq := c.psql.Insert("meals").Columns("id", "ate_at", "name").Values(mealID, m.Created, mealID)
			_, err := c.execContext(ctx, iq)
			if err != nil {
				return err
			}
		}

		q := c.psql.Insert("meal_photo").Columns("meal", "gphotos_id").Values(mealID, m.PhotoID)
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

func (c *Client) GetAllMeals(ctx context.Context) ([]Meal, error) {
	ctx, span := c.tracer.Start(ctx, "GetAllMeals")
	defer span.End()
	q := c.psql.Select("id", "name", "ate_at").From("meals").OrderBy("ate_at DESC")
	var results []Meal
	err := c.selectContext(ctx, q, &results)
	return results, err
}

func (c *Client) GetPhotosForMeal(ctx context.Context, meal string) ([]Photo, error) {
	ctx, span := c.tracer.Start(ctx, "GetPhotosForMeal")
	defer span.End()
	return c.getPhotos(ctx, func(q sq.SelectBuilder) sq.SelectBuilder {
		return q.LeftJoin("meal_photo on meal_photo.gphotos_id = gphotos_photos.id").
			Where(sq.Eq{"meal": meal})
	})
}
