package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/davecgh/go-spew/spew"
	"github.com/lib/pq"
	"github.com/nickysemenza/gourd/internal/common"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// IngredientByName retrieves an ingredient by name, creating it if it does not exist.
func (c *Client) IngredientByName(ctx context.Context, name string) (*Ingredient, error) {
	ctx, span := c.tracer.Start(ctx, "IngredientByName")
	defer span.End()

	switch name {
	case "full-fat Greek yoghurt", "full fat yogurt or whole milk":
		name = "greek yogurt"
	case "A small handful of coriander leave":
		name = "coriander leaves"
	case "double cream":
		name = "heavy cream"
	}
	ingredient := &Ingredient{}
	err := c.db.GetContext(ctx, ingredient, `SELECT * FROM ingredients
	WHERE lower(name) = lower($1) LIMIT 1`, name)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = c.db.ExecContext(ctx, `INSERT INTO ingredients (id, name) VALUES ($1, $2)`, common.ID("i"), name)
		if err != nil {
			return nil, err
		}
		return c.IngredientByName(ctx, name)
	}
	return ingredient, err
}

// nolint: funlen
func (c *Client) updateRecipe(ctx context.Context, tx *sql.Tx, r *RecipeDetail) error {
	q := c.psql.
		Update(recipeDetailsTable).Where(sq.Eq{"id": r.Id}).Set("name", r.Name).
		Set("unit", r.Unit)

	_, err := c.execTx(ctx, tx, q)
	if err != nil {
		return fmt.Errorf("failed to update recipe details: %w", err)
	}

	if err := c.AssignIds(ctx, r); err != nil {
		return fmt.Errorf("failed to assign ids: %w", err)
	}

	if len(r.Sections) == 0 {
		return nil
	}

	// sections
	sectionInsert := c.psql.Insert(sectionsTable).Columns("id", "recipe_detail_id", "duration_timerange")
	for _, s := range r.Sections {
		sectionInsert = sectionInsert.Values(s.Id, s.RecipeDetailId, s.TimeRange)
	}

	_, err = c.execTx(ctx, tx, sectionInsert)
	if err != nil {
		return fmt.Errorf("failed to insert sections: %w", err)
	}

	instructionsInsert := c.psql.Insert(sInstructionsTable).Columns("id", "section_id", "instruction")
	ingredientsInsert := c.psql.Insert(sIngredientsTable).Columns("id", "section_id", "ingredient_id", "recipe_id",
		"amounts", "adjective", "optional", "original", "sub_for_ingredient_id")

	var hasInstructions, hasIngredients bool
	for _, s := range r.Sections {
		for _, i := range s.Instructions {
			hasInstructions = true
			instructionsInsert = instructionsInsert.Values(i.Id, i.SectionId, i.Instruction)
		}

		for _, i := range s.Ingredients {
			hasIngredients = true

			amounts, err := json.Marshal(i.Amounts)
			if err != nil {
				return fmt.Errorf("failed to insert sections: %w", err)
			}
			ingredientsInsert = ingredientsInsert.Values(i.Id, i.SectionId, i.IngredientId, i.RecipeId,
				string(amounts), i.Adjective, i.Optional, i.Original, i.SubsFor)
		}

	}
	if hasInstructions {
		if _, err = c.execTx(ctx, tx, instructionsInsert); err != nil {
			return fmt.Errorf("failed to insert section instructions %w", err)
		}
	}
	if hasIngredients {
		if _, err = c.execTx(ctx, tx, ingredientsInsert); err != nil {
			return fmt.Errorf("failed to insert section ingredients: %w", err)
		}
	}
	return nil
}

