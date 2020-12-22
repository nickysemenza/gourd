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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
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

// AssignIds adds ids where missing.
func (c *Client) AssignIds(ctx context.Context, r *RecipeDetail) error {
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
	res, err := c.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	span.SetAttributes(label.Int64("rows_affected", rows))
	return res, nil
}
