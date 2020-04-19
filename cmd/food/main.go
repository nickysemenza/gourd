package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nickysemenza/food/db"
	"github.com/nickysemenza/food/manager"
	"github.com/nickysemenza/food/server"
	"github.com/spf13/viper"
)

const ()

func main() {

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5555)
	viper.SetDefault("DB_USER", "food")
	viper.SetDefault("DB_PASSWORD", "food")
	viper.SetDefault("DB_DBNAME", "food")

	viper.AutomaticEnv()

	dbConn, err := sqlx.Open("postgres", db.ConnnectionString(
		viper.GetString("DB_HOST"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_DBNAME"),
		viper.GetInt64("DB_PORT")))

	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	if err = dbConn.Ping(); err != nil {
		panic(err)
	}
	dbClient := db.New(dbConn)
	m := manager.New(dbClient)
	s := server.Server{Manager: m, HTTPPort: 4242, DB: dbClient}

	spew.Dump(s.Run())

}
