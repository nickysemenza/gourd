package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/nickysemenza/gourd/internal/db/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/guregu/null.v4/zero"
)

func (a *API) recipeDetailFromModel(ctx context.Context, d *models.RecipeDetail) (*RecipeDetail, error) {
	ctx, span := a.tracer.Start(ctx, "recipeDetailFromModel")
	defer span.End()

	span.SetAttributes(attribute.String("name", d.Name))
	span.AddEvent("model", trace.WithAttributes(attribute.String("model", spew.Sdump(d))))
	sections := make([]RecipeSection, 0)
	for _, section := range d.R.RecipeSections {
		rs := RecipeSection{
			Id:           section.ID,
			Ingredients:  []SectionIngredient{},
			Instructions: []SectionInstruction{},
		}
		if section.DurationTimerange.Valid {
			err := json.Unmarshal([]byte(section.DurationTimerange.JSON), &rs.Duration)
			if err != nil {
				return nil, err
			}
		}
		for _, instruction := range section.R.SectionRecipeSectionInstructions {
			rs.Instructions = append(rs.Instructions, SectionInstruction{
				Id:          instruction.ID,
				Instruction: instruction.Instruction.String,
			})
		}
		for _, ingredient := range section.R.SectionRecipeSectionIngredients {
			si := SectionIngredient{
				Id:        ingredient.ID,
				Adjective: ingredient.Adjective.Ptr(),
				Amounts:   []Amount{},
				Optional:  ingredient.Optional.Ptr(),
				Original:  ingredient.Original.Ptr(),
			}

			if ingredient.SubForIngredientID.Valid {
				l(ctx).Warnf("SubForIngredientID not implemented for %s", ingredient.ID)
			}
			rawAmounts := []Amount{}
			err := ingredient.Amounts.Unmarshal(&rawAmounts)
			if err != nil {
				return nil, err
			}
			for _, amt := range rawAmounts {
				si.Amounts = append(si.Amounts, Amount{
					Unit:   amt.Unit,
					Value:  amt.Value,
					Source: zero.StringFrom("db").Ptr(),
				})
			}
			switch {
			case ingredient.RecipeID.Valid:
				si.Kind = IngredientKindRecipe
				foo, err := a.recipeFromModel(ctx, ingredient.R.Recipe)
				if err != nil {
					return nil, err
				}
				si.Recipe = &foo.Detail
			case ingredient.IngredientID.Valid:
				si.Kind = IngredientKindIngredient
				var err error
				si.Ingredient, err = a.ingredientFromModel(ctx, ingredient.R.Ingredient, true, false)
				if err != nil {
					return nil, err
				}
				if err := a.enhanceMulti(ctx, &si); err != nil {
					return nil, err
				}

			default:
				return nil, fmt.Errorf("ingredient is not valid")
			}

			rs.Ingredients = append(rs.Ingredients, si)
		}
		sections = append(sections, rs)

	}

	rd := RecipeDetail{
		CreatedAt:       d.CreatedAt,
		Id:              d.ID,
		IsLatestVersion: d.IsLatestVersion.Bool,
		Name:            d.Name,
		Quantity:        d.Quantity.Int,
		Sections:        sections,
		Servings:        d.Servings.Ptr(),
		Sources:         []RecipeSource{},
		Tags:            d.Tags,
		Unit:            d.Unit.String,
		Version:         d.Version,
	}
	if rd.Tags == nil {
		rd.Tags = []string{}
	}
	if d.Source.Valid {
		if err := d.Source.Unmarshal(&rd.Sources); err != nil {
			return nil, err
		}
	}
	return &rd, nil

}
func (a *API) recipeFromModel(ctx context.Context, recipe *models.Recipe) (*RecipeWrapper, error) {
	ctx, span := a.tracer.Start(ctx, "recipeFromModel")
	defer span.End()
	if recipe == nil || len(recipe.R.RecipeDetails) == 0 {
		return nil, nil
	}

	rw := RecipeWrapper{
		Id: recipe.ID,
	}
	other := []RecipeDetail{}
	for _, d := range recipe.R.RecipeDetails {
		rd, err := a.recipeDetailFromModel(ctx, d)
		if err != nil {
			return nil, err
		}
		if rd.IsLatestVersion {
			rw.Detail = *rd
			if d.DeletedAt.Valid {
				return nil, nil
			}
		} else {
			other = append(other, *rd)
		}
	}
	rw.OtherVersions = &other

	gp := models.GphotosPhotoSlice{}
	linkedMeals := []Meal{}
	for _, m := range recipe.R.MealRecipes {
		for _, x := range m.R.Meal.R.MealGphotos {
			gp = append(gp, x.R.Gphoto)
		}
		linkedMeals = append(linkedMeals, Meal{
			Id:    m.R.Meal.ID,
			Name:  m.R.Meal.Name,
			AteAt: m.R.Meal.AteAt,
		})
	}
	images, err := a.imagesFromModel(ctx, recipe.R.NotionRecipes, gp)
	if err != nil {
		return nil, err
	}
	rw.LinkedPhotos = &images

	rw.LinkedMeals = &linkedMeals

	return &rw, nil
}
