package model

type Recipe struct {
	UUID         string   `json:"uuid"`
	Name         string   `json:"name"`
	TotalMinutes int      `json:"total_minutes"`
	Unit         string   `json:"unit"`
	SectionIDs   []string `json:"sections"`
}

type Section struct {
	UUID            string   `json:"uuid"`
	RecipeUUID      string   `json:"recipe_uuid"`
	Minutes         int      `json:"minutes"`
	InstructionsIDs []string `json:"instructions"`
	IngredientsIDs  []string `json:"ingredients"`
}

type SectionIngredient struct {
	UUID         string  `json:"uuid"`
	IngredientID string  `json:"info"`
	Grams        float64 `json:"grams"`
}
