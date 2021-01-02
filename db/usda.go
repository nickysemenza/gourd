package db

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"gopkg.in/guregu/null.v3/zero"
)

type Food struct {
	Description string   `db:"description"`
	DataType    string   `db:"data_type"`
	FdcID       int      `db:"fdc_id"`
	CategoryID  zero.Int `db:"food_category_id"`
}

type Foods []Food

// func (e Foods) String(i int) string {
// 	return e[i].Description
// }

// func (e Foods) Len() int {
// 	return len(e)
// }

// func (c *Client) GetFoods(ctx context.Context) (Foods, error) {
// 	query, args, err := c.psql.Select("description, data_type, fdc_id").From("usda_food").ToSql()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var res Foods
// 	err = c.db.SelectContext(ctx, &res, query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return res, nil
// }

func (c *Client) GetFood(ctx context.Context, fdcID int) (*Food, error) {
	q := c.psql.Select(
		"food_category_id",
		"data_type",
		"description",
		"fdc_id",
	).From("usda_food").Where(sq.Eq{"fdc_id": fdcID})

	f := &Food{}
	return f, c.getContext(ctx, q, f)
}

func (c *Client) SearchFoods(ctx context.Context, searchQuery string, dataType string, foodCategoryID *int) ([]Food, error) {
	q := c.psql.Select(
		"food_category_id",
		"data_type",
		"description",
		"fdc_id",
	).From("usda_food").Where(sq.ILike{"description": fmt.Sprintf("%%%s%%", searchQuery)})
	if foodCategoryID != nil {
		q = q.Where(sq.Eq{"food_category_id": &foodCategoryID})
	}
	if dataType != "" {
		q = q.Where(sq.Eq{"data_type": dataType})
	}

	res := []Food{}

	if err := c.selectContext(ctx, q, &res); err != nil {
		return nil, err
	}
	return res, nil
}

type FoodNutrient struct {
	FdcID      int      `db:"fdc_id"`
	NutrientID int      `json:"nutrient" db:"nutrient_id"`
	Amount     float64  `json:"amount" db:"amount"`
	DataPoints zero.Int `json:"data_points" db:"data_points"`
}
type FoodNutrients []FoodNutrient

func (c *Client) GetFoodNutrients(ctx context.Context, fdcID ...int) (FoodNutrients, error) {
	q := c.psql.Select(
		"nutrient_id",
		"amount",
		"data_points",
	).From("usda_food_nutrient").Where(sq.Eq{"fdc_id": fdcID})

	fns := []FoodNutrient{}
	err := c.selectContext(ctx, q, &fns)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return fns, nil
}

type Nutrient struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UnitName string `json:"unitName"`
}
type Nutrients []Nutrient

func (r Nutrients) ById() map[int]Nutrient {
	m := make(map[int]Nutrient)
	for _, x := range r {
		m[x.ID] = x
	}
	return m
}

func (c *Client) GetNutrients(ctx context.Context, nutrientID ...int) (Nutrients, error) {
	q := c.psql.Select("id", "name", "unit_name AS unitName").From("usda_nutrient").Where(sq.Eq{"id": nutrientID})

	ns := []Nutrient{}
	err := c.selectContext(ctx, q, &ns)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return ns, nil
}

func (c *Client) GetCategory(ctx context.Context, categoryID int64) (*FoodCategory, error) {
	q := c.psql.Select(
		"code",
		"description",
	).From("usda_food_category").Where(sq.Eq{"id": categoryID})

	x := &FoodCategory{}
	err := c.getContext(ctx, q, x)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return x, nil
}
func (c *Client) GetBrandInfo(ctx context.Context, fdcID int) (*BrandedFood, error) {
	q := c.psql.Select(
		"brand_owner",
		"ingredients",
		"serving_size",
		"serving_size_unit",
		"household_serving_fulltext",
		"branded_food_category",
	).From("usda_branded_food").Where(sq.Eq{"fdc_id": fdcID})

	x := &BrandedFood{}
	err := c.getContext(ctx, q, x)
	if err != nil {
		return nil, fmt.Errorf("faileo select: %w", err)
	}
	return x, nil
}

type BrandedFood struct {
	BrandOwner          *string `db:"brand_owner"`
	Ingredients         *string `db:"ingredients"`
	ServingSize         float64 `db:"serving_size"`
	ServingSizeUnit     string  `db:"serving_size_unit"`
	HouseholdServing    *string `db:"household_serving_fulltext"`
	BrandedFoodCategory *string `db:"branded_food_category"`
}

type FoodCategory struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}
