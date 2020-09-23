//go:generate oapi-codegen --package api --generate server,spec -o api-server.gen.go openapi.yaml
//go:generate oapi-codegen --package api --generate types -o api-types.gen.go openapi.yaml
// todo go:generate oapi-codegen --package api --generate client -o api-client.gen.go openapi.yaml

package api

import (
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/manager"
	"gopkg.in/guregu/null.v3/zero"
)

type API struct {
	*manager.Manager
}

func NewAPI(m *manager.Manager) *API {
	return &API{Manager: m}
}

func transformRecipe(dbr db.Recipe) Recipe {
	return Recipe{
		Id:           dbr.UUID,
		Name:         dbr.Name,
		Source:       dbr.Source.Ptr(),
		TotalMinutes: dbr.TotalMinutes.Ptr(),
	}
}
func transformRecipeFull(dbr *db.Recipe) RecipeDetail {
	return RecipeDetail{Recipe: transformRecipe(*dbr), Sections: transformRecipeSections(*&dbr.Sections)}
}
func transformIngredient(dbr db.Ingredient) Ingredient {
	return Ingredient{Id: dbr.UUID, Name: dbr.Name}
}
func (i *Ingredient) toDB() *db.Ingredient {
	return &db.Ingredient{UUID: i.Id, Name: i.Name}
}

func (r *RecipeDetail) toDB() *db.Recipe {
	dbr := db.Recipe{
		UUID:         r.Recipe.Id,
		Name:         r.Recipe.Name,
		Source:       zero.StringFromPtr(r.Recipe.Source),
		TotalMinutes: zero.IntFromPtr(r.Recipe.TotalMinutes),
	}

	for _, s := range r.Sections {
		dbs := db.Section{
			Minutes: zero.IntFrom(s.Minutes),
		}
		for _, i := range s.Instructions {
			dbs.Instructions = append(dbs.Instructions, db.SectionInstruction{
				Instruction: i.Instruction,
			})
		}
		for _, i := range s.Ingredients {
			si := db.SectionIngredient{
				Grams:     zero.FloatFromPtr(i.Grams),
				Amount:    zero.FloatFromPtr(i.Amount),
				Unit:      zero.StringFromPtr(i.Unit),
				Adjective: zero.StringFromPtr(i.Adjective),
				Optional:  zero.BoolFromPtr(i.Optional),
			}
			if i.Kind == "recipe" {
				si.RecipeUUID = zero.StringFrom(i.Recipe.Id)
			} else {
				si.IngredientUUID = zero.StringFrom(i.Ingredient.Id)
			}

			dbs.Ingredients = append(dbs.Ingredients, si)
		}

		dbr.Sections = append(dbr.Sections, dbs)
	}

	spew.Dump(dbr)
	return &dbr

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
			item := SectionIngredient{
				Id:        i.UUID,
				Grams:     i.Grams.Ptr(),
				Amount:    i.Amount.Ptr(),
				Unit:      i.Unit.Ptr(),
				Adjective: i.Adjective.Ptr(),
				Optional:  i.Optional.Ptr(),
			}
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
			Minutes:      d.Minutes.Int64,
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
	ctx := c.Request().Context()
	var r RecipeDetail
	if err := c.Bind(&r); err != nil {
		err = fmt.Errorf("invalid format for input: %w", err)
		return sendErr(c, http.StatusBadRequest, err)
	}
	uuid, err := a.DB().InsertRecipe(ctx, r.toDB())
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}
	r2, err := a.Manager.DB().GetRecipeByUUIDFull(ctx, uuid)
	return c.JSON(http.StatusCreated, transformRecipeFull(r2))
}

// Info for a specific recipe
// (GET /recipes/{recipeId})
func (a *API) GetRecipeById(c echo.Context, recipeId string) error {
	ctx := c.Request().Context()
	r, err := a.Manager.DB().GetRecipeByUUIDFull(ctx, recipeId)
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, transformRecipeFull(r))
}

func (a *API) CreateIngredients(c echo.Context) error {
	ctx := c.Request().Context()
	var i Ingredient
	if err := c.Bind(&i); err != nil {
		err = fmt.Errorf("invalid format for input: %w", err)
		return sendErr(c, http.StatusBadRequest, err)
	}
	ing, err := a.DB().IngredientByName(ctx, i.Name)
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, transformIngredient(*ing))

}
func (a *API) ListIngredients(c echo.Context, params ListIngredientsParams) error {
	ctx := c.Request().Context()
	items := []IngredientDetail{}

	paginationParams, listMeta := parsePagination(params.Offset, params.Limit)
	ing, count, err := a.Manager.DB().GetIngredients(ctx, "", paginationParams...)
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}
	for _, i := range ing {
		items = append(items, IngredientDetail{Ingredient: transformIngredient(i)})
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
