// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

import (
	"time"
)

// Amount defines model for Amount.
type Amount struct {

	// if it was explicit, inferred, etc
	Source *string `json:"source,omitempty"`

	// unit
	Unit string `json:"unit"`

	// value
	Value float64 `json:"value"`
}

// AuthResp defines model for AuthResp.
type AuthResp struct {
	Jwt  string                 `json:"jwt"`
	User map[string]interface{} `json:"user"`
}

// BrandedFood defines model for BrandedFood.
type BrandedFood struct {
	BrandOwner          *string `json:"brand_owner,omitempty"`
	BrandedFoodCategory *string `json:"branded_food_category,omitempty"`
	HouseholdServing    *string `json:"household_serving,omitempty"`
	Ingredients         *string `json:"ingredients,omitempty"`
	ServingSize         float64 `json:"serving_size"`
	ServingSizeUnit     string  `json:"serving_size_unit"`
}

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
}

// Food defines model for Food.
type Food struct {

	// branded_food
	BrandedInfo *BrandedFood `json:"branded_info,omitempty"`

	// food category, set for some
	Category *FoodCategory `json:"category,omitempty"`
	DataType FoodDataType  `json:"data_type"`

	// Food description
	Description string `json:"description"`

	// FDC Id
	FdcId int `json:"fdc_id"`

	// todo
	Nutrients []FoodNutrient `json:"nutrients"`

	// portion datapoints
	Portions *[]FoodPortion `json:"portions,omitempty"`

	// mappings of equivalent units
	UnitMappings []UnitMapping `json:"unit_mappings"`
}

// FoodCategory defines model for FoodCategory.
type FoodCategory struct {

	// Food description
	Code string `json:"code"`

	// Food description
	Description string `json:"description"`
}

// FoodDataType defines model for FoodDataType.
type FoodDataType string

// List of FoodDataType
const (
	FoodDataType_agricultural_acquisition FoodDataType = "agricultural_acquisition"
	FoodDataType_branded_food             FoodDataType = "branded_food"
	FoodDataType_foundation_food          FoodDataType = "foundation_food"
	FoodDataType_market_acquisition       FoodDataType = "market_acquisition"
	FoodDataType_sample_food              FoodDataType = "sample_food"
	FoodDataType_sr_legacy_food           FoodDataType = "sr_legacy_food"
	FoodDataType_sub_sample_food          FoodDataType = "sub_sample_food"
	FoodDataType_survey_fndds_food        FoodDataType = "survey_fndds_food"
)

// FoodNutrient defines model for FoodNutrient.
type FoodNutrient struct {
	Amount     float64 `json:"amount"`
	DataPoints int     `json:"data_points"`

	// todo
	Nutrient Nutrient `json:"nutrient"`
}

// FoodNutrientUnit defines model for FoodNutrientUnit.
type FoodNutrientUnit string

// List of FoodNutrientUnit
const (
	FoodNutrientUnit_G      FoodNutrientUnit = "G"
	FoodNutrientUnit_IU     FoodNutrientUnit = "IU"
	FoodNutrientUnit_KCAL   FoodNutrientUnit = "KCAL"
	FoodNutrientUnit_MG     FoodNutrientUnit = "MG"
	FoodNutrientUnit_MG_ATE FoodNutrientUnit = "MG_ATE"
	FoodNutrientUnit_SP_GR  FoodNutrientUnit = "SP_GR"
	FoodNutrientUnit_UG     FoodNutrientUnit = "UG"
	FoodNutrientUnit_kJ     FoodNutrientUnit = "kJ"
)

// FoodPortion defines model for FoodPortion.
type FoodPortion struct {
	Amount             float64 `json:"amount"`
	GramWeight         float64 `json:"gram_weight"`
	Id                 int     `json:"id"`
	Modifier           string  `json:"modifier"`
	PortionDescription string  `json:"portion_description"`
}

