//go:generate oapi-codegen --package api --generate types,server,spec -o api.gen.go openapi.yaml

package api

import (
	"fmt"
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

func transformRecipe(dbr db.Recipe) Recipe {
	return Recipe{Id: dbr.UUID, Name: dbr.Name, Source: dbr.Source.Ptr()}
}
func transformIngredient(dbr db.Ingredient) Ingredient {
	return Ingredient{Id: dbr.UUID, Name: dbr.Name}
}

func transformRecipeSections(dbs []db.Section) []RecipeSection {
	var s []RecipeSection
	for _, d := range dbs {
		var ing []SectionIngredient
		var ins []SectionInstruction

		for _, i := range d.Instructions {
			ins = append(ins, SectionInstruction{Id: i.UUID, Instruction: i.Instruction})
		}
		for _, i := range d.Ingredients {
			g := float32(i.Grams.Float64)
			item := SectionIngredient{Id: i.UUID, Grams: &g}
			if i.RawRecipe != nil {
				item.Kind = "recipe"
				r := transformRecipe(*i.RawRecipe)
				item.Recipe = &r
			} else {
				item.Kind = "ingredient"
				i := transformIngredient(*i.RawIngredient)
				item.Ingredient = &i
			}
			ing = append(ing, item)
		}

		s = append(s, RecipeSection{
			Id:           d.UUID,
			Minutes:      int(d.Minutes.Int64),
			Ingredients:  ing,
			Instructions: ins,
		})
	}
	return s
}

// List all recipes
// (GET /recipes)
func (a *API) ListRecipes(c echo.Context, params ListRecipesParams) error {
	ctx := c.Request().Context()
	items := []Recipe{}

	paginationParams, listMeta := parsePagination(params.Offset, params.Limit)
	recipes, count, err := a.Manager.DB().GetRecipes(ctx, "", paginationParams...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}
	for _, r := range recipes {
		items = append(items, transformRecipe(r))
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
	var r RecipeDetail
	if err := c.Bind(&r); err != nil {
		err = fmt.Errorf("invalid format for input: %w", err)
		return sendErr(c, http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, r)
}

// Info for a specific recipe
// (GET /recipes/{recipeId})
func (a *API) GetRecipeById(c echo.Context, recipeId string) error {
	ctx := c.Request().Context()
	r, err := a.Manager.DB().GetRecipeByUUIDFull(ctx, recipeId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}
	detail := RecipeDetail{Recipe: transformRecipe(*r), Sections: transformRecipeSections(*&r.Sections)}
	return c.JSON(http.StatusOK, detail)
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
		items = append(items, transformIngredient(i))
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

func sendErr(ctx echo.Context, code int, err error) error {
	return ctx.JSON(code, Error{
		Message: err.Error(),
	})
}
