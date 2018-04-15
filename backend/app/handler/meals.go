package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nickysemenza/food/backend/app/model"
	"log"
	"net/http"
)

//GetAllMeals gets all meals, with their related recipes
func GetAllMeals(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	var meals []model.Meal
	db.Order("time DESC").Preload("RecipeMeal.Recipe").Find(&meals)
	c.JSON(http.StatusOK, meals)
}

//GetMealByID retrieves a meal
func GetMealByID(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	id := c.Params.ByName("id")
	var meal model.Meal
	db.Preload("RecipeMeal.Recipe.Sections.Ingredients.Item").Preload("RecipeMeal.Recipe.Sections.Instructions").First(&meal, id)
	c.JSON(http.StatusOK, meal)
}

func UpdateMealByID(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	id := c.Params.ByName("id")

	var updatedMeal model.Meal

	if err := c.BindJSON(&updatedMeal); err != nil {
		log.Println(err)
	}

	updatedMeal.CreateOrUpdate(db)

	db.Preload("RecipeMeal.Recipe").First(&updatedMeal, id)
	c.JSON(http.StatusOK, updatedMeal)
}
