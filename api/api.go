//go:generate oapi-codegen --package api --generate types,server,spec -o api.gen.go openapi.yaml

package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
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
func (a *API) ListRecipes(c echo.Context, params ListRecipesParams) error {
	ctx := c.Request().Context()
	items := []Recipe{}

	paginationParams, listMeta := parsePagination(params.Offset, params.Limit)
	ing, count, err := a.Manager.DB().GetRecipes(ctx, "", paginationParams...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}
	for _, i := range ing {
		items = append(items, Recipe{Id: i.UUID, Name: i.Name, Source: i.Source.Ptr()})
	}
	listMeta.setTotalCount(count)

	resp := PaginatedRecipes{
		Recipes: &items,
		Meta:    listMeta,
	}
	return c.JSON(http.StatusOK, resp)
}

// Create a recipe
// (POST /recipes)
func (a *API) CreateRecipes(c echo.Context) error {
	panic("not implemented") // TODO: Implement
}

// Info for a specific recipe
// (GET /recipes/{recipeId})
func (a *API) GetRecipeById(c echo.Context, recipeId string) error {
	panic("not implemented") // TODO: Implement
}

func (a *API) ListIngredients(c echo.Context, params ListIngredientsParams) error {
	ctx := c.Request().Context()
	items := []Ingredient{}

	paginationParams, listMeta := parsePagination(params.Offset, params.Limit)
	ing, count, err := a.Manager.DB().GetIngredients(ctx, "", paginationParams...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}
	for _, i := range ing {
		items = append(items, Ingredient{Id: i.UUID, Name: i.Name})
	}
	listMeta.TotalCount = int(count)

	resp := PaginatedIngredients{
		Ingredients: &items,
		Meta:        listMeta,
	}
	return c.JSON(http.StatusOK, resp)
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

func (l *List) setTotalCount(count uint64) {
	c := int(count)
	l.TotalCount = c
	l.PageCount = c / l.Limit
}
