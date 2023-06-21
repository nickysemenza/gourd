package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/nickysemenza/gourd/internal/db/models"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (a *API) recipeByDetailID(ctx context.Context, detailId string) (*RecipeWrapper, error) {
	ctx, span := a.tracer.Start(ctx, "recipeByDetailID")
	defer span.End()

	detail, err := models.FindRecipeDetail(ctx, a.db.DB(), detailId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to find recipe detail with detail id %s: %w", detailId, err)
		}
		return nil, err
	}
	return a.recipeByWrapperID(ctx, detail.RecipeID)
}
func (a *API) recipebyWrapperWhere(ctx context.Context, where ...QueryMod) (*RecipeWrapper, error) {
	recipe, err := models.Recipes(
		append(
			recipeQueryMods,
			where...,
		)...,
	).One(ctx, a.db.DB())
	if err != nil {
		return nil, err
	}
	return a.recipeFromModel(ctx, recipe)
}

func (a *API) recipeByWrapperID(ctx context.Context, wrapperId string) (*RecipeWrapper, error) {
	rw, err := a.recipebyWrapperWhere(ctx, Where("id = ?", wrapperId))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to find recipe with id %s: %w", wrapperId, err)
	}
	return rw, err
}

func (a *API) recipeByExactName(ctx context.Context, name string) (*RecipeWrapper, error) {
	rw, err := a.recipebyWrapperWhere(ctx, InnerJoin("recipe_details rd on rd.recipe_id = recipes.id"),
		Where("lower(name) = lower(?)", name))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to find recipe with name %s: %w", name, err)
	}
	return rw, err
}

func (a *API) recipeDetailsWhereV2(ctx context.Context, where QueryMod) (models.RecipeDetailSlice, error) {
	ctx, span := a.tracer.Start(ctx, "recipeDetailsWhereV2")
	defer span.End()
	recipes, err := models.
		RecipeDetails(
			where,
			OrderBy("version DESC"),
		).
		All(ctx, a.db.DB())
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

var recipeQueryMods = []QueryMod{
	Load(models.RecipeRels.RecipeDetails),
	// sections -> ingredients -> parents way 1
	Load(Rels(models.RecipeRels.RecipeDetails,
		models.RecipeDetailRels.RecipeSections,
		models.RecipeSectionRels.SectionRecipeSectionIngredients,
		models.RecipeSectionIngredientRels.Ingredient,
		models.IngredientRels.ParentIngredient,
	)),
	// sections -> ingredients -> parents way 2
	Load(Rels(models.RecipeRels.RecipeDetails,
		models.RecipeDetailRels.RecipeSections,
		models.RecipeSectionRels.SectionRecipeSectionIngredients,
		models.RecipeSectionIngredientRels.Ingredient,
		models.IngredientRels.ParentIngredientIngredients,
	)),

	// sections -> ingredients -> mappings
	Load(Rels(models.RecipeRels.RecipeDetails,
		models.RecipeDetailRels.RecipeSections,
		models.RecipeSectionRels.SectionRecipeSectionIngredients,
		models.RecipeSectionIngredientRels.Ingredient,
		models.IngredientRels.IngredientUnits,
	)),

	// sections -> ingredients -> recipes
	Load(Rels(models.RecipeRels.RecipeDetails,
		models.RecipeDetailRels.RecipeSections,
		models.RecipeSectionRels.SectionRecipeSectionIngredients,
		models.RecipeSectionIngredientRels.Recipe,
		models.RecipeRels.RecipeDetails,
	)),

	// sections -> instructions
	Load(Rels(models.RecipeRels.RecipeDetails,
		models.RecipeDetailRels.RecipeSections,
		models.RecipeSectionRels.SectionRecipeSectionInstructions)),
	// has images via notion recipe
	Load(Rels(models.RecipeRels.NotionRecipes,
		models.NotionRecipeRels.PageNotionImages,
		models.NotionImageRels.Image,
	)),
	Load(Rels(models.RecipeRels.MealRecipes,
		models.MealRecipeRels.Meal,
		models.MealRels.MealGphotos,
		models.MealGphotoRels.Gphoto,
		models.GphotosPhotoRels.Image,
	)),
}

func (a *API) RecipeListV2(ctx context.Context, pagination Items, mods ...QueryMod) ([]RecipeWrapper, int64, error) {
	ctx, span := a.tracer.Start(ctx, "RecipeListV2")
	defer span.End()
	filters := []QueryMod{
		Limit(pagination.Limit),
		Offset(pagination.Offset),
	}
	recipes, err := models.Recipes(
		append(
			recipeQueryMods,
			append(mods, filters...)...,
		)...,
	).
		All(ctx, a.db.DB())

	if err != nil {
		return nil, 0, err
	}
	count, err := models.Recipes(
		append(
			recipeQueryMods,
			append(mods, filters...)...,
		)...,
	).Count(ctx, a.db.DB())
	if err != nil {
		return nil, 0, err
	}
	items := []RecipeWrapper{}
	for _, recipe := range recipes {
		rw, err := a.recipeFromModel(ctx, recipe)
		if err != nil {
			return nil, 0, err
		}
		if rw != nil {
			items = append(items, *rw)
		}
	}

	return items, count, nil
}
