package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nickysemenza/food/db"
	"github.com/nickysemenza/food/graph/generated"
	"github.com/nickysemenza/food/graph/model"
	"github.com/vektah/gqlparser/gqlerror"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
	"gopkg.in/guregu/null.v3/zero"
)

func (r *foodResolver) Category(ctx context.Context, obj *model.Food) (*model.FoodCategory, error) {
	if !obj.CategoryID.Valid {
		return nil, nil
	}
	return r.DB.GetCategory(ctx, obj.CategoryID.Int64)
}

func (r *foodResolver) Nutrients(ctx context.Context, obj *model.Food) ([]*model.FoodNutrient, error) {
	return r.DB.GetFoodNutrients(ctx, obj.FdcID)
}

func (r *foodNutrientResolver) Nutrient(ctx context.Context, obj *model.FoodNutrient) (*model.Nutrient, error) {
	return r.DB.GetNutrient(ctx, obj.NutrientID)
}

func (r *foodNutrientResolver) DataPoints(ctx context.Context, obj *model.FoodNutrient) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *ingredientResolver) Recipes(ctx context.Context, obj *model.Ingredient) ([]*model.Recipe, error) {
	dbr, err := r.DB.GetRecipesWithIngredient(ctx, obj.UUID)
	if err != nil {
		return nil, err
	}
	recipes := []*model.Recipe{}
	for _, x := range dbr {
		each := x
		recipes = append(recipes, fromRecipe(&each))
	}
	return recipes, nil
}

func (r *mutationResolver) CreateRecipe(ctx context.Context, recipe *model.NewRecipe) (*model.Recipe, error) {
	uuid, err := r.DB.InsertRecipe(ctx, &db.Recipe{Name: recipe.Name})
	if err != nil {
		return nil, err
	}
	return r.Query().Recipe(ctx, uuid)
}

func (r *mutationResolver) UpdateRecipe(ctx context.Context, recipe *model.RecipeInput) (*model.Recipe, error) {
	uuid := recipe.UUID
	res, err := r.Resolver.DB.GetRecipeByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if res == nil {
		graphql.AddError(ctx, gqlerror.Errorf("no recipe found with uuid %s", uuid))
		return nil, nil
	}

	dbr := &db.Recipe{
		UUID: uuid,
		Name: recipe.Name,
	}
	if recipe.TotalMinutes != nil {
		dbr.TotalMinutes = zero.IntFrom(int64(*recipe.TotalMinutes))
	}
	for _, s := range recipe.Sections {
		dbs := db.Section{
			// RecipeUUID: uuid,
			Minutes: zero.IntFrom(int64(s.Minutes)),
		}
		for _, i := range s.Instructions {
			dbs.Instructions = append(dbs.Instructions, db.SectionInstruction{
				Instruction: i.Instruction,
			})
		}
		for _, i := range s.Ingredients {
			dbsi := db.SectionIngredient{
				// Name:      i.Name,
				Grams:     zero.FloatFrom(i.Grams),
				Amount:    zero.FloatFromPtr(i.Amount),
				Unit:      zero.StringFromPtr(i.Unit),
				Adjective: zero.StringFromPtr(i.Adjective),
				Optional:  zero.BoolFromPtr(i.Optional),
				//TODO(nicky) kind:
			}
			switch i.Kind {
			case model.SectionIngredientKindIngredient:
				dbsi.IngredientUUID = zero.StringFrom(i.InfoUUID)
			case model.SectionIngredientKindRecipe:
				dbsi.RecipeUUID = zero.StringFrom(i.InfoUUID)
			default:
				return nil, fmt.Errorf("invalid kind: %s", i.Kind)
			}
			dbs.Ingredients = append(dbs.Ingredients, dbsi)
		}
		dbr.Sections = append(dbr.Sections, dbs)
	}

	if err := r.DB.UpdateRecipe(ctx, dbr); err != nil {
		return nil, err
	}
	return r.Query().Recipe(ctx, uuid)
}

func (r *mutationResolver) CreateIngredient(ctx context.Context, name string) (*model.Ingredient, error) {
	ing, err := r.DB.IngredientByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return fromIngredient(ing), nil
}

