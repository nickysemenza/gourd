package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
)

// Client is a database client
type Client struct {
	db *sqlx.DB
}

// Recipe represents a recipe
type Recipe struct {
	UUID         string         `db:"uuid"`
	Name         string         `db:"name"`
	TotalMinutes sql.NullInt64  `db:"total_minutes"`
	Equipment    sql.NullString `db:"equipment"`
	Source       sql.NullString `db:"source"`
	Servings     sql.NullInt64  `db:"servings"`
	Quantity     sql.NullString `db:"quantity"`
	Unit         sql.NullString `db:"unit"`
	Sections     []Section
}

// Section represents a Section
type Section struct {
	UUID         string        `db:"uuid"`
	RecipeUUID   string        `db:"recipe"`
	Minutes      sql.NullInt64 `db:"minutes"`
	Sort         sql.NullInt64 `db:"sort"`
	Ingredients  []SectionIngredient
	Instructions []SectionInstruction
}

// SectionIngredient is a foo
type SectionIngredient struct {
	UUID        string          `db:"uuid"`
	Sort        sql.NullInt64   `db:"sort"`
	Name        string          // todo: use this to load an Ingredient
	Grams       sql.NullFloat64 `db:"grams"`
	SectionUUID string          `db:"section"`
	Amount      sql.NullFloat64 `db:"amount"`
	Unit        sql.NullString  `db:"unit"`
	Adjective   sql.NullString  `db:"adjective"`
	Optional    bool            `db:"optional"`

	// one of the following:
	RecipeUUID     string `db:"recipe"`
	IngredientUUID string `db:"ingredient"`
}

// SectionInstruction represents a SectionInstruction
type SectionInstruction struct {
	UUID        string        `db:"uuid"`
	Sort        sql.NullInt64 `db:"sort"`
	Instruction string        `db:"instruction"`
	SectionUUID string        `db:"section"`
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// New creates a new Client
func New(db *sqlx.DB) *Client {
	return &Client{db: db}
}

// ConnnectionString returns a DSN
func ConnnectionString(host, user, password, dbname string, port int64) string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}
func setUUID(val string) string {
	if val != "" {
		return val
	}
	u, _ := uuid.NewV4()
	return u.String()
}

// AssignUUIDs adds uuids where missing
func (c *Client) AssignUUIDs(ctx context.Context, r *Recipe) error {
	r.UUID = setUUID(r.UUID)
	for x := range r.Sections {
		r.Sections[x].UUID = setUUID(r.Sections[x].UUID)
		r.Sections[x].RecipeUUID = r.UUID
		for y := range r.Sections[x].Ingredients {
			r.Sections[x].Ingredients[y].UUID = setUUID(r.Sections[x].Ingredients[y].UUID)
			r.Sections[x].Ingredients[y].SectionUUID = r.Sections[x].UUID
			// ing, err := m.IngredientByName(ctx, r.Sections[x].Ingredients[y].Name)
			// if err != nil {
			// return err
			// }
			// r.Sections[x].Ingredients[y].IngredientUUID = ing.UUID
		}
		for y := range r.Sections[x].Instructions {
			r.Sections[x].Instructions[y].UUID = setUUID(r.Sections[x].Instructions[y].UUID)
			r.Sections[x].Instructions[y].SectionUUID = r.Sections[x].UUID
		}
	}
	return nil
}

// GetRecipeByUUID asdasd
func (c *Client) GetRecipeByUUID(ctx context.Context, uuid string) (*Recipe, error) {

	query, args, err := psql.Select("*").From("recipes").Where(sq.Eq{"uuid": uuid}).ToSql()
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

	query, args, err = psql.Select("*").From("recipe_sections").Where(sq.Eq{"recipe": uuid}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	// var sections *[]Section
	err = c.db.SelectContext(ctx, &r.Sections, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	// for _, s := range *sections {
	// 	r.Sections = append(r.Sections, s)
	// }

	return r, nil

}

func (c *Client) updateRecipe(ctx context.Context, tx *sql.Tx, r *Recipe) error {
	query, args, err := psql.
		Update("recipes").Where(sq.Eq{"uuid": r.UUID}).Set("name", r.Name).
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
	c.AssignUUIDs(ctx, r)
	//update recipe_section_instructions
	//update recipe_section_ingredients
	//update recipe_sections

	_, err = c.db.ExecContext(ctx, `DELETE FROM recipe_sections WHERE recipe = $1`, r.UUID)
	if err != nil {
		return err
	}

	if len(r.Sections) == 0 {
		return nil
	}

	sectionInsert := psql.Insert("recipe_sections").Columns("uuid", "recipe", "minutes")

	for _, s := range r.Sections {
		sectionInsert = sectionInsert.Values(s.UUID, s.RecipeUUID, s.Minutes)
	}
	query, args, err = sectionInsert.ToSql()

	spew.Dump(sectionInsert)
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil

}

// UpdateRecipe updates a recipe
func (c *Client) UpdateRecipe(ctx context.Context, r *Recipe) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	err = c.updateRecipe(ctx, tx, r)
	if err != nil {
		return err
	}
	return tx.Commit()

}

// InsertRecipe inserts a recipe
func (c *Client) InsertRecipe(ctx context.Context, r *Recipe) error {
	query, args, err := psql.
		Insert("recipes").Columns("uuid", "name").Values(r.UUID, r.Name).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	err = c.updateRecipe(ctx, tx, r)
	if err != nil {
		return err
	}
	return tx.Commit()

}
