package api

import (
	"context"
	"strings"
)

type IngredientID string
type IngredientUsage struct {
	Amounts    []Amount
	RequiredBy []ResSummary // chain of dependencies
	Multiplier float64
}
type ResSummary struct {
	Name     string
	DetailID string
}

func (r *RecipeDetail) summary() ResSummary {
	return ResSummary{
		Name:     r.Name,
		DetailID: r.Id,
	}
}

type bar struct {
	RecipeId   string
	Multiplier float64
}

type UsageValue struct {
	Ings         []IngredientUsage
	Sum          []Amount
	IngredientID IngredientID // redundant
	Name         string
}

type UsageSummary map[IngredientID]UsageValue

func (a *API) IngredientUsage(ctx context.Context, details []bar) (UsageSummary, error) {
	u := make(map[IngredientID]UsageValue)

	for _, d := range details {
		ru, err := a.singleIngredientUsage(ctx, d)
		if err != nil {
			return nil, err
		}
		for k, v := range ru {
			pVal, ok := u[k]
			if !ok {
				pVal.IngredientID = v.IngredientID
				pVal.Name = v.Name
			}
			pVal.Ings = append(pVal.Ings, v.Ings...)
			u[k] = pVal
		}
	}

	// sum the things
	for k, v := range u {

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
					v.Sum = append(v.Sum, Amount{Unit: fa.Unit, Value: fa.Value * i.Multiplier})
				}
			}

		}

		u[k] = v

	}

	return u, nil

}

func (a *API) singleIngredientUsage(ctx context.Context, d bar) (UsageSummary, error) {
	res, err := a.recipeById(ctx, d.RecipeId)
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
					RequiredBy: []ResSummary{res.Detail.summary()},
				}
				id := IngredientID(detail.Ingredient.Id)
				pVal, ok := u[id]
				if !ok {
					pVal = UsageValue{
						IngredientID: id,
						Name:         detail.Ingredient.Name,
					}
				}
				pVal.Ings = append(pVal.Ings, usage)
				u[id] = pVal
			} else if i.Kind == IngredientKindRecipe {
				ru, err := a.singleIngredientUsage(ctx, bar{
					RecipeId:   i.Recipe.Id,
					Multiplier: 1, // todo
				})
				if err != nil {
					return nil, err
				}

				for k, v := range ru {
					pVal, ok := u[k]
					if !ok {
						pVal = UsageValue{
							IngredientID: v.IngredientID,
							Name:         v.Name,
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
