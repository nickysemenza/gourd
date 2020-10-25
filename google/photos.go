package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/image"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/api/global"
	"gopkg.in/guregu/null.v3/zero"
)

// https://developers.google.com/photos/library/reference/rest/v1/mediaItems/batchGet#query-parameters
const maxPhotoBatchGet = 50

func (c *Client) getPhotosClient(ctx context.Context) (*gphotos.Client, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("bad token: %w", err)
	}
	tc := c.oc.Client(ctx, token)
	return gphotos.NewClient(tc)
}
func (c *Client) batchGet(ctx context.Context, ids []string) ([]photoslibrary.MediaItem, error) {
	ctx, span := global.Tracer("google").Start(ctx, "google.batchGet")
	defer span.End()
	if size := len(ids); size > maxPhotoBatchGet {
		return nil, fmt.Errorf("requested %d, limit is %d", size, maxPhotoBatchGet)
	}
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("bad token: %w", err)
	}
	tc := c.oc.Client(ctx, token)

	url := "https://photoslibrary.googleapis.com/v1/mediaItems:batchGet?mediaItemIds=" + strings.Join(ids, "&mediaItemIds=")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	res, err := tc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var batchResult BatchGetResult
	json.Unmarshal(body, &batchResult)
	var items []photoslibrary.MediaItem
	for _, b := range batchResult.MediaItemResults {
		items = append(items, b.MediaItem)
	}
	span.SetAttribute("result-raw", batchResult)
	return items, nil
}
func (c *Client) GetMediaItems(ctx context.Context, ids []string) (map[string]photoslibrary.MediaItem, error) {
	ctx, span := global.Tracer("google").Start(ctx, "google.GetMediaItems")
	defer span.End()
	chunks := chunkBy(ids, maxPhotoBatchGet)
	urls := map[string]photoslibrary.MediaItem{}
	wg := sync.WaitGroup{}
	var m sync.Mutex
	for _, chunk := range chunks {
		wg.Add(1)
		go func(chunk []string) {
			items, err := c.batchGet(ctx, chunk)
			if err != nil {
				log.Error(err)
			}
			m.Lock()
			for _, item := range items {
				urls[item.Id] = item
			}
			m.Unlock()
			wg.Done()
		}(chunk)
	}

	wg.Wait()

	return urls, nil

}
func (c *Client) GetTest(ctx context.Context) (interface{}, error) {
	ids := []string{"AIbigFqimylf7SUbAvFeiDPDBg_K_rH5DYtsZUMAiD2yMhJDeHIadDYJnc2Q7vnqKT4DQJeB5IQ7qNEk1Iu0-9k9lfolG6i9-A",
		"AIbigFrtZdWe-gFN1KOuPBhPlFsSNftIy2tyH0yW3JxQPALG-qPg1BsByn12LwoUM_om-DI_rB7OLwhZ8UpzPBxStrlbb9_SQQ"}
	return c.GetMediaItems(ctx, ids)

}
func (c *Client) GetAvailableAlbums(ctx context.Context) ([]photoslibrary.Album, error) {
	client, err := c.getPhotosClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("bad client: %w", err)
	}
	var albums []photoslibrary.Album
	err = client.Albums.List().PageSize(maxPhotoBatchGet).Pages(ctx, func(r *photoslibrary.ListAlbumsResponse) error {
		for _, a := range r.Albums {
			albums = append(albums, *a)
		}
		return nil
	})
	return albums, err
}

func (c *Client) SyncAlbums(ctx context.Context) error {
	ctx, span := global.Tracer("google").Start(ctx, "google.SyncAlbums")
	defer span.End()

	client, err := c.getPhotosClient(ctx)
	if err != nil {
		return fmt.Errorf("bad client: %w", err)
	}

	dbPhotos, err := c.db.GetAllPhotos(ctx)
	if err != nil {
		return err
	}
	albums, err := c.db.GetAlbums(ctx)
	if err != nil {
		return err
	}
	log.Infof("syncing %d album(s)", len(albums))
	numBlurHashCalculated := 0
	for _, album := range albums {
		var photos []db.Photo
		err = client.MediaItems.Search(&photoslibrary.SearchMediaItemsRequest{
			AlbumId:  album.ID,
			PageSize: maxPhotoBatchGet,
		}).Pages(ctx, func(r *photoslibrary.SearchMediaItemsResponse) error {
			for _, m := range r.MediaItems {
				t, err := time.Parse(time.RFC3339, m.MediaMetadata.CreationTime)
				if err != nil {
					return err
				}
				var bh string
				if dbPhoto, ok := dbPhotos[m.Id]; ok && dbPhoto.BlurHash.Valid && dbPhoto.BlurHash.String != "" {
					bh = dbPhoto.BlurHash.String
				} else {
					bh, err = image.GetBlurHash(ctx, m.BaseUrl)
					if err != nil {
						return err
					}
					numBlurHashCalculated++
				}

				photos = append(photos, db.Photo{
					AlbumID:  album.ID,
					PhotoID:  m.Id,
					Created:  t,
					BlurHash: zero.StringFrom(bh),
				})
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = c.db.UpsertPhotos(ctx, photos)
		if err != nil {
			return err
		}
		log.Infof(
			"synced %d photos from album %s. Calcualted %d blur hashes",
			len(photos), album, numBlurHashCalculated)
	}
	return nil
}

type BatchGetResult struct {
	MediaItemResults []struct {
		MediaItem photoslibrary.MediaItem `json:"mediaItem"`
	} `json:"mediaItemResults"`
}
