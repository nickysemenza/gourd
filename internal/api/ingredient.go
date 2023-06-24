package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/internal/db/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.opentelemetry.io/otel/attribute"
)

func (a *API) CreateIngredients(c echo.Context) error {
	ctx := c.Request().Context()
	var i Ingredient
	if err := c.Bind(&i); err != nil {
		err = fmt.Errorf("invalid format for input: %w", err)
		return sendErr(c, http.StatusBadRequest, err)
	}
	ing, err := a.ingredientByName(ctx, i.Name)
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}

	detail, err := a.ingredientById(ctx, IngredientID(ing.ID), true)
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, detail.Ingredient)
}

func (a *API) IngredientIdByName(ctx context.Context, name string) (string, error) {
	ing, err := a.ingredientByName(ctx, name)
	if err != nil {
		return "", err
	}
	return ing.ID, nil
}

func (a *API) enhanceWithFDC(ctx context.Context, fdcId int64, detail *IngredientDetail) (err error) {
	ctx, span := a.tracer.Start(ctx, "enhanceWithFDC")
	defer span.End()

	food, err := a.grabFood(ctx, int(fdcId))
	if err != nil {
		return
	}
	if food == nil {
		return
	}

	detail.Food = food
	span.SetAttributes(attribute.Int64("fdc_id", fdcId))

	detail.UnitMappings = append(detail.UnitMappings, food.UnitMappings...)
	return
}

func (a *API) ListIngredients(c echo.Context, params ListIngredientsParams) error {
	ctx := c.Request().Context()

	ctx, span := a.tracer.Start(ctx, "ListIngredients")
	defer span.End()

	listMeta := parsePagination(params.Offset, params.Limit)

	var mods []qm.QueryMod
	if params.IngredientId != nil && len(*params.IngredientId) > 0 {
		qm.Where("ingredients.id = ?", params.IngredientId)
	}

	items, count, err := a.IngredientListV2(ctx, listMeta, mods...)
	if err != nil {
		return handleErr(c, err)
	}

	listMeta.setTotalCount(uint64(count))

	resp := PaginatedIngredients{
		Ingredients: &items,
		Meta:        listMeta,
	}
	return c.JSON(http.StatusOK, resp)
}
func (a *API) convertIngredientToRecipe(ctx context.Context, ingredientId string) (*models.RecipeDetail, error) {
	ctx, span := a.tracer.Start(ctx, "convertIngredientToRecipe")
	defer span.End()
	convertedPrefix := "[converted]"

	tx, err := a.db.DB().BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	i, err := models.FindIngredient(ctx, tx, ingredientId)
	if err != nil {
		return nil, err
	}
	name := i.Name
	if strings.Contains(name, convertedPrefix) {
		return nil, fmt.Errorf("%s has already been converted to recipe", name)
	}

	dbVersion, err := a.insertRecipeWrapper(ctx, tx, &RecipeWrapperInput{Detail: RecipeDetailInput{Name: name}})
	if err != nil {
		return nil, err
	}

	sectionIngredients, err := models.RecipeSectionIngredients(
		models.RecipeSectionIngredientWhere.IngredientID.EQ(null.StringFrom(ingredientId)),
	).All(ctx, tx)
	if err != nil {
		return nil, err
	}
	affected, err := sectionIngredients.UpdateAll(ctx, tx,
		models.M{
			models.RecipeSectionIngredientColumns.IngredientID: nil,
			models.RecipeSectionIngredientColumns.RecipeID:     dbVersion.RecipeID,
		})
	if err != nil {
		return nil, err
	}

	l(ctx).Infof("updated %d section ingredients to point to recipe", affected)

	i.Name = convertedPrefix + " " + i.Name
	_, err = i.Update(ctx, tx, boil.Whitelist(models.IngredientColumns.Name))
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return dbVersion, nil

}
func (a *API) ConvertIngredientToRecipe(c echo.Context, ingredientId string) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "ConvertIngredientToRecipe")
	defer span.End()

	detail, err := a.convertIngredientToRecipe(ctx, ingredientId)

	if err != nil {
		return handleErr(c, err)
	}
	tr, err := a.recipeByDetailID(ctx, detail.ID)
	if err != nil {
		return handleErr(c, err)
	}

	return c.JSON(http.StatusCreated, tr)
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
		return handleErr(c, err)
	}

	detail, err := a.ingredientById(ctx, IngredientID(ingredientId), true)
	if err != nil {
		return handleErr(c, err)
	}

	return c.JSON(http.StatusCreated, detail.Ingredient)
}

func (a *API) GetIngredientById(c echo.Context, ingredientId string) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "GetIngredientById")
	defer span.End()

	ing, err := a.ingredientById(ctx, IngredientID(ingredientId), true)
	if err != nil {
		return handleErr(c, err)
	}
	if ing == nil {
		return sendErr(c, http.StatusNotFound, fmt.Errorf("no ingredient with id %s", ingredientId))
	}
	return c.JSON(http.StatusOK, ing)
}

func (a *API) Scrape(ctx context.Context, url string) (*RecipeWrapper, error) {
	ctx, span := a.tracer.Start(ctx, "Scrape")
	defer span.End()

	r, err := a.FetchAndTransform(ctx, url, a.IngredientIdByName)
	if err != nil {
		return nil, err
	}
	return a.CreateRecipe(ctx, r)
}
func (a *API) ScrapeRecipe(c echo.Context) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "ScrapeRecipe")
	defer span.End()

	var i ScrapeRecipeJSONBody
	if err := c.Bind(&i); err != nil {
		err = fmt.Errorf("invalid format for input: %w", err)
		return sendErr(c, http.StatusBadRequest, err)
	}

	r, err := a.Scrape(ctx, i.Url)
	if err != nil {
		return handleErr(c, err)
	}
	return c.JSON(http.StatusCreated, r)
}
