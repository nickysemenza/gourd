package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nickysemenza/food/backend/app/config"
	"github.com/nickysemenza/food/backend/app/model"
	"log"
	"net/http"
)

//GetAllMeals gets all meals, with their related recipes
func GetAllMeals(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	var meals []model.Meal
	e.DB.Order("time DESC").Preload("RecipeMeal.Recipe").Find(&meals)
	respondSuccess(w, meals)
	return nil
}

//GetMealByID retrieves a meal
func GetMealByID(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	var meal model.Meal
	e.DB.Preload("RecipeMeal.Recipe").First(&meal, id)
	respondSuccess(w, meal)
	return nil
}

func UpdateMealByID(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]

	decoder := json.NewDecoder(r.Body)
	var updatedMeal model.Meal
	err := decoder.Decode(&updatedMeal)
	if err != nil {
		log.Println(err)
	}

	updatedMeal.CreateOrUpdate(e.DB)

	e.DB.Preload("RecipeMeal.Recipe").First(&updatedMeal, id)
	respondSuccess(w, updatedMeal)
	return nil
}
