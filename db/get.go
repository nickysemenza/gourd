package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/nickysemenza/gourd/graph/model"
	log "github.com/sirupsen/logrus"
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
	cacheKey := fmt.Sprintf("i:%s", uuid)

	tr := global.Tracer("db")
	ctx, span := tr.Start(ctx, "GetIngredientByUUID")
	defer span.End()

	cval, hit := c.cache.Get(cacheKey)
	log.WithField("key", cacheKey).WithField("hit", hit).Debug("cache:ingredients")
	if hit {
		ing, ok := cval.(Ingredient)
		if ok {
			return &ing, nil
		}
	}

	query, args, err := c.psql.Select("*").From(ingredientsTable).Where(sq.Eq{"uuid": uuid}).ToSql()
	if err != nil {
		return nil, err
	}
	ingredient := &Ingredient{}
	err = c.db.GetContext(ctx, ingredient, query, args...)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	c.cache.SetWithTTL(cacheKey, *ingredient, 0, time.Second)
	return ingredient, nil
}

// GetRecipeByUUID gets a recipe by UUID, shallowly.
func (c *Client) GetRecipeByUUID(ctx context.Context, uuid string) (*Recipe, error) {
	return c.getRecipe(ctx, c.psql.Select("*").From(recipesTable).Where(sq.Eq{"uuid": uuid}))
}

// GetRecipeByUUID gets a recipe by name, shallowly.
func (c *Client) GetRecipeByName(ctx context.Context, name string) (*Recipe, error) {
	return c.getRecipe(ctx, c.psql.Select("*").From(recipesTable).Where(sq.Eq{"name": name}))
}

func (c *Client) getRecipe(ctx context.Context, sb sq.SelectBuilder) (*Recipe, error) {
	query, args, err := sb.ToSql()
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
func (c *Client) GetRecipes(ctx context.Context, searchQuery string) ([]Recipe, error) {
	q := c.psql.Select("*").From(recipesTable)
	if searchQuery != "" {
		q = q.Where(sq.ILike{"name": fmt.Sprintf("%%%s%%", searchQuery)})
	}
	query, args, err := q.ToSql()
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

// GetRecipesWithIngredient gets all recipes with an ingredeitn
// todo: consolidate into getrecipes
func (c *Client) GetRecipesWithIngredient(ctx context.Context, ingredient string) ([]Recipe, error) {
	query, args, err := c.psql.Select(getRecipeColumns()...).From(recipesTable).
		Join("recipe_sections on recipe_sections.recipe = recipes.uuid").
		Join("recipe_section_ingredients on recipe_sections.uuid = recipe_section_ingredients.section").
		Where(sq.Eq{"ingredient": ingredient}).
		ToSql()
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
			if i.IngredientUUID.String != "" {
				ing, err := c.GetIngredientByUUID(ctx, i.IngredientUUID.String)
				if err != nil {
					return nil, err
				}
				r.Sections[x].Ingredients[y].RawIngredient = ing
			}
			if i.RecipeUUID.String != "" {
				rec, err := c.GetRecipeByUUID(ctx, i.RecipeUUID.String)
				if err != nil {
					return nil, err
				}
				r.Sections[x].Ingredients[y].RawRecipe = rec
			}
		}
	}

	return r, nil
}

// GetIngredients returns all ingredients.
func (c *Client) GetIngredients(ctx context.Context, searchQuery string) ([]Ingredient, error) {
	q := c.psql.Select("*").From(ingredientsTable)
	if searchQuery != "" {
		q = q.Where(sq.ILike{"name": fmt.Sprintf("%%%s%%", searchQuery)})
	}
	query, args, err := q.ToSql()
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

func (c *Client) GetMeals(ctx context.Context, recipe string) ([]*model.Meal, error) {
	query, args, err := c.psql.Select("meal_uuid AS uuid", "name", "notion_link AS notionURL").From("meals").
		LeftJoin("meal_recipe on meals.uuid = meal_recipe.meal_uuid").
		Where(sq.Eq{"recipe_uuid": recipe}).ToSql()
	if err != nil {
		return nil, err
	}
	var res []*model.Meal
	err = c.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
