package model

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
