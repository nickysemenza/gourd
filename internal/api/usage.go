package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/maps"
)

// IngredientID is an alias
type IngredientID string

func (r *RecipeDetail) summary(multiplier float64) EntitySummary {
	return EntitySummary{
		Name:       r.Name,
		Id:         r.Id,
		Kind:       IngredientKindRecipe,
		Multiplier: multiplier,
	}
}

// SumRecipes sums the amounts of ingredients in the recipes
func (a *API) SumRecipes(c echo.Context) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "SumRecipes")
	defer span.End()

	var r SumRecipesJSONRequestBody
	if err := c.Bind(&r); err != nil {
		err = fmt.Errorf("invalid format for input: %w", err)
		return sendErr(c, http.StatusBadRequest, err)
	}
	ingredientSums, err := a.ingredientUsage(ctx, r.Inputs...)
	if err != nil {
		return handleErr(c, err)
	}

	s := SumsResponse{
		Sums: maps.Values(ingredientSums),
	}

	for _, eachRecipe := range r.Inputs {
		recipeSpecific, err := a.ingredientUsage(ctx, eachRecipe)
		if err != nil {
			return handleErr(c, err)
		}
		if s.ByRecipe == nil {
			s.ByRecipe = make(map[string][]UsageValue)
		}
		s.ByRecipe[eachRecipe.Id] = maps.Values(recipeSpecific)
	}
	return c.JSON(http.StatusOK, s)

}

// UsageSummary is a map of ingredient ids to their value
type UsageSummary map[IngredientID]UsageValue

func (u UsageSummary) add(ingId IngredientID, usageVal UsageValue) {
	pVal, ok := u[ingId]
	if !ok {
		pVal.Meta = usageVal.Meta
	}
	pVal.Ings = append(pVal.Ings, usageVal.Ings...)
	u[ingId] = pVal
}

// ingredientUsage calculate usage for provided recipes
func (a *API) ingredientUsage(ctx context.Context, inputRecipes ...EntitySummary) (UsageSummary, error) {
	ctx, span := a.tracer.Start(ctx, "IngredientUsage")
	defer span.End()
	summary := make(UsageSummary)

	for _, inputRecipe := range inputRecipes {
		ru, err := a.singleIngredientUsage(ctx, inputRecipe)
		if err != nil {
			return nil, err
		}
		for ingID, usageVal := range ru {
			summary.add(ingID, usageVal)
		}
	}

	l(ctx).Infof("summary has %d ingredients", len(summary))
	// sum the things
	for ingredientID, v := range summary {
		if v.Sum == nil {
			v.Sum = []Amount{}
		}
		for _, i := range v.Ings {
			// for each of the 'usages' (times they appear in recipes) of the current ingredient,
			// try this with both grams and then non-grams:
			//	 iterate through each of the current sums for the
			//	 ingredient, and add the current usage to the sum if the unit matches,
			// 	 otherwise create a new one.
			for _, amt := range []*Amount{firstAmount(i.Amounts, true), firstAmount(i.Amounts, false)} {
				if amt == nil {
					continue
				}
				amt2, err := a.NormalizeAmount(ctx, *amt)
				if err != nil {
					return nil, err
				}
				added := false
				for x, each := range v.Sum {

					if each.Unit == amt2.Unit {
						each.Value += (amt2.Value * i.Multiplier)
						v.Sum[x] = each
						added = true
						break
					}
				}
				if !added {
					v.Sum = append(v.Sum, Amount{Unit: amt2.Unit, Value: amt2.Value * i.Multiplier})
				}
				break
			}

		}

		summary[ingredientID] = v

	}

	return summary, nil

}

func (a *API) singleIngredientUsage(ctx context.Context, inputRecipe EntitySummary) (UsageSummary, error) {
	ctx, span := a.tracer.Start(ctx, "singleIngredientUsage")
	defer span.End()

	if inputRecipe.Kind != IngredientKindRecipe {
		return nil, fmt.Errorf("bad kind: %s", inputRecipe.Kind)
	}
	recipe, err := a.recipeById(ctx, inputRecipe.Id)
	if err != nil {
		return nil, err
	}
	l(ctx).Infof("singleUsage for %s (%v)", recipe.Detail.Name, inputRecipe)
	inputRecipe.Multiplier = Coalesce(inputRecipe.Multiplier, 1)
	totalSum := make(UsageSummary)
	for _, section := range recipe.Detail.Sections {
		for _, si := range section.Ingredients {
			// for each of the ingredient line items in the recipe
			switch si.Kind {
			case IngredientKindIngredient:
				ing := si.Ingredient.Ingredient
				si.Amounts = removeCalculatedAmounts(si.Amounts)

				totalSum.add(IngredientID(ing.Id), UsageValue{
					Meta: EntitySummary{
						Kind: IngredientKindIngredient,
						Id:   ing.Id,
						Name: ing.Name,
					},
					Ings: []IngredientUsage{{
						Amounts:    si.Amounts,
						Multiplier: inputRecipe.Multiplier,
						RequiredBy: []EntitySummary{recipe.Detail.summary(inputRecipe.Multiplier)},
					}},
				})

			case IngredientKindRecipe:
				subRecipeMultiplier := 1.0
				for _, a := range si.Amounts {
					if a.Unit == "recipe" {
						// special case adjusting multiplier to 0.5 from "1/2 recipe foo"
						subRecipeMultiplier = a.Value
						break
					}
				}
				subRecipeMultiplier *= inputRecipe.Multiplier

				// recurse into the sub-recipe
				usageSummaryForRecipe, err := a.singleIngredientUsage(ctx, EntitySummary{
					Id:         si.Recipe.Id,
					Multiplier: subRecipeMultiplier,
					Kind:       IngredientKindRecipe,
				})
				if err != nil {
					return nil, err
				}

				for ingID, eachIngUsage := range usageSummaryForRecipe {

					for x := range eachIngUsage.Ings {
						// prepend bc it's like a call stack
						eachIngUsage.Ings[x].RequiredBy = append(
							[]EntitySummary{recipe.Detail.summary(subRecipeMultiplier)},
							eachIngUsage.Ings[x].RequiredBy...,
						)
					}
					totalSum.add(ingID, eachIngUsage)

				}
			}
		}
	}
	return totalSum, nil
}

// IsGram checks if the unit is gram
func (a Amount) IsGram() bool {
	return a.Unit == "g" || strings.HasPrefix(a.Unit, "gr")
}

// IsMoneyKCal checks if the unit is money or kcalories
func (a Amount) IsMoneyKCal() bool {
	return a.Unit == "$" || a.Unit == "kcal"
}
func firstAmount(a []Amount, grams bool) *Amount {
	for _, s := range a {
		if s.IsGram() && grams {
			return &s
		}
		if !s.IsGram() && !s.IsMoneyKCal() && !grams {
			return &s
		}
	}
	return nil
}

func removeCalculatedAmounts(a []Amount) []Amount {
	out := make([]Amount, 0, len(a))
	for _, s := range a {
		if !s.IsMoneyKCal() {
			out = append(out, s)
		}
	}
	return out
}
