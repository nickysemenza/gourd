package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/models"
	"github.com/nickysemenza/gourd/rs_client"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/guregu/null.v4/zero"
)

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
func (a *API) addDetailsToIngredients(ctx context.Context, ing []db.Ingredient) ([]IngredientDetail, error) {
	ctx, span := a.tracer.Start(ctx, "addDetailsToIngredients")
	defer span.End()

	var ingredientIds []string
	for _, i := range ing {
		ingredientIds = append(ingredientIds, i.Id)
	}

	span.AddEvent("ingredient-params", trace.WithAttributes(attribute.StringSlice("id", ingredientIds)))

	parent, _, err := a.DB().GetIngrientsParent(ctx, ingredientIds...)
	if err != nil {
		return nil, err
	}
	for _, i := range parent {
		ingredientIds = append(ingredientIds, i.Id)
	}

	span.AddEvent("ingredient-plus-parent", trace.WithAttributes(attribute.StringSlice("id", ingredientIds)))

	linkedRecipes, err := a.DB().GetRecipeDetailsWithIngredient(ctx, ingredientIds...)
	if err != nil {
		return nil, err
	}

	items := make([]IngredientDetail, len(ing))
	for x, i := range ing {
		// assemble

		detail, err := a.makeDetail(ctx, i, parent, linkedRecipes)
		if err != nil {
			return nil, err
		}
		unitMappings, err := a.DB().GetIngredientUnits(ctx, []string{i.Id, i.Parent.String})
		if err != nil {
			return nil, err
		}

		for _, m := range unitMappings {
			detail.UnitMappings = append(detail.UnitMappings, UnitMapping{
				Amount{Unit: m.UnitA, Value: m.AmountA},
				Amount{Unit: m.UnitB, Value: m.AmountB},
				zero.StringFrom(fmt.Sprintf("%s (%s)", m.Source, i.Id)).Ptr(),
			})
		}
		span.AddEvent("mappings", trace.WithAttributes(attribute.String("mappings", spew.Sdump(unitMappings))))

		if i.FdcID.Valid {
			err := a.enhanceWithFDC(ctx, i.FdcID.Int64, detail)
			if err != nil {
				return nil, fmt.Errorf("enhanceWithFDC: %w", err)
			}
		}
		items[x] = *detail
	}

	return items, nil
}
func (a *API) getFoodById2(ctx context.Context, fdcId int) (*Food, error) {
	ctx, span := a.tracer.Start(ctx, "getFoodById2")
	defer span.End()

	foodRec, err := models.UsdaFoods(
		qm.Where("fdc_id = ?", fdcId),
		qm.Load(
			qm.Rels(models.UsdaFoodRels.FDCUsdaFoodNutrients, models.UsdaFoodNutrientRels.Nutrient),
		),
		qm.Load(
			qm.Rels(models.UsdaFoodRels.FDCUsdaBrandedFood),
		),
		// qm.Load(
		// 	qm.Rels(models.UsdaFoodRels.FoodCategory),
		// ),
		qm.Load(
			qm.Rels(models.UsdaFoodRels.FDCUsdaFoodPortions),
		),
	).One(ctx, a.db.DB())

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get food with id %d: %w", fdcId, err)
	}

	if foodRec == nil {
		return nil, nil
	}

	f := Food{
		FdcId:       foodRec.FDCID,
		Description: foodRec.Description.String,
		DataType:    FoodDataType(foodRec.DataType.String),
	}

	nutrients := []FoodNutrient{}
	for _, fn := range foodRec.R.FDCUsdaFoodNutrients {
		n := fn.R.Nutrient
		if n == nil {
			continue
		}
		nutrients = append(nutrients, FoodNutrient{
			Amount:     float64(fn.Amount.Float32),
			DataPoints: fn.DataPoints.Int,
			Nutrient: Nutrient{
				Id:       n.ID,
				Name:     n.Name.String,
				UnitName: FoodNutrientUnit(n.UnitName.String),
			},
		})

	}
	f.Nutrients = nutrients

	if foodRec.R.FDCUsdaBrandedFood != nil {
		brandInfo := foodRec.R.FDCUsdaBrandedFood
		f.BrandedInfo = &BrandedFood{
			BrandOwner:          brandInfo.BrandOwner.Ptr(),
			BrandedFoodCategory: brandInfo.BrandedFoodCategory.Ptr(),
			HouseholdServing:    brandInfo.HouseholdServingFulltext.Ptr(),
			Ingredients:         brandInfo.Ingredients.Ptr(),
			ServingSize:         float64(brandInfo.ServingSize.Float32),
			ServingSizeUnit:     brandInfo.ServingSizeUnit.String,
		}
	}
	// todo: FoodCategory{} but not super necessary

	portions := []FoodPortion{}
	for _, p := range foodRec.R.FDCUsdaFoodPortions {
		portions = append(portions, FoodPortion{
			Amount:             float64(p.Amount.Float32),
			GramWeight:         float64(p.GramWeight.Float32),
			Id:                 p.ID,
			Modifier:           p.Modifier.String,
			PortionDescription: p.PortionDescription.String,
		})
	}
	f.Portions = &portions

	m, err := a.UnitMappingsFromFood(ctx, &f)
	if err != nil {
		return nil, err
	}
	f.UnitMappings = m
	return &f, nil
}

