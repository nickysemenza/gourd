package utils

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
			}

			var c model.Recipe
			json.Unmarshal(raw, &c)
			u.Env.DB.Create(&c)
		}
	}
}

func (u Utils) Export(path string) {
	recipes := []model.Recipe{}
	u.Env.DB.Preload("Sections.Instructions").Preload("Sections.Ingredients.Item").Find(&recipes)

	for _, r := range recipes {
		jsonData, _ := json.Marshal(r)
		err := ioutil.WriteFile(path+r.Slug+".json", jsonData, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
