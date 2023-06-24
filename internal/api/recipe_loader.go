package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/nickysemenza/gourd/internal/db/models"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.mitsakis.org/workerpool"
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
	return a.recipeFromModel(ctx, recipe, true)
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

	recipes, count, err := countAndQuery[models.RecipeSlice](ctx, a.db.DB(), models.Recipes, qmWithPagination(recipeQueryMods, pagination, mods...)...)
	if err != nil {
		return nil, 0, err
	}

	p, err := workerpool.NewPoolWithResults(
		8,
		func(job workerpool.Job[*models.Recipe], workerID int) (*RecipeWrapper, error) {
			return a.recipeFromModel(ctx, job.Payload, false)
		})
	if err != nil {
		return nil, 0, err
	}
	go func() {
		for _, nRecipe := range recipes {
			p.Submit(nRecipe)
		}
		p.StopAndWait()
	}()

	items := []RecipeWrapper{}

	for result := range p.Results {
		if result.Error != nil {
			l(ctx).Error(result.Error)
			return nil, 0, result.Error
		} else {
			res := result.Value
			items = append(items, *res)

		}
	}

	return items, count, nil
}
