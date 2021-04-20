package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/rs_client"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
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

	sameAs, _, err := a.Manager.DB().GetIngrientsSameAs(ctx, ingredientIds...)
	if err != nil {
		return nil, err
	}
	for _, i := range sameAs {
		ingredientIds = append(ingredientIds, i.Id)
	}
	linkedRecipes, err := a.Manager.DB().GetRecipeDetailsWithIngredient(ctx, ingredientIds...)
	if err != nil {
		return nil, err
	}
	unitMappings, err := a.Manager.DB().GetIngredientUnits(ctx, ingredientIds)
	if err != nil {
		return nil, err
	}

	items := make([]IngredientDetail, len(ing))
	for x, i := range ing {
		// assemble
		ctx, span2 := a.tracer.Start(ctx, "addDetailsToIngredients: enhance w/ fdc")
		defer span2.End()
		detail := makeDetail(i, sameAs, linkedRecipes)

		for _, x := range append(sameAs.IdsBySameAs(i.Id), i.Id) {
			if val, ok := unitMappings[x]; ok {
				for _, m := range val {
					detail.UnitMappings = append(detail.UnitMappings, UnitMapping{
						Amount{Unit: m.UnitA, Value: m.AmountA},
						Amount{Unit: m.UnitB, Value: m.AmountB},
						fmt.Sprintf("%s (%s)", m.Source, x),
					})
				}
			}
		}

		if i.FdcID.Valid {
			food, err := a.getFoodById(ctx, int(i.FdcID.Int64))
			if err != nil {
				return nil, err
			}
			detail.Food = food
			span2.SetAttributes(attribute.Int("fdc_id", food.FdcId))

			m, err := UnitMappingsFromFood(ctx, food)
			if err != nil {
				return nil, err
			}
			detail.UnitMappings = append(detail.UnitMappings, m...)

		}

		items[x] = detail
		span2.End()
	}
	return items, nil
}
func UnitMappingsFromFood(ctx context.Context, food *Food) ([]UnitMapping, error) {
	// todo: store these in DB instead of inline parsing ?
	m := []UnitMapping{}
	if food.BrandedInfo != nil && food.BrandedInfo.HouseholdServing != nil {
		var res []Amount
		err := rs_client.Parse(ctx, *food.BrandedInfo.HouseholdServing, rs_client.Amount, &res)
		if err != nil {
			return nil, err
		}
		if len(res) > 0 {
			m = append(m, UnitMapping{
				Amount{Unit: food.BrandedInfo.ServingSizeUnit, Value: food.BrandedInfo.ServingSize},
				res[0],
				"fdc hs"})
		}
	}
	if food.Portions != nil {
		for _, p := range *food.Portions {

			if p.PortionDescription != "" {
				var res []Amount
				err := rs_client.Parse(ctx, p.PortionDescription, rs_client.Amount, &res)
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
					"fdc p1"})
			} else {
				m = append(m, UnitMapping{
					Amount{Unit: p.Modifier, Value: p.Amount},
					Amount{Unit: "grams", Value: p.GramWeight},
					"fdc p2"})
			}

		}
	}
	for _, n := range food.Nutrients {
		if n.Nutrient.UnitName == FoodNutrientUnit_KCAL {
			m = append(m, UnitMapping{
				Amount{Unit: "kcal", Value: n.Amount},
				Amount{Unit: "grams", Value: 100},
				"fdc n"})
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
	ing, count, err := a.Manager.DB().GetIngredients(ctx, "", ids, paginationParams...)
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

	detail := IngredientDetail{
		Ingredient:   transformIngredient(i),
		Children:     same,
		Recipes:      recipes,
		UnitMappings: []UnitMapping{},
	}

	return detail
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
