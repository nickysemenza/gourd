package app

import (
	"context"
	"encoding/json"
	"fmt"

	//load jpeg and png for processing
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/nickysemenza/food/app/model"
)

//Import imports a folder of recipes in json format
//TODO: move Import/Export to not be App methods
func (a App) Import(ctx context.Context, path string) {
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
				eachRecipe.CreateOrUpdate(ctx, true)
			}
		}
	}
	log.Printf("Exported %d recipes from %s", len(files), path)
}

//Export exports a folder of recipes in json format
func (a App) Export(ctx context.Context, path string) {
	recipes := []model.Recipe{}
	db := model.GetDBFromContext(ctx)
	db.Preload("Sections.Instructions").Preload("Sections.Ingredients.Item").Find(&recipes)
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
