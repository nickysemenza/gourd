package app

import (
	"encoding/json"
	"fmt"
	"github.com/nickysemenza/food/backend/app/handler"
	"github.com/nickysemenza/food/backend/app/model"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Utils struct {
	*handler.Env
}

func (u Utils) Import(path string) {
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
				eachRecipe.CreateOrUpdate(u.Env.DB, true)
			}
		}
	}
	log.Printf("Exported %d recipes from %s", len(files), path)
}

func (u Utils) Export(path string) {
	recipes := []model.Recipe{}
	u.Env.DB.Preload("Sections.Instructions").Preload("Sections.Ingredients.Item").Find(&recipes)
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
