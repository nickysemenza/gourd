package main

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nickysemenza/food/db"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
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
	err = dbConn.Ping()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	c := db.New(dbConn)
	res, err := c.GetRecipeByUUID(ctx, "fb1d53ef-47e0-4de2-bc68-9773f5353089")
	spew.Dump(res, err)
	y, _ := yaml.Marshal(res)
	fmt.Printf("%s", y)

	// spew.Dump(c.InsertRecipe(ctx, &db.Recipe{UUID: "cz", Name: "azz"}))

	// m := manager.New(db)
	// res, err := m.LoadFromFile(ctx, "recipes/chocolate-chip-cookies.yaml")
	// // m.AssignUUIDs(res)
	// spew.Dump(res, err)

	// err = m.SaveRecipe(ctx, res)
	// spew.Dump(err)

}
