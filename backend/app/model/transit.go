package model

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func Export(db *gorm.DB) {
	recipes := []Recipe{}
	db.Preload("Sections.Instructions").Preload("Sections.Ingredients.Item").Find(&recipes)

	pwd, _ := os.Getwd()
	pwd += "/recipes2/"

	for _, r := range recipes {
		jsonData, _ := json.Marshal(r)
		err := ioutil.WriteFile(pwd+r.Slug+".json", jsonData, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func LegacyImport(db *gorm.DB) {
	pwd, _ := os.Getwd()
	pwd += "/recipes/"

	files, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			raw, err := ioutil.ReadFile(pwd + f.Name())
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			legacyImport(db, raw, strings.TrimSuffix(f.Name(), ".json"))
		}
	}

}
func legacyImport(db *gorm.DB, bytes []byte, slug string) {
	var c LegacyData
	json.Unmarshal(bytes, &c)
	//fmt.Println(c)

	toImport := Recipe{}
	toImport.Slug = slug
	toImport.Equipment = strings.Join(c.Equipment, ", ")
	toImport.Unit = c.Unit
	toImport.Title = c.Title
	toImport.Source = c.Source
	toImport.Servings = uint(c.Servings)
	toImport.TotalMinutes = uint(c.TotalMinutes)
	toImport.Quantity = uint(c.Quantity)
	toImport.CreatedAt = time.Now()
	toImport.UpdatedAt = time.Now()
	for _, s := range c.Sections {
		eachS := Section{}
		s.Minutes = c.TotalMinutes
		for _, b := range s.Instructions {
			eachS.Instructions = append(eachS.Instructions, SectionInstruction{Name: b})
		}
		for _, b := range s.Ingredients {
			eachI := Ingredient{}
			db.FirstOrCreate(&eachI, Ingredient{Name: b.Name})
			eachS.Ingredients = append(eachS.Ingredients, SectionIngredient{
				Item:       eachI,
				Grams:      b.Grams,
				Amount:     b.Measurement.Amount,
				Unit:       b.Measurement.Unit,
				Substitute: b.Substitute,
				Modifier:   b.Modifier,
				Optional:   false,
			})
		}
		toImport.Sections = append(toImport.Sections, eachS)
	}
	b, err := json.MarshalIndent(toImport, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(b))
	db.Create(&toImport)
}

type LegacyData struct {
	Equipment []string `json:"equipment"`
	Quantity  int      `json:"quantity"`
	Sections  []struct {
		Ingredients []struct {
			Measurement struct {
				Amount float32 `json:"amount"`
				Unit   string  `json:"unit"`
			} `json:"measurement"`
			Modifier   string  `json:"modifier,omitempty"`
			Name       string  `json:"name"`
			Substitute string  `json:"substitute"`
			Grams      float32 `json:"grams,omitempty"`
		} `json:"ingredients,omitempty"`
		Instructions []string `json:"instructions"`
		Minutes      int      `json:"minutes"`
	} `json:"sections"`
	Servings     int    `json:"servings"`
	Source       string `json:"source"`
	Title        string `json:"title"`
	TotalMinutes int    `json:"totalMinutes"`
	Unit         string `json:"unit"`
}
