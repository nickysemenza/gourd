package manager

import (
	"context"
	"io/ioutil"

	"github.com/nickysemenza/gourd/db"
	yaml "gopkg.in/yaml.v2"
)

// Manager manages recipes
type Manager struct {
	db *db.Client
}

func New(db *db.Client) *Manager {
	return &Manager{db: db}
}

func (m *Manager) DB() *db.Client {
	return m.db
}

func (m *Manager) LoadFromFile(ctx context.Context, filename string) (*Recipe, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	r := &Recipe{}
	err = yaml.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (m *Manager) GetRecipe(ctx context.Context, uuid string) (*Recipe, error) {
	r, err := m.db.GetRecipeByUUIDFull(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return FromRecipe(r), nil
}
