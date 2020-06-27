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

type Food struct {
	FdcID       int          `json:"fdc_id" db:"fdc_id"`
	Description string       `json:"description" db:"description"`
	DataType    FoodDataType `json:"data_type" db:"data_type"`
	CategoryID  zero.Int     `json:"category" db:"food_category_id"`
	NutrientIDs []int        `json:"nutrients"`
}
