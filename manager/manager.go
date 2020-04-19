package manager

import (
	"context"
	"io/ioutil"

	"github.com/nickysemenza/food/db"
	"gopkg.in/yaml.v1"
)

// Manager manages recipes
type Manager struct {
	db *db.Client
}

// New creates a new Manager
func New(db *db.Client) *Manager {
	return &Manager{db: db}
}

// LoadFromFile loads a recipe from a file
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
