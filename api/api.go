//go:generate ../bin/oapi-codegen --package api --generate server,spec -o api-server.gen.go openapi.yaml
//go:generate ../bin/oapi-codegen --package api --generate types,skip-prune -o api-types.gen.go openapi.yaml
// todo go:generate ../bin/oapi-codegen --package api --generate client -o api-client.gen.go openapi.yaml

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	sq "github.com/Masterminds/squirrel"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/auth"
	"github.com/nickysemenza/gourd/common"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/google"
	"github.com/nickysemenza/gourd/image"
	"github.com/nickysemenza/gourd/notion"
	"github.com/nickysemenza/gourd/photos"
	"github.com/nickysemenza/gourd/rs_client"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/guregu/null.v4/zero"
)

type API struct {
	db         *db.Client
	Google     *google.Client
	Photos     *photos.Photos
	Auth       *auth.Auth
	R          *rs_client.Client
	Notion     *notion.Client
	ImageStore image.Store
	tracer     trace.Tracer
}

func New(db *db.Client, g *google.Client, auth *auth.Auth, r *rs_client.Client, notion *notion.Client, imageStore image.Store) *API {
	return &API{
		db:         db,
		Google:     g,
		Auth:       auth,
		Photos:     photos.New(db, g, imageStore),
		R:          r,
		Notion:     notion,
		ImageStore: imageStore,
		tracer:     otel.Tracer("api"),
	}
}

