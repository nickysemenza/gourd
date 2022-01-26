package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/models"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gopkg.in/guregu/null.v4/zero"
)

// from usda_food_category
var catmap = map[int]FoodCategory{
	1:  {Code: "0100", Description: " Dairy and Egg Products"},
	2:  {Code: "0200", Description: " Spices and Herbs"},
	3:  {Code: "0300", Description: " Baby Foods"},
	4:  {Code: "0400", Description: " Fats and Oils"},
	5:  {Code: "0500", Description: " Poultry Products"},
	6:  {Code: "0600", Description: " Soups, Sauces, and Gravies"},
	7:  {Code: "0700", Description: " Sausages and Luncheon Meats"},
	8:  {Code: "0800", Description: " Breakfast Cereals"},
	9:  {Code: "0900", Description: " Fruits and Fruit Juices"},
	10: {Code: "1000", Description: " Pork Products"},
	11: {Code: "1100", Description: " Vegetables and Vegetable Products"},
	12: {Code: "1200", Description: " Nut and Seed Products"},
	13: {Code: "1300", Description: " Beef Products"},
	14: {Code: "1400", Description: " Beverages"},
	15: {Code: "1500", Description: " Finfish and Shellfish Products"},
	16: {Code: "1600", Description: " Legumes and Legume Products"},
	17: {Code: "1700", Description: " Lamb, Veal, and Game Products"},
	18: {Code: "1800", Description: " Baked Products"},
	19: {Code: "1900", Description: " Sweets"},
	20: {Code: "2000", Description: " Cereal Grains and Pasta"},
	21: {Code: "2100", Description: " Fast Foods"},
	22: {Code: "2200", Description: " Meals, Entrees, and Side Dishes"},
	23: {Code: "2500", Description: " Snacks"},
	24: {Code: "3500", Description: " American Indian/Alaska Native Foods"},
	25: {Code: "3600", Description: " Restaurant Foods"},
	26: {Code: "4500", Description: " Branded Food Products Database"},
	27: {Code: "2600", Description: " Quality Control Materials"},
	28: {Code: "1410", Description: " Alcoholic Beverages"},
}

func (a *API) foodFromRec(ctx context.Context, foodRec *models.UsdaFood) (*Food, error) {
	if foodRec == nil {
		return nil, nil
	}

	f := Food{
		FdcId:       foodRec.FDCID,
		Description: foodRec.Description.String,
		DataType:    FoodDataType(foodRec.DataType.String),
	}
	if foodRec.FoodCategoryID.Valid {
		cat, ok := catmap[foodRec.FoodCategoryID.Int]
		if ok {
			f.Category = &cat
		}
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
func (a *API) GetFoodById(c echo.Context, fdcId int) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "GetFoodById")
	defer span.End()

	f, err := a.getFoodById(ctx, fdcId)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	if f == nil {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusOK, *f)
}
func (a *API) GetFoodsByIds(c echo.Context, params GetFoodsByIdsParams) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "GetFoodsByIds")
	defer span.End()

	foodRecs, err := models.UsdaFoods(
		qm.Where("fdc_id = any(?)", pq.Array(params.FdcId)),
		qm.Load(qm.Rels(models.UsdaFoodRels.FDCUsdaFoodNutrients,
			models.UsdaFoodNutrientRels.Nutrient)),
		qm.Load(qm.Rels(models.UsdaFoodRels.FDCUsdaBrandedFood)),
		qm.Load(qm.Rels(models.UsdaFoodRels.FDCUsdaFoodPortions)),
	).All(ctx, a.db.DB())
	if err != nil {
		return err
	}
	items := []Food{}
	for _, foodRec := range foodRecs {
		f, err := a.foodFromRec(ctx, foodRec)
		if err != nil {
			return err
		}
		items = append(items, *f)
	}

	listMeta := Items{PageCount: 1}
	listMeta.setTotalCount(uint64(len(items)))

	resp := PaginatedFoods{
		Foods: &items,
		Meta:  listMeta,
	}

	return c.JSON(http.StatusOK, resp)
}
func (a *API) SearchFoods(c echo.Context, params SearchFoodsParams) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "SearchFoods")
	defer span.End()

	paginationParams, listMeta := parsePagination(params.Offset, params.Limit)
	dataTypes := []string{}
	if params.DataTypes != nil {
		for _, x := range *params.DataTypes {
			if x != "" {
				dataTypes = append(dataTypes, string(x))
			}
		}
	}
	foods, count, err := a.DB().SearchFoods(ctx, string(params.Name), dataTypes, nil, paginationParams...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}
	listMeta.setTotalCount(count)

	items, err := a.buildPaginatedFood(ctx, foods)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}
	resp := PaginatedFoods{
		Foods: &items,
		Meta:  listMeta,
	}

	return c.JSON(http.StatusOK, resp)

}

