package handler

import (
	"github.com/jinzhu/gorm"
	"github.com/nickysemenza/food/backend/app/model"
	"net/http"
)

func GetAllRecipes(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	projects := []model.Recipe{}
	db.Preload("Sections.Instructions").Preload("Sections.Ingredients.Item").Find(&projects)
	respondJSON(w, http.StatusOK, projects)
}
