package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/nickysemenza/gourd/internal/common"
	"github.com/nickysemenza/gourd/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var ingredientQueryMods = []QueryMod{
	Load(models.IngredientRels.ParentIngredient),
	Load(models.IngredientRels.ParentIngredientIngredients),
	Load(models.IngredientRels.IngredientUnits),
	Load(Rels(
		models.IngredientRels.RecipeSectionIngredients,
		models.RecipeSectionIngredientRels.Section,
		models.RecipeSectionRels.RecipeDetail,
	)),
}

func (a *API) ingredientByName(ctx context.Context, name string) (*models.Ingredient, error) {
	ctx, span := a.tracer.Start(ctx, "ingredientByName")
	defer span.End()

	switch name {
	case "full-fat Greek yoghurt", "full fat yogurt or whole milk":
		name = "greek yogurt"
	case "A small handful of coriander leave":
		name = "coriander leaves"
	case "double cream":
		name = "heavy cream"
	}

	ingredient, err := models.Ingredients(Where("lower(name) = lower(?)", name)).One(ctx, a.db.DB())
	if errors.Is(err, sql.ErrNoRows) {
		ingredient = &models.Ingredient{
			Name: name,
			ID:   common.ID("i"),
		}
		err := ingredient.Insert(ctx, a.db.DB(), boil.Infer())
		if err != nil {
			return nil, err
		}
		return a.ingredientByName(ctx, name)
	}
	return ingredient, err
}

func (a *API) ingredientById(ctx context.Context, ingredientId IngredientID, withDetail bool) (*IngredientDetail, error) {
	ctx, span := a.tracer.Start(ctx, "ingredientByIdV2")
	defer span.End()
	ingredient, err := models.Ingredients(
		append(ingredientQueryMods, Where("id = ?", ingredientId))...,
	).One(ctx, a.db.DB())

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to find ingredient with id %s: %w", ingredientId, err)
		}
		return nil, err
	}
	rw, err := a.ingredientFromModel(ctx, ingredient, withDetail, true, true)
	return rw, err
}

func (a *API) IngredientListV2(ctx context.Context, pagination Items, mods ...QueryMod) ([]IngredientDetail, int64, error) {
	ctx, span := a.tracer.Start(ctx, "IngredientListV2")
	defer span.End()
	filters := []QueryMod{
		Where("parent_ingredient_id IS NULL"),
		Limit(pagination.Limit),
		Offset(pagination.Offset),
	}
	ingredients, count, err := countAndQuery[models.IngredientSlice](ctx, a.db.DB(), models.Ingredients, qmWithPagination(ingredientQueryMods, pagination, filters...)...)
	if err != nil {
		return nil, 0, err
	}

	if err != nil {
		return nil, 0, err
	}
	items := []IngredientDetail{}
	for _, recipe := range ingredients {
		rw, err := a.ingredientFromModel(ctx, recipe, true, true, true)
		if err != nil {
			return nil, 0, err
		}
		if rw != nil {
			items = append(items, *rw)
		}
	}

	return items, count, nil
}