func (a *API) getFoodById(ctx context.Context, fdcId int) (*Food, error) {
	ctx, span := a.tracer.Start(ctx, "getFoodById")
	defer span.End()

	foodRec, err := models.UsdaFoods(
		qm.Where("fdc_id = ?", fdcId),
		qm.Load(qm.Rels(models.UsdaFoodRels.FDCUsdaFoodNutrients,
			models.UsdaFoodNutrientRels.Nutrient)),
		qm.Load(qm.Rels(models.UsdaFoodRels.FDCUsdaBrandedFood)),
		qm.Load(qm.Rels(models.UsdaFoodRels.FDCUsdaFoodPortions)),
	).One(ctx, a.db.DB())

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get food with id %d: %w", fdcId, err)
	}

	f, err := a.foodFromRec(ctx, foodRec)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (a *API) buildPaginatedFood(ctx context.Context, foods []db.Food) ([]Food, error) {
	ctx, span := a.tracer.Start(ctx, "buildPaginatedFood")
	defer span.End()

	ids := []int{}
	for _, food := range foods {
		ids = append(ids, food.FdcID)
	}
	foodRecs, err := models.UsdaFoods(
		qm.Where("fdc_id = any(?)", pq.Array(ids)),
		qm.Load(qm.Rels(models.UsdaFoodRels.FDCUsdaFoodNutrients,
			models.UsdaFoodNutrientRels.Nutrient)),
		qm.Load(qm.Rels(models.UsdaFoodRels.FDCUsdaBrandedFood)),
		qm.Load(qm.Rels(models.UsdaFoodRels.FDCUsdaFoodPortions)),
	).All(ctx, a.db.DB())
	if err != nil {
		return nil, err
	}
	items := []Food{}
	for _, foodRec := range foodRecs {
		f, err := a.foodFromRec(ctx, foodRec)
		if err != nil {
			return nil, err
		}
		items = append(items, *f)
	}
	return items, nil

}

func (a *API) AssociateFoodWithIngredient(c echo.Context, ingredientId string, params AssociateFoodWithIngredientParams) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "AssociateFoodWithIngredient")
	defer span.End()
	err := a.DB().AssociateFoodWithIngredient(ctx, ingredientId, params.FdcId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, nil)
}

func (a *API) LoadIngredientMappings(c echo.Context) error {
	ctx := c.Request().Context()
	var r IngredientMappingsPayload
	if err := c.Bind(&r); err != nil {
		err = fmt.Errorf("invalid format for input: %w", err)
		return handleErr(c, err)
	}

	if err := a.loadIngredientMappings(ctx, r.IngredientMappings); err != nil {
		return handleErr(c, err)
	}
	return c.JSON(http.StatusAccepted, nil)
}

func (a *API) loadIngredientMappings(ctx context.Context, mapping []IngredientMapping) error {
	ctx, span := a.tracer.Start(ctx, "loadIngredientMappings")
	defer span.End()

	for _, m := range mapping {

		ing, err := a.DB().IngredientByName(ctx, m.Name)
		if err != nil {
			return err
		}

		if m.FdcId != nil {
			err = a.DB().AssociateFoodWithIngredient(ctx, ing.Id, *m.FdcId)
			if err != nil {
				return err
			}
		}

		childIds := []string{}
		for _, alias := range m.Aliases {
			ing, err := a.DB().IngredientByName(ctx, alias)
			if err != nil {
				return err
			}
			childIds = append(childIds, ing.Id)
		}
		err = a.DB().MergeIngredients(ctx, ing.Id, childIds)
		if err != nil {
			return err
		}

		actualPairs := len(m.UnitMappings)
		loadedPairs := int64(0)
		for _, u := range m.UnitMappings {
			u := db.IngredientUnitMapping{
				IngredientId: ing.Id,
				UnitA:        u.A.Unit,
				AmountA:      u.A.Value,
				UnitB:        u.B.Unit,
				AmountB:      u.B.Value,
				Source:       zero.StringFromPtr(u.Source).String,
			}
			num, err := a.DB().AddIngredientUnit(ctx, u)
			if err != nil {
				return err
			}
			loadedPairs += num
		}
		log.Printf("loaded %s (%v), fdc: %d=>%s, %d/%d unit pairs", m.Name, strings.Join(m.Aliases, ", "), m.FdcId, ing.Id, loadedPairs, actualPairs)

	}
	return nil
}
