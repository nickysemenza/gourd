package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nickysemenza/food/app/model"
)

//GetAllMeals gets all meals, with their related recipes
func GetAllMeals(c *gin.Context) {
	db := model.GetDBFromContext(c.MustGet("ctx").(context.Context))
	var meals []model.Meal
	db.Order("time DESC").Preload("RecipeMeal.Recipe").Find(&meals)
	c.JSON(http.StatusOK, meals)
}

//GetMealByID retrieves a meal
func GetMealByID(c *gin.Context) {
	db := model.GetDBFromContext(c.MustGet("ctx").(context.Context))
	id := c.Params.ByName("id")
	var meal model.Meal
	db.Preload("RecipeMeal.Recipe.Sections.Ingredients.Item").Preload("RecipeMeal.Recipe.Sections.Instructions").First(&meal, id)
	c.JSON(http.StatusOK, meal)
}

func UpdateMealByID(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	id := c.Params.ByName("id")

	var updatedMeal model.Meal

	if err := c.BindJSON(&updatedMeal); err != nil {
		log.Println(err)
	}

	updatedMeal.CreateOrUpdate(ctx)

	db := model.GetDBFromContext(ctx)
	db.Preload("RecipeMeal.Recipe").First(&updatedMeal, id) //TODO: move this into model
	c.JSON(http.StatusOK, updatedMeal)
}
