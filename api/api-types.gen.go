// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

import (
	"time"
)

// AuthResp defines model for AuthResp.
type AuthResp struct {
	Jwt  string                 `json:"jwt"`
	User map[string]interface{} `json:"user"`
}

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
}

// GooglePhoto defines model for GooglePhoto.
type GooglePhoto struct {

	// public image
	BaseUrl string `json:"base_url"`

	// when it was taken
	Created time.Time `json:"created"`

	// id
	Id string `json:"id"`
}

// Ingredient defines model for Ingredient.
type Ingredient struct {

	// UUID
	Id string `json:"id"`

	// Ingredient name
	Name string `json:"name"`
}

// IngredientDetail defines model for IngredientDetail.
type IngredientDetail struct {

	// Ingredients that are equivalent
	Children *[]Ingredient `json:"children,omitempty"`

	// An Ingredient
	Ingredient Ingredient `json:"ingredient"`

	// Recipes referencing this ingredient
	Recipes *[]Recipe `json:"recipes,omitempty"`
}

// List defines model for List.
type List struct {

	// How many items were requested for this page
	Limit int `json:"limit"`

	// todo
	Offset int `json:"offset"`

	// Total number of pages available
	PageCount int `json:"page_count"`

	// What number page this is
	PageNumber int `json:"page_number"`

	// Total number of items across all pages
	TotalCount int `json:"total_count"`
}

// Meal defines model for Meal.
type Meal struct {

	// when it was taken
	AteAt time.Time `json:"ate_at"`

	// id
	Id string `json:"id"`

	// public image
	Name   string        `json:"name"`
	Photos []GooglePhoto `json:"photos"`
}

// PaginatedIngredients defines model for PaginatedIngredients.
type PaginatedIngredients struct {
	Ingredients *[]IngredientDetail `json:"ingredients,omitempty"`
	Meta        *List               `json:"meta,omitempty"`
}

// PaginatedMeals defines model for PaginatedMeals.
type PaginatedMeals struct {
	Meals *[]Meal `json:"meals,omitempty"`
	Meta  *List   `json:"meta,omitempty"`
}

// PaginatedPhotos defines model for PaginatedPhotos.
type PaginatedPhotos struct {
	Meta   *List          `json:"meta,omitempty"`
	Photos *[]GooglePhoto `json:"photos,omitempty"`
}

// PaginatedRecipes defines model for PaginatedRecipes.
type PaginatedRecipes struct {
	Meta    *List     `json:"meta,omitempty"`
	Recipes *[]Recipe `json:"recipes,omitempty"`
}

// Recipe defines model for Recipe.
type Recipe struct {

	// UUID
	Id string `json:"id"`

	// recipe name
	Name string `json:"name"`

	// serving quantity
	Quantity int64 `json:"quantity"`

	// num servings
	Servings *int64 `json:"servings,omitempty"`

	// book or website? deprecated?
	Source *string `json:"source,omitempty"`

	// todo
	TotalMinutes *int64 `json:"total_minutes,omitempty"`

	// serving unit
	Unit string `json:"unit"`
}

// RecipeDetail defines model for RecipeDetail.
type RecipeDetail struct {

	// A recipe
	Recipe Recipe `json:"recipe"`

	// sections of the recipe
	Sections []RecipeSection `json:"sections"`
}

// RecipeSection defines model for RecipeSection.
type RecipeSection struct {

	// UUID
	Id string `json:"id"`

	// x
	Ingredients []SectionIngredient `json:"ingredients"`

	// x
	Instructions []SectionInstruction `json:"instructions"`

	// How many minutes the step takes, approximately (todo - make this a range)
	Minutes int64 `json:"minutes"`
}

// SectionIngredient defines model for SectionIngredient.
type SectionIngredient struct {

	// adjective
	Adjective string `json:"adjective"`

	// amount
	Amount *float64 `json:"amount,omitempty"`

	// weight in grams
	Grams float64 `json:"grams"`

	// UUID
	Id string `json:"id"`

	// An Ingredient
	Ingredient *Ingredient `json:"ingredient,omitempty"`

	// what kind of ingredient
	Kind string `json:"kind"`

	// optional
	Optional *bool `json:"optional,omitempty"`

	// A recipe
	Recipe *Recipe `json:"recipe,omitempty"`

	// unit
	Unit string `json:"unit"`
}

// SectionInstruction defines model for SectionInstruction.
type SectionInstruction struct {

	// UUID
	Id string `json:"id"`

	// instruction
	Instruction string `json:"instruction"`
}

// LimitParam defines model for limitParam.
type LimitParam int

// OffsetParam defines model for offsetParam.
type OffsetParam int

// AuthLoginParams defines parameters for AuthLogin.
type AuthLoginParams struct {

	// Google code
	Code string `json:"code"`
}

// ListIngredientsParams defines parameters for ListIngredients.
type ListIngredientsParams struct {

	// The number of items to skip before starting to collect the result set.
	Offset *OffsetParam `json:"offset,omitempty"`

	// The numbers of items to return.
	Limit *LimitParam `json:"limit,omitempty"`
}

// CreateIngredientsJSONBody defines parameters for CreateIngredients.
type CreateIngredientsJSONBody Ingredient

// ListMealsParams defines parameters for ListMeals.
type ListMealsParams struct {

	// The number of items to skip before starting to collect the result set.
	Offset *OffsetParam `json:"offset,omitempty"`

	// The numbers of items to return.
	Limit *LimitParam `json:"limit,omitempty"`
}

// ListPhotosParams defines parameters for ListPhotos.
type ListPhotosParams struct {

	// The number of items to skip before starting to collect the result set.
	Offset *OffsetParam `json:"offset,omitempty"`

	// The numbers of items to return.
	Limit *LimitParam `json:"limit,omitempty"`
}

// ListRecipesParams defines parameters for ListRecipes.
type ListRecipesParams struct {

	// The number of items to skip before starting to collect the result set.
	Offset *OffsetParam `json:"offset,omitempty"`

	// The numbers of items to return.
	Limit *LimitParam `json:"limit,omitempty"`
}

// CreateRecipesJSONBody defines parameters for CreateRecipes.
type CreateRecipesJSONBody RecipeDetail

// CreateIngredientsRequestBody defines body for CreateIngredients for application/json ContentType.
type CreateIngredientsJSONRequestBody CreateIngredientsJSONBody

// CreateRecipesRequestBody defines body for CreateRecipes for application/json ContentType.
type CreateRecipesJSONRequestBody CreateRecipesJSONBody