func (a *API) transformRecipe(ctx context.Context, dbr db.RecipeDetail, includeOtherVersions bool) RecipeDetail {
	sections, err := a.transformRecipeSections(ctx, dbr.Sections)
	if err != nil {
		panic(err)
	}
	rd := RecipeDetail{
		Id:              dbr.Id,
		Name:            dbr.Name,
		Version:         dbr.Version,
		IsLatestVersion: dbr.LatestVersion,
		Quantity:        dbr.Quantity.Int64,
		Servings:        &dbr.Servings.Int64,
		Unit:            dbr.Unit.String,
		Sections:        sections,
		// CreatedAt:       &dbr.CreatedAt,
	}
	if dbr.Source.Valid {
		if err := json.Unmarshal([]byte(dbr.Source.String), &rd.Sources); err != nil {
			panic(err)
		}
	}
	if includeOtherVersions {
		recipes, err := a.DB().GetRecipeDetailWhere(ctx, sq.And{sq.Eq{"recipe": dbr.RecipeId}, sq.NotEq{"id": dbr.Id}})
		if err != nil {
			panic(err)
		}
		details := a.transformRecipes(ctx, recipes, false)
		rd.OtherVersions = &details
	}
	return rd
}
func (a *API) transformRecipes(ctx context.Context, dbr db.RecipeDetails, includeOtherVersions bool) []RecipeDetail {
	r := make([]RecipeDetail, len(dbr))
	for x, d := range dbr {
		r[x] = a.transformRecipe(ctx, d, includeOtherVersions)
	}
	return r

}
func (a *API) transformRecipeFull(ctx context.Context, dbr *db.RecipeDetail) *RecipeWrapper {
	return &RecipeWrapper{
		Id:     dbr.RecipeId,
		Detail: a.transformRecipe(ctx, *dbr, true),
	}
}
func transformIngredient(dbr db.Ingredient) Ingredient {
	return Ingredient{
		Id:     dbr.Id,
		Name:   dbr.Name,
		FdcId:  dbr.FdcID.Ptr(),
		Parent: dbr.Parent.Ptr(),
	}
}
func (a *API) sectionIngredientTODB(ctx context.Context, i SectionIngredientInput) (*db.SectionIngredient, error) {
	si := db.SectionIngredient{
		Id:        common.ID("si"),
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

	if i.Name == nil && i.TargetId == nil {
		return nil, nil
	}

	iOrRiD := zero.StringFromPtr(i.TargetId).ValueOrZero()
	name := zero.StringFromPtr(i.Name).ValueOrZero()

	switch i.Kind {
	case IngredientKindRecipe:

		eq := sq.Eq{"recipe": iOrRiD}
		if iOrRiD != "" {
			eq = sq.Eq{"name": name}
		}
		rs, err := a.DB().GetRecipeDetailWhere(ctx, eq)
		if err != nil {
			return nil, err
		}

		r := rs.First()
		if r != nil {
			iOrRiD = r.RecipeId
		} else {
			r, err := a.DB().InsertRecipe(ctx, &db.RecipeDetail{Name: name})
			if err != nil {
				return nil, err
			}
			iOrRiD = r.RecipeId
		}
		si.RecipeId = zero.StringFrom(iOrRiD)
	case IngredientKindIngredient:

		// missing id, need to find/create
		if iOrRiD == "" {
			ing, err := a.DB().IngredientByName(ctx, name)
			if err != nil {
				return nil, err
			}
			iOrRiD = ing.Id
		}

		si.IngredientId = zero.StringFrom(iOrRiD)
	case "":
		// empty table row, drop it
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown kind: %s", i.Kind)

	}
	return &si, nil
}
func (a *API) recipeWrappertoDB(ctx context.Context, r *RecipeWrapperInput) (*db.RecipeDetail, error) {
	dbr := db.RecipeDetail{
		// Id:   r.Detail.Id,
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
			// i. = common.ID("s_ing")
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
				si2.SubsFor = zero.StringFromPtr(i.TargetId)
				dbs.Ingredients = append(dbs.Ingredients, si2)
			}
		}

		dbr.Sections = append(dbr.Sections, dbs)
	}

	return &dbr, nil

}
func hasGrams(amounts []Amount) bool {
	for _, amt := range amounts {
		if amt.Unit == "grams" || amt.Unit == "g" || amt.Unit == "gram" {
			return true
		}
	}
	return false
}
func (a *API) transformRecipeSections(ctx context.Context, dbs []db.Section) ([]RecipeSection, error) {
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
			for _, amt := range i.Amounts {
				item.Amounts = append(item.Amounts, Amount{
					Unit:   amt.Unit,
					Value:  amt.Value,
					Source: zero.StringFrom("db").Ptr(),
				})

			}

			if i.RawRecipe != nil {
				item.Kind = "recipe"
				r := a.transformRecipe(ctx, *i.RawRecipe, false)
				item.Recipe = &r
			} else {
				item.Kind = "ingredient"
				foo, err := a.addDetailsToIngredients(ctx, []db.Ingredient{*i.RawIngredient})
				if err != nil {
					return nil, err
				}
				// i := transformIngredient(*i.RawIngredient)
				item.Ingredient = &foo[0]

				targets := []UnitConversionRequestTarget{UnitConversionRequestTargetCalories, UnitConversionRequestTargetMoney, UnitConversionRequestTargetVolume}
				if !hasGrams(item.Amounts) {
					targets = append(targets, UnitConversionRequestTargetWeight)
				}
				for _, t := range targets {
					err = a.enhance(ctx, t, &item)
					if err != nil {
						return nil, err
					}
				}
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
func (a *API) enhance(ctx context.Context, with UnitConversionRequestTarget, item *SectionIngredient) error {
	req := UnitConversionRequest{
		Input:        item.Amounts,
		Target:       &with,
		UnitMappings: item.Ingredient.UnitMappings,
	}
	var res Amount
	err := a.R.Convert(
		ctx, req, &res,
	)
	if err != nil {
		return fmt.Errorf("enhance: %w", err)
	}
	if res.Value != 0 {
		res.Value = math.Round(res.Value*100) / 100
		res.Source = zero.StringFrom("calculated").Ptr()
		item.Amounts = append(item.Amounts, res)
	}
	return nil
}

// Items all recipes
// (GET /recipes)
func (a *API) ListRecipes(c echo.Context, params ListRecipesParams) error {
	ctx := c.Request().Context()
	ctx, span := a.tracer.Start(ctx, "ListRecipes")
	defer span.End()

	paginationParams, listMeta := parsePagination(params.Offset, params.Limit)
	recipes, count, err := a.DB().GetRecipes(ctx, "", paginationParams...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}

	var recipeIDs []string
	for _, r := range recipes {
		recipeIDs = append(recipeIDs, r.Id)
	}

	details, err := a.DB().GetRecipeDetailWhere(ctx, sq.Eq{"recipe": recipeIDs})
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}
	byId := details.ByRecipeId()
	items := []Recipe{}

	for _, r := range recipeIDs {
		meals, err := a.getLinkedMeals(ctx, r)
		if err != nil {
			return sendErr(c, http.StatusBadRequest, err)
		}
		photosDB, err := a.DB().GetPhotosWithRecipe(ctx, r)
		if err != nil {
			return sendErr(c, http.StatusBadRequest, err)
		}
		photos := []Photo{}
		for _, p := range photosDB {
			url := a.ImageStore.GetImageURL(ctx, p.ID)
			photos = append(photos, Photo{
				Id:      p.ID,
				Source:  PhotoSource(p.Source),
				BaseUrl: url,
			})
		}

		items = append(items, Recipe{
			Id:           r,
			LinkedMeals:  &meals,
			LinkedPhotos: &photos,
			Versions:     a.transformRecipes(ctx, byId[r], true),
		})
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
	var r RecipeWrapperInput
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
func (a *API) CreateRecipeDetails(ctx context.Context, recipes ...RecipeDetailInput) error {
	for _, r := range recipes {
		_, err := a.CreateRecipe(ctx, &RecipeWrapperInput{Detail: r})
		if err != nil {
			return err
		}
	}
	return nil
}
func (a *API) CreateRecipe(ctx context.Context, r *RecipeWrapperInput) (*RecipeWrapper, error) {
	ctx, span := a.tracer.Start(ctx, "CreateRecipe")
	defer span.End()

	span.AddEvent("got wrapper", trace.WithAttributes(attribute.String("recipe", spew.Sdump(r))))
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

	return a.transformRecipeFull(ctx, r2), nil
}

// Info for a specific recipe
// (GET /recipes/{recipeId})
func (a *API) GetRecipeById(c echo.Context, recipeId string) error {
	ctx := c.Request().Context()
	r, err := a.DB().GetRecipeDetailByIdFull(ctx, recipeId)
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}
	if r == nil {
		return sendErr(c, http.StatusNotFound, fmt.Errorf("could not find recipe with detail %s", recipeId))
	}
	apiR := a.transformRecipeFull(ctx, r)

	return c.JSON(http.StatusOK, apiR)
}

func (a *API) GetRecipesByIds(c echo.Context, params GetRecipesByIdsParams) error {

	ctx := c.Request().Context()
	list := []RecipeWrapper{}

	for _, recipeId := range params.RecipeId {
		r, err := a.DB().GetRecipeDetailByIdFull(ctx, recipeId)
		if err != nil {
			return sendErr(c, http.StatusBadRequest, err)
		}
		if r == nil {
			return sendErr(c, http.StatusNotFound, fmt.Errorf("could not find recipe with detail %s", recipeId))
		}
		apiR := a.transformRecipeFull(ctx, r)
		list = append(list, *apiR)
	}
	result := PaginatedRecipeWrappers{Recipes: &list}
	return c.JSON(http.StatusOK, result)
}

func (a *API) AuthLogin(c echo.Context, params AuthLoginParams) error {
	ctx := c.Request().Context()
	jwt, rawUser, err := a.ProcessGoogleAuth(ctx, params.Code)
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

	dbAlbums, err := a.DB().GetAlbums(ctx)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}

	albums, err := a.Photos.GetAvailableAlbums(ctx)
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

	recipes, recipesCount, err := a.DB().GetRecipesDetails(ctx, string(params.Name))
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	ingredients, ingredientsCount, err := a.DB().GetIngredients(ctx, string(params.Name), nil)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}

	listMeta.setTotalCount(recipesCount + ingredientsCount)

	var resRecipes []RecipeWrapper
	var resIngredients []Ingredient

	for _, x := range recipes {
		r := a.transformRecipe(ctx, x, true)
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

func (a *API) NotionTest(c echo.Context) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "Notion")
	defer span.End()

	res, err := a.Notion.Dump(ctx)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, res)
}
