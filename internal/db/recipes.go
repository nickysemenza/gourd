package db

import (
	"context"
	"database/sql"

	"github.com/nickysemenza/gourd/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
)

// MergeIngredients sets the provided ingredients `parent` to the first one.
// TODO: prevent cyclic loop?
func (c *Client) MergeIngredients(ctx context.Context, tx *sql.Tx, ingredientID string, ids []string) error {
	ctx, span := c.tracer.Start(ctx, "MergeIngredients")
	defer span.End()

	_, err := models.Ingredients(
		qm.Where("id = any(?)", types.Array(ids))).
		UpdateAll(ctx, tx, models.M{models.IngredientColumns.ParentIngredientID: ingredientID})

	return err
}
func (c *Client) RecipeIngredientDependencies(ctx context.Context) ([]RecipeIngredientDependency, error) {
	res := []RecipeIngredientDependency{}
	q := `SELECT distinct
	recipe_details.name AS recipe_name,
	recipe_details.id AS recipe_id,
-- 	coalesce(alts.id, ingredients.id) AS ingredient_id,
--  	coalesce(alts.name, ingredients.name) AS ingredient_name,
--  	r2.name AS ingredient_recipe_name,
-- 	r2.id AS ingredient_recipe_id,
	coalesce(alts.id, ingredients.id, r2.id) AS ingredient_id,
 	coalesce(alts.name, ingredients.name, r2.name) AS ingredient_name,
	(case when r2.id is null then 'ingredient' else 'recipe' end)  as ingredient_kind
	

FROM
	recipe_details
	LEFT JOIN recipe_sections ON recipe_details.id = recipe_sections.recipe_detail_id
	LEFT JOIN recipe_section_ingredients ON recipe_section_ingredients.section_id = recipe_sections.id
	LEFT JOIN recipe_details r2  ON r2.recipe_id = recipe_section_ingredients.recipe_id
	LEFT JOIN ingredients ON recipe_section_ingredients.ingredient_id = ingredients.id
	LEFT JOIN ingredients alts ON ingredients.parent_ingredient_id = alts.id
WHERE
	(recipe_details.is_latest_version = TRUE AND (r2.is_latest_version IS NULL OR r2.is_latest_version = TRUE)
	)
	AND (ingredients.id IS NOT NULL OR r2.id IS NOT NULL)`
	err := c.db.SelectContext(ctx, &res, q)
	return res, err
}
