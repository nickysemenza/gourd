package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"go.opentelemetry.io/otel/api/global"
)

// GetRecipeSections finds the sections.
func (c *Client) GetRecipeSections(ctx context.Context, recipeUUID string) ([]Section, error) {
	query, args, err := c.psql.Select("*").From(sectionsTable).Where(sq.Eq{"recipe": recipeUUID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	var sections []Section
	err = c.db.SelectContext(ctx, &sections, query, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return sections, nil
}

// GetSectionInstructions finds the instructions for a section.
func (c *Client) GetSectionInstructions(ctx context.Context, sectionUUID string) ([]SectionInstruction, error) {
	query, args, err := c.psql.Select("*").From(sInstructionsTable).Where(sq.Eq{"section": sectionUUID}).ToSql()
	if err != nil {
		return nil, err
	}
	var res []SectionInstruction
	err = c.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetSectionIngredients finds the ingredients for a section.
func (c *Client) GetSectionIngredients(ctx context.Context, sectionUUID string) ([]SectionIngredient, error) {
	query, args, err := c.psql.Select("*").From(sIngredientsTable).Where(sq.Eq{"section": sectionUUID}).ToSql()
	if err != nil {
		return nil, err
	}
	var res []SectionIngredient
	err = c.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetIngredientByUUID finds an ingredient.
func (c *Client) GetIngredientByUUID(ctx context.Context, uuid string) (*Ingredient, error) {
	tr := global.Tracer("db")
	ctx, span := tr.Start(ctx, "GetIngredientByUUID")
	defer span.End()

	query, args, err := c.psql.Select("*").From(ingredientsTable).Where(sq.Eq{"uuid": uuid}).ToSql()
	if err != nil {
		return nil, err
	}
	ingredient := &Ingredient{}
	err = c.db.GetContext(ctx, ingredient, query, args...)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return ingredient, nil
}

// GetRecipeByUUID gets a recipe by UUID, shallowly.
func (c *Client) GetRecipeByUUID(ctx context.Context, uuid string) (*Recipe, error) {
	query, args, err := c.psql.Select("*").From(recipesTable).Where(sq.Eq{"uuid": uuid}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	r := &Recipe{}
	err = c.db.GetContext(ctx, r, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return r, nil
}

// GetRecipes returns all recipes, shallowly.
func (c *Client) GetRecipes(ctx context.Context) ([]Recipe, error) {
	query, args, err := c.psql.Select("*").From(recipesTable).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	r := []Recipe{}
	err = c.db.SelectContext(ctx, &r, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return r, nil
}

// GetRecipeByUUIDFull gets a recipe by UUID, with all dependencies.
func (c *Client) GetRecipeByUUIDFull(ctx context.Context, uuid string) (*Recipe, error) {
	r, err := c.GetRecipeByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return r, nil
	}

	r.Sections, err = c.GetRecipeSections(ctx, uuid)
	if err != nil {
		return nil, err
	}

	for x, s := range r.Sections {
		r.Sections[x].Instructions, err = c.GetSectionInstructions(ctx, s.UUID)
		if err != nil {
			return nil, err
		}
		r.Sections[x].Ingredients, err = c.GetSectionIngredients(ctx, s.UUID)
		if err != nil {
			return nil, err
		}

		for y, i := range r.Sections[x].Ingredients {
			ing, err := c.GetIngredientByUUID(ctx, i.IngredientUUID.String)
			if err != nil {
				return nil, err
			}
			if ing != nil {
				r.Sections[x].Ingredients[y].Name = ing.Name
			}
		}
	}

	return r, nil
}

// GetIngredients returns all ingredients.
func (c *Client) GetIngredients(ctx context.Context) ([]Ingredient, error) {
	query, args, err := c.psql.Select("*").From(ingredientsTable).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	i := []Ingredient{}
	err = c.db.SelectContext(ctx, &i, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return i, nil
}
