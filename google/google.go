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

type Client struct {
	oc oauth2.Config
	db *db.Client

	_token *oauth2.Token
}

func New(db *db.Client, clientID, clientSecret, redirectURL string) *Client {
	oc := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "postmessage", // https://github.com/anthonyjgrove/react-google-login/issues/66#issuecomment-648521442
		Scopes: []string{
			"https://www.googleapis.com/auth/photoslibrary.readonly",
			gauth.UserinfoProfileScope,
			gauth.UserinfoEmailScope,
		},
	}
	return &Client{oc: oc, db: db}
}
func (c *Client) GetURL() string {
	return c.oc.AuthCodeURL("state", oauth2.AccessTypeOffline)
}
func (c *Client) Finish(ctx context.Context, code string) error {
	token, err := c.oc.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("bad token exchage: %w", err)
	}

	ui, err := c.getUserInfo(ctx, token)
	if err != nil {
		return err
	}
	if ui.Email != "14nicholasse@gmail.com" {
		return fmt.Errorf("user not allowed")
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

	if c._token != nil {
		return c._token, nil
	}

	var token oauth2.Token
	res, err := c.db.GetKV(ctx, "gphotos-oauth2-token")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(res), &token)
	if err != nil {
		return nil, err
	}
	c._token = &token
	return &token, nil
}

func (c *Client) GetUserInfo(ctx context.Context) (*gauth.Userinfo, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, err
	}
	return c.getUserInfo(ctx, token)

}
func (c *Client) getUserInfo(ctx context.Context, token *oauth2.Token) (*gauth.Userinfo, error) {
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
