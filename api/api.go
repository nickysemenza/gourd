//go:generate oapi-codegen --package api --generate server,spec -o api-server.gen.go openapi.yaml
//go:generate oapi-codegen --package api --generate types -o api-types.gen.go openapi.yaml
// todo go:generate oapi-codegen --package api --generate client -o api-client.gen.go openapi.yaml

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	sq "github.com/Masterminds/squirrel"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/common"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/manager"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/guregu/null.v4/zero"
)

type API struct {
	*manager.Manager
	tracer trace.Tracer
}

func NewAPI(m *manager.Manager) *API {
	return &API{
		Manager: m,
		tracer:  otel.Tracer("api"),
	}
}

func transformRecipe(dbr db.RecipeDetail) RecipeDetail {
	sections, err := transformRecipeSections(dbr.Sections)
	if err != nil {
		panic(err)
	}
	rd := RecipeDetail{
		Id:              dbr.Id,
		Name:            dbr.Name,
		Version:         &dbr.Version,
		IsLatestVersion: &dbr.LatestVersion,
		Quantity:        dbr.Quantity.Int64,
		Servings:        &dbr.Servings.Int64,
		Unit:            dbr.Unit.String,
		Sections:        sections,
	}
	if dbr.Source.Valid {
		if err := json.Unmarshal([]byte(dbr.Source.String), &rd.Sources); err != nil {
			panic(err)
		}
	}
	return rd
}
func transformRecipes(dbr db.RecipeDetails) []RecipeDetail {
	r := make([]RecipeDetail, len(dbr))
	for x, d := range dbr {
		r[x] = transformRecipe(d)
	}
	return r

}
func transformRecipeFull(dbr *db.RecipeDetail) *RecipeWrapper {
	return &RecipeWrapper{
		Id:     dbr.RecipeId,
		Detail: transformRecipe(*dbr),
	}
}
func transformIngredient(dbr db.Ingredient) Ingredient {
	return Ingredient{
		Id:     dbr.Id,
		Name:   dbr.Name,
		FdcId:  dbr.FdcID.Ptr(),
		SameAs: dbr.SameAs.Ptr(),
	}
}
func (a *API) sectionIngredientTODB(ctx context.Context, i SectionIngredient) (*db.SectionIngredient, error) {
	si := db.SectionIngredient{
		Id:        i.Id,
		Adjective: zero.StringFromPtr(i.Adjective),
		Original:  zero.StringFromPtr(i.Original),
		Optional:  zero.BoolFromPtr(i.Optional),
		Amounts:   []db.Amount{},
	}
	for _, amt := range i.Amounts {
		si.Amounts = append(si.Amounts, db.Amount{
			Unit:  amt.Unit,
			Value: amt.Value,
		})
	}
	switch i.Kind {
	case IngredientKind_recipe:
		if i.Recipe == nil {
			return nil, nil
		}
		var eq sq.Eq
		id := i.Recipe.Id
		if id == "" {
			eq = sq.Eq{"name": i.Recipe.Name}
		} else {
			eq = sq.Eq{"recipe": i.Recipe.Id}
		}
		rs, err := a.DB().GetRecipeDetailWhere(ctx, eq)
		if err != nil {
			return nil, err
		}

		r := rs.First()
		if r != nil {
			id = r.RecipeId
		} else {
			r, err := a.DB().InsertRecipe(ctx, &db.RecipeDetail{Name: i.Recipe.Name})
			if err != nil {
				return nil, err
			}
			id = r.RecipeId
		}
		si.RecipeId = zero.StringFrom(id)
	case IngredientKind_ingredient:
		if i.Ingredient == nil || (i.Ingredient.Name == "" && i.Ingredient.Id == "") {
			return nil, nil
		}
		id := i.Ingredient.Id

		// missing id, need to find/create
		if id == "" {
			ing, err := a.DB().IngredientByName(ctx, i.Ingredient.Name)
			if err != nil {
				return nil, err
			}
			id = ing.Id
		}

		si.IngredientId = zero.StringFrom(id)
	case "":
		// empty table row, drop it
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown kind: %s", i.Kind)

	}
	return &si, nil
}
func (a *API) recipeWrappertoDB(ctx context.Context, r *RecipeWrapper) (*db.RecipeDetail, error) {
	dbr := db.RecipeDetail{
		Id:   r.Detail.Id,
		Name: r.Detail.Name,
	}
	source, err := json.Marshal(r.Detail.Sources)
	if err != nil {
		return nil, err
	}
	dbr.Source = zero.StringFrom(string(source))

	for _, s := range r.Detail.Sections {
		if len(s.Ingredients) == 0 && len(s.Instructions) == 0 {
			// empty section, drop it
			continue
		}
		data, err := json.Marshal(s.Duration)
		if err != nil {
			return nil, err
		}

		dbs := db.Section{
			TimeRange: string(data),
		}
		for _, i := range s.Instructions {
			if i.Instruction == "" {
				continue
			}
			dbs.Instructions = append(dbs.Instructions, db.SectionInstruction{
				Instruction: i.Instruction,
			})
		}
		for _, i := range s.Ingredients {
			i.Id = common.ID("ing")
			si, err := a.sectionIngredientTODB(ctx, i)
			if err != nil {
				return nil, err
			}
			if si == nil {
				continue
			}
			dbs.Ingredients = append(dbs.Ingredients, *si)
			if i.Substitutes == nil {
				continue
			}
			for _, sub := range *i.Substitutes {
				si, err := a.sectionIngredientTODB(ctx, sub)
				if err != nil {
					return nil, err
				}
				if si == nil {
					continue
				}
				si2 := *si
				si2.SubsFor = zero.StringFrom(i.Id)
				dbs.Ingredients = append(dbs.Ingredients, si2)
				spew.Dump(si2)
			}
		}

		dbr.Sections = append(dbr.Sections, dbs)
	}

	return &dbr, nil

}
func transformRecipeSections(dbs []db.Section) ([]RecipeSection, error) {
	s := []RecipeSection{}
	for _, d := range dbs {
		ing := []SectionIngredient{}
		ingSubs := map[string][]SectionIngredient{}
		ins := []SectionInstruction{}

		for _, i := range d.Instructions {
			ins = append(ins, SectionInstruction{Id: i.Id, Instruction: i.Instruction})
		}
		for _, i := range d.Ingredients {
			item := SectionIngredient{
				Id:        i.Id,
				Adjective: i.Adjective.Ptr(),
				Original:  i.Original.Ptr(),
				Optional:  i.Optional.Ptr(),
				Amounts:   []Amount{},
			}
			hasGrams := false
			for _, amt := range i.Amounts {
				item.Amounts = append(item.Amounts, Amount{
					Unit:   amt.Unit,
					Value:  amt.Value,
					Source: zero.StringFrom("db").Ptr(),
				})

				if amt.Unit == "grams" || amt.Unit == "g" || amt.Unit == "gram" {
					hasGrams = true
				}
			}
			if !hasGrams {
				item.Amounts = append(item.Amounts, Amount{
					Unit:   "gram",
					Value:  1.1,
					Source: zero.StringFrom("calculated-todo").Ptr(),
				})
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
			if i.SubsFor.Valid {
				ingSubs[i.SubsFor.String] = append(ingSubs[i.SubsFor.String], item)
			} else {
				ing = append(ing, item)
			}
		}
		for x, i := range ing {
			if subs, ok := ingSubs[i.Id]; ok {
				ing[x].Substitutes = &subs
			}
		}
		rs := RecipeSection{
			Id:           d.Id,
			Ingredients:  ing,
			Instructions: ins,
		}
		err := json.Unmarshal([]byte(d.TimeRange), &rs.Duration)
		if err != nil {
			return nil, err
		}
		s = append(s, rs)
	}
	return s, nil
}

// Items all recipes
// (GET /recipes)
func (a *API) ListRecipes(c echo.Context, params ListRecipesParams) error {
	ctx := c.Request().Context()
	ctx, span := a.tracer.Start(ctx, "ListRecipes")
	defer span.End()

	paginationParams, listMeta := parsePagination(params.Offset, params.Limit)
	recipes, count, err := a.Manager.DB().GetRecipes(ctx, "", paginationParams...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}

	var recipeIDs []string
	for _, r := range recipes {
		recipeIDs = append(recipeIDs, r.Id)
	}

	details, err := a.Manager.DB().GetRecipeDetailWhere(ctx, sq.Eq{"recipe": recipeIDs})
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}
	byId := details.ByRecipeId()
	items := []Recipe{}
	for _, r := range recipeIDs {
		items = append(items, Recipe{Id: r, Versions: transformRecipes(byId[r])})
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
	var r RecipeWrapper
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
func (a *API) CreateRecipeDetails(ctx context.Context, recipes ...RecipeDetail) error {
	for _, r := range recipes {
		_, err := a.CreateRecipe(ctx, &RecipeWrapper{Detail: r})
		if err != nil {
			return err
		}
	}
	return nil
}
func (a *API) CreateRecipe(ctx context.Context, r *RecipeWrapper) (*RecipeWrapper, error) {
	ctx, span := a.tracer.Start(ctx, "CreateRecipe")
	defer span.End()

	span.AddEvent("got wrapper", trace.WithAttributes(attribute.String("id", r.Id), attribute.String("recipe", spew.Sdump(r))))
	dbVersion, err := a.recipeWrappertoDB(ctx, r)
	if err != nil {
		return nil, err
	}
	r2, err := a.DB().InsertRecipe(ctx, dbVersion)
	if err != nil {
		return nil, err
	}

	if r2 == nil {
		return nil, fmt.Errorf("failed to create recipe with name %s", r.Detail.Name)
	}

	return transformRecipeFull(r2), nil
}

// Info for a specific recipe
// (GET /recipes/{recipeId})
func (a *API) GetRecipeById(c echo.Context, recipeId string) error {
	ctx := c.Request().Context()
	r, err := a.Manager.DB().GetRecipeDetailByIdFull(ctx, recipeId)
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}
	if r == nil {
		return sendErr(c, http.StatusNotFound, fmt.Errorf("could not find recipe with detail %s", recipeId))
	}
	apiR := transformRecipeFull(r)

	return c.JSON(http.StatusOK, apiR)
}

func (a *API) GetRecipesByIds(c echo.Context, params GetRecipesByIdsParams) error {

	ctx := c.Request().Context()
	list := []RecipeWrapper{}

	for _, recipeId := range params.RecipeId {
		r, err := a.Manager.DB().GetRecipeDetailByIdFull(ctx, recipeId)
		if err != nil {
			return sendErr(c, http.StatusBadRequest, err)
		}
		if r == nil {
			return sendErr(c, http.StatusNotFound, fmt.Errorf("could not find recipe with detail %s", recipeId))
		}
		apiR := transformRecipeFull(r)
		list = append(list, *apiR)
	}
	result := PaginatedRecipeWrappers{Recipes: &list}
	return c.JSON(http.StatusOK, result)
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

type Test struct {
	Albums []GooglePhotosAlbum `json:"albums,omitempty"`
}

func (a *API) ListAllAlbums(c echo.Context) error {
	ctx := c.Request().Context()

	var resp Test

	dbAlbums, err := a.Manager.DB().GetAlbums(ctx)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}

	albums, err := a.Manager.Photos.GetAvailableAlbums(ctx)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}

	for _, a := range albums {
		gpa := GooglePhotosAlbum{
			Id:         a.Id,
			ProductUrl: a.ProductUrl,
			Title:      a.Title,
		}

		for _, dbA := range dbAlbums {
			if dbA.ID == gpa.Id {
				gpa.Usecase = dbA.Usecase
			}
		}

		resp.Albums = append(resp.Albums, gpa)
	}

	return c.JSON(http.StatusOK, resp)
}

