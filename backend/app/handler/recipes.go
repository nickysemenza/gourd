package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nickysemenza/food/backend/app/model"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

func GetAllRecipes(e *Env, w http.ResponseWriter, r *http.Request) error {
	var recipes []model.Recipe
	e.DB.Select([]string{"slug"}).Find(&recipes)
	var slugs []string
	for _, r := range recipes {
		slugs = append(slugs, r.Slug)
	}
	respondSuccess(w, slugs)
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
	recipe := model.Recipe{}
	vars := mux.Vars(r)
	slug := vars["slug"]
	if err := e.DB.Where("slug = ?", slug).Preload("Sections.Instructions").Preload("Sections.Ingredients.Item").First(&recipe).Error; err != nil {
		return StatusError{Code: 404, Err: errors.New("recipe " + slug + " not found")}
	}

	decoder := json.NewDecoder(r.Body)
	var updatedRecipe model.Recipe
	err := decoder.Decode(&updatedRecipe)
	if err != nil {
		panic(err)
	}

	for x := range updatedRecipe.Sections {
		eachSection := &updatedRecipe.Sections[x]
		for y := range eachSection.Ingredients {
			eachSectionIngredient := &eachSection.Ingredients[y]
			eachItem := &eachSectionIngredient.Item
			if eachItem.ID == 0 {
				//	new ingredient!
				//	find by name, to see if we have existing
				eachItem.FindOrCreateUsingName(e.DB)
				//eachItem = model.GetIngredientByName(e.DB, eachItem.Name)
				log.Printf("[ingredient] %s does not have an ID, giving it %d: ", eachItem.Name, eachItem.ID)
			} else {
				//	get fresh obj via eachIngredient.ID
				fresh := *eachItem
				fresh.GetFresh(e.DB)
				//	if eachIngredient.Name != fresh.Name IT WAS MUTATED AAH!
				if eachItem.Name != fresh.Name {
					log.Printf("[ingredient] name of %d was muted! %s->%s", eachItem.ID, eachItem.Name, fresh.Name)
					//we want to preserve the original eachItem; create new w/ eachItem.Name

					// find by name, or create new
					newItem := model.Ingredient{Name: eachItem.Name}
					newItem.FindOrCreateUsingName(e.DB)
					eachSectionIngredient.Item = newItem
				}
			}

		}
	}

	e.DB.Save(&updatedRecipe)
	respondSuccess(w, updatedRecipe)
	return nil
}
