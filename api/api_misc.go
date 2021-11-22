package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/models"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (a *API) ingredientFromModel(_ context.Context, ingredient *models.Ingredient) (*IngredientDetail, error) {
	if ingredient == nil {
		return nil, nil
	}
	i := IngredientDetail{
		Name: ingredient.Name,
	}
	return &i, nil
}
func (a *API) recipeFromModel(ctx context.Context, recipe *models.Recipe) (*RecipeWrapper, error) {
	if recipe == nil {
		return nil, nil
	}
	sections := []RecipeSection{}

	// d := recipe.R.RecipeDetails[0]
	rw := RecipeWrapper{
		Id: recipe.ID,
		// Detail: rd,
	}
	other := []RecipeDetail{}

	for _, d := range recipe.R.RecipeDetails {

		for _, section := range d.R.RecipeSections {
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
					var err error
					si.Ingredient, err = a.ingredientFromModel(ctx, ingredient.R.RecipeSectionIngredientIngredient)
					if err != nil {
						return nil, err
					}
				} else {
					si.Kind = IngredientKindRecipe
					// si.Recipe = &recipeFromDB(ingredient.R.RecipeSectionIngredientRecipe).Detail
				}
				s.Ingredients = append(s.Ingredients, si)
			}
			sections = append(sections, s)

		}

		// if detail.IsLatestVersion.Bool {

		// 	others, err := models.Recipes(
		// 		Where("recipes.id = ?", recipe.ID),
		// 		Load(models.RecipeRels.RecipeDetails,
		// 			Where("recipe_details.is_latest_version = ?", false)),
		// 	).
		// 		All(ctx, a.db.DB())
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	for _, o := range others {
		// 		r, err := a.recipeFromModel(ctx, o)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 		other = append(other, r.Detail)
		// 	}
		// }

		rd := RecipeDetail{
			Id:        d.ID,
			CreatedAt: d.CreatedAt,
			Name:      d.Name,
			Quantity:  int64(d.Quantity.Int),
			// Servings:  int64(d.Servings.Int),
			// Sources:   d.Source,
			Unit:            d.Unit.String,
			Version:         int64(d.Version),
			Sections:        sections,
			IsLatestVersion: d.IsLatestVersion.Bool,
			// OtherVersions:   &other,
		}
		if rd.IsLatestVersion {
			rw.Detail = rd
		} else {
			other = append(other, rd)
		}
	}
	rw.Detail.OtherVersions = &other

	items := []Photo{}
	for _, notionRecipe := range recipe.R.NotionRecipes {
		for _, notionImage := range notionRecipe.R.PageNotionImages {
			url := a.ImageStore.GetImageURL(ctx, notionImage.Image)
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

	return &rw, nil
}
func (a *API) RecipeListV2(ctx context.Context, limit, offset uint64) ([]RecipeWrapper, error) {
	recipes, err := models.Recipes(
		// Load(models.RecipeRels.RecipeDetails, Where("recipe_details.is_latest_version = ?", true)),
		// has many sections, has many ingredients, which can be ingredients or recipes
		Load(Rels(models.RecipeRels.RecipeDetails,
			models.RecipeDetailRels.RecipeSections,
			models.RecipeSectionRels.SectionRecipeSectionIngredients,
			models.RecipeSectionIngredientRels.RecipeSectionIngredientIngredient,
			// TODO
			// models.RecipeSectionIngredientRels.RecipeSectionIngredientRecipe,
		)),
		// has many sections, has many instructions
		Load(Rels(models.RecipeRels.RecipeDetails,
			models.RecipeDetailRels.RecipeSections,
			models.RecipeSectionRels.SectionRecipeSectionInstructions)),
		// has images via notion recipe
		Load(Rels(models.RecipeRels.NotionRecipes,
			models.NotionRecipeRels.PageNotionImages,
			models.NotionImageRels.NotionImageImage,
		)),
		Limit(int(limit)),
		Offset(int(offset)),
	).
		All(ctx, a.db.DB())
	if err != nil {
		return nil, err
	}
	// spew.Dump(recipes)
	items := []RecipeWrapper{}
	for _, recipe := range recipes {
		rw, err := a.recipeFromModel(ctx, recipe)
		if err != nil {
			return nil, err
		}
		items = append(items, *rw)
	}
	return items, nil
}
func (a *API) Misc(c echo.Context) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "Misc")
	defer span.End()

	items, err := a.RecipeListV2(ctx, 10, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	// s := spew.Sdump(recipes)
	// // s = strings.ReplaceAll(s, "\n", "<br/>")
	// s = fmt.Sprintf("<html>%s</html>", s)
	return c.JSON(http.StatusOK, items)
}
