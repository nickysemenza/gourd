package app

import (
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/nickysemenza/food/backend/app/model"
)

//Import imports a folder of recipes in json format
func (a App) Import(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			raw, err := ioutil.ReadFile(path + f.Name())
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			} else {
				var eachRecipe model.Recipe
				json.Unmarshal(raw, &eachRecipe)
				log.Printf("importing: %s (%s)\n", eachRecipe.Title, eachRecipe.Slug)
				eachRecipe.CreateOrUpdate(a.Env.DB, true)
			}
		}
	}
	log.Printf("Exported %d recipes from %s", len(files), path)
}

//Export exports a folder of recipes in json format
func (a App) Export(path string) {
	recipes := []model.Recipe{}
	a.Env.DB.Preload("Sections.Instructions").Preload("Sections.Ingredients.Item").Find(&recipes)
	for _, r := range recipes {
		jsonData, _ := json.Marshal(r)
		err := ioutil.WriteFile(path+r.Slug+".json", jsonData, 0644)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("exporting: %s (%s)\n", r.Title, r.Slug)
		}
	}
	log.Printf("Exported %d recipes to %s", len(recipes), path)
}
