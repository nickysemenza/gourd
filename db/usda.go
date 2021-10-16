package db

import (
	"context"
	"fmt"
	"strings"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"gopkg.in/guregu/null.v4/zero"
)

type Food struct {
	Description string   `db:"description"`
	DataType    string   `db:"data_type"`
	FdcID       int      `db:"fdc_id"`
	CategoryID  zero.Int `db:"food_category_id"`
}

type Foods []Food

func (c *Client) GetFood(ctx context.Context, fdcID int) (*Food, error) {
	ctx, span := c.tracer.Start(ctx, "GetFood")
	defer span.End()
	q := c.psql.Select(
		"food_category_id",
		"data_type",
		"description",
		"fdc_id",
	).From("usda_food").Where(sq.Eq{"fdc_id": fdcID})

	f := &Food{}
	return f, c.getContext(ctx, q, f)
}
func (c *Client) getFoods(ctx context.Context, addons func(q sq.SelectBuilder, count bool) sq.SelectBuilder) ([]Food, uint64, error) {

	favoriteBrands := []string{
		// brands
		"Guittard Chocolate Co.",
		"Bob's Red Mill Natural Foods, Inc.",
		"The King Arthur Flour Company, Inc.",

		// stores
		"'Whole Foods Market, Inc.",
		"WHOLE FOODS MARKET",
		"TRADER JOE'S",
		"Target Stores",
		"Kikkoman Sales USA, Inc.",
	}
	for x := range favoriteBrands {
		// escape single quote
		favoriteBrands[x] = strings.ReplaceAll(favoriteBrands[x], "'", "''")
	}

	ctx, span := c.tracer.Start(ctx, "getFoods")
	defer span.End()

	q := c.psql.Select(
		"food_category_id",
		"data_type",
		"description",
		"usda_food.fdc_id as fdc_id",
	).From("usda_food").
		LeftJoin("usda_branded_food	on usda_food.fdc_id = usda_branded_food.fdc_id").
		OrderBy(fmt.Sprintf(`array_position(array['%s'], brand_owner)`, strings.Join(favoriteBrands, "','"))).
		OrderBy("length(ingredients) asc") // shorter ingredient list = more likely to be 'pure'

	cq := c.psql.Select("count(*)").From("usda_food")

	q = addons(q, false)
	cq = addons(cq, true)

	cq = cq.RemoveLimit().RemoveOffset()

	res := []Food{}
	var count uint64

	numQueries := 2
	var wg sync.WaitGroup
	wg.Add(numQueries)
	errs := make(chan error, numQueries)

	go func() {
		if err := c.selectContext(ctx, q, &res); err != nil {
			errs <- err
		}
		wg.Done()
	}()
	go func() {
		if err := c.getContext(ctx, cq, &count); err != nil {
			errs <- err
		}
		wg.Done()
	}()

	wg.Wait()
	close(errs)
	if len(errs) > 0 {
		return nil, 0, <-errs
	}

	return res, count, nil
}
func (c *Client) SearchFoods(ctx context.Context, searchQuery string, dataType []string, foodCategoryID *int, opts ...SearchOption) ([]Food, uint64, error) {
	ctx, span := c.tracer.Start(ctx, "GetIngrientsParent")
	defer span.End()
	return c.getFoods(ctx, func(q sq.SelectBuilder, count bool) sq.SelectBuilder {
		w := `description ILIKE '%' || $1 || '%'`

		q = q.Where(w, searchQuery)
		q = newSearchQuery(opts...).apply(q)

		if foodCategoryID != nil {
			q = q.Where(sq.Eq{"food_category_id": &foodCategoryID})
		}
		if len(dataType) > 0 {
			q = q.Where(sq.Eq{"data_type": dataType})
		}
		if !count {
			q = q.OrderBy("length(description) ASC", "fdc_id DESC")
		}
		return q
	})
}
func (c *Client) FoodsByIds(ctx context.Context, ids []int) ([]Food, uint64, error) {
	ctx, span := c.tracer.Start(ctx, "FoodsByIds")
	defer span.End()
	return c.getFoods(ctx, func(q sq.SelectBuilder, count bool) sq.SelectBuilder {
		q = q.Where(sq.Eq{"usda_food.fdc_id": ids})
		if !count {
			q = q.OrderBy("usda_food.fdc_id DESC")
		}
		return q
	})
}

type FoodNutrient struct {
	FdcID      int      `db:"fdc_id"`
	NutrientID int      `json:"nutrient" db:"nutrient_id"`
	Amount     float64  `json:"amount" db:"amount"`
	DataPoints zero.Int `json:"data_points" db:"data_points"`
}
type FoodNutrients []FoodNutrient

func (c *Client) GetFoodNutrients(ctx context.Context, fdcID ...int) (FoodNutrients, error) {
	ctx, span := c.tracer.Start(ctx, "GetFoodNutrients")
	defer span.End()
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
	ctx, span := c.tracer.Start(ctx, "GetNutrients")
	defer span.End()

	q := c.psql.Select("id", "name", "unit_name AS unitName").From("usda_nutrient").Where(sq.Eq{"id": nutrientID})

	ns := []Nutrient{}
	err := c.selectContext(ctx, q, &ns)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return ns, nil
}

func (c *Client) GetCategory(ctx context.Context, categoryID int64) (*FoodCategory, error) {
	ctx, span := c.tracer.Start(ctx, "GetCategory")
	defer span.End()

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
	ctx, span := c.tracer.Start(ctx, "GetBrandInfo")
	defer span.End()

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

type FoodPortion struct {
	Id                 int         `db:"id"`
	FdcID              int         `db:"fdc_id"`
	Amount             zero.Float  `db:"amount"`
	GramWeight         float64     `db:"gram_weight"`
	Modifier           zero.String `db:"modifier"`
	PortionDescription zero.String `db:"portion_description"`
}
type FoodPortions []FoodPortion

func (r FoodPortions) ByFdcId() map[int][]FoodPortion {
	m := make(map[int][]FoodPortion)
	for _, x := range r {
		m[x.FdcID] = append(m[x.FdcID], x)
	}
	return m
}

func (c *Client) GetFoodPortions(ctx context.Context, fdcId ...int) (FoodPortions, error) {
	ctx, span := c.tracer.Start(ctx, "GetFoodPortions")
	defer span.End()

	q := c.psql.Select("id", "fdc_id", "amount", "portion_description", "modifier", "gram_weight").
		From("usda_food_portion").Where(sq.Eq{"fdc_id": fdcId})

	ns := []FoodPortion{}
	err := c.selectContext(ctx, q, &ns)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return ns, nil
}
func (c *Client) AssociateFoodWithIngredient(ctx context.Context, ingredient string, fdcId int) error {
	ctx, span := c.tracer.Start(ctx, "AssociateFoodWithIngredient")
	defer span.End()

	q := c.psql.Update("ingredients").Set("fdc_id", fdcId).Where(sq.Eq{"id": ingredient})
	_, err := c.execContext(ctx, q)
	return err
}