// GooglePhoto defines model for GooglePhoto.
type GooglePhoto struct {

	// public image
	BaseUrl string `json:"base_url"`

	// blur hash
	BlurHash *string `json:"blur_hash,omitempty"`

	// when it was taken
	Created time.Time `json:"created"`

	// height px
	Height int64 `json:"height"`

	// id
	Id string `json:"id"`

	// width px
	Width int64 `json:"width"`
}

// GooglePhotosAlbum defines model for GooglePhotosAlbum.
type GooglePhotosAlbum struct {

	// id
	Id string `json:"id"`

	// product_url
	ProductUrl string `json:"product_url"`

	// title
	Title string `json:"title"`

	// usecase
	Usecase string `json:"usecase"`
}

// Ingredient defines model for Ingredient.
type Ingredient struct {

	// FDC id equivalent to this ingredient
	FdcId *int64 `json:"fdc_id,omitempty"`

	// id
	Id string `json:"id"`

	// Ingredient name
	Name string `json:"name"`

	// ingredient ID for a similar (likely a different spelling)
	SameAs *string `json:"same_as,omitempty"`
}

// IngredientDetail defines model for IngredientDetail.
type IngredientDetail struct {

	// Ingredients that are equivalent
	Children []IngredientDetail `json:"children"`

	// A top level food
	Food *Food `json:"food,omitempty"`

	// An Ingredient
	Ingredient Ingredient `json:"ingredient"`

	// Recipes referencing this ingredient
	Recipes []RecipeDetail `json:"recipes"`

	// mappings of equivalent units
	UnitMappings []UnitMapping `json:"unit_mappings"`
}

// IngredientKind defines model for IngredientKind.
type IngredientKind string

// List of IngredientKind
const (
	IngredientKind_ingredient IngredientKind = "ingredient"
	IngredientKind_recipe     IngredientKind = "recipe"
)

// Items defines model for Items.
type Items struct {

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
	Name    string        `json:"name"`
	Photos  []GooglePhoto `json:"photos"`
	Recipes *[]MealRecipe `json:"recipes,omitempty"`
}

// MealRecipe defines model for MealRecipe.
type MealRecipe struct {

	// when it was taken
	Multiplier float64 `json:"multiplier"`

	// A revision of a recipe
	Recipe RecipeDetail `json:"recipe"`
}

// MealRecipeUpdate defines model for MealRecipeUpdate.
type MealRecipeUpdate struct {

	// todo
	Action string `json:"action"`

	// multiplier
	Multiplier float64 `json:"multiplier"`

	// Recipe Id
	RecipeId string `json:"recipe_id"`
}

// Nutrient defines model for Nutrient.
type Nutrient struct {

	// todo
	Id int `json:"id"`

	// todo
	Name     string           `json:"name"`
	UnitName FoodNutrientUnit `json:"unit_name"`
}

// PaginatedFoods defines model for PaginatedFoods.
type PaginatedFoods struct {
	Foods *[]Food `json:"foods,omitempty"`

	// A generic list (for pagination use)
	Meta *Items `json:"meta,omitempty"`
}

// PaginatedIngredients defines model for PaginatedIngredients.
type PaginatedIngredients struct {
	Ingredients *[]IngredientDetail `json:"ingredients,omitempty"`

	// A generic list (for pagination use)
	Meta *Items `json:"meta,omitempty"`
}

// PaginatedMeals defines model for PaginatedMeals.
type PaginatedMeals struct {
	Meals *[]Meal `json:"meals,omitempty"`

	// A generic list (for pagination use)
	Meta *Items `json:"meta,omitempty"`
}

// PaginatedPhotos defines model for PaginatedPhotos.
type PaginatedPhotos struct {

	// A generic list (for pagination use)
	Meta   *Items         `json:"meta,omitempty"`
	Photos *[]GooglePhoto `json:"photos,omitempty"`
}

// PaginatedRecipeWrappers defines model for PaginatedRecipeWrappers.
type PaginatedRecipeWrappers struct {

	// A generic list (for pagination use)
	Meta    *Items           `json:"meta,omitempty"`
	Recipes *[]RecipeWrapper `json:"recipes,omitempty"`
}

