package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/nickysemenza/food/graph/generated"
	"github.com/nickysemenza/food/graph/model"
)

func (r *mutationResolver) CreateRecipe(ctx context.Context, input *model.NewRecipe) (*model.Recipe, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Recipes(ctx context.Context) ([]*model.Recipe, error) {
	return []*model.Recipe{&model.Recipe{}, &model.Recipe{}}, nil
}

func (r *queryResolver) Recipe(ctx context.Context, uuid string) (*model.Recipe, error) {
	res, err := r.Resolver.Manager.GetRecipe(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return &model.Recipe{UUID: res.UUID, Name: res.Name, TotalMinutes: int(res.TotalMinutes), Unit: res.Unit}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
