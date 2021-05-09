package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/db"
	log "github.com/sirupsen/logrus"
)

func (a *API) addDetailToFood(ctx context.Context, f *Food, categoryId int64) error {
	ctx, span := a.tracer.Start(ctx, "addDetailToFood")
	defer span.End()

	fatalErrors := make(chan error)
	wgDone := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(3)

	fdcId := f.FdcId

	fNutrients := make([]FoodNutrient, 0)
	go func() {
		nutrientRows, err := a.DB().GetFoodNutrients(ctx, fdcId)
		if err != nil {
			fatalErrors <- err
		}

		nutrientIDs := make([]int, len(nutrientRows))
		for x, nr := range nutrientRows {
			nutrientIDs[x] = nr.NutrientID
		}
		nutrients, err := a.DB().GetNutrients(ctx, nutrientIDs...)
		if err != nil {
			fatalErrors <- err
		}

		nutrientsById := nutrients.ById()

		for _, nr := range nutrientRows {
			nDetail := nutrientsById[nr.NutrientID]
			fNutrients = append(fNutrients, FoodNutrient{
				Amount:     nr.Amount,
				DataPoints: int(nr.DataPoints.Int64),
				Nutrient: Nutrient{
					Id:       nDetail.ID,
					Name:     nDetail.Name,
					UnitName: FoodNutrientUnit(nDetail.UnitName),
				},
			})
		}
		wg.Done()
	}()

	var brandInfoRes *BrandedFood

	go func() {
		brandInfo, err := a.DB().GetBrandInfo(ctx, fdcId)
		if err != nil {
			fatalErrors <- err
		}
		if brandInfo != nil && brandInfo.BrandOwner != nil && *brandInfo.BrandOwner != "" {
			brandInfoRes = &BrandedFood{
				BrandOwner:          brandInfo.BrandOwner,
				BrandedFoodCategory: brandInfo.BrandedFoodCategory,
				HouseholdServing:    brandInfo.HouseholdServing,
				Ingredients:         brandInfo.Ingredients,
				ServingSize:         brandInfo.ServingSize,
				ServingSizeUnit:     brandInfo.ServingSizeUnit,
			}
		}
		wg.Done()
	}()
	category, err := a.DB().GetCategory(ctx, categoryId)
	if err != nil {
		return err
	}
	if category != nil && category.Code != "" {
		f.Category = &FoodCategory{
			Code:        category.Code,
			Description: category.Description,
		}
	}

	apiPortions := make([]FoodPortion, 0)
	go func() {
		portions, err := a.DB().GetFoodPortions(ctx, fdcId)
		if err != nil {
			fatalErrors <- err
		}

		if portionsById, ok := portions.ByFdcId()[fdcId]; ok {

			for _, p := range portionsById {
				apiPortions = append(apiPortions, FoodPortion{
					Amount:             p.Amount.Float64,
					GramWeight:         p.GramWeight,
					Id:                 p.Id,
					Modifier:           p.Modifier.String,
					PortionDescription: p.PortionDescription.String,
				})
			}
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(wgDone)
	}()

	// Wait until either WaitGroup is done or an error is received through the channel
	select {
	case <-wgDone:
		// carry on
		break
	case err := <-fatalErrors:
		// close(fatalErrors)
		err = fmt.Errorf("err on chann: %w", err)
		log.Print(err)
		return err
	}
	f.Nutrients = fNutrients
	f.Portions = &apiPortions
	f.BrandedInfo = brandInfoRes
	m, err := UnitMappingsFromFood(ctx, f)
	if err != nil {
		return err
	}
	f.UnitMappings = append(f.UnitMappings, m...)

	return nil
}
func (a *API) getFoodById(ctx context.Context, fdcId int) (*Food, error) {
	ctx, span := a.tracer.Start(ctx, "getFoodById")
	defer span.End()
	food, err := a.DB().GetFood(ctx, fdcId)
	if err != nil {
		return nil, err
	}
	if food.FdcID == 0 {
		return nil, fmt.Errorf("could not find food with fdc id %d", fdcId)
	}

	f := Food{
		Description: food.Description,
		DataType:    FoodDataType(food.DataType),
		FdcId:       food.FdcID,
	}

	err = a.addDetailToFood(ctx, &f, food.CategoryID.Int64)
	return &f, err
}
func (a *API) GetFoodById(c echo.Context, fdcId int) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "GetFoodById")
	defer span.End()

	f, err := a.getFoodById(ctx, fdcId)
	if err != nil {
		return sendErr(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, *f)
}
func (a *API) GetFoodsByIds(c echo.Context, params GetFoodsByIdsParams) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "GetFoodsByIds")
	defer span.End()

	foods, count, err := a.Manager.DB().FoodsByIds(ctx, params.FdcId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}
	listMeta := &Items{PageCount: 1}
	listMeta.setTotalCount(count)
	// spew.Dump(foods, count)
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
	foods, count, err := a.Manager.DB().SearchFoods(ctx, string(params.Name), dataTypes, nil, paginationParams...)
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

