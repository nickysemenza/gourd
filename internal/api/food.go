package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/internal/common"
	"github.com/nickysemenza/gourd/internal/db/models"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"gopkg.in/guregu/null.v4/zero"
)

func (a *API) grabFood(ctx context.Context, id int) (*TempFood, error) {
	ctx, span := a.tracer.Start(ctx, "grabFood")
	defer span.End()
	var result *TempFood
	err := a.R.Send(ctx, "debug/get_usda?name="+url.QueryEscape(fmt.Sprintf("%d", id)), nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *API) GetFoodById(c echo.Context, fdcId int) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "GetFoodById")
	defer span.End()

	f, err := a.grabFood(ctx, fdcId)
	if err != nil {
		return handleErr(c, err)
	}
	if f == nil {
		return handleErr(c, fmt.Errorf("failed to load food %v: %w", fdcId, common.ErrNotFound))
	}

	return c.JSON(http.StatusOK, *f)
}
func (a *API) GetFoodsByIds(c echo.Context, params GetFoodsByIdsParams) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "GetFoodsByIds")
	defer span.End()

	items := []TempFood{}
	for _, id := range params.FdcId {
		res, err := a.grabFood(ctx, id)
		if err != nil {
			return handleErr(c, err)
		}
		if res != nil {
			items = append(items, *res)
		}
	}

	listMeta := Items{PageCount: 1}
	listMeta.setTotalCount(int64(len(items)))

	resp := PaginatedFoods{
		Foods: &items,
		Meta:  listMeta,
	}

	return c.JSON(http.StatusOK, resp)
}
func (a *API) SearchFoods(c echo.Context, params SearchFoodsParams) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "SearchFoods")
	defer span.End()

	var byItem []TempFood
	err := a.R.Send(ctx, "debug/search_usda?name="+url.QueryEscape(string(params.Name)), nil, &byItem)
	if err != nil {
		return handleErr(c, err)
	}
	resp := FoodSearchResult{
		Foods: byItem,
	}

	return c.JSON(http.StatusOK, resp)

}

func (a *API) AssociateFoodWithIngredient(c echo.Context, ingredientId string, params AssociateFoodWithIngredientParams) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "AssociateFoodWithIngredient")
	defer span.End()
	err := a.DB().AssociateFoodWithIngredient(ctx, ingredientId, params.FdcId)
	if err != nil {
		return handleErr(c, err)
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

	if err := a.insertIngredientMappings(ctx, r.IngredientMappings); err != nil {
		return handleErr(c, err)
	}
	return c.JSON(http.StatusAccepted, nil)
}

type IngredientMappings []IngredientMapping

func (m IngredientMappings) Valdiate() error {
	var names, aliases []string
	for _, each := range m {
		names = append(names, each.Name)
		aliases = append(aliases, each.Aliases...)
	}
	for _, x := range names {
		if slices.Contains(aliases, x) {
			return fmt.Errorf("Validate: %s is a top level and an alias", x)
		}
	}
	return nil
}

func (a *API) insertIngredientMappings(ctx context.Context, mapping IngredientMappings) error {
	ctx, span := a.tracer.Start(ctx, "insertIngredientMappings")
	defer span.End()

	if err := mapping.Valdiate(); err != nil {
		return err
	}

	tx := a.tx(ctx)

	for _, m := range mapping {

		ing, err := a.ingredientByName(ctx, m.Name)
		if err != nil {
			return err
		}

		if m.FdcId != nil {
			food, err := a.grabFood(ctx, *m.FdcId)
			if err != nil {
				return err
			}
			if food == nil {
				l(ctx).Warnf("food %d not found,", *m.FdcId)
			}
			err = a.DB().AssociateFoodWithIngredient(ctx, ing.ID, *m.FdcId)
			if err != nil {
				return err
			}
		}

		childIds := []string{}
		for _, alias := range m.Aliases {
			ing, err := a.ingredientByName(ctx, alias)
			if err != nil {
				return err
			}
			childIds = append(childIds, ing.ID)
		}
		err = a.DB().MergeIngredients(ctx, tx, ing.ID, childIds)
		if err != nil {
			return err
		}

		actualPairs := len(m.UnitMappings)
		loadedPairs := int64(0)
		for _, u := range m.UnitMappings {
			u := models.IngredientUnit{
				IngredientID: ing.ID,
				UnitA:        u.A.Unit,
				AmountA:      common.DecimalFromFloat(u.A.Value),
				UnitB:        u.B.Unit,
				AmountB:      common.DecimalFromFloat(u.B.Value),
				Source:       zero.StringFromPtr(u.Source).String,
			}
			num, err := a.DB().AddIngredientUnit(ctx, u)
			if err != nil {
				return err
			}
			loadedPairs += num
		}
		log.Printf("loaded %s (%v), fdc: %d=>%s, %d/%d unit pairs", m.Name, strings.Join(m.Aliases, ", "), m.FdcId, ing.ID, loadedPairs, actualPairs)

	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
