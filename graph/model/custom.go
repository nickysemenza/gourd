package model

import "gopkg.in/guregu/null.v3/zero"

// Recipe has section IDs
type Recipe struct {
	UUID         string   `json:"uuid"`
	Name         string   `json:"name"`
	TotalMinutes int      `json:"total_minutes"`
	Unit         string   `json:"unit"`
	SectionIDs   []string `json:"sections"`
}

// SectionIngredient has ingredients  and instructions by ID
type Section struct {
	UUID            string   `json:"uuid"`
	RecipeUUID      string   `json:"recipe_uuid"`
	Minutes         int      `json:"minutes"`
	InstructionsIDs []string `json:"instructions"`
	IngredientsIDs  []string `json:"ingredients"`
}

// SectionIngredient has ingredients by ID
type SectionIngredient struct {
	UUID         string  `json:"uuid"`
	IngredientID string  `json:"info"`
	RecipeID     string  `json:"info_recipe"`
	Grams        float64 `json:"grams"`
	Amount       float64 `json:"amount"`
	Unit         string  `json:"unit"`
	Adjective    string  `json:"adjective"`
	Optional     bool    `json:"optional"`
}

type Ingredient struct {
	UUID       string   `json:"uuid"`
	Name       string   `json:"name"`
	RecipesIDs []string `json:"recipes"`
}

func (i Ingredient) IsIngredientInfo() {}
func (r Recipe) IsIngredientInfo()     {}

type FoodNutrient struct {
	NutrientID int      `json:"nutrient" db:"nutrient_id"`
	Amount     float64  `json:"amount" db:"amount"`
	DataPoints zero.Int `json:"data_points" db:"data_points"`
}

type Nutrient struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	UnitName string `json:"unit_name" db:"unit_name"`
}

type Food struct {
	FdcID       int          `json:"fdc_id" db:"fdc_id"`
	Description string       `json:"description" db:"description"`
	DataType    FoodDataType `json:"data_type" db:"data_type"`
	CategoryID  zero.Int     `json:"category" db:"food_category_id"`
	NutrientIDs []int        `json:"nutrients"`
}

type FoodCategory struct {
	Code        string `json:"code" db:"code"`
	Description string `json:"description" db:"description"`
}
