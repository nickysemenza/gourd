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

// func (c *Client) SearchFoods(ctx context.Context, searchQuery string, dataType []string, foodCategoryID *int, opts ...SearchOption) ([]Food, uint64, error) {
// 	ctx, span := c.tracer.Start(ctx, "GetIngrientsParent")
// 	defer span.End()
// 	return c.getFoods(ctx, func(q sq.SelectBuilder, count bool) sq.SelectBuilder {
// 		w := `description ILIKE '%' || $1 || '%'`

// 		q = q.Where(w, searchQuery)
// 		q = newSearchQuery(opts...).apply(q)

// 		if foodCategoryID != nil {
// 			q = q.Where(sq.Eq{"food_category_id": &foodCategoryID})
// 		}
// 		if len(dataType) > 0 {
// 			q = q.Where(sq.Eq{"data_type": dataType})
// 		}
// 		if !count {
// 			q = q.OrderBy("length(description) ASC", "fdc_id DESC")
// 		}
// 		return q
// 	})
// }

func (c *Client) AssociateFoodWithIngredient(ctx context.Context, ingredient string, fdcId int) error {
	ctx, span := c.tracer.Start(ctx, "AssociateFoodWithIngredient")
	defer span.End()

	q := c.psql.Update("ingredients").Set("fdc_id", fdcId).Where(sq.Eq{"id": ingredient})
	_, err := c.execContext(ctx, q)
	return err
}