// PaginatedRecipes defines model for PaginatedRecipes.
type PaginatedRecipes struct {

	// A generic list (for pagination use)
	Meta    *Items    `json:"meta,omitempty"`
	Recipes *[]Recipe `json:"recipes,omitempty"`
}

// Recipe defines model for Recipe.
type Recipe struct {

	// id
	Id string `json:"id"`

	// all the versions of the recipe
	Versions []RecipeDetail `json:"versions"`
}

// RecipeDependency defines model for RecipeDependency.
type RecipeDependency struct {

	// id
	IngredientId   string         `json:"ingredient_id"`
	IngredientKind IngredientKind `json:"ingredient_kind"`

	// id
	IngredientName string `json:"ingredient_name"`

	// recipe_id
	RecipeId string `json:"recipe_id"`

	// id
	RecipeName string `json:"recipe_name"`
}

// RecipeDetail defines model for RecipeDetail.
type RecipeDetail struct {

	// id
	Id string `json:"id"`

	// whether or not it is the most recent version
	IsLatestVersion *bool `json:"is_latest_version,omitempty"`

	// recipe name
	Name string `json:"name"`

	// serving quantity
	Quantity int64 `json:"quantity"`

	// sections of the recipe
	Sections []RecipeSection `json:"sections"`

	// num servings
	Servings *int64 `json:"servings,omitempty"`

	// book or websites
	Sources *[]RecipeSource `json:"sources,omitempty"`

	// serving unit
	Unit string `json:"unit"`

	// version of the recipe
	Version *int64 `json:"version,omitempty"`
}

// RecipeSection defines model for RecipeSection.
type RecipeSection struct {

	// A range of time or a specific duration of time (in seconds)
	Duration *TimeRange `json:"duration,omitempty"`

	// id
	Id string `json:"id"`

	// x
	Ingredients []SectionIngredient `json:"ingredients"`

	// x
	Instructions []SectionInstruction `json:"instructions"`
}

// RecipeSource defines model for RecipeSource.
type RecipeSource struct {

	// page number/section (if book)
	Page *string `json:"page,omitempty"`

	// title (if book)
	Title *string `json:"title,omitempty"`

	// url
	Url *string `json:"url,omitempty"`
}

// RecipeWrapper defines model for RecipeWrapper.
type RecipeWrapper struct {

	// A revision of a recipe
	Detail RecipeDetail `json:"detail"`

	// id
	Id string `json:"id"`
}

// SearchResult defines model for SearchResult.
type SearchResult struct {

	// The ingredients
	Ingredients *[]Ingredient `json:"ingredients,omitempty"`

	// A generic list (for pagination use)
	Meta *Items `json:"meta,omitempty"`

	// The recipes
	Recipes *[]RecipeWrapper `json:"recipes,omitempty"`
}

// SectionIngredient defines model for SectionIngredient.
type SectionIngredient struct {

	// adjective
	Adjective *string `json:"adjective,omitempty"`

	// the various measures
	Amounts []Amount `json:"amounts"`

	// id
	Id string `json:"id"`

	// An Ingredient
	Ingredient *IngredientDetail `json:"ingredient,omitempty"`
	Kind       IngredientKind    `json:"kind"`

	// optional
	Optional *bool `json:"optional,omitempty"`

	// raw line item (pre-import/scrape)
	Original *string `json:"original,omitempty"`

	// A revision of a recipe
	Recipe *RecipeDetail `json:"recipe,omitempty"`

	// x
	Substitutes *[]SectionIngredient `json:"substitutes,omitempty"`
}

// SectionInstruction defines model for SectionInstruction.
type SectionInstruction struct {

	// id
	Id string `json:"id"`

	// instruction
	Instruction string `json:"instruction"`
}

// TimeRange defines model for TimeRange.
type TimeRange struct {

	// The maximum amount of seconds (if a range)
	Max int `json:"max"`

	// The minimum amount of seconds (or the total, if not a range)
	Min int `json:"min"`
}

