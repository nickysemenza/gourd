//go:generate oapi-codegen --package api --generate server,spec -o api-server.gen.go openapi.yaml
//go:generate oapi-codegen --package api --generate types -o api-types.gen.go openapi.yaml
// todo go:generate oapi-codegen --package api --generate client -o api-client.gen.go openapi.yaml

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	sq "github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/manager"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/guregu/null.v3/zero"
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
	return Ingredient{Id: dbr.Id, Name: dbr.Name}
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
			si := db.SectionIngredient{
				Grams:     zero.FloatFrom(i.Grams),
				Amount:    zero.FloatFromPtr(i.Amount),
				Unit:      zero.StringFromPtr(i.Unit),
				Adjective: zero.StringFromPtr(i.Adjective),
				Original:  zero.StringFromPtr(i.Original),
				Optional:  zero.BoolFromPtr(i.Optional),
			}
			switch i.Kind {
			case "recipe":
				if i.Recipe == nil {
					continue
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
			case "ingredient":
				if i.Ingredient == nil || (i.Ingredient.Name == "" && i.Ingredient.Id == "") {
					continue
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
				continue
			default:
				return nil, fmt.Errorf("unknown kind: %s", i.Kind)

			}

			dbs.Ingredients = append(dbs.Ingredients, si)
		}

		dbr.Sections = append(dbr.Sections, dbs)
	}

	return &dbr, nil

}
func transformRecipeSections(dbs []db.Section) ([]RecipeSection, error) {
	s := []RecipeSection{}
	for _, d := range dbs {
		ing := []SectionIngredient{}
		ins := []SectionInstruction{}

		for _, i := range d.Instructions {
			ins = append(ins, SectionInstruction{Id: i.Id, Instruction: i.Instruction})
		}
		for _, i := range d.Ingredients {
			item := SectionIngredient{
				Id:        i.Id,
				Grams:     i.Grams.Float64,
				Amount:    i.Amount.Ptr(),
				Unit:      i.Unit.Ptr(),
				Adjective: i.Adjective.Ptr(),
				Original:  i.Original.Ptr(),
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

// List all recipes
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

func (a *API) IngredientIdByName(ctx context.Context, name, kind string) (string, error) {
	ing, err := a.DB().IngredientByName(ctx, name)
	if err != nil {
		return "", err
	}
	return ing.Id, nil
}
func (a *API) ListIngredients(c echo.Context, params ListIngredientsParams) error {
	ctx := c.Request().Context()

	ctx, span := a.tracer.Start(ctx, "ListIngredients")
	defer span.End()

	items := []IngredientDetail{}

	paginationParams, listMeta := parsePagination(params.Offset, params.Limit)
	ing, count, err := a.Manager.DB().GetIngredients(ctx, "", paginationParams...)
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}
	var ingredientIds []string
	for _, i := range ing {
		ingredientIds = append(ingredientIds, i.Id)
	}

	sameAs, _, err := a.Manager.DB().GetIngrientsSameAs(ctx, ingredientIds...)
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}
	for _, i := range sameAs {
		ingredientIds = append(ingredientIds, i.Id)
	}
	linkedRecipes, err := a.Manager.DB().GetRecipeDetailsWithIngredient(ctx, ingredientIds...)
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}

	for _, i := range ing {
		// assemble
		items = append(items, makeDetail(i, sameAs, linkedRecipes))
	}
	listMeta.TotalCount = int(count)

	resp := PaginatedIngredients{
		Ingredients: &items,
		Meta:        listMeta,
	}
	return c.JSON(http.StatusOK, resp)
}

func makeDetail(i db.Ingredient, sameAs db.Ingredients, linkedRecipes db.RecipeDetails) IngredientDetail {
	// find linked ingredients
	same := []IngredientDetail{}
	for _, x := range sameAs.BySameAs()[i.Id] {
		same = append(same, makeDetail(x, sameAs, linkedRecipes))
	}

	// find linked recipes
	recipes := []RecipeDetail{}
	for _, x := range linkedRecipes.ByIngredientId()[i.Id] {
		recipes = append(recipes, transformRecipe(x))
	}
	return IngredientDetail{
		Ingredient: transformIngredient(i),
		Children:   &same,
		Recipes:    &recipes,
	}
}

func (a *API) fromDBPhoto(ctx context.Context, photos []db.Photo, getURLs bool) ([]GooglePhoto, []string, error) {
	ctx, span := a.tracer.Start(ctx, "fromDBPhoto")
	defer span.End()
	items := []GooglePhoto{}
	var ids []string
	for _, p := range photos {
		gp := GooglePhoto{Id: p.PhotoID, Created: p.Created}
		if p.BlurHash.Valid {
			s := p.BlurHash.String
			gp.BlurHash = &s
		}
		items = append(items, gp)
		ids = append(ids, p.PhotoID)
	}

	if getURLs {
		results, err := a.Manager.Photos.GetMediaItems(ctx, ids)
		if err != nil {
			return nil, nil, err
		}
		for x, item := range items {
			val, ok := results[item.Id]
			if !ok {
				continue
			}
			items[x].BaseUrl = val.BaseUrl
			items[x].Width = val.MediaMetadata.Width
			items[x].Height = val.MediaMetadata.Height
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
	urls, err := a.Manager.Photos.GetMediaItems(ctx, gphotoIDs)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	for x, item := range items {
		for y, photo := range item.Photos {
			val, ok := urls[photo.Id]
			if !ok {
				continue
			}
			items[x].Photos[y].BaseUrl = val.BaseUrl
			items[x].Photos[y].Width = val.MediaMetadata.Width
			items[x].Photos[y].Height = val.MediaMetadata.Height
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
	l.PageCount = int(math.Ceil(float64(c) / float64(l.Limit)))
}

func sendErr(ctx echo.Context, code int, err error) error {
	trace.SpanFromContext(ctx.Request().Context()).AddEvent(fmt.Sprintf("error: %s", err))
	return ctx.JSON(code, Error{
		Message: err.Error(),
	})
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
	ingredients, ingredientsCount, err := a.Manager.DB().GetIngredients(ctx, string(params.Name))
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

func (a *API) ConvertIngredientToRecipe(c echo.Context, ingredientId string) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "ConvertIngredientToRecipe")
	defer span.End()

	detail, err := a.DB().IngredientToRecipe(ctx, ingredientId)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, transformRecipe(*detail))
}

func (a *API) MergeIngredients(c echo.Context, ingredientId string) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "MergeIngredients")
	defer span.End()

	var r MergeIngredientsJSONRequestBody
	if err := c.Bind(&r); err != nil {
		err = fmt.Errorf("invalid format for input: %w", err)
		return sendErr(c, http.StatusBadRequest, err)
	}

	err := a.DB().MergeIngredients(ctx, ingredientId, r.IngredientIds)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}

	ing, err := a.DB().GetIngredientById(ctx, ingredientId)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, transformIngredient(*ing))
}
