package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nickysemenza/food/app/model"
	"github.com/pkg/errors"
)

//GetAllRecipes gets all recipes: GET /recipes
func GetAllRecipes(c *gin.Context) {
	var recipes []model.Recipe
	db := model.GetDBFromContext(c.MustGet("ctx").(context.Context))
	db.Preload("Images.Sizes").Preload("Categories").Find(&recipes)
	c.JSON(http.StatusOK, recipes)
}

//GetRecipe gets a recipe by its slug: GET /recipes/{slug}
func GetRecipe(c *gin.Context) {
	slug := c.Params.ByName("slug")

	// span := c.MustGet("tracing-context").(opentracing.Span)
	ctx := c.MustGet("ctx").(context.Context)

	recipe, err := model.GetRecipeFromSlug(ctx, slug)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.New("recipe "+slug+" not found"))
	} else {
		c.JSON(http.StatusOK, recipe)
	}

}

//PutRecipe updates or creates: PUT /recipes/{slug}
func PutRecipe(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	var updatedRecipe model.Recipe

	if err := c.BindJSON(&updatedRecipe); err != nil {
		log.Println(err)
	}

	updatedRecipe.CreateOrUpdate(ctx, false)

	slug := updatedRecipe.Slug

	if recipe, err := model.GetRecipeFromSlug(ctx, slug); err != nil {
		c.JSON(404, errors.New("recipe "+slug+" not found"))
	} else {
		c.JSON(http.StatusOK, recipe)
	}

}

//CreateRecipe Creates a new recipe from a Slug and Title
func CreateRecipe(c *gin.Context) {
	db := model.GetDBFromContext(c.MustGet("ctx").(context.Context))
	//decode the data from JSON encoded request body
	var parsed struct {
		Slug  string `json:"slug"`
		Title string `json:"title"`
	}
	if err := c.BindJSON(&parsed); err != nil {
		log.Println(err)
	}

	//see if one exists
	recipe := model.Recipe{}
	if !db.Where("slug = ?", parsed.Slug).First(&recipe).RecordNotFound() {
		c.JSON(500, "slug exists already")

	}
	recipe.Slug = parsed.Slug
	recipe.Title = parsed.Title
	db.Save(&recipe)
	c.JSON(http.StatusOK, "added!")

}

//AddNote adds a Note to a Recipe based on Slug, and Note Body
func AddNote(c *gin.Context) {
	db := model.GetDBFromContext(c.MustGet("ctx").(context.Context))
	//find the recipe we are adding a note to
	recipe := model.Recipe{}
	slug := c.Params.ByName("slug")
	if err := db.Where("slug = ?", slug).First(&recipe).Error; err != nil {
		c.JSON(404, errors.New("recipe "+slug+" not found"))
	}

	//decode the note from JSON encoded request body
	var parsed struct {
		Note string `json:"note"`
	}

	if err := c.BindJSON(&parsed); err != nil {
		log.Println(err)
	}

	//add a new RecipeNote Model, save it
	note := model.RecipeNote{
		Body:     parsed.Note,
		RecipeID: recipe.ID,
	}
	db.Save(&note)

	c.JSON(http.StatusOK, note)

}

//GetAllCategories gets all categories that exist
func GetAllCategories(c *gin.Context) {
	db := model.GetDBFromContext(c.MustGet("ctx").(context.Context))
	var categories []model.Category
	db.Find(&categories)
	c.JSON(http.StatusOK, categories)

}
