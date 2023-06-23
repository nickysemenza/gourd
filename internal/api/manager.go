package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/internal/db"
)

func (a *API) DB() *db.Client {
	return a.db
}

func (a *API) ProcessGoogleAuth(ctx context.Context, code string) (jwt string, rawUser map[string]interface{}, err error) {
	err = a.Google.Finish(ctx, code)
	if err != nil {
		return
	}
	user, err := a.Google.GetUserInfo(ctx)
	if err != nil {
		return
	}
	rawUser = map[string]interface{}{"raw": user}

	jwt, err = a.Auth.GetJWT(user)
	return
}

func (a *API) DoSync(c echo.Context, params DoSyncParams) error {
	ctx := c.Request().Context()
	ctx, span := a.tracer.Start(ctx, "DoSync")
	defer span.End()
	err := a.Sync(ctx, params.LookbackDays)
	if err != nil {
		return handleErr(c, err)
	}
	return c.JSON(http.StatusOK, nil)
}
func (a *API) Sync(ctx context.Context, lookbackDays int) error {
	now := time.Now()
	ctx, span := a.tracer.Start(ctx, "Sync")
	defer span.End()

	if err := a.syncRecipeFromNotion(ctx, time.Hour*24*time.Duration(lookbackDays)); err != nil {
		return fmt.Errorf("notion: %w", err)
	}
	l(ctx).Infof("synced recipes from notion")
	if err := a.DB().SyncNotionMealFromNotionRecipe(ctx); err != nil {
		return fmt.Errorf("notion meal: %w", err)
	}
	// if a.GPhotos != nil {
	// 	if err := a.GPhotos.SyncAlbums(ctx); err != nil {
	// 		return fmt.Errorf("gphotos: %w", err)
	// 	}
	// }
	// if err := a.DB().SyncMealsFromGPhotos(ctx); err != nil {
	// 	return fmt.Errorf("gphotos meal: %w", err)
	// }
	l(ctx).Infof("sync complete in %s", time.Since(now))
	return nil
}
