package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nickysemenza/food/backend/app/model"
	"github.com/pkg/errors"
	"net/http"
)

func GetAllRecipes(e *Env, w http.ResponseWriter, r *http.Request) error {
	var recipes []model.Recipe
	e.DB.Find(&recipes)
	respondSuccess(w, recipes)
	return nil
}
func ErrorTest(e *Env, w http.ResponseWriter, r *http.Request) error {
	return StatusError{Code: 201, Err: errors.New("sad..")}
}
func GetRecipe(e *Env, w http.ResponseWriter, r *http.Request) error {
	recipe := model.Recipe{}
	vars := mux.Vars(r)
	slug := vars["slug"]
	if err := e.DB.Where("slug = ?", slug).Preload("Sections.Instructions").Preload("Sections.Ingredients.Item").First(&recipe).Error; err != nil {
		return StatusError{Code: 404, Err: errors.New("recipe " + slug + " not found")}
	}
	respondSuccess(w, recipe)
	return nil
}

func PutRecipe(e *Env, w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	var updatedRecipe model.Recipe
	err := decoder.Decode(&updatedRecipe)
	if err != nil {
		panic(err)
	}

	updatedRecipe.CreateOrUpdate(e.DB, false)

	respondSuccess(w, updatedRecipe)
	return nil
}
