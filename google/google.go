package google

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nickysemenza/gourd/db"
	"go.opentelemetry.io/otel/api/global"
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
func (c *Client) GetURL() string {
	return c.oc.AuthCodeURL("state", oauth2.AccessTypeOffline)
}
func (c *Client) Finish(ctx context.Context, code string) error {
	token, err := c.oc.Exchange(ctx, code, oauth2.AccessTypeOffline)
	if err != nil {
		return fmt.Errorf("bad token exchage: %w", err)
	}

	tokenStr, err := json.Marshal(token)
	if err != nil {
		return err
	}
	if err := c.db.SetKV(ctx, "gphotos-oauth2-token", string(tokenStr)); err != nil {
		return err
	}

	return nil
}
func (c *Client) getToken(ctx context.Context) (*oauth2.Token, error) {
	ctx, span := global.Tracer("google").Start(ctx, "google.getToken")
	defer span.End()
	var token oauth2.Token
	res, err := c.db.GetKV(ctx, "gphotos-oauth2-token")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(res), &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (c *Client) GetUserInfo(ctx context.Context, token *oauth2.Token) (*gauth.Userinfo, error) {
	oauth2Service, err := gauth.NewService(ctx, option.WithTokenSource(c.oc.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}
	return gauth.NewUserinfoV2MeService(oauth2Service).Get().Context(ctx).Do()
}

// https://gist.github.com/mustafaturan/7a29e8251a7369645fb6c2965f8c2daf
func chunkBy(items []string, chunkSize int) (chunks [][]string) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}
