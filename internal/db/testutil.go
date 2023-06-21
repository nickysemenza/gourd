//go:build integration
// +build integration

package db

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq" // for pg
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// NewTestDB makes a test DB.
func NewTestDB(t *testing.T) *Client {
	viper.SetDefault("DATABASE_URL", "postgres://gourd:gourd@localhost:5555/food")
	viper.AutomaticEnv()

	dsn := viper.GetString("DATABASE_URL") + "?sslmode=disable"
	dbConn, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	ctx := context.Background()
	err = AutoMigrate(ctx, dbConn, false)
	require.NoError(t, err)
	err = AutoMigrate(ctx, dbConn, true)
	require.NoError(t, err)

	d, err := New(dbConn)
	require.NoError(t, err)
	return d
}
