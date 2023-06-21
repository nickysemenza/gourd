package db

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (c *Client) AssociateFoodWithIngredient(ctx context.Context, ingredient string, fdcId int) error {
	ctx, span := c.tracer.Start(ctx, "AssociateFoodWithIngredient")
	defer span.End()

	q := c.psql.Update("ingredients").Set("fdc_id", fdcId).Where(sq.Eq{"id": ingredient})
	_, err := c.execContext(ctx, q)
	return err
}
