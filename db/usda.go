package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/nickysemenza/food/graph/model"
)

type Food struct {
	Description string `db:"description"`
	DataType    string `db:"data_type"`
	FdcID       string `db:"fdc_id"`
}

type Foods []Food

func (e Foods) String(i int) string {
	return e[i].Description
}

func (e Foods) Len() int {
	return len(e)
}

func (c *Client) GetFoods(ctx context.Context) (Foods, error) {
	query, args, err := c.psql.Select("description, data_type, fdc_id").From("food").ToSql()
	if err != nil {
		return nil, err
	}
	var res Foods
	err = c.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetFood(ctx context.Context, fdcID int) (*model.Food, error) {
	query, args, err := c.psql.Select(
		"food_category_id",
		"data_type",
		"description",
		"fdc_id",
	).From("food").Where(sq.Eq{"fdc_id": fdcID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	f := &model.Food{}
	err = c.db.GetContext(ctx, f, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return f, nil
}

func (c *Client) GetFoodNutrients(ctx context.Context, fdcID int) ([]*model.FoodNutrient, error) {
	query, args, err := c.psql.Select(
		"nutrient_id",
		"amount",
		"data_points",
	).From("food_nutrient").Where(sq.Eq{"fdc_id": fdcID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	fns := []*model.FoodNutrient{}
	err = c.db.SelectContext(ctx, &fns, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	// fns2 := []*model.FoodNutrient{}
	// for _, x := range fns {
	// 	fns2 = append(fns2, &x)
	// }
	return fns, nil
}

func (c *Client) GetNutrient(ctx context.Context, nutrientID int) (*model.Nutrient, error) {
	query, args, err := c.psql.Select(
		"id",
		"name",
		"unit_name",
	).From("nutrient").Where(sq.Eq{"id": nutrientID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	x := &model.Nutrient{}
	err = c.db.GetContext(ctx, x, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return x, nil
}

func (c *Client) GetCategory(ctx context.Context, categoryID int64) (*model.FoodCategory, error) {
	query, args, err := c.psql.Select(
		"code",
		"description",
	).From("food_category").Where(sq.Eq{"id": categoryID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	x := &model.FoodCategory{}
	err = c.db.GetContext(ctx, x, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return x, nil
}
