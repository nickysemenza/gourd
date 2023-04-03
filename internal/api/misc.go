package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nickysemenza/gourd/internal/db/models"
	"github.com/sirupsen/logrus"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gopkg.in/guregu/null.v4/zero"
)

func (a *API) ingredientFromModel(_ context.Context, ingredient *models.Ingredient) *IngredientDetail {
	if ingredient == nil {
		return nil
	}
	i := IngredientDetail{
		Recipes:      []RecipeDetail{},
		UnitMappings: []UnitMapping{},
	}
	return &i
}

func (a *API) imageFromModel(ctx context.Context, p *models.Image) (*Photo, error) {
	if p == nil {
		return nil, nil
	}
	url, err := a.ImageStore.GetImageURL(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	return &Photo{
		Id:       p.ID,
		TakenAt:  p.TakenAt.Ptr(),
		BlurHash: &p.BlurHash,
		Width:    300,
		Height:   400,
		BaseUrl:  url,
		Source:   PhotoSource(p.Source),
	}, nil

}

func (a *API) imagesFromModel(ctx context.Context, nr models.NotionRecipeSlice, gp models.GphotosPhotoSlice) ([]Photo, error) {
	items := []Photo{}
	for _, notionRecipe := range nr {
		for _, notionImage := range notionRecipe.R.PageNotionImages {
			photo, err := a.imageFromModel(ctx, notionImage.R.Image)
			if err != nil {
				return nil, err
			}
			items = append(items, *photo)
		}
	}

	for _, gPhoto := range gp {

		image, err := a.imageFromModel(ctx, gPhoto.R.Image)
		if err != nil {
			return nil, err
		}
		items = append(items, *image)
	}

	return items, nil
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
		sections := make([]RecipeSection, 0)
		for _, section := range d.R.RecipeSections {
			s := RecipeSection{
				Id:           section.ID,
				Ingredients:  []SectionIngredient{},
				Instructions: []SectionInstruction{},
			}
			if section.DurationTimerange.Valid {
				err := json.Unmarshal([]byte(section.DurationTimerange.JSON), &s.Duration)
				if err != nil {
					return nil, err
				}
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
					Amounts:   []Amount{},
				}
				if ingredient.IngredientID.Valid {
					si.Kind = IngredientKindIngredient
					si.Ingredient = a.ingredientFromModel(ctx, ingredient.R.Ingredient)
				} else {
					si.Kind = IngredientKindRecipe
					foo, err := a.recipeFromModel(ctx, ingredient.R.Recipe)
					if err != nil {
						return nil, err
					}
					si.Recipe = &foo.Detail
				}
				amounts := []Amount{}
				err := ingredient.Amounts.Unmarshal(&amounts)
				if err != nil {
					return nil, err
				}
				for _, amt := range amounts {
					si.Amounts = append(si.Amounts, Amount{
						Unit:   amt.Unit,
						Value:  amt.Value,
						Source: zero.StringFrom("db").Ptr(),
					})
					//todo: unit conversion restuests
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
			Sources:         []RecipeSource{},
			Unit:            d.Unit.String,
			Version:         int64(d.Version),
			Sections:        sections,
			IsLatestVersion: d.IsLatestVersion.Bool,
			Tags:            d.Tags,
		}
		if rd.Tags == nil {
			rd.Tags = []string{}
		}
		if d.Source.Valid {
			if err := d.Source.Unmarshal(&rd.Sources); err != nil {
				return nil, err
			}
		}

		if rd.IsLatestVersion {
			rw.Detail = rd
			if d.DeletedAt.Valid {
				return nil, nil
			}
		} else {
			other = append(other, rd)
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
func (a *API) RecipeListV2(ctx context.Context, limit, offset uint64) ([]RecipeWrapper, error) {
	ctx, span := a.tracer.Start(ctx, "RecipeListV2")
	defer span.End()
	recipes, err := models.Recipes(
		// Load(models.RecipeRels.RecipeDetails, Where("recipe_details.is_latest_version = ?", true)),
		// has many sections, has many ingredients, which can be ingredients or recipes
		Load(models.RecipeRels.RecipeDetails),
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
		Load(Rels(models.RecipeRels.MealRecipes,
			models.MealRecipeRels.Meal,
			models.MealRels.MealGphotos,
			models.MealGphotoRels.Gphoto,
			models.GphotosPhotoRels.Image,
		)),
		OrderBy("recipes.created_at DESC"),
		Limit(int(limit)),
		Offset(int(offset)),
	).
		All(ctx, a.db.DB())
	if err != nil {
		return nil, err
	}
	items := []RecipeWrapper{}
	for _, recipe := range recipes {
		rw, err := a.recipeFromModel(ctx, recipe)
		if err != nil {
			return nil, err
		}
		if rw != nil {
			items = append(items, *rw)
		}
	}

	return items, nil
}
func (a *API) indexRecipeDetails(ctx context.Context, details ...RecipeDetail) {
	ctx, span := a.tracer.Start(ctx, "indexRecipeDetails")
	defer span.End()
	if err := a.R.Send(ctx, "index_recipe_detail", details, nil); err != nil {
		logrus.Error(err)
	}
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
		Load(Rels(models.RecipeDetailRels.Recipe,
			models.RecipeRels.MealRecipes,
			models.MealRecipeRels.Meal,
			models.MealRels.MealGphotos,
			models.MealGphotoRels.Gphoto,
			models.GphotosPhotoRels.Image,
		)),
	).
		One(ctx, a.db.DB())
	if err != nil {
		return nil, err
	}
	gp := models.GphotosPhotoSlice{}
	for _, m := range rd.R.Recipe.R.MealRecipes {
		for _, x := range m.R.Meal.R.MealGphotos {
			gp = append(gp, x.R.Gphoto)
		}

	}

	return a.imagesFromModel(ctx, rd.R.Recipe.R.NotionRecipes, gp)

}

func (a *API) Misc(c echo.Context) error {
	ctx, span := a.tracer.Start(c.Request().Context(), "Misc")
	defer span.End()

	// items, err := a.imagesFromRecipeDetailId(ctx, "rd_08c6db27")
	// items, err := a.Notion.PageById(ctx, "f6a5d0759d4a4becb95adf696b1cccb0")
	items, err := a.IngredientUsage(ctx, EntitySummary{Id: "rd_2dfbb24c", Kind: IngredientKindRecipe})

	if err != nil {
		return handleErr(c, err)
	}
	// s := spew.Sdump(recipes)
	// // s = strings.ReplaceAll(s, "\n", "<br/>")
	// s = fmt.Sprintf("<html>%s</html>", s)
	return c.JSON(http.StatusOK, items)
}
