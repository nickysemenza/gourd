package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type IngredientID string

func (r *RecipeDetail) summary(multiplier float64) EntitySummary {
	return EntitySummary{
		Name:       r.Name,
		Id:         r.Id,
		Kind:       IngredientKindRecipe,
		Multiplier: multiplier,
	}
}

// type UsageValue struct {
// 	Ings []IngredientUsage
// 	Sum  []Amount
// 	Ing  EntitySummary
// }

func (a *API) SumRecipes(c echo.Context) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "SumRecipes")
	defer span.End()

	var r SumRecipesJSONRequestBody
	if err := c.Bind(&r); err != nil {
		err = fmt.Errorf("invalid format for input: %w", err)
		return sendErr(c, http.StatusBadRequest, err)
	}
	res, err := a.IngredientUsage(ctx, r.Inputs)
	if err != nil {
		return handleErr(c, err)
	}

	type Sumresp struct {
		Sums []UsageValue `json:"sums"`
	}
	s := Sumresp{}

	for _, v := range res {
		s.Sums = append(s.Sums, v)
	}
	return c.JSON(http.StatusOK, s)

}

type UsageSummary map[IngredientID]UsageValue

func (a *API) IngredientUsage(ctx context.Context, inputRecipes []EntitySummary) (UsageSummary, error) {
	ctx, span := a.tracer.Start(ctx, "IngredientUsage")
	defer span.End()
	u := make(map[IngredientID]UsageValue)

	for _, inputRecipe := range inputRecipes {
		ru, err := a.singleIngredientUsage(ctx, inputRecipe)
		if err != nil {
			return nil, err
		}
		for k, v := range ru {
			pVal, ok := u[k]
			if !ok {
				pVal.Ing = v.Ing
			}
			pVal.Ings = append(pVal.Ings, v.Ings...)
			u[k] = pVal
		}
	}

	// sum the things
	for k, v := range u {
		if v.Sum == nil {
			v.Sum = []Amount{}
		}
		for _, i := range v.Ings {
			fa := firstAmount(i.Amounts, true)
			fb := firstAmount(i.Amounts, false)
			added := false
			if fa != nil {

				for x, a := range v.Sum {
					if a.Unit == fa.Unit {
						a.Value += (fa.Value * i.Multiplier)
						v.Sum[x] = a
						added = true
						break
					}
				}
				if !added {
					v.Sum = append(v.Sum, Amount{Unit: fa.Unit, Value: fa.Value * i.Multiplier})
				}
			} else if fb != nil {
				for x, a := range v.Sum {
					if a.Unit == fb.Unit {
						a.Value += (fb.Value * i.Multiplier)
						v.Sum[x] = a
						added = true
						break
					}
				}
				if !added {
					v.Sum = append(v.Sum, Amount{Unit: fb.Unit, Value: fb.Value * i.Multiplier})
				}
			}

		}

		u[k] = v

	}

	return u, nil

}

func (a *API) singleIngredientUsage(ctx context.Context, inputRecipe EntitySummary) (UsageSummary, error) {
	ctx, span := a.tracer.Start(ctx, "singleIngredientUsage")
	defer span.End()

	if inputRecipe.Kind != IngredientKindRecipe {
		return nil, fmt.Errorf("bad kind: %s", inputRecipe.Kind)
	}
	res, err := a.recipeById(ctx, inputRecipe.Id)
	if err != nil {
		return nil, err
	}
	inputRecipe.Multiplier = Coalesce(inputRecipe.Multiplier, 1)
	totalSum := make(UsageSummary)
	for _, s := range res.Detail.Sections {
		for _, i := range s.Ingredients {
			if i.Kind == IngredientKindIngredient {
				detail := i.Ingredient
				usage := IngredientUsage{
					Amounts:    i.Amounts,
					Multiplier: inputRecipe.Multiplier,
					RequiredBy: []EntitySummary{res.Detail.summary(inputRecipe.Multiplier)},
				}
				id := IngredientID(detail.Ingredient.Id)
				pVal, ok := totalSum[id]
				if !ok {
					pVal = UsageValue{
						Ing: EntitySummary{
							Kind: IngredientKindIngredient,
							Id:   string(id),
							Name: detail.Ingredient.Name,
						},
					}
				}
				// usage.Amounts = removeCalculatedAmounts(usage.Amounts)

				pVal.Ings = append(pVal.Ings, usage)
				totalSum[id] = pVal
			} else if i.Kind == IngredientKindRecipe {
				var subRecipeMultiplier float64 = 1.0
				for _, a := range i.Amounts {
					if a.Unit == "recipe" {
						subRecipeMultiplier = a.Value
						break
					}
				}
				subRecipeMultiplier *= inputRecipe.Multiplier
				usageSummaryForRecipe, err := a.singleIngredientUsage(ctx, EntitySummary{
					Id:         i.Recipe.Id,
					Multiplier: subRecipeMultiplier,
					Kind:       IngredientKindRecipe,
				})
				if err != nil {
					return nil, err
				}

				for ingID, eachIngUsage := range usageSummaryForRecipe {
					totalIngUsage, ok := totalSum[ingID]
					if !ok {
						totalIngUsage = UsageValue{
							Ing: eachIngUsage.Ing,
						}
					}
					for x := range eachIngUsage.Ings {
						eachIngUsage.Ings[x].RequiredBy = append(eachIngUsage.Ings[x].RequiredBy, res.Detail.summary(subRecipeMultiplier))
					}
					totalIngUsage.Ings = append(totalIngUsage.Ings, eachIngUsage.Ings...)
					totalSum[ingID] = totalIngUsage
				}
			}
		}
	}
	return totalSum, nil
}
func (a Amount) IsGram() bool {
	return a.Unit == "g" || strings.HasPrefix(a.Unit, "gr")
}
func (a Amount) IsCalculated() bool {
	if a.Source != nil && *a.Source == "calculated" {
		return true
	}
	return a.Unit == "$" || a.Unit == "kcal"
}
func firstAmount(a []Amount, grams bool) *Amount {
	for _, s := range a {
		if s.IsGram() && grams {
			return &s
		}
		if !s.IsGram() && !s.IsCalculated() && !grams {
			return &s
		}
	}
	return nil
}
func removeCalculatedAmounts(a []Amount) []Amount {
	var out []Amount
	for _, s := range a {
		if !s.IsCalculated() {
			out = append(out, s)
		}
	}
	return out
}