func (a *API) buildPaginatedFood(ctx context.Context, foods []db.Food) ([]Food, error) {
	items := make([]Food, len(foods))
	var wg sync.WaitGroup
	fatalErrors := make(chan error)
	wgDone := make(chan bool)

	for x := range foods {
		wg.Add(1)
		go func(i int) {
			food := foods[i]
			f := Food{
				Description: food.Description,
				DataType:    FoodDataType(food.DataType),
				FdcId:       food.FdcID,
				Nutrients:   make([]FoodNutrient, 0),
			}
			err := a.addDetailToFood(ctx, &f, food.CategoryID.Int64)
			if err != nil {
				fatalErrors <- err
			}
			items[i] = f
			wg.Done()
		}(x)
	}
	go func() {
		wg.Wait()
		close(wgDone)
	}()

	// Wait until either WaitGroup is done or an error is received through the channel
	select {
	case <-wgDone:
		// carry on
		break
	case err := <-fatalErrors:
		// close(fatalErrors)
		return nil, err
	}

	return items, nil

}

func (a *API) AssociateFoodWithIngredient(c echo.Context, ingredientId string, params AssociateFoodWithIngredientParams) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "AssociateFoodWithIngredient")
	defer span.End()
	err := a.Manager.DB().AssociateFoodWithIngredient(ctx, ingredientId, params.FdcId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, nil)
}

func (a *API) LoadIngredientMappings(ctx context.Context, mapping []IngredientMapping) error {
	for _, m := range mapping {

		ing, err := a.DB().IngredientByName(ctx, m.Name)
		if err != nil {
			return err
		}

		err = a.Manager.DB().AssociateFoodWithIngredient(ctx, ing.Id, m.FdcID)
		if err != nil {
			return err
		}
		log.Printf("loaded %s (%v), fdc: %d=>%s, %d unit pairs", m.Name, strings.Join(m.Aliases, ", "), m.FdcID, ing.Id, len(m.UnitMappings))

		sameAsIds := []string{}
		for _, alias := range m.Aliases {
			ing, err := a.DB().IngredientByName(ctx, alias)
			if err != nil {
				return err
			}
			sameAsIds = append(sameAsIds, ing.Id)
		}
		err = a.Manager.DB().MergeIngredients(ctx, ing.Id, sameAsIds)
		if err != nil {
			return err
		}

		for _, u := range m.UnitMappings {
			u := db.IngredientUnitMapping{
				IngredientId: ing.Id,
				UnitA:        u.A.Unit,
				AmountA:      u.A.Value,
				UnitB:        u.B.Unit,
				AmountB:      u.B.Value,
				Source:       u.Source,
			}
			err = a.Manager.DB().AddIngredientUnit(ctx, u)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