func (c *Client) insertRecipe(ctx context.Context, tx *sql.Tx, r *RecipeDetail) (recipeDetail *RecipeDetail, err error) {
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now()
	}
	// if we have an existing recipe with the same Id or name, this one is a n+1 version of that one
	version := int64(1)
	var modifying *RecipeDetail
	parentID := ""
	if r.Id != "" {
		res, err := c.GetRecipeDetailWhere(ctx, sq.Eq{"id": r.Id})
		if err != nil {
			return nil, fmt.Errorf("failed to find prior recipe: %w", err)
		}
		modifying = res.First()
	}
	if modifying == nil {
		res, err := c.GetRecipeDetailWhere(ctx, sq.Eq{"name": r.Name})
		if err != nil {
			return nil, fmt.Errorf("failed to find prior recipe: %w", err)
		}
		modifying = res.First()
	}

	if modifying != nil {
		latestVersion, err := c.GetRecipeDetailWhere(ctx, sq.Eq{"recipe_id": modifying.RecipeId})
		if err != nil {
			return nil, fmt.Errorf("failed to find prior recipe: %w", err)
		}
		version = latestVersion.First().Version + 1
		parentID = latestVersion.First().RecipeId

		if len(r.Sections) == 0 {
			log.Infof("no sections, modifying, so just returning")
			return modifying, nil
		}
	}
	r.Version = version
	r.Id = common.ID("rd")

	if parentID == "" {
		parentID = common.ID("r")
		_, err = c.execTx(ctx, tx, c.psql.Insert(recipesTable).Columns("id", "created_at").Values(parentID, r.CreatedAt))
		if err != nil {
			return nil, fmt.Errorf("failed to insert parent recipe: %w", err)
		}
	} else {
		// there is a parent and therefore other children, which are no longer latest
		_, err = c.execTx(ctx, tx, c.psql.
			Update(recipeDetailsTable).
			Set("is_latest_version", false).
			Where(sq.Eq{"recipe_id": parentID}),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update other verions to not be latest: %w", err)
		}
	}
	r.RecipeId = parentID

	if r.Id == "" {
		// todo? needed?
		r.Id = common.ID("rd")
	}
	log.Println("inserting", r.Id, r.Name)

	date := time.Now()
	if !r.CreatedAt.IsZero() {
		date = r.CreatedAt
	}

	_, err = c.execTx(ctx, tx, c.psql.
		Insert(recipeDetailsTable).
		Columns("id", "recipe_id", "name", "version", "is_latest_version", "source", "created_at", "tags").
		Values(r.Id, r.RecipeId, r.Name, r.Version, true, r.Source, date, pq.Array(r.Tags)))
	if err != nil {
		return nil, fmt.Errorf("failed to insert new recipe details row: %w", err)
	}

	err = c.updateRecipe(ctx, tx, r)
	if err != nil {
		return nil, err
	}

	return r, nil

}

// InsertRecipe inserts a recipe.
func (c *Client) InsertRecipe(ctx context.Context, r *RecipeDetail) (*RecipeDetail, error) {
	ctx, span := c.tracer.Start(ctx, "InsertRecipe")
	defer span.End()

	span.AddEvent("got detail", trace.WithAttributes(attribute.String("id", r.Id), attribute.String("recipe", spew.Sdump(r))))

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.insertRecipe(ctx, tx, r)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return c.GetRecipeDetailByIdFull(ctx, res.Id)
}

func (c *Client) IngredientToRecipe(ctx context.Context, ingredientID string) (*RecipeDetail, error) {
	ctx, span := c.tracer.Start(ctx, "IngredientToRecipe")
	defer span.End()

	i, err := c.GetIngredientById(ctx, ingredientID)
	if err != nil {
		return nil, err
	}
	if i == nil {
		return nil, fmt.Errorf("failed to find ingredient with id %s", ingredientID)
	}

	if strings.Contains(i.Name, "[converted]") {
		return nil, fmt.Errorf("%s has already been converted to recipe", i.Name)
	}

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	newDetail, err := c.insertRecipe(ctx, tx, &RecipeDetail{Name: i.Name})
	if err != nil {
		return nil, fmt.Errorf("failed to insert new recipe: %w", err)
	}

	if _, err = c.execTx(ctx,
		tx,
		c.psql.
			Update(sIngredientsTable).
			Set("ingredient_id", nil).
			Set("recipe_id", newDetail.RecipeId).
			Where(sq.Eq{"ingredient_id": ingredientID})); err != nil {
		return nil, fmt.Errorf("failed to update references to transformed ingredient: %w", err)
	}

	if _, err = c.execTx(ctx,
		tx,
		c.psql.
			Update(ingredientsTable).
			Set("name", fmt.Sprintf("[converted] %s", i.Name)).
			Where(sq.Eq{"id": i.Id})); err != nil {
		return nil, fmt.Errorf("failed to deprecated ingredient: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return c.GetRecipeDetailByIdFull(ctx, newDetail.Id)
}

// MergeIngredients sets the provided ingredients `parent` to the first one.
// TODO: prevent cyclic loop?
func (c *Client) MergeIngredients(ctx context.Context, ingredientID string, ids []string) error {
	ctx, span := c.tracer.Start(ctx, "MergeIngredients")
	defer span.End()

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err = c.execTx(ctx,
		tx,
		c.psql.
			Update(ingredientsTable).
			Set("parent_ingredient_id", ingredientID).
			Where(sq.Eq{"id": ids})); err != nil {
		return fmt.Errorf("failed to update ingredient: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
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
