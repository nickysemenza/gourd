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
func (a *API) imagesFromModel(ctx context.Context, nr models.NotionRecipeSlice) []Photo {
	items := []Photo{}
	for _, notionRecipe := range nr {
		for _, notionImage := range notionRecipe.R.PageNotionImages {
			url := a.ImageStore.GetImageURL(ctx, notionImage.ImageID)
			items = append(items, Photo{
				Id:       notionImage.BlockID,
				Created:  notionImage.LastSeen,
				BlurHash: &notionImage.R.Image.BlurHash,
				Width:    300,
				Height:   400,
				BaseUrl:  url,
				Source:   PhotoSourceNotion,
			})
		}
	}
	return items
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
				Id:           section.ID,
				Ingredients:  []SectionIngredient{},
				Instructions: []SectionInstruction{},
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
				if ingredient.IngredientID.Valid {
					si.Kind = IngredientKindIngredient
					var err error
					si.Ingredient, err = a.ingredientFromModel(ctx, ingredient.R.Ingredient)
					if err != nil {
						return nil, err
					}
				} else {
					si.Kind = IngredientKindRecipe
					foo, err := a.recipeFromModel(ctx, ingredient.R.Recipe)
					if err != nil {
						return nil, err
					}
					si.Recipe = &foo.Detail
				}
				s.Ingredients = append(s.Ingredients, si)
			}
			sections = append(sections, s)

		}

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

	items := a.imagesFromModel(ctx, recipe.R.NotionRecipes)
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
			models.RecipeSectionIngredientRels.Ingredient,
		)),
		Load(Rels(models.RecipeRels.RecipeDetails,
			models.RecipeDetailRels.RecipeSections,
			models.RecipeSectionRels.SectionRecipeSectionIngredients,
			models.RecipeSectionIngredientRels.Recipe,
			models.RecipeRels.RecipeDetails,
		)),
		// has many sections, has many instructions
		Load(Rels(models.RecipeRels.RecipeDetails,
			models.RecipeDetailRels.RecipeSections,
			models.RecipeSectionRels.SectionRecipeSectionInstructions)),
		// has images via notion recipe
		Load(Rels(models.RecipeRels.NotionRecipes,
			models.NotionRecipeRels.PageNotionImages,
			models.NotionImageRels.Image,
		)),
		OrderBy("recipes.created_at DESC"),
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

func (a *API) imagesFromRecipeDetailId(ctx context.Context, id string) ([]Photo, error) {
	rd, err := models.RecipeDetails(
		Where("recipe_details.id = ?", id),
		Load(Rels(
			models.RecipeDetailRels.Recipe,
			models.RecipeRels.NotionRecipes,
			models.NotionRecipeRels.PageNotionImages,
			models.NotionImageRels.Image,
		)),
	).
		One(ctx, a.db.DB())
	if err != nil {
		return nil, err
	}
	return a.imagesFromModel(ctx, rd.R.Recipe.R.NotionRecipes), nil
	// rw, err := a.recipeFromModel(ctx, recipe.R.Recipe)
	// if err != nil {
	// 	return nil, err
	// }
	// return rw, nil
}

func (a *API) Misc(c echo.Context) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "Misc")
	defer span.End()

	// items, err := a.imagesFromRecipeDetailId(ctx, "rd_08c6db27")
	items, err := a.Notion.PageById(ctx, "57946cb8ce6c40fdb685a25d14bf9feb")

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	// s := spew.Sdump(recipes)
	// // s = strings.ReplaceAll(s, "\n", "<br/>")
	// s = fmt.Sprintf("<html>%s</html>", s)
	return c.JSON(http.StatusOK, items)
}