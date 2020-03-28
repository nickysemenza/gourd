package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
	// Sections     []Section `db:"sections"`
}

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

// GetRecipeByUUID asdasd
func (c *Client) GetRecipeByUUID(ctx context.Context, uuid string) (*Recipe, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

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

	return r, nil

}

func (c *Client) updateRecipe(ctx context.Context, tx *sql.Tx, r Recipe) error {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("recipes").Where(sq.Eq{"uuid": r.UUID}).Set("name", r.Name).
		Set("total_minutes", r.TotalMinutes).
		Set("unit", r.Unit).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	return err
}

// UpdateRecipe updates a recipe
func (c *Client) UpdateRecipe(ctx context.Context, r Recipe) error {
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
func (c *Client) InsertRecipe(ctx context.Context, r Recipe) error {
	query, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
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
