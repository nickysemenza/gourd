//go:generate oapi-codegen --package api --generate server,spec -o api-server.gen.go openapi.yaml
//go:generate oapi-codegen --package api --generate types -o api-types.gen.go openapi.yaml
// todo go:generate oapi-codegen --package api --generate client -o api-client.gen.go openapi.yaml

package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/manager"
	"go.opentelemetry.io/otel/api/global"
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
func transformRecipeFull(dbr *db.Recipe) *RecipeDetail {
	return &RecipeDetail{Recipe: transformRecipe(*dbr), Sections: transformRecipeSections(*&dbr.Sections)}
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
				Grams:     zero.FloatFrom(i.Grams),
				Amount:    zero.FloatFromPtr(i.Amount),
				Unit:      zero.StringFrom(i.Unit),
				Adjective: zero.StringFrom(i.Adjective),
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
	s := make([]RecipeSection, 0)
	for _, d := range dbs {
		ing := make([]SectionIngredient, 1)
		ins := make([]SectionInstruction, 1)

		for _, i := range d.Instructions {
			ins = append(ins, SectionInstruction{Id: i.UUID, Instruction: i.Instruction})
		}
		for _, i := range d.Ingredients {
			item := SectionIngredient{
				Id:        i.UUID,
				Grams:     i.Grams.Float64,
				Amount:    i.Amount.Ptr(),
				Unit:      i.Unit.String,
				Adjective: i.Adjective.String,
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
	recipe, err := a.CreateRecipe(ctx, &r)
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, recipe)
}
func (a *API) CreateRecipe(ctx context.Context, r *RecipeDetail) (*RecipeDetail, error) {
	id := r.Recipe.Id
	var err error
	if id == "" {
		id, err = a.DB().InsertRecipe(ctx, r.toDB())
		if err != nil {
			return nil, err
		}
		r.Recipe.Id = id
	}
	if err := a.DB().UpdateRecipe(ctx, r.toDB()); err != nil {
		return nil, err
	}
	r2, err := a.Manager.DB().GetRecipeByUUIDFull(ctx, id)
	if err != nil {
		return nil, err
	}
	return transformRecipeFull(r2), nil
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

func (a *API) IngredientUUIDByName(ctx context.Context, name, kind string) (string, error) {
	ing, err := a.DB().IngredientByName(ctx, name)
	if err != nil {
		return "", err
	}
	return ing.UUID, nil
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

func (a *API) fromDBPhoto(ctx context.Context, photos []db.Photo, getURLs bool) ([]GooglePhoto, []string, error) {
	ctx, span := global.Tracer("api").Start(ctx, "google.fromDBPhoto")
	defer span.End()
	items := []GooglePhoto{}
	var ids []string
	for _, p := range photos {
		items = append(items, GooglePhoto{Id: p.PhotoID, Created: p.Created})
		ids = append(ids, p.PhotoID)
	}

	if getURLs {
		urls, err := a.Manager.Google.GetBaseURLs(ctx, ids)
		if err != nil {
			return nil, nil, err
		}
		for x, item := range items {
			url, ok := urls[item.Id]
			if !ok {
				continue
			}
			items[x].BaseUrl = url
		}
	}
	return items, ids, nil
}
func (a *API) ListPhotos(c echo.Context, params ListPhotosParams) error {
	ctx := c.Request().Context()
	photos, err := a.Manager.DB().GetPhotos(ctx)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	items, _, err := a.fromDBPhoto(ctx, photos, true)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}

	resp := PaginatedPhotos{
		Photos: &items,
	}
	return c.JSON(http.StatusOK, resp)
}

func (a *API) ListMeals(c echo.Context, params ListMealsParams) error {
	ctx := c.Request().Context()

	items := []Meal{}

	meals, err := a.DB().GetAllMeals(ctx)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	var gphotoIDs []string
	for _, m := range meals {
		meal := Meal{Id: m.ID, Name: m.Name, AteAt: m.AteAt}

		photos, err := a.DB().GetPhotosForMeal(ctx, m.ID)
		if err != nil {
			return sendErr(c, http.StatusInternalServerError, err)
		}

		photos2, gIDs, err := a.fromDBPhoto(ctx, photos, false)
		if err != nil {
			return sendErr(c, http.StatusInternalServerError, err)
		}
		meal.Photos = photos2
		gphotoIDs = append(gphotoIDs, gIDs...)

		items = append(items, meal)
	}
	urls, err := a.Manager.Google.GetBaseURLs(ctx, gphotoIDs)
	if err != nil {
		return err
	}
	for x, item := range items {
		for y, photo := range item.Photos {
			url, ok := urls[photo.Id]
			if !ok {
				continue
			}
			items[x].Photos[y].BaseUrl = url
		}
	}

	resp := PaginatedMeals{
		Meals: &items,
	}
	return c.JSON(http.StatusOK, resp)
}
func (a *API) AuthLogin(c echo.Context, params AuthLoginParams) error {
	ctx := c.Request().Context()
	jwt, rawUser, err := a.Manager.ProcessGoogleAuth(ctx, params.Code)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}

	resp := AuthResp{
		Jwt:  jwt,
		User: rawUser,
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
