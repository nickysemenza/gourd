// +build integration

package db

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq" // for pg
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// NewDB makes a test DB.
func NewDB(t *testing.T) *Client {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5555)
	viper.SetDefault("DB_USER", "gourd")
	viper.SetDefault("DB_PASSWORD", "gourd")
	viper.SetDefault("DB_DBNAME", "food")
	viper.AutomaticEnv()

	dbConn, err := sql.Open("postgres", ConnnectionString(
		viper.GetString("DB_HOST"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_DBNAME"),
		viper.GetInt64("DB_PORT")))
	require.NoError(t, err)
	d, err := New(dbConn)
	require.NoError(t, err)
	return d
}
