package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nickysemenza/food/backend/app/config"
	"github.com/nickysemenza/food/backend/app/model"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

//GetAllRecipes gets all recipes: GET /recipes
func GetAllRecipes(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	var recipes []model.Recipe
	e.DB.Preload("Images.Sizes").Preload("Categories").Find(&recipes)
	respondSuccess(w, recipes)
	return nil
}

//GetRecipe gets a recipe by its slug: GET /recipes/{slug}
func GetRecipe(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	recipe := model.Recipe{}
	slug := mux.Vars(r)["slug"]

	if err := recipe.GetFromSlug(e.DB, slug); err != nil {
		return StatusError{Code: 404, Err: errors.New("recipe " + slug + " not found")}
	}

	respondSuccess(w, recipe)
	return nil
}

//PutRecipe updates or creates: PUT /recipes/{slug}
func PutRecipe(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	var updatedRecipe model.Recipe
	err := decoder.Decode(&updatedRecipe)
	if err != nil {
		log.Println(err)
	}

	updatedRecipe.CreateOrUpdate(e.DB, false)

	slug := updatedRecipe.Slug
	recipe := model.Recipe{}

	if err := recipe.GetFromSlug(e.DB, slug); err != nil {
		return StatusError{Code: 404, Err: errors.New("recipe " + slug + " not found")}
	}

	respondSuccess(w, recipe)
	return nil
}

//CreateRecipe Creates a new recipe from a Slug and Title
func CreateRecipe(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	//decode the data from JSON encoded request body
	decoder := json.NewDecoder(r.Body)
	var parsed struct {
		Slug  string `json:"slug"`
		Title string `json:"title"`
	}
	err := decoder.Decode(&parsed)
	if err != nil {
		log.Println(err)
	}

	//see if one exists
	recipe := model.Recipe{}
	if !e.DB.Where("slug = ?", parsed.Slug).First(&recipe).RecordNotFound() {
		respondError(w, 500, "slug exists already")
		return nil
	}
	recipe.Slug = parsed.Slug
	recipe.Title = parsed.Title
	e.DB.Save(&recipe)
	respondSuccess(w, "added!")
	return nil
}

//AddNote adds a Note to a Recipe based on Slug, and Note Body
func AddNote(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	//find the recipe we are adding a note to
	recipe := model.Recipe{}
	slug := mux.Vars(r)["slug"]
	if err := e.DB.Where("slug = ?", slug).First(&recipe).Error; err != nil {
		return StatusError{Code: 404, Err: errors.New("recipe " + slug + " not found")}
	}

	//decode the note from JSON encoded request body
	decoder := json.NewDecoder(r.Body)
	var parsed struct {
		Note string `json:"note"`
	}
	err := decoder.Decode(&parsed)
	if err != nil {
		log.Println(err)
	}
	//add a new RecipeNote Model, save it
	note := model.RecipeNote{
		Body:     parsed.Note,
		RecipeID: recipe.ID,
	}
	e.DB.Save(&note)

	respondSuccess(w, note)
	return nil
}

//GetAllCategories gets all categories that exist
func GetAllCategories(e *config.Env, w http.ResponseWriter, r *http.Request) error {
	var categories []model.Category
	e.DB.Find(&categories)
	respondSuccess(w, categories)
	return nil
}
