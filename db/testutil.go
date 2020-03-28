package db

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // for pg
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func newDB(t *testing.T) *Client {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5555)
	viper.SetDefault("DB_USER", "food")
	viper.SetDefault("DB_PASSWORD", "food")
	viper.SetDefault("DB_DBNAME", "food")
	viper.AutomaticEnv()

	dbConn, err := sqlx.Open("postgres", ConnnectionString(
		viper.GetString("DB_HOST"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_DBNAME"),
		viper.GetInt64("DB_PORT")))
	require.NoError(t, err)
	return New(dbConn)

}
