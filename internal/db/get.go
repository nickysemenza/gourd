package db

import (
	"context"
)



func (c *Client) AddIngredientUnit(ctx context.Context, m IngredientUnitMapping) (int64, error) {
	q := c.psql.Insert("ingredient_units").
		Columns("ingredient_id", "unit_a", "amount_a", "unit_b", "amount_b", "source").
		Values(m.IngredientId, m.UnitA, m.AmountA, m.UnitB, m.AmountB, m.Source).Suffix("ON CONFLICT (ingredient_id, unit_a, amount_a, unit_b, amount_b) DO NOTHING")
	r, err := c.execContext(ctx, q)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}
