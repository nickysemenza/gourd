package handler

import (
	"github.com/gorilla/mux"
	"github.com/nickysemenza/food/backend/app/model"
	"github.com/pkg/errors"
	"net/http"
)

func GetAllRecipes(e *Env, w http.ResponseWriter, r *http.Request) error {
	recipes := []model.Recipe{}
	e.DB.Select([]string{"slug"}).Find(&recipes)
	var slugs []string
	for _, r := range recipes {
		slugs = append(slugs, r.Slug)
	}
	respondJSON(w, http.StatusOK, slugs)
	return nil
}
func ErrorTest(e *Env, w http.ResponseWriter, r *http.Request) error {
	return StatusError{Code: 201, Err: errors.New("sad..")}
}
func GetRecipe(e *Env, w http.ResponseWriter, r *http.Request) error {
	recipe := model.Recipe{}
	vars := mux.Vars(r)
	slug := vars["slug"]
	e.DB.Where("slug = ?", slug).Preload("Sections.Instructions").Preload("Sections.Ingredients.Item").First(&recipe)
	respondJSON(w, http.StatusOK, recipe)
	return nil
}
