package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/dgraph-io/ristretto"
	"github.com/jmoiron/sqlx"
	"github.com/nickysemenza/gourd/common"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/guregu/null.v3/zero"
)

const (
	sIngredientsTable  = "recipe_section_ingredients"
	sInstructionsTable = "recipe_section_instructions"
	sectionsTable      = "recipe_sections"
	recipesTable       = "recipes"
	recipeDetailsTable = "recipe_details"
	ingredientsTable   = "ingredients"
)

// Client is a database client
type Client struct {
	db     *sqlx.DB
	psql   sq.StatementBuilderType
	cache  *ristretto.Cache
	tracer trace.Tracer
}

func getRecipeDetailColumns() []string {
	return []string{
		"recipe_details.id",
		"name", "version",
		"total_minutes", "equipment",
		"source", "servings",
		"quantity", "recipe_details.unit"}
}

// RecipeDetail represents a recipe
type Recipe struct {
	Id     string `db:"id"`
	Detail RecipeDetail
}

// RecipeDetail represents a recipe
type RecipeDetail struct {
	Id           string      `db:"id"`
	RecipeId     string      `db:"recipe"`
	Name         string      `db:"name"`
	TotalMinutes zero.Int    `db:"total_minutes"`
	Equipment    zero.String `db:"equipment"`
	Source       zero.String `db:"source"`
	Servings     zero.Int    `db:"servings"`
	Quantity     zero.Int    `db:"quantity"`
	Unit         zero.String `db:"unit"`
	Version      int64       `db:"version"`
	Sections     []Section
}

// Section represents a Section
type Section struct {
	Id             string   `db:"id"`
	RecipeDetailId string   `db:"recipe_detail"`
	Minutes        zero.Int `db:"minutes"`
	Sort           zero.Int `db:"sort"`
	Ingredients    []SectionIngredient
	Instructions   []SectionInstruction
}

// SectionIngredient is a foo
type SectionIngredient struct {
	Id        string      `db:"id"`
	SectionId string      `db:"section"`
	Sort      zero.Int    `db:"sort"`
	Grams     zero.Float  `db:"grams"`
	Amount    zero.Float  `db:"amount"`
	Unit      zero.String `db:"unit"`
	Adjective zero.String `db:"adjective"`
	Optional  zero.Bool   `db:"optional"`

	// one of the following is required for get and update:
	RecipeId     zero.String `db:"recipe"`
	IngredientId zero.String `db:"ingredient"`

	// one of these is populated via gets
	RawRecipe     *RecipeDetail
	RawIngredient *Ingredient

	// deprecated
	// Name        string      // todo: use this to load an Ingredient
}

// SectionInstruction represents a SectionInstruction
type SectionInstruction struct {
	Id          string   `db:"id"`
	Sort        zero.Int `db:"sort"`
	Instruction string   `db:"instruction"`
	SectionId   string   `db:"section"`
}

// Ingredient is a globally-scoped ingredient
type Ingredient struct {
	Id     string      `json:"id"`
	Name   string      `json:"name"`
	FdcID  zero.Int    `db:"fdc_id"`
	SameAs zero.String `db:"same_as"`
}

// New creates a new Client.
func New(dbConn *sql.DB) (*Client, error) {
	dbx := sqlx.NewDb(dbConn, "postgres")
	if err := dbx.Ping(); err != nil {
		return nil, err
	}

	// nolint:gomnd
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
		Metrics:     true,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		db:     dbx,
		psql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		cache:  cache,
		tracer: otel.Tracer("db"),
	}, nil
}

// ConnnectionString returns a DSN.
func ConnnectionString(host, user, password, dbname string, port int64) string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}
func setId(val string) string {
	if val != "" {
		return val
	}
	return common.UUID()
}

// AssignIds adds ids where missing.
func (c *Client) AssignIds(ctx context.Context, r *RecipeDetail) error {
	// r.Id = setId(r.Id)
	for x := range r.Sections {
		r.Sections[x].Id = common.UUID()
		r.Sections[x].RecipeDetailId = r.Id
		for y := range r.Sections[x].Ingredients {
			r.Sections[x].Ingredients[y].Id = common.UUID()
			r.Sections[x].Ingredients[y].SectionId = r.Sections[x].Id

		}
		for y := range r.Sections[x].Instructions {
			r.Sections[x].Instructions[y].Id = common.UUID()
			r.Sections[x].Instructions[y].SectionId = r.Sections[x].Id
		}
	}
	return nil
}

// IngredientByName retrieves an ingredient by name, creating it if it does not exist.
func (c *Client) IngredientByName(ctx context.Context, name string) (*Ingredient, error) {
	ingredient := &Ingredient{}
	err := c.db.GetContext(ctx, ingredient, `SELECT * FROM ingredients
	WHERE name = $1 LIMIT 1`, name)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = c.db.ExecContext(ctx, `INSERT INTO ingredients (id, name) VALUES ($1, $2)`, setId(""), name)
		if err != nil {
			return nil, err
		}
		return c.IngredientByName(ctx, name)
	}
	return ingredient, err
}

