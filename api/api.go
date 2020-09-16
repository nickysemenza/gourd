//go:generate oapi-codegen --package api --generate types,client,chi-server,spec -o api.gen.go openapi.yaml

package api

import (
	"encoding/json"
	"net/http"
)

type API struct {
}

func NewAPI() *API {
	return &API{}
}

// List all recipes
// (GET /recipes)
func (a *API) ListRecipes(w http.ResponseWriter, r *http.Request, params ListRecipesParams) {

	// var x = 4

	var items = []Recipe{{Id: "foo", Name: "bar", Quantity: 52, Unit: "cookies"}}
	// var items2 = Recipes{Recipe{Id: 1}}

	resp := PaginatedRecipes{
		Recipes: &items,
		// Data2: &items2,
		Meta: &List{PageNumber: 2, Limit: 3},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Create a recipe
// (POST /recipes)
func (a *API) CreateRecipes(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

// Info for a specific recipe
// (GET /recipes/{recipeId})
func (a *API) GetRecipeById(w http.ResponseWriter, r *http.Request, recipeId string) {
	panic("not implemented") // TODO: Implement
}
