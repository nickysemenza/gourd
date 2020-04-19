package graph

import (
	"github.com/nickysemenza/food/db"
	"github.com/nickysemenza/food/graph/model"
)

func fromRecipe(res *db.Recipe) *model.Recipe {
	if res == nil {
		return nil
	}
	return &model.Recipe{
		UUID:         res.UUID,
		Name:         res.Name,
		TotalMinutes: int(res.TotalMinutes.Int64),
		Unit:         res.Unit.String,
	}

}
