package main

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/davecgh/go-spew/spew"
	"github.com/nickysemenza/food/manager"
)

const (
	host     = "localhost"
	port     = 5555
	user     = "food"
	password = "food"
	dbname   = "food"
)

func main() {
	fmt.Println("hello")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sqlx.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	m := manager.New(db)
	res, err := m.LoadFromFile(ctx, "recipes/chocolate-chip-cookies.yaml")
	// m.AssignUUIDs(res)
	spew.Dump(res, err)

	err = m.SaveRecipe(ctx, res)
	spew.Dump(err)

}
