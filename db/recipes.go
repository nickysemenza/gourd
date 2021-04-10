package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/nickysemenza/gourd/common"
)

// IngredientByName retrieves an ingredient by name, creating it if it does not exist.
func (c *Client) IngredientByName(ctx context.Context, name string) (*Ingredient, error) {
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

//nolint: funlen
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
	sectionInsert := c.psql.Insert(sectionsTable).Columns("id", "recipe_detail", "duration_timerange")
	for _, s := range r.Sections {
		sectionInsert = sectionInsert.Values(s.Id, s.RecipeDetailId, s.TimeRange)
	}

	_, err = c.execTx(ctx, tx, sectionInsert)
	if err != nil {
		return fmt.Errorf("failed to insert sections: %w", err)
	}

	instructionsInsert := c.psql.Insert(sInstructionsTable).Columns("id", "section", "instruction")
	ingredientsInsert := c.psql.Insert(sIngredientsTable).Columns("id", "section", "ingredient", "recipe",
		"grams", "amount", "unit", "adjective", "optional", "original", "substitutes_for")

	var hasInstructions, hasIngredients bool
	for _, s := range r.Sections {
		for _, i := range s.Instructions {
			hasInstructions = true
			instructionsInsert = instructionsInsert.Values(i.Id, i.SectionId, i.Instruction)
		}

		for _, i := range s.Ingredients {
			hasIngredients = true

			ingredientsInsert = ingredientsInsert.Values(i.Id, i.SectionId, i.IngredientId, i.RecipeId,
				i.Grams, i.Amount, i.Unit, i.Adjective, i.Optional, i.Original, i.SubsFor)
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
		latestVersion, err := c.GetRecipeDetailWhere(ctx, sq.Eq{"recipe": modifying.RecipeId})
		if err != nil {
			return nil, fmt.Errorf("failed to find prior recipe: %w", err)
		}
		version = latestVersion.First().Version + 1
		parentID = latestVersion.First().RecipeId
	}
	r.Version = version
	r.Id = common.ID("rd")

	if parentID == "" {
		parentID = common.ID("r")
		_, err = c.execTx(ctx, tx, c.psql.Insert(recipesTable).Columns("id").Values(parentID))
		if err != nil {
			return nil, fmt.Errorf("failed to insert parent recipe: %w", err)
		}
	} else {
		// there is a parent and therefore other children, which are no longer latest
		_, err = c.execTx(ctx, tx, c.psql.
			Update(recipeDetailsTable).
			Set("is_latest_version", false).
			Where(sq.Eq{"recipe": parentID}),
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

	_, err = c.execTx(ctx, tx, c.psql.
		Insert(recipeDetailsTable).
		Columns("id", "recipe", "name", "version", "is_latest_version", "source").
		Values(r.Id, r.RecipeId, r.Name, r.Version, true, r.Source))
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
			Set("ingredient", nil).
			Set("recipe", newDetail.RecipeId).
			Where(sq.Eq{"ingredient": ingredientID})); err != nil {
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

// MergeIngredients sets the provided ingredients `same_as` to the first one.
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
			Set("same_as", ingredientID).
			Where(sq.Eq{"id": ids})); err != nil {
		return fmt.Errorf("failed to update ingredient: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
