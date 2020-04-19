package graph

import (
	"github.com/nickysemenza/food/db"
	"github.com/nickysemenza/food/manager"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver holds dependencies for GQL
type Resolver struct {
	Manager *manager.Manager
	DB      *db.Client
}
