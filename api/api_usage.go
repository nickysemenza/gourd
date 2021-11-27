package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type IngredientID string

func (r *RecipeDetail) summary() EntitySummary {
	return EntitySummary{
		Name: r.Name,
		Id:   r.Id,
		Kind: IngredientKindRecipe,
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

func (a *API) IngredientUsage(ctx context.Context, details []EntitySummary) (UsageSummary, error) {
	ctx, span := a.tracer.Start(ctx, "IngredientUsage")
	defer span.End()
	u := make(map[IngredientID]UsageValue)

	for _, d := range details {
		ru, err := a.singleIngredientUsage(ctx, d)
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

func (a *API) singleIngredientUsage(ctx context.Context, d EntitySummary) (UsageSummary, error) {
	ctx, span := a.tracer.Start(ctx, "singleIngredientUsage")
	defer span.End()

	if d.Kind != IngredientKindRecipe {
		return nil, fmt.Errorf("bad kind: %s", d.Kind)
	}
	res, err := a.recipeById(ctx, d.Id)
	if err != nil {
		return nil, err
	}
	if d.Multiplier == 0 {
		d.Multiplier = 1
	}
	u := make(map[IngredientID]UsageValue)
	for _, s := range res.Detail.Sections {
		for _, i := range s.Ingredients {
			if i.Kind == IngredientKindIngredient {
				detail := i.Ingredient
				usage := IngredientUsage{
					Amounts:    i.Amounts,
					Multiplier: d.Multiplier,
					RequiredBy: []EntitySummary{res.Detail.summary()},
				}
				id := IngredientID(detail.Ingredient.Id)
				pVal, ok := u[id]
				if !ok {
					pVal = UsageValue{
						Ing: EntitySummary{
							Kind: IngredientKindIngredient,
							Id:   string(id),
							Name: detail.Ingredient.Name,
						},
					}
				}
				pVal.Ings = append(pVal.Ings, usage)
				u[id] = pVal
			} else if i.Kind == IngredientKindRecipe {
				ru, err := a.singleIngredientUsage(ctx, EntitySummary{
					Id:         i.Recipe.Id,
					Multiplier: 1, // todo
					Kind:       IngredientKindRecipe,
				})
				if err != nil {
					return nil, err
				}

				for k, v := range ru {
					pVal, ok := u[k]
					if !ok {
						pVal = UsageValue{
							Ing: v.Ing,
						}
					}
					for x := range v.Ings {
						v.Ings[x].RequiredBy = append(v.Ings[x].RequiredBy, res.Detail.summary())
					}
					pVal.Ings = append(pVal.Ings, v.Ings...)
					u[k] = pVal
				}
			}
		}
	}
	return u, nil
}
func (a Amount) IsGram() bool {
	return a.Unit == "g" || strings.HasPrefix(a.Unit, "gr")
}
func firstAmount(a []Amount, grams bool) *Amount {
	for _, s := range a {
		if s.IsGram() && grams {
			return &s
		}
		if !s.IsGram() && s.Unit != "$" && s.Unit != "kcal" && !grams {
			return &s
		}
	}
	return nil
}
