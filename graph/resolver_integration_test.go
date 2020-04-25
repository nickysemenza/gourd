package graph

import (
	"fmt"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/nickysemenza/food/db"
	"github.com/nickysemenza/food/graph/generated"
	"github.com/nickysemenza/food/manager"
	"github.com/stretchr/testify/require"
)

func TestCreateUpdateList(t *testing.T) {
	tdb := db.NewDB(t)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{
			Manager: manager.New(tdb),
			DB:      tdb,
		}},
	))
	srv.Use(Observability{})

	// create a recipe
	var resp struct {
		CreateRecipe struct {
			UUID string
		}
	}
	c := client.New(srv)
	err := c.Post(`mutation{
		createRecipe(recipe: {name: "`+fmt.Sprintf("rr-%d", time.Now().Unix())+`"}) {uuid}
	  }`, &resp)

	require.NoError(t, err)
	newUUID := resp.CreateRecipe.UUID
	require.NotEmpty(t, newUUID)

	// update recipe name
	var resp2 struct {
		UpdateRecipe struct {
			UUID string
		}
	}
	newName := fmt.Sprintf("name2-%d", time.Now().Unix())
	err = c.Post(`mutation{
		updateRecipe(recipe: {uuid: "`+newUUID+`",name: "`+newName+`"}) {uuid}
	  }`, &resp2)

	require.NoError(t, err)

	// ensure recipe is in the getlist
	var resp3 struct {
		Recipes []struct {
			UUID         string
			Name         string
			Sections     interface{}
			TotalMinutes interface{} `json:"total_minutes"`
			Unit         interface{}
		}
	}
	err = c.Post(`query {
		recipes {
		  uuid
		  name
		  total_minutes
		  unit
		  sections {
			minutes
			ingredients {
			  uuid
			  info {
				name
			  }
			  grams
			}
			instructions {
			  instruction
			  uuid
			}
		  }
		}
	  }
	  `, &resp3)
	require.NoError(t, err)

	found := false
	for _, x := range resp3.Recipes {
		if x.UUID == newUUID {
			found = true
			require.Equal(t, x.Name, newName)
		}
	}
	require.True(t, found)

	// t.Fatal()

}