func (r *queryResolver) Recipes(ctx context.Context, searchQuery string) ([]*model.Recipe, error) {
	dbr, err := r.DB.GetRecipes(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	recipes := []*model.Recipe{}
	for _, x := range dbr {
		each := x
		recipes = append(recipes, fromRecipe(&each))
	}
	return recipes, nil
}

func (r *queryResolver) Recipe(ctx context.Context, uuid string) (*model.Recipe, error) {
	tr := global.Tracer("graph")
	ctx, span := tr.Start(ctx, "Recipe")
	defer span.End()
	span.SetAttributes(core.KeyValue{Key: "uuid", Value: core.String(uuid)})
	res, err := r.Resolver.DB.GetRecipeByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if res == nil {
		graphql.AddError(ctx, gqlerror.Errorf("no recipe found with uuid %s", uuid))
		return nil, nil
	}

	return fromRecipe(res), nil
}

func (r *queryResolver) Ingredients(ctx context.Context, searchQuery string) ([]*model.Ingredient, error) {
	dbr, err := r.DB.GetIngredients(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	ingredients := []*model.Ingredient{}
	for _, x := range dbr {
		each := x
		ingredients = append(ingredients, fromIngredient(&each))
	}
	return ingredients, nil
}

func (r *queryResolver) Ingredient(ctx context.Context, uuid string) (*model.Ingredient, error) {
	ing, err := r.DB.GetIngredientByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return fromIngredient(ing), nil
}

func (r *queryResolver) Food(ctx context.Context, fdcID int) (*model.Food, error) {
	return r.DB.GetFood(ctx, fdcID)
}

func (r *queryResolver) Foods(ctx context.Context, searchQuery string, dataType *model.FoodDataType, foodCategoryID *int) ([]*model.Food, error) {
	return r.DB.SearchFoods(ctx, searchQuery, dataType, foodCategoryID)
}

func (r *recipeResolver) Sections(ctx context.Context, obj *model.Recipe) ([]*model.Section, error) {
	sections, err := r.DB.GetRecipeSections(ctx, obj.UUID)
	if err != nil {
		return nil, err
	}
	s := []*model.Section{}
	for _, dbs := range sections {
		s = append(s, &model.Section{UUID: dbs.UUID, RecipeUUID: dbs.RecipeUUID, Minutes: int(dbs.Minutes.Int64)})
	}
	return s, nil
}

func (r *sectionResolver) Instructions(ctx context.Context, obj *model.Section) ([]*model.SectionInstruction, error) {
	res, err := r.DB.GetSectionInstructions(ctx, obj.UUID)
	if err != nil {
		return nil, err
	}
	i := []*model.SectionInstruction{}
	for _, item := range res {
		i = append(i, &model.SectionInstruction{UUID: item.UUID, Instruction: item.Instruction})
	}
	return i, nil
}

func (r *sectionResolver) Ingredients(ctx context.Context, obj *model.Section) ([]*model.SectionIngredient, error) {
	res, err := r.DB.GetSectionIngredients(ctx, obj.UUID)
	if err != nil {
		return nil, err
	}
	i := []*model.SectionIngredient{}
	for _, item := range res {
		var info model.IngredientInfo
		var kind model.SectionIngredientKind
		if item.RecipeUUID.Valid {
			info, err = r.Resolver.Query().Recipe(ctx, item.RecipeUUID.String)
			if err != nil {
				return nil, err
			}
			kind = model.SectionIngredientKindRecipe
		} else if item.IngredientUUID.Valid {
			// info = model.Ingredient{Name: "todo"}
			info, err = r.Resolver.Query().Ingredient(ctx, item.IngredientUUID.String)
			if err != nil {
				return nil, err
			}
			kind = model.SectionIngredientKindIngredient
		}
		i = append(i, &model.SectionIngredient{
			UUID:      item.UUID,
			Grams:     item.Grams.Float64,
			Amount:    item.Amount.Float64,
			Unit:      item.Unit.String,
			Adjective: item.Adjective.String,
			Optional:  item.Optional.Bool,
			Info:      info,
			Kind:      kind,
			// IngredientID: item.IngredientUUID.String,
			// RecipeID:     item.RecipeUUID.String,
		})
	}
	return i, nil
}

// Food returns generated.FoodResolver implementation.
func (r *Resolver) Food() generated.FoodResolver { return &foodResolver{r} }

// FoodNutrient returns generated.FoodNutrientResolver implementation.
func (r *Resolver) FoodNutrient() generated.FoodNutrientResolver { return &foodNutrientResolver{r} }

// Ingredient returns generated.IngredientResolver implementation.
func (r *Resolver) Ingredient() generated.IngredientResolver { return &ingredientResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Recipe returns generated.RecipeResolver implementation.
func (r *Resolver) Recipe() generated.RecipeResolver { return &recipeResolver{r} }

// Section returns generated.SectionResolver implementation.
func (r *Resolver) Section() generated.SectionResolver { return &sectionResolver{r} }

type foodResolver struct{ *Resolver }
type foodNutrientResolver struct{ *Resolver }
type ingredientResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type recipeResolver struct{ *Resolver }
type sectionResolver struct{ *Resolver }
