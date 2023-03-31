package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/internal/common"
	"github.com/nickysemenza/gourd/internal/db"
	"github.com/nickysemenza/gourd/internal/db/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slices"
	"gopkg.in/guregu/null.v4/zero"
)

// Items all recipes
// (GET /recipes)
func (a *API) ListRecipes(c echo.Context, params ListRecipesParams) error {
	ctx := c.Request().Context()
	ctx, span := a.tracer.Start(ctx, "ListRecipes")
	defer span.End()

	paginationParams, listMeta := parsePagination(params.Offset, params.Limit)
	limit, offset := db.GetPagination(paginationParams...)

	count, err := models.Recipes().Count(ctx, a.db.DB())
	if err != nil {
		return handleErr(c, err)
	}

	listMeta.setTotalCount(uint64(count))

	items, err := a.RecipeListV2(ctx, limit, offset)
	if err != nil {
		return handleErr(c, err)
	}

	resp := PaginatedRecipeWrappers{
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
		return handleErr(c, err)
	}
	recipe, err := a.CreateRecipe(ctx, &r)
	if err != nil {
		return handleErr(c, err)
	}

	return c.JSON(http.StatusCreated, recipe)
}
func (a *API) CreateRecipeDetails(ctx context.Context, recipes ...RecipeDetailInput) (detailIDs []string, err error) {
	for _, r := range recipes {
		d, err := a.CreateRecipe(ctx, &RecipeWrapperInput{Detail: r})
		if err != nil {
			return nil, err
		}
		detailIDs = append(detailIDs, d.Detail.Id)
	}
	return
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

	w, err := a.transformRecipeFull(ctx, r2)
	if err != nil {
		return nil, err
	}
	a.indexRecipeDetails(ctx, append(w.OtherVersions, w.Detail)...)
	return w, nil
}

func (a *API) recipeById(ctx context.Context, recipeId string) (*RecipeWrapper, error) {
	r, err := a.DB().GetRecipeDetailByIdFull(ctx, recipeId)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, fmt.Errorf("could not find recipe with detail %s", recipeId)
	}
	full, err := a.transformRecipeFull(ctx, r)
	if err != nil {
		return nil, err
	}
	p, err := a.imagesFromRecipeDetailId(ctx, recipeId)
	if err != nil {
		return nil, err
	}
	full.LinkedPhotos = &p

	if slices.Contains(r.Tags, "notion") {
		// add notion links as sources as well
		notionRecipes, err := models.NotionRecipes(qm.Where("recipe_id = ?", r.RecipeId)).All(ctx, a.db.DB())
		if err != nil {
			return nil, err
		}
		notionSOurces := []RecipeSource{}
		for _, n := range notionRecipes {
			url := fmt.Sprintf("https://notion.so/%s", strings.ReplaceAll(n.PageID, "-", ""))
			notionSOurces = append(notionSOurces, RecipeSource{Url: &url})
		}
		full.Detail.Sources = append(full.Detail.Sources, notionSOurces...)

	}
	return full, nil
}

// Info for a specific recipe
// (GET /recipes/{recipeId})
func (a *API) GetRecipeById(c echo.Context, recipeId string) error {
	ctx := c.Request().Context()

	apiR, err := a.recipeById(ctx, recipeId)
	if err != nil {
		return handleErr(c, err)
	}

	return c.JSON(http.StatusOK, apiR)
}

func (a *API) GetRecipesByIds(c echo.Context, params GetRecipesByIdsParams) error {

	ctx := c.Request().Context()
	list := []RecipeWrapper{}

	for _, recipeId := range params.RecipeId {
		apiR, err := a.recipeById(ctx, recipeId)
		if err != nil {
			return handleErr(c, err)
		}
		list = append(list, *apiR)
	}
	result := PaginatedRecipeWrappers{Recipes: &list}
	return c.JSON(http.StatusOK, result)
}

