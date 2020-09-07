package graph

import (
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/graph/model"
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

func fromIngredient(res *db.Ingredient) *model.Ingredient {
	if res == nil {
		return nil
	}
	i := &model.Ingredient{
		UUID:   res.UUID,
		Name:   res.Name,
		SameAs: res.SameAs.String,
	}
	if res.FdcID.Valid {
		i.FdcID = res.FdcID.Int64
	}

	return i
}
