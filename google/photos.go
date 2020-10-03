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
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/api/global"
)

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
	return items, nil
}
func (c *Client) GetBaseURLs(ctx context.Context, ids []string) (map[string]string, error) {
	ctx, span := global.Tracer("google").Start(ctx, "google.GetBaseURLs")
	defer span.End()
	chunks := chunkBy(ids, maxPhotoBatchGet)
	urls := map[string]string{}
	// for _, chunk := range chunks {
	// 	items, err := c.batchGet(ctx, chunk)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	for _, item := range items {
	// 		urls[item.Id] = item.BaseUrl
	// 	}
	// }
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
				urls[item.Id] = item.BaseUrl
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
	return c.GetBaseURLs(ctx, ids)

}

func (c *Client) SyncAlbums(ctx context.Context) error {
	client, err := c.getPhotosClient(ctx)
	if err != nil {
		return fmt.Errorf("bad client: %w", err)
	}

	albums, err := c.db.GetAlbums(ctx)
	if err != nil {
		return err
	}
	for _, album := range albums {
		var photos []db.Photo
		err = client.MediaItems.Search(&photoslibrary.SearchMediaItemsRequest{
			AlbumId:  album,
			PageSize: 50,
		}).Pages(ctx, func(r *photoslibrary.SearchMediaItemsResponse) error {
			for _, m := range r.MediaItems {
				t, err := time.Parse(time.RFC3339, m.MediaMetadata.CreationTime)
				if err != nil {
					return err
				}
				photos = append(photos, db.Photo{
					AlbumID: album,
					PhotoID: m.Id,
					Created: t,
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
		log.Infof("synced %d photos from album %s", len(photos), album)
	}
	return nil
}

type BatchGetResult struct {
	MediaItemResults []struct {
		MediaItem photoslibrary.MediaItem `json:"mediaItem"`
	} `json:"mediaItemResults"`
}
