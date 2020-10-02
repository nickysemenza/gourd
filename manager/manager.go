package manager

import (
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/google"
)

// Manager manages recipes
type Manager struct {
	db     *db.Client
	Google *google.Client
}

func New(db *db.Client) *Manager {
	return &Manager{db: db}
}

func (m *Manager) DB() *db.Client {
	return m.db
}
