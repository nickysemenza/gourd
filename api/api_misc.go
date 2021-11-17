package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/models"
	"github.com/sirupsen/logrus"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func ingredientFromModel(ingredient *models.Ingredient) *IngredientDetail {
	if ingredient == nil {
		return nil
	}
	i := IngredientDetail{
		Name: ingredient.Name,
	}
	return &i
}
func (a *API) recipeFromModel(recipe *models.Recipe) *RecipeWrapper {
	if recipe == nil {
		return nil
	}
	sections := []RecipeSection{}

	for _, section := range recipe.R.RecipeDetail.R.RecipeSections {
		s := RecipeSection{
			Id: section.ID,
			// Duration: ,
		}
		for _, instruction := range section.R.SectionRecipeSectionInstructions {
			s.Instructions = append(s.Instructions, SectionInstruction{
				Id:          instruction.ID,
				Instruction: instruction.Instruction.String,
			})
		}
		for _, ingredient := range section.R.SectionRecipeSectionIngredients {
			si := SectionIngredient{
				Id:        ingredient.ID,
				Adjective: ingredient.Adjective.Ptr(),
				// Instruction: ingredient.Instruction.String,
			}
			if ingredient.Ingredient.Valid {
				si.Kind = IngredientKindIngredient
				si.Ingredient = ingredientFromModel(ingredient.R.RecipeSectionIngredientIngredient)
			} else {
				si.Kind = IngredientKindRecipe
				// si.Recipe = &recipeFromDB(ingredient.R.RecipeSectionIngredientRecipe).Detail
			}
			s.Ingredients = append(s.Ingredients, si)
		}
		sections = append(sections, s)

	}

	rd := RecipeDetail{
		Id:        recipe.R.RecipeDetail.ID,
		CreatedAt: recipe.R.RecipeDetail.CreatedAt,
		Name:      recipe.R.RecipeDetail.Name,
		Quantity:  int64(recipe.R.RecipeDetail.Quantity.Int),
		// Servings:  int64(recipe.R.RecipeDetail.Servings.Int),
		// Sources:   recipe.R.RecipeDetail.Source,
		Unit:            recipe.R.RecipeDetail.Unit.String,
		Version:         int64(recipe.R.RecipeDetail.Version),
		Sections:        sections,
		IsLatestVersion: recipe.R.RecipeDetail.IsLatestVersion.Bool,
	}

	rw := RecipeWrapper{
		Id:     recipe.ID,
		Detail: rd,
	}

	items := []Photo{}
	for _, notionRecipe := range recipe.R.NotionRecipes {
		for _, notionImage := range notionRecipe.R.PageNotionImages {
			url := a.ImageStore.GetImageURL(context.Background(), notionImage.Image)
			items = append(items, Photo{
				Id:       notionImage.BlockID,
				Created:  notionImage.LastSeen,
				BlurHash: &notionImage.R.NotionImageImage.BlurHash,
				Width:    300,
				Height:   400,
				BaseUrl:  url,
				Source:   PhotoSourceNotion,
			})
		}
	}
	rw.LinkedPhotos = &items

	return &rw
}
func (a *API) Misc(c echo.Context) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "Misc")
	defer span.End()

	recipes, err := models.Recipes(
		Load(Rels(models.RecipeRels.RecipeDetail,
			models.RecipeDetailRels.RecipeSections,
			models.RecipeSectionRels.SectionRecipeSectionIngredients,
			models.RecipeSectionIngredientRels.RecipeSectionIngredientIngredient,
			// models.RecipeSectionIngredientRels.RecipeSectionIngredientRecipe,
		)),
		Load(Rels(models.RecipeRels.RecipeDetail,
			models.RecipeDetailRels.RecipeSections,
			models.RecipeSectionRels.SectionRecipeSectionInstructions)),
		Load(Rels(models.RecipeRels.NotionRecipes,
			models.NotionRecipeRels.PageNotionImages,
			models.NotionImageRels.NotionImageImage,
		)),
	).All(ctx, a.db.DB())
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	// spew.Dump(recipes)
	items := []RecipeWrapper{}
	for _, recipe := range recipes {
		rw := a.recipeFromModel(recipe)
		items = append(items, *rw)
	}
	// s := spew.Sdump(recipes)
	// // s = strings.ReplaceAll(s, "\n", "<br/>")
	// s = fmt.Sprintf("<html>%s</html>", s)
	return c.JSON(http.StatusOK, items)
}