func (a *API) Search(c echo.Context, params SearchParams) error {
	ctx := c.Request().Context()

	ctx, span := a.tracer.Start(ctx, "Search")
	defer span.End()

	_, listMeta := parsePagination(params.Offset, params.Limit)

	recipes, recipesCount, err := a.Manager.DB().GetRecipesDetails(ctx, string(params.Name))
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	ingredients, ingredientsCount, err := a.Manager.DB().GetIngredients(ctx, string(params.Name), nil)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}

	listMeta.setTotalCount(recipesCount + ingredientsCount)

	var resRecipes []RecipeWrapper
	var resIngredients []Ingredient

	for _, x := range recipes {
		r := transformRecipe(x)
		resRecipes = append(resRecipes, RecipeWrapper{Detail: r, Id: x.RecipeId})
	}
	for _, x := range ingredients {
		i := transformIngredient(x)
		resIngredients = append(resIngredients, i)
	}

	return c.JSON(http.StatusOK, SearchResult{Recipes: &resRecipes, Ingredients: &resIngredients})
}

type Test2 struct {
	Items []RecipeDependency `json:"items,omitempty"`
}

func (a *API) RecipeDependencies(c echo.Context) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "RecipeDependencies")
	defer span.End()

	res := []RecipeDependency{}
	dbRows, err := a.DB().RecipeIngredientDependencies(ctx)
	for _, r := range dbRows {
		res = append(res, RecipeDependency{
			IngredientId:   r.IngredientId,
			IngredientKind: IngredientKind(r.IngredientKind),
			IngredientName: r.IngredientName,
			RecipeId:       r.RecipeId,
			RecipeName:     r.RecipeName,
		})
	}
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, Test2{res})
}
