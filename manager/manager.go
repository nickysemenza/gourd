package manager

import (
	"context"

	"github.com/nickysemenza/gourd/auth"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/google"
)

// Manager manages recipes
type Manager struct {
	db     *db.Client
	Google *google.Client
	Auth   *auth.Auth
}

func New(db *db.Client, g *google.Client, auth *auth.Auth) *Manager {
	return &Manager{db: db,
		Google: g,
		Auth:   auth,
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
