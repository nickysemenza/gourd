// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type NewRecipe struct {
	Name string `json:"name"`
}

type RecipeInput struct {
	UUID         string          `json:"uuid"`
	Name         string          `json:"name"`
	TotalMinutes *int            `json:"total_minutes"`
	Unit         *string         `json:"unit"`
	Sections     []*SectionInput `json:"sections"`
}

type SectionIngredientInput struct {
	Name  string  `json:"name"`
	Grams float64 `json:"grams"`
}

type SectionInput struct {
	Minutes      int                        `json:"minutes"`
	Instructions []*SectionInstructionInput `json:"instructions"`
	Ingredients  []*SectionIngredientInput  `json:"ingredients"`
}

type SectionInstruction struct {
	UUID        string `json:"uuid"`
	Instruction string `json:"instruction"`
}

type SectionInstructionInput struct {
	Instruction string `json:"instruction"`
}