//nolint: funlen
func (c *Client) updateRecipe(ctx context.Context, tx *sql.Tx, r *RecipeDetail) error {
	query, args, err := c.psql.
		Update(recipeDetailsTable).Where(sq.Eq{"id": r.Id}).Set("name", r.Name).
		Set("total_minutes", r.TotalMinutes).
		Set("unit", r.Unit).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	if err := c.AssignIds(ctx, r); err != nil {
		return err
	}
	// update recipe_section_instructions
	// update recipe_section_ingredients
	// update recipe_sections

	// c.psql.Delete(sectionInstructionsTable).Where(sq.Eq{""})

	if len(r.Sections) == 0 {
		return nil
	}

	// sections
	sectionInsert := c.psql.Insert(sectionsTable).Columns("id", "recipe_detail", "minutes")
	for _, s := range r.Sections {
		sectionInsert = sectionInsert.Values(s.Id, s.RecipeDetailId, s.Minutes)
	}
	query, args, err = sectionInsert.ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	for _, s := range r.Sections {
		if len(s.Instructions) > 0 {
			instructionsInsert := c.psql.Insert(sInstructionsTable).Columns("id", "section", "instruction")
			for _, i := range s.Instructions {
				instructionsInsert = instructionsInsert.Values(i.Id, i.SectionId, i.Instruction)
			}
			if _, err = instructionsInsert.RunWith(tx).ExecContext(ctx); err != nil {
				return err
			}
		}

		if len(s.Ingredients) > 0 {
			ingredientsInsert := c.psql.Insert(sIngredientsTable).Columns("id", "section", "ingredient", "recipe",
				"grams", "amount", "unit", "adjective", "optional")
			for _, i := range s.Ingredients {
				ingredientsInsert = ingredientsInsert.Values(i.Id, i.SectionId, i.IngredientId, i.RecipeId,
					i.Grams, i.Amount, i.Unit, i.Adjective, i.Optional)
			}
			if _, err = ingredientsInsert.RunWith(tx).ExecContext(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Client) insertRecipe(ctx context.Context, r *RecipeDetail) (*RecipeDetail, error) {
	ctx, span := c.tracer.Start(ctx, "insertRecipe")
	defer span.End()
	r.Id = setId(r.Id)
	log.Println("inserting", r.Id, r.Name)
	query, args, err := c.psql.
		Insert(recipeDetailsTable).Columns("id", "recipe", "name", "version").Values(r.Id, r.RecipeId, r.Name, r.Version).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	err = c.updateRecipe(ctx, tx, r)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return c.GetRecipeDetailByIdFull(ctx, r.Id)
}

// InsertRecipe inserts a recipe.
func (c *Client) InsertRecipe(ctx context.Context, r *RecipeDetail) (*RecipeDetail, error) {
	ctx, span := c.tracer.Start(ctx, "InsertRecipe")
	defer span.End()

	// if we have an existing recipe with the same Id or name, this one is a n+1 version of that one
	version := int64(1)
	var modifying *RecipeDetail
	parentID := ""
	var err error
	if r.Id != "" {
		modifying, err = c.GetRecipeDetailWhere(ctx, sq.Eq{"id": r.Id})
		if err != nil {
			return nil, fmt.Errorf("failed to find prior recipe: %w", err)
		}
	}
	if modifying == nil {
		modifying, err = c.GetRecipeDetailWhere(ctx, sq.Eq{"name": r.Name})
		if err != nil {
			return nil, fmt.Errorf("failed to find prior recipe: %w", err)
		}
	}

	if modifying != nil {
		latestVersion, err := c.GetRecipeDetailWhere(ctx, sq.Eq{"recipe": modifying.RecipeId})
		if err != nil {
			return nil, fmt.Errorf("failed to find prior recipe: %w", err)
		}
		version = latestVersion.Version + 1
		parentID = latestVersion.RecipeId
	}
	r.Version = version
	r.Id = common.UUID()

	if parentID == "" {
		parentID = common.UUID()
		_, err = c.execContext(ctx, c.psql.
			Insert(recipesTable).Columns("id").Values(parentID))
		if err != nil {
			return nil, fmt.Errorf("failed to insert parent recipe: %w", err)
		}
	}
	r.RecipeId = parentID

	return c.insertRecipe(ctx, r)

}

func (c *Client) getContext(ctx context.Context, q sq.SelectBuilder, dest interface{}) error {
	ctx, span := c.tracer.Start(ctx, "getContext")
	defer span.End()

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	err = c.db.GetContext(ctx, dest, query, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to GetContext: (%s %s) %w", query, args, err)
	}
	return nil
}
func (c *Client) selectContext(ctx context.Context, q sq.SelectBuilder, dest interface{}) error {
	ctx, span := c.tracer.Start(ctx, "selectContext")
	defer span.End()

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	err = c.db.SelectContext(ctx, dest, query, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to SelectContext: %w", err)
	}
	return nil
}

// nolint: unparam
func (c *Client) execContext(ctx context.Context, q sq.Sqlizer) (sql.Result, error) {
	ctx, span := c.tracer.Start(ctx, "execContext")
	defer span.End()

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	return c.db.ExecContext(ctx, query, args...)
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

	newRecipe, err := c.InsertRecipe(ctx, &RecipeDetail{Name: i.Name})
	if err != nil {
		return nil, err
	}

	if _, err = c.execContext(ctx,
		c.psql.
			Update(sIngredientsTable).
			Set("ingredient", nil).
			Set("recipe", newRecipe.RecipeId).
			Where(sq.Eq{"ingredient": ingredientID})); err != nil {
		return nil, fmt.Errorf("failed to update references to transformed ingredient: %w", err)
	}

	if _, err = c.execContext(ctx, c.psql.
		Update(ingredientsTable).
		Set("name", fmt.Sprintf("[deprecated] %s", i.Name)).
		Where(sq.Eq{"id": i.Id})); err != nil {
		return nil, fmt.Errorf("failed to deprecated ingredient: %w", err)
	}

	return newRecipe, nil
}