// UnitMapping defines model for UnitMapping.
type UnitMapping struct {

	// amount and unit
	A Amount `json:"a"`

	// amount and unit
	B Amount `json:"b"`

	// source of the mapping
	Source string `json:"source"`
}

// LimitParam defines model for limitParam.
type LimitParam int

// NameParam defines model for nameParam.
type NameParam string

// OffsetParam defines model for offsetParam.
type OffsetParam int

// AuthLoginParams defines parameters for AuthLogin.
type AuthLoginParams struct {

	// Google code
	Code string `json:"code"`
}

// GetFoodsByIdsParams defines parameters for GetFoodsByIds.
type GetFoodsByIdsParams struct {

	// ids
	FdcId []int `json:"fdc_id"`
}

// SearchFoodsParams defines parameters for SearchFoods.
type SearchFoodsParams struct {

	// The number of items to skip before starting to collect the result set.
	Offset *OffsetParam `json:"offset,omitempty"`

	// The numbers of items to return.
	Limit *LimitParam `json:"limit,omitempty"`

	// The search query (name).
	Name NameParam `json:"name"`

	// The data types
	DataTypes *[]FoodDataType `json:"data_types,omitempty"`
}

// ListIngredientsParams defines parameters for ListIngredients.
type ListIngredientsParams struct {

	// The number of items to skip before starting to collect the result set.
	Offset *OffsetParam `json:"offset,omitempty"`

	// The numbers of items to return.
	Limit *LimitParam `json:"limit,omitempty"`

	// ids
	IngredientId *[]string `json:"ingredient_id,omitempty"`
}

// CreateIngredientsJSONBody defines parameters for CreateIngredients.
type CreateIngredientsJSONBody Ingredient

// AssociateFoodWithIngredientParams defines parameters for AssociateFoodWithIngredient.
type AssociateFoodWithIngredientParams struct {

	// The FDC id of the food to link to the ingredient
	FdcId int `json:"fdc_id"`
}

// MergeIngredientsJSONBody defines parameters for MergeIngredients.
type MergeIngredientsJSONBody struct {
	IngredientIds []string `json:"ingredient_ids"`
}

// ListMealsParams defines parameters for ListMeals.
type ListMealsParams struct {

	// The number of items to skip before starting to collect the result set.
	Offset *OffsetParam `json:"offset,omitempty"`

	// The numbers of items to return.
	Limit *LimitParam `json:"limit,omitempty"`
}

// UpdateRecipesForMealJSONBody defines parameters for UpdateRecipesForMeal.
type UpdateRecipesForMealJSONBody MealRecipeUpdate

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
type CreateRecipesJSONBody RecipeWrapper

// GetRecipesByIdsParams defines parameters for GetRecipesByIds.
type GetRecipesByIdsParams struct {

	// detail ids
	RecipeId []string `json:"recipe_id"`
}

// SearchParams defines parameters for Search.
type SearchParams struct {

	// The number of items to skip before starting to collect the result set.
	Offset *OffsetParam `json:"offset,omitempty"`

	// The numbers of items to return.
	Limit *LimitParam `json:"limit,omitempty"`

	// The search query (name).
	Name NameParam `json:"name"`
}

// CreateIngredientsRequestBody defines body for CreateIngredients for application/json ContentType.
type CreateIngredientsJSONRequestBody CreateIngredientsJSONBody

// MergeIngredientsRequestBody defines body for MergeIngredients for application/json ContentType.
type MergeIngredientsJSONRequestBody MergeIngredientsJSONBody

// UpdateRecipesForMealRequestBody defines body for UpdateRecipesForMeal for application/json ContentType.
type UpdateRecipesForMealJSONRequestBody UpdateRecipesForMealJSONBody

// CreateRecipesRequestBody defines body for CreateRecipes for application/json ContentType.
type CreateRecipesJSONRequestBody CreateRecipesJSONBody