func (a *API) transformRecipe(ctx context.Context, dbr db.RecipeDetail, includeOtherVersions bool) (*RecipeWrapper, error) {
	ctx, span := a.tracer.Start(ctx, "transformRecipe")
	defer span.End()
	sections, err := a.transformRecipeSections(ctx, dbr.Sections)
	if err != nil {
		return nil, err
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
		CreatedAt:       dbr.CreatedAt,
		Tags:            []string{},
		Sources:         []RecipeSource{},
	}
	if dbr.Source.String != "" {
		if err := json.Unmarshal([]byte(dbr.Source.String), &rd.Sources); err != nil {
			return nil, err
		}
	}
	w := RecipeWrapper{
		Id:     dbr.RecipeId,
		Detail: rd,
	}
	if includeOtherVersions {
		recipes, err := a.DB().GetRecipeDetailWhere(ctx, sq.And{sq.Eq{"recipe_id": dbr.RecipeId}, sq.NotEq{"id": dbr.Id}})
		if err != nil {
			return nil, err
		}
		details, err := a.transformRecipes(ctx, recipes, false)
		if err != nil {
			return nil, err
		}
		// w.OtherVersions = &details
		for _, d := range details {
			w.OtherVersions = append(w.OtherVersions, d.Detail)
		}
	}

	return &w, nil
}
func (a *API) transformRecipes(ctx context.Context, dbr db.RecipeDetails, includeOtherVersions bool) ([]RecipeWrapper, error) {
	ctx, span := a.tracer.Start(ctx, "transformRecipes")
	defer span.End()
	r := make([]RecipeWrapper, len(dbr))
	for x, d := range dbr {
		each, err := a.transformRecipe(ctx, d, includeOtherVersions)
		if err != nil {
			return nil, err
		}
		r[x] = *each
	}
	return r, nil

}
func (a *API) transformRecipeFull(ctx context.Context, dbr *db.RecipeDetail) (*RecipeWrapper, error) {
	d, err := a.transformRecipe(ctx, *dbr, true)
	if err != nil {
		return nil, err
	}

	return d, nil
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

		eq := sq.Eq{"recipe_id": iOrRiD}
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
		Name:      r.Detail.Name,
		CreatedAt: zero.TimeFromPtr(r.Detail.Date).Time,
	}
	if r.Detail.Sources != nil {
		source, err := json.Marshal(r.Detail.Sources)
		if err != nil {
			return nil, err
		}
		dbr.Source = zero.StringFrom(string(source))
	}
	dbr.Tags = r.Detail.Tags

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

func (a *API) transformRecipeSections(ctx context.Context, dbs []db.Section) ([]RecipeSection, error) {
	ctx, span := a.tracer.Start(ctx, "transformRecipeSections")
	defer span.End()
	s := []RecipeSection{}
	for _, d := range dbs {
		ing := []SectionIngredient{}
		ingSubs := map[string][]SectionIngredient{}
		ins := []SectionInstruction{}

		for _, i := range d.Instructions {
			ins = append(ins, SectionInstruction{Id: i.Id, Instruction: i.Instruction})
		}
		for _, i := range d.Ingredients {
			si := SectionIngredient{
				Id:        i.Id,
				Adjective: i.Adjective.Ptr(),
				Original:  i.Original.Ptr(),
				Optional:  i.Optional.Ptr(),
				Amounts:   []Amount{},
			}
			for _, amt := range i.Amounts {
				si.Amounts = append(si.Amounts, Amount{
					Unit:   amt.Unit,
					Value:  amt.Value,
					Source: zero.StringFrom("db").Ptr(),
				})

			}

			if i.RawRecipe != nil {
				si.Kind = "recipe"
				r, err := a.transformRecipe(ctx, *i.RawRecipe, false)
				if err != nil {
					return nil, err
				}
				si.Recipe = &r.Detail
			} else {
				si.Kind = "ingredient"
				foo, err := a.addDetailsToIngredients(ctx, []db.Ingredient{*i.RawIngredient})
				if err != nil {
					return nil, err
				}
				// i := transformIngredient(*i.RawIngredient)
				si.Ingredient = &foo[0]

				targets := []UnitConversionRequestTarget{UnitConversionRequestTargetCalories, UnitConversionRequestTargetMoney, UnitConversionRequestTargetVolume}
				if !hasGrams(si.Amounts) {
					targets = append(targets, UnitConversionRequestTargetWeight)
				}
				for _, t := range targets {
					err = a.enhance(ctx, t, &si)
					if err != nil {
						return nil, err
					}
				}
			}
			if i.SubsFor.Valid {
				ingSubs[i.SubsFor.String] = append(ingSubs[i.SubsFor.String], si)
			} else {
				ing = append(ing, si)
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
	if err := a.R.Send(ctx, "convert", req, &res); err != nil {
		return fmt.Errorf("enhance: %w", err)
	}
	if res.Value != 0 {
		res.Value = math.Round(res.Value*100) / 100
		res.Source = zero.StringFrom("calculated").Ptr()
		item.Amounts = append(item.Amounts, res)
	}
	return nil
}
