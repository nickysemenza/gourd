package manager

import (
	"context"

	"github.com/nickysemenza/gourd/auth"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/google"
	"github.com/nickysemenza/gourd/notion"
	"github.com/nickysemenza/gourd/photos"
	"github.com/nickysemenza/gourd/rs_client"
)

// Manager manages recipes
type Manager struct {
	db     *db.Client
	Google *google.Client
	Photos *photos.Photos
	Auth   *auth.Auth
	R      *rs_client.Client
	Notion *notion.Client
}

func New(db *db.Client, g *google.Client, auth *auth.Auth, r *rs_client.Client, notion *notion.Client) *Manager {
	return &Manager{db: db,
		Google: g,
		Auth:   auth,
		Photos: photos.New(db, g),
		R:      r,
		Notion: notion,
	}
}

func (m *Manager) DB() *db.Client {
	return m.db
}

func (m *Manager) ProcessGoogleAuth(ctx context.Context, code string) (jwt string, rawUser map[string]interface{}, err error) {
	err = m.Google.Finish(ctx, code)
	if err != nil {
		return
	}
	user, err := m.Google.GetUserInfo(ctx)
	if err != nil {
		return
	}
	rawUser = map[string]interface{}{"raw": user}

	jwt, err = m.Auth.GetJWT(user)
	return
}

func (m *Manager) SyncNotionToMeals(ctx context.Context) error {
	nMeals, err := m.Notion.Dump(ctx)
	if err != nil {
		return err
	}

	for _, nMeal := range nMeals {
		if nMeal.Time == nil {
			continue
		}
		_, err := m.db.MealIDInRange(ctx, *nMeal.Time, nMeal.Title, &nMeal.PageID)
		if err != nil {
			return err
		}
		// nm
	}
	return nil
}
