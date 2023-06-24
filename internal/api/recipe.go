package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"

	"database/sql"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/internal/common"
	"github.com/nickysemenza/gourd/internal/db/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

	listMeta := parsePagination(params.Offset, params.Limit)

	items, count, err := a.RecipeListV2(ctx, listMeta)
	if err != nil {
		return handleErr(c, err)
	}

	listMeta.setTotalCount(uint64(count))

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
func (a *API) createRecipeDetails(ctx context.Context, recipes ...RecipeDetailInput) (detailIDs []string, err error) {
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

	tx, err := a.db.DB().BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	dbVersion, err := a.insertRecipeWrapper(ctx, tx, r)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	w, err := a.recipeById(ctx, dbVersion.ID)
	if err != nil {
		return nil, err
	}

	a.indexRecipeDetails(ctx, w)
	return w, nil
}

func (a *API) recipeById(ctx context.Context, recipeDetailId string) (*RecipeWrapper, error) {
	full, err := a.recipeByDetailID(ctx, recipeDetailId)
	if err != nil {
		return nil, err
	}
	if slices.Contains(full.Detail.Tags, "notion") {
		// add notion links as sources as well
		notionRecipes, err := models.NotionRecipes(qm.Where("recipe_id = ?", full.Id)).All(ctx, a.db.DB())
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

func (a *API) sectionIngredientTODB(ctx context.Context, tx *sql.Tx, i SectionIngredientInput) (*models.RecipeSectionIngredient, error) {
	si := models.RecipeSectionIngredient{
		ID:        common.ID("rs_ing"),
		Adjective: null.StringFromPtr(i.Adjective),
		Original:  null.StringFromPtr(i.Original),
		Optional:  null.BoolFromPtr(i.Optional),
	}
	amounts := []Amount{}
	for _, amt := range i.Amounts {
		amounts = append(amounts, Amount{
			Unit:  amt.Unit,
			Value: amt.Value,
		})
	}

	if err := si.Amounts.Marshal(amounts); err != nil {
		return nil, err
	}

	if i.Name == nil && i.TargetId == nil {
		return nil, nil
	}

	iOrRiD := zero.StringFromPtr(i.TargetId).ValueOrZero()
	name := zero.StringFromPtr(i.Name).ValueOrZero()

	switch i.Kind {
	case IngredientKindRecipe:

		where := qm.Where("recipe_id = ?", iOrRiD)
		if iOrRiD == "" {
			where = qm.Where("name = ?", name)
		}
		rs, err := a.recipeDetailsWhereV2(ctx, where)
		if err != nil {
			return nil, err
		}
		l(ctx).Infof("TEST123 %v", rs)

		// find or create
		{
			var r *models.RecipeDetail

			if len(rs) > 0 {
				r = rs[0]
				l(ctx).Infof("ingredient %s is a recipe, found %v", name, r.ID)

			} else {
				r, err = a.insertRecipeWrapper(ctx, tx, &RecipeWrapperInput{Detail: RecipeDetailInput{Name: name}})
				if err != nil {
					return nil, err
				}
				l(ctx).Infof("ingredient %s is a recipe, created %v", name, r.ID)
			}
			iOrRiD = r.RecipeID

			si.RecipeID = null.StringFrom(iOrRiD)
		}
	case IngredientKindIngredient:

		// missing id, need to find/create
		if iOrRiD == "" {
			ing, err := a.ingredientByName(ctx, name)
			if err != nil {
				return nil, err
			}
			iOrRiD = ing.ID
		}

		si.IngredientID = null.StringFrom(iOrRiD)
	case "":
		l(ctx).Warnf("empty kind")
		// empty table row, drop it
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown kind: %s", i.Kind)

	}
	return &si, nil
}
func (a *API) insertRecipeWrapper(ctx context.Context, tx *sql.Tx, wrapperIsThisNeeded *RecipeWrapperInput) (*models.RecipeDetail, error) {
	ctx, span := a.tracer.Start(ctx, "insertRecipeWrapper")
	defer span.End()
	span.AddEvent("got wrapper",
		trace.WithAttributes(attribute.String("recipe", spew.Sdump(wrapperIsThisNeeded))))

	if wrapperIsThisNeeded == nil {
		return nil, fmt.Errorf("no wrapper")
	}

	// first, see if there is one with the same name?
	// if so, we'll update it
	insert := wrapperIsThisNeeded.Detail
	l(ctx).Infof("insertRecipeWrapper: %s", insert.Name)

	w := *wrapperIsThisNeeded
	wrapperId := ""
	if w.Id != nil {
		wrapperId = *w.Id
	}
	var modifying *RecipeWrapper
	version := 1
	if wrapperId != "" {
		res, err := a.recipeByWrapperID(ctx, wrapperId)
		if err != nil {
			return nil, fmt.Errorf("failed to find prior recipe with wrapper id %s: %w", wrapperId, err)
		}
		modifying = res
	}

	if modifying == nil {
		res, err := a.recipeByExactName(ctx, insert.Name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// no problem, we'll insert a new one
			} else {
				return nil, err
			}
		} else {
			modifying = res
		}

	}

	if modifying != nil {
		// todo: make sure above 2 only return latest version
		version = modifying.Detail.Version + 1

		if len(insert.Sections) == 0 {
			l(ctx).Infof("no sections, modifying, so just returning")
			// TODO
			// return modifying.Detail, nil
		}
		// todo: check for equality
	}

	dbr := models.RecipeDetail{
		CreatedAt:       zero.TimeFromPtr(insert.Date).Time,
		ID:              common.ID("rd"),
		IsLatestVersion: null.BoolFrom(true),
		Name:            insert.Name,
		Quantity:        null.IntFrom(insert.Quantity),
		Servings:        null.IntFromPtr(insert.Servings),
		Tags:            insert.Tags,
		Unit:            null.StringFrom(insert.Unit),
		Version:         version,
	}
	if insert.Sources != nil {
		err := dbr.Source.Marshal(insert.Sources)
		if err != nil {
			return nil, err
		}
	}

	var dbWrapper models.Recipe

	if modifying == nil {
		dbWrapper = models.Recipe{
			ID: common.ID("r"),
		}
		if err := dbWrapper.Insert(ctx, tx, boil.Infer()); err != nil {
			return nil, err
		}
	} else {
		// modifying
		dbWrapper = models.Recipe{
			ID: modifying.Id,
		}

		currentMostRecent, err := models.RecipeDetails(
			models.RecipeDetailWhere.IsLatestVersion.EQ(null.BoolFrom(true)),
			models.RecipeDetailWhere.RecipeID.EQ(modifying.Id),
		).One(ctx, tx)
		if err != nil {
			return nil, err
		}
		currentMostRecent.IsLatestVersion = null.BoolFrom(false)
		touched, err := currentMostRecent.Update(ctx, tx, boil.Whitelist(models.RecipeDetailColumns.IsLatestVersion))
		if err != nil {
			return nil, err
		}
		l(ctx).Infof("updated %d rows that were previously the most recent version", touched)
	}

	if err := dbWrapper.AddRecipeDetails(ctx, tx, true, &dbr); err != nil {
		return nil, err
	}

	/// add sections
	for _, s := range insert.Sections {
		if len(s.Ingredients) == 0 && len(s.Instructions) == 0 {
			// empty section, drop it
			continue
		}
		data, err := json.Marshal(s.Duration)
		if err != nil {
			return nil, err
		}

		dbs := &models.RecipeSection{
			ID:                common.ID("rs"),
			DurationTimerange: null.JSONFrom(data),
		}
		if err := dbr.AddRecipeSections(ctx, tx, true, dbs); err != nil {
			return nil, err
		}

		for _, i := range s.Instructions {
			if i.Instruction == "" {
				continue
			}
			if err := dbs.AddSectionRecipeSectionInstructions(ctx, tx, true, &models.RecipeSectionInstruction{
				ID:          common.ID("rs_ins"),
				Instruction: null.StringFrom(i.Instruction),
			}); err != nil {
				return nil, err
			}
		}
		for _, i := range s.Ingredients {
			// i. = common.ID("s_ing")
			si, err := a.sectionIngredientTODB(ctx, tx, i)
			if err != nil {
				return nil, err
			}
			if si == nil {
				continue
			}
			if err := dbs.AddSectionRecipeSectionIngredients(ctx, tx, true, si); err != nil {
				return nil, err
			}
			if i.Substitutes == nil {
				continue
			}
			for _, sub := range *i.Substitutes {
				si, err := a.sectionIngredientTODB(ctx, tx, sub)
				if err != nil {
					return nil, err
				}
				if si == nil {
					continue
				}
				si2 := si
				si2.SubForIngredientID = null.StringFromPtr(i.TargetId)
				if err := dbs.AddSectionRecipeSectionIngredients(ctx, tx, true, si2); err != nil {
					return nil, err
				}
			}
		}

	}

	return &dbr, nil

}
func (a *API) enhanceMulti(ctx context.Context, si *SectionIngredient) (err error) {
	ctx, span := a.tracer.Start(ctx, "enhanceMulti")
	defer span.End()

	targets := []UnitConversionRequestTarget{
		UnitConversionRequestTargetCalories,
		UnitConversionRequestTargetMoney,
		UnitConversionRequestTargetVolume,
	}
	if !hasGrams(si.Amounts) {
		targets = append(targets, UnitConversionRequestTargetWeight)
	}
	for _, t := range targets {
		err = a.enhance(ctx, t, si)
		if err != nil {
			return err
		}
	}
	return nil
}
func (a *API) enhance(ctx context.Context, with UnitConversionRequestTarget, item *SectionIngredient) error {
	ctx, span := a.tracer.Start(ctx, "enhance")
	defer span.End()
	req := UnitConversionRequest{
		Input:        item.Amounts,
		Target:       &with,
		UnitMappings: item.Ingredient.UnitMappings,
	}
	var res Amount
	if err := a.R.Send(ctx, "convert", req, &res); err != nil {
		return fmt.Errorf("enhance: %w", err)
	}
	if res.Value > 0 { // todo: if success
		res.Value = math.Round(res.Value*100) / 100
		res.Source = zero.StringFrom("calculated").Ptr()
		item.Amounts = append(item.Amounts, res)
	}
	return nil
}
