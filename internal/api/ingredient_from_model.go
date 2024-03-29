package api

import (
	"context"
	"fmt"

	"github.com/nickysemenza/gourd/internal/db/models"
	"go.opentelemetry.io/otel/attribute"
	"gopkg.in/guregu/null.v4/zero"
)

func (a *API) ingredientFromModel(ctx context.Context, ingredient *models.Ingredient, withDetail, recurseIngredients, recurseRecipes bool) (*IngredientWrapper, error) {
	ctx, span := a.tracer.Start(ctx, "ingredientFromModel")
	defer span.End()
	if ingredient == nil {
		l(ctx).Warnf("ingredientFromModel called with nil ingredient")
		return nil, nil
	}
	span.SetAttributes(attribute.String("ingredient-name", ingredient.Name))
	detail := &IngredientWrapper{
		Ingredient: Ingredient{
			Id:     ingredient.ID,
			Name:   ingredient.Name,
			FdcId:  ingredient.FDCID.Ptr(),
			Parent: ingredient.ParentIngredientID.Ptr(),
		},
		Recipes:      []RecipeDetail{},
		UnitMappings: []UnitMapping{},
	}

	l(ctx).Infof("ingredientFromModel %s %s: recurseRecipes=%v recurseRecipes=%v", ingredient.Name, ingredient.ID, recurseIngredients, recurseRecipes)

	var children []IngredientWrapper
	if recurseIngredients {
		for _, x := range ingredient.R.ParentIngredientIngredients {
			l(ctx).Infof("checking child ingredient %s", x.ID)
			res, err := a.ingredientFromModel(ctx, x, withDetail, false, false)
			if err != nil {
				return nil, err
			}
			if res != nil {
				children = append(children, *res)
			}
		}
		if parent := ingredient.R.ParentIngredient; parent != nil {
			l(ctx).Infof("checking parent ingredient %s", parent.ID)
			res, err := a.ingredientFromModel(ctx, parent, withDetail, false, false)
			if err != nil {
				return nil, err
			}
			if res != nil {
				children = append(children, *res)
			}
		}
	}
	if recurseRecipes {
		for _, s := range ingredient.R.RecipeSectionIngredients {
			if s.R.Section != nil && s.R.Section.R.RecipeDetail != nil && s.R.Section.R.RecipeDetail.ID != "" {
				wrapper, err := a.recipeByDetailID(ctx, s.R.Section.R.RecipeDetail.ID)
				if err != nil {
					panic(err)
				}
				detail.Recipes = append(detail.Recipes, wrapper.Detail)
			}

		}
	}
	detail.Children = &children

	for _, m := range ingredient.R.IngredientUnits {
		fB, _ := m.AmountA.Float64()
		fA, _ := m.AmountB.Float64()
		detail.UnitMappings = append(detail.UnitMappings, UnitMapping{
			Amount{Unit: m.UnitA, Value: fA},
			Amount{Unit: m.UnitB, Value: fB},
			zero.StringFrom(fmt.Sprintf("%s (%s)", m.Source, ingredient.ID)).Ptr(),
		})
	}

	if ingredient.FDCID.Valid && withDetail {
		err := a.enhanceWithFDC(ctx, int64(ingredient.FDCID.Int), detail)
		if err != nil {
			return nil, fmt.Errorf("enhanceWithFDC: %w", err)
		}

	}

	return detail, nil
}