func (a *API) enhanceWithFDC(ctx context.Context, fdcId int64, detail *IngredientDetail) (err error) {
	ctx, span := a.tracer.Start(ctx, "enhanceWithFDC")
	defer span.End()

	food, err := a.getFoodById2(ctx, int(fdcId))
	if err != nil {
		return
	}
	if food == nil {
		return
	}

	detail.Food = food
	span.SetAttributes(attribute.Int("fdc_id", food.FdcId))
	// ??
	// var m []UnitMapping
	// m, err = a.UnitMappingsFromFood(ctx, food)
	// if err != nil {
	// 	return
	// }
	detail.UnitMappings = append(detail.UnitMappings, food.UnitMappings...)
	return
}
func (a *API) UnitMappingsFromFood(ctx context.Context, food *Food) ([]UnitMapping, error) {
	// todo: store these in DB instead of inline parsing ?
	m := []UnitMapping{}
	if food.BrandedInfo != nil && food.BrandedInfo.HouseholdServing != nil {
		var res []Amount
		err := a.R.Call(ctx, *food.BrandedInfo.HouseholdServing, rs_client.ParseAmount, &res)
		if err != nil {
			return nil, err
		}
		if len(res) > 0 {
			m = append(m, UnitMapping{
				Amount{Unit: food.BrandedInfo.ServingSizeUnit, Value: food.BrandedInfo.ServingSize},
				res[0],
				zero.StringFrom("fdc hs").Ptr()})
		}
	}
	if food.Portions != nil {
		for _, p := range *food.Portions {

			if p.PortionDescription != "" {
				var res []Amount
				err := a.R.Call(ctx, p.PortionDescription, rs_client.ParseAmount, &res)
				if err != nil {
					err := fmt.Errorf("failed to parse '%s' :%w", p.PortionDescription, err)
					log.Error(err)
					continue
					// return nil, err
				}
				if len(res) == 0 {
					continue
				}
				m = append(m, UnitMapping{
					res[0],
					Amount{Unit: "grams", Value: p.GramWeight},
					zero.StringFrom("fdc p1").Ptr()})
			} else {
				m = append(m, UnitMapping{
					Amount{Unit: p.Modifier, Value: p.Amount},
					Amount{Unit: "grams", Value: p.GramWeight},
					zero.StringFrom("fdc p2").Ptr()})
			}

		}
	}
	for _, n := range food.Nutrients {
		if n.Nutrient.UnitName == FoodNutrientUnitKCAL {
			m = append(m, UnitMapping{
				Amount{Unit: "kcal", Value: n.Amount},
				Amount{Unit: "grams", Value: 100},
				zero.StringFrom("fdc n").Ptr()})
		}
	}
	return m, nil

}
func (a *API) ListIngredients(c echo.Context, params ListIngredientsParams) error {
	ctx := c.Request().Context()

	ctx, span := a.tracer.Start(ctx, "ListIngredients")
	defer span.End()

	paginationParams, listMeta := parsePagination(params.Offset, params.Limit)
	var ids []string
	if params.IngredientId != nil {
		ids = *params.IngredientId
	}
	ing, count, err := a.DB().GetIngredients(ctx, "", ids, paginationParams...)
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}

	items, err := a.addDetailsToIngredients(ctx, ing)
	if err != nil {
		return sendErr(c, http.StatusBadRequest, err)
	}

	listMeta.TotalCount = int(count)

	resp := PaginatedIngredients{
		Ingredients: &items,
		Meta:        listMeta,
	}
	return c.JSON(http.StatusOK, resp)
}

func (a *API) makeDetail(ctx context.Context, i db.Ingredient, parent db.Ingredients, linkedRecipes db.RecipeDetails) (*IngredientDetail, error) {
	ctx, span := a.tracer.Start(ctx, "makeDetail")
	defer span.End()

	// find linked ingredients
	same := []IngredientDetail{}
	for _, x := range parent.ByParent()[i.Id] {
		d, err := a.makeDetail(ctx, x, parent, linkedRecipes)
		if err != nil {
			return nil, err
		}
		same = append(same, *d)
	}

	// find linked recipes
	recipes := []RecipeDetail{}
	for _, x := range linkedRecipes.ByIngredientId()[i.Id] {
		tr, err := a.transformRecipe(ctx, x, false)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, *tr)
	}

	detail := IngredientDetail{
		Ingredient:   transformIngredient(i),
		Children:     &same,
		Recipes:      recipes,
		UnitMappings: []UnitMapping{},
	}

	return &detail, nil
}

func (a *API) ConvertIngredientToRecipe(c echo.Context, ingredientId string) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "ConvertIngredientToRecipe")
	defer span.End()

	detail, err := a.DB().IngredientToRecipe(ctx, ingredientId)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	tr, err := a.transformRecipe(ctx, *detail, true)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
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
		return sendErr(c, http.StatusInternalServerError, err)
	}

	ing, err := a.DB().GetIngredientById(ctx, ingredientId)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, transformIngredient(*ing))
}

func (a *API) GetIngredientById(c echo.Context, ingredientId string) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "GetIngredientById")
	defer span.End()

	ing, err := a.DB().GetIngredientById(ctx, ingredientId)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	if ing == nil {
		return sendErr(c, http.StatusNotFound, fmt.Errorf("no ingredient with id %s", ingredientId))
	}
	foo, err := a.addDetailsToIngredients(ctx, []db.Ingredient{*ing})
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, foo[0])
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
		return sendErr(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, r)
}
