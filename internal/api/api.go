//go:generate ../../bin/oapi-codegen --package api --generate server,spec -o api-server.gen.go openapi.yaml
//go:generate ../../bin/oapi-codegen --package api --generate types,skip-prune -o api-types.gen.go openapi.yaml
//go:generate ../../bin/oapi-codegen --package api --generate client -o api-client.gen.go openapi.yaml

package api

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/internal/auth"
	"github.com/nickysemenza/gourd/internal/clients/google"
	"github.com/nickysemenza/gourd/internal/clients/gphotos"
	"github.com/nickysemenza/gourd/internal/clients/notion"
	"github.com/nickysemenza/gourd/internal/clients/rs_client"
	"github.com/nickysemenza/gourd/internal/db"
	"github.com/nickysemenza/gourd/internal/image"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type API struct {
	db         *db.Client
	Google     *google.Client
	GPhotos    *gphotos.Photos
	Auth       *auth.Auth
	R          *rs_client.Client
	Notion     *notion.Client
	ImageStore image.Store
	tracer     trace.Tracer
}

type ServerOption func(*API)

func WithNotionClient(n *notion.Client) ServerOption {
	return func(s *API) {
		s.Notion = n
	}
}

func WithGoogleClient(g *google.Client) ServerOption {
	return func(s *API) {
		s.Google = g
	}
}
func WithAuthClient(a *auth.Auth) ServerOption {
	return func(s *API) {
		s.Auth = a
	}
}

func New(db *db.Client,
	r *rs_client.Client,
	imageStore image.Store,
	opts ...ServerOption) *API {
	a := API{
		db:         db,
		R:          r,
		ImageStore: imageStore,
		tracer:     otel.Tracer("api"),
	}

	for _, o := range opts {
		o(&a)
	}
	if a.Google != nil {
		a.GPhotos = gphotos.New(a.db, a.Google, a.ImageStore)
	}
	return &a
}

func hasGrams(amounts []Amount) bool {
	for _, amt := range amounts {
		if amt.Unit == "grams" || amt.Unit == "g" || amt.Unit == "gram" {
			return true
		}
	}
	return false
}

func (a *API) AuthLogin(c echo.Context, params AuthLoginParams) error {
	ctx := c.Request().Context()
	jwt, rawUser, err := a.ProcessGoogleAuth(ctx, params.Code)
	if err != nil {
		return handleErr(c, err)
	}

	resp := AuthResp{
		Jwt:  jwt,
		User: rawUser,
	}
	return c.JSON(http.StatusOK, resp)
}

func (a *API) ListAllAlbums(c echo.Context) error {
	ctx := c.Request().Context()

	var resp struct {
		Albums []GooglePhotosAlbum `json:"albums,omitempty"`
	}

	dbAlbums, err := a.DB().GetAlbums(ctx)
	if err != nil {
		return handleErr(c, err)
	}

	albums, err := a.GPhotos.GetAvailableAlbums(ctx)
	if err != nil {
		return handleErr(c, err)
	}

	for _, a := range albums {
		gpa := GooglePhotosAlbum{
			Id:         a.Id,
			ProductUrl: a.ProductUrl,
			Title:      a.Title,
		}

		for _, dbA := range dbAlbums {
			if dbA.ID == gpa.Id {
				gpa.Usecase = dbA.Usecase
			}
		}

		resp.Albums = append(resp.Albums, gpa)
	}

	return c.JSON(http.StatusOK, resp)
}

func (a *API) RecipeDependencies(c echo.Context) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "RecipeDependencies")
	defer span.End()

	res := []RecipeDependency{}
	dbRows, err := a.DB().RecipeIngredientDependencies(ctx)
	for _, r := range dbRows {
		res = append(res, RecipeDependency{
			IngredientId:   r.IngredientId,
			IngredientKind: IngredientKind(r.IngredientKind),
			IngredientName: r.IngredientName,
			RecipeId:       r.RecipeId,
			RecipeName:     r.RecipeName,
		})
	}
	if err != nil {
		return handleErr(c, err)
	}
	return c.JSON(http.StatusOK, struct {
		Items []RecipeDependency `json:"items,omitempty"`
	}{res})
}

type Res struct {
	Notion  notion.Recipe
	Details *RecipeDetailInput
}

func (a *API) getEnhancedNotion(ctx context.Context, lookback time.Duration, pageID string) ([]Res, error) {
	res, err := a.Notion.GetAll(ctx, lookback, pageID)
	if err != nil {
		return nil, err
	}

	var res2 []Res
	for _, r := range res {
		i, err := a.notionRecipeToInput(ctx, r)
		if err != nil {
			return nil, err
		}
		res2 = append(res2, Res{
			Notion:  r,
			Details: i,
		})
	}
	return res2, nil
}
func (a *API) NotionTest(c echo.Context) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "Notion")
	defer span.End()

	daysSince := 45
	if daysSinceStr := c.QueryParam("days_since"); daysSinceStr != "" {
		var err error
		daysSince, err = strconv.Atoi(daysSinceStr)
		if err != nil {
			return handleErr(c, err)
		}
	}
	res2, err := a.getEnhancedNotion(ctx, time.Hour*24*time.Duration(daysSince), c.QueryParam("page_id"))
	if err != nil {
		return handleErr(c, err)
	}

	return c.JSON(http.StatusOK, res2)
}

func (a *API) GetConfig(c echo.Context) error {
	res := ConfigData{
		GoogleClientId: a.Google.GetClientID(),
		GoogleScopes:   "profile email https://www.googleapis.com/auth/photoslibrary.readonly",
	}
	return c.JSON(http.StatusOK, res)
}

func (a *API) tx(ctx context.Context) *sql.Tx {
	tx, err := a.db.DB().BeginTx(ctx, nil)
	if err != nil {
		panic("failed to start tx: " + err.Error())
	}
	return tx
}
