//go:generate oapi-codegen --package api --generate types,client,chi-server,spec -o api.gen.go openapi.yaml

package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/manager"
)

type API struct {
	*manager.Manager
}

func NewAPI(m *manager.Manager) *API {
	return &API{Manager: m}
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

	success(w, resp)
}
func writeErr(ctx context.Context, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	e := Error{Message: err.Error()}
	json.NewEncoder(w).Encode(e)
}
func success(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
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

func (a *API) ListIngredients(w http.ResponseWriter, r *http.Request, params ListIngredientsParams) {

	ctx := r.Context()
	items := []Ingredient{}

	paginationParams, listMeta := parsePagination(params.Offset, params.Limit)
	ing, count, err := a.Manager.DB().GetIngredients(ctx, "", paginationParams...)
	if err != nil {
		writeErr(ctx, w, err)
		return
	}
	for _, i := range ing {
		items = append(items, Ingredient{Id: i.UUID, Name: i.Name})
	}
	listMeta.TotalCount = int(count)

	resp := PaginatedIngredients{
		Ingredients: &items,
		Meta:        listMeta,
	}
	success(w, resp)
}

func parsePagination(o *OffsetParam, l *LimitParam) ([]db.SearchOption, *List) {
	offset := 0
	limit := 20
	if o != nil {
		offset = int(*o)
	}
	if l != nil {
		limit = int(*l)
	}
	return []db.SearchOption{db.WithOffset(uint64(offset)), db.WithLimit(uint64(limit))}, &List{Offset: offset, Limit: limit, PageNumber: (offset/limit + 1)}
}
