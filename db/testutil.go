//go:build integration
// +build integration

package db

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq" // for pg
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// NewTestDB makes a test DB.
func NewTestDB(t *testing.T, kind Kind) *Client {
	viper.SetDefault("DATABASE_URL", "postgres://gourd:gourd@localhost:5555/food")
	viper.SetDefault("DATABASE_URL_USDA", "postgres://gourd:gourd@localhost:5556/usda")
	viper.AutomaticEnv()

	var dsn string
	if kind == USDA {
		dsn = viper.GetString("DATABASE_URL_USDA") + "?sslmode=disable"
	} else {
		dsn = viper.GetString("DATABASE_URL") + "?sslmode=disable"
	}
	dbConn, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	if kind == Gourd {
		err = AutoMigrate(dbConn, false)
		require.NoError(t, err)
		err = AutoMigrate(dbConn, true)
		require.NoError(t, err)
	}

	d, err := New(dbConn, kind)
	require.NoError(t, err)
	return d
}
