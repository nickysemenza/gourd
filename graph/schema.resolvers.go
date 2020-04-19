package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/davecgh/go-spew/spew"
	"github.com/nickysemenza/food/graph/generated"
	"github.com/nickysemenza/food/graph/model"
	"github.com/vektah/gqlparser/gqlerror"
)

func (r *mutationResolver) CreateRecipe(ctx context.Context, input *model.NewRecipe) (*model.Recipe, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Recipes(ctx context.Context) ([]*model.Recipe, error) {
	dbr, err := r.DB.GetRecipes(ctx)
	if err != nil {
		return nil, err
	} 
	recipes := []*model.Recipe{}
	for _, x := range dbr {
		recipes = append(recipes,fromRecipe(&x))
	}
	return recipes, nil
}

func (r *queryResolver) Recipe(ctx context.Context, uuid string) (*model.Recipe, error) {
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

func (r *recipeResolver) Sections(ctx context.Context, obj *model.Recipe) ([]*model.Section, error) {
	sections, err := r.DB.GetRecipeSections(ctx, obj.UUID)
	if err != nil {
		return nil, err
	}
	s := []*model.Section{{}, {}}
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
		i = append(i, &model.SectionIngredient{UUID: item.UUID, Grams: item.Grams.Float64, IngredientID: item.IngredientUUID.String})
	}
	return i, nil
}

func (r *sectionIngredientResolver) Info(ctx context.Context, obj *model.SectionIngredient) (*model.Ingredient, error) {
	spew.Dump(obj)
	ing, err := r.DB.GetIngredientByUUID(ctx, obj.IngredientID)
	if err != nil {
		return nil, err
	}
	if ing == nil {
		return nil, nil
	}
	return &model.Ingredient{Name: ing.Name, UUID: ing.UUID}, nil

}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Recipe returns generated.RecipeResolver implementation.
func (r *Resolver) Recipe() generated.RecipeResolver { return &recipeResolver{r} }

// Section returns generated.SectionResolver implementation.
func (r *Resolver) Section() generated.SectionResolver { return &sectionResolver{r} }

// SectionIngredient returns generated.SectionIngredientResolver implementation.
func (r *Resolver) SectionIngredient() generated.SectionIngredientResolver {
	return &sectionIngredientResolver{r}
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type recipeResolver struct{ *Resolver }
type sectionResolver struct{ *Resolver }
type sectionIngredientResolver struct{ *Resolver }
