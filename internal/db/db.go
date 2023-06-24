package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	servertiming "github.com/mitchellh/go-server-timing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Client is a database client
type Client struct {
	db     *sqlx.DB
	psql   sq.StatementBuilderType
	tracer trace.Tracer
}

func (c *Client) DB() *sqlx.DB {
	return c.db
}

type RecipeIngredientDependency struct {
	RecipeName     string `db:"recipe_name"`
	RecipeId       string `db:"recipe_id"`
	IngredientName string `db:"ingredient_name"`
	IngredientId   string `db:"ingredient_id"`
	IngredientKind string `db:"ingredient_kind"`
}

// New creates a new Client.
func New(dbConn *sql.DB) (*Client, error) {
	dbx := sqlx.NewDb(dbConn, "postgres")
	if err := dbx.Ping(); err != nil {
		return nil, err
	}

	return &Client{
		db:     dbx,
		psql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		tracer: otel.Tracer("db"),
	}, nil
}

// ConnnectionString returns a DSN.
func ConnnectionString(host, user, password, dbname string, port int64) string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
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
		span.RecordError(err)

		return fmt.Errorf("failed to GetContext: (%s %s) %w", query, args, err)
	}
	return nil
}
func (c *Client) selectContext(ctx context.Context, q sq.SelectBuilder, dest interface{}) error {
	ctx, span := c.tracer.Start(ctx, "selectContext")
	defer span.End()

	timing := servertiming.FromContext(ctx)
	defer timing.NewMetric("selectContext").Start().Stop()

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	err = c.db.SelectContext(ctx, dest, query, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		span.RecordError(err)

		return fmt.Errorf("failed to SelectContext: %w", err)
	}
	return nil
}

// nolint: unparam
func (c *Client) execContext(ctx context.Context, q sq.Sqlizer) (sql.Result, error) {
	ctx, span := c.tracer.Start(ctx, "execContext")
	defer span.End()

	timing := servertiming.FromContext(ctx)
	defer timing.NewMetric("execContext").Start().Stop()

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	res, err := c.execTx(ctx, tx, q)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return res, nil
}

func (c *Client) execTx(ctx context.Context, tx *sql.Tx, q sq.Sqlizer) (sql.Result, error) {
	ctx, span := c.tracer.Start(ctx, "execTx")
	defer span.End()

	query, args, err := q.ToSql()
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	span.SetAttributes(attribute.Int64("rows_affected", rows))
	return res, nil
}

type NotionRecipeMeta struct {
	Tags []string `json:"tags"`
}
