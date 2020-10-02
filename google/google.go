package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"github.com/nickysemenza/gourd/db"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// https://developers.google.com/photos/library/reference/rest/v1/mediaItems/batchGet#query-parameters
const maxPhotoBatchGet = 50

type Client struct {
	oc oauth2.Config
	db *db.Client
}

func New(db *db.Client, clientID, clientSecret, redirectURL string) *Client {
	oc := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/photoslibrary.readonly",
			gauth.UserinfoProfileScope,
			gauth.UserinfoEmailScope,
		},
	}
	return &Client{oc, db}
}
func (p *Client) GetURL() string {
	return p.oc.AuthCodeURL("state", oauth2.AccessTypeOffline)
}
func (p *Client) Finish(ctx context.Context, code string) error {
	token, err := p.oc.Exchange(ctx, code, oauth2.AccessTypeOffline)
	if err != nil {
		return fmt.Errorf("bad token exchage: %w", err)
	}

	tokenStr, err := json.Marshal(token)
	if err != nil {
		return err
	}
	if err := p.db.SetKV(ctx, "gphotos-oauth2-token", string(tokenStr)); err != nil {
		return err
	}

	return nil
}
func (p *Client) getToken(ctx context.Context) (*oauth2.Token, error) {
	var token oauth2.Token
	res, err := p.db.GetKV(ctx, "gphotos-oauth2-token")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(res), &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
func (p *Client) getPhotosClient(ctx context.Context) (*gphotos.Client, error) {
	token, err := p.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("bad token: %w", err)
	}
	tc := p.oc.Client(ctx, token)
	return gphotos.NewClient(tc)
}
func (p *Client) batchGet(ctx context.Context, ids []string) ([]photoslibrary.MediaItem, error) {
	if size := len(ids); size > maxPhotoBatchGet {
		return nil, fmt.Errorf("requested %d, limit is %d", size, maxPhotoBatchGet)
	}
	token, err := p.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("bad token: %w", err)
	}
	tc := p.oc.Client(ctx, token)

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
func (p *Client) GetBaseURLs(ctx context.Context, ids []string) (map[string]string, error) {
	chunks := chunkBy(ids, maxPhotoBatchGet)
	urls := map[string]string{}
	for _, chunk := range chunks {
		items, err := p.batchGet(ctx, chunk)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			urls[item.Id] = item.BaseUrl
		}
	}
	return urls, nil

}
func (p *Client) GetTest(ctx context.Context) (interface{}, error) {
	ids := []string{"AIbigFqimylf7SUbAvFeiDPDBg_K_rH5DYtsZUMAiD2yMhJDeHIadDYJnc2Q7vnqKT4DQJeB5IQ7qNEk1Iu0-9k9lfolG6i9-A",
		"AIbigFrtZdWe-gFN1KOuPBhPlFsSNftIy2tyH0yW3JxQPALG-qPg1BsByn12LwoUM_om-DI_rB7OLwhZ8UpzPBxStrlbb9_SQQ"}
	return p.GetBaseURLs(ctx, ids)

}
func (p *Client) GetUserInfo(ctx context.Context, token *oauth2.Token) (*gauth.Userinfo, error) {
	oauth2Service, err := gauth.NewService(ctx, option.WithTokenSource(p.oc.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}
	return gauth.NewUserinfoV2MeService(oauth2Service).Get().Context(ctx).Do()
}

func (p *Client) SyncAlbums(ctx context.Context) error {
	client, err := p.getPhotosClient(ctx)
	if err != nil {
		return fmt.Errorf("bad client: %w", err)
	}

	albums, err := p.db.GetAlbums(ctx)
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
		err = p.db.UpsertPhotos(ctx, photos)
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

// https://gist.github.com/mustafaturan/7a29e8251a7369645fb6c2965f8c2daf
func chunkBy(items []string, chunkSize int) (chunks [][]string) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}
