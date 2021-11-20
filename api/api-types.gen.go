// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.2 DO NOT EDIT.
package api

import (
	"time"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for FoodDataType.
const (
	FoodDataTypeAgriculturalAcquisition FoodDataType = "agricultural_acquisition"

	FoodDataTypeBrandedFood FoodDataType = "branded_food"

	FoodDataTypeFoundationFood FoodDataType = "foundation_food"

	FoodDataTypeMarketAcquisition FoodDataType = "market_acquisition"

	FoodDataTypeSampleFood FoodDataType = "sample_food"

	FoodDataTypeSrLegacyFood FoodDataType = "sr_legacy_food"

	FoodDataTypeSubSampleFood FoodDataType = "sub_sample_food"

	FoodDataTypeSurveyFnddsFood FoodDataType = "survey_fndds_food"
)

// Defines values for FoodNutrientUnit.
const (
	FoodNutrientUnitG FoodNutrientUnit = "G"

	FoodNutrientUnitIU FoodNutrientUnit = "IU"

	FoodNutrientUnitKCAL FoodNutrientUnit = "KCAL"

	FoodNutrientUnitKJ FoodNutrientUnit = "kJ"

	FoodNutrientUnitMG FoodNutrientUnit = "MG"

	FoodNutrientUnitMGATE FoodNutrientUnit = "MG_ATE"

	FoodNutrientUnitSPGR FoodNutrientUnit = "SP_GR"

	FoodNutrientUnitUG FoodNutrientUnit = "UG"
)

// Defines values for IngredientKind.
const (
	IngredientKindIngredient IngredientKind = "ingredient"

	IngredientKindRecipe IngredientKind = "recipe"
)

// Defines values for MealRecipeUpdateAction.
const (
	MealRecipeUpdateActionAdd MealRecipeUpdateAction = "add"

	MealRecipeUpdateActionRemove MealRecipeUpdateAction = "remove"
)

// Defines values for PhotoSource.
const (
	PhotoSourceGoogle PhotoSource = "google"

	PhotoSourceNotion PhotoSource = "notion"
)

// Defines values for UnitConversionRequestTarget.
const (
	UnitConversionRequestTargetCalories UnitConversionRequestTarget = "calories"

	UnitConversionRequestTargetMoney UnitConversionRequestTarget = "money"

	UnitConversionRequestTargetOther UnitConversionRequestTarget = "other"

	UnitConversionRequestTargetVolume UnitConversionRequestTarget = "volume"

	UnitConversionRequestTargetWeight UnitConversionRequestTarget = "weight"
)

// amount and unit
type Amount struct {
	// if it was explicit, inferred, etc
	Source *string `json:"source,omitempty"`

	// unit
	Unit string `json:"unit"`

	// value
	UpperValue *float64 `json:"upper_value,omitempty"`

	// value
	Value float64 `json:"value"`
}

// todo
type AuthResp struct {
	Jwt  string                 `json:"jwt"`
	User map[string]interface{} `json:"user"`
}

// branded_food
type BrandedFood struct {
	BrandOwner          *string `json:"brand_owner,omitempty"`
	BrandedFoodCategory *string `json:"branded_food_category,omitempty"`
	HouseholdServing    *string `json:"household_serving,omitempty"`
	Ingredients         *string `json:"ingredients,omitempty"`
	ServingSize         float64 `json:"serving_size"`
	ServingSizeUnit     string  `json:"serving_size_unit"`
}

// A generic error message
type Error struct {
	Message string `json:"message"`
}

// A top level food
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

// food category, set for some
type FoodCategory struct {
	// Food description
	Code string `json:"code"`

	// Food description
	Description string `json:"description"`
}

// FoodDataType defines model for FoodDataType.
type FoodDataType string

// todo
type FoodNutrient struct {
	Amount     float64 `json:"amount"`
	DataPoints int     `json:"data_points"`

	// todo
	Nutrient Nutrient `json:"nutrient"`
}

// FoodNutrientUnit defines model for FoodNutrientUnit.
type FoodNutrientUnit string

// food_portion
type FoodPortion struct {
	Amount             float64 `json:"amount"`
	GramWeight         float64 `json:"gram_weight"`
	Id                 int     `json:"id"`
	Modifier           string  `json:"modifier"`
	PortionDescription string  `json:"portion_description"`
}

// an album containing `Photo`
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

// An Ingredient
type Ingredient struct {
	// FDC id equivalent to this ingredient
	FdcId *int64 `json:"fdc_id,omitempty"`

	// id
	Id string `json:"id"`

	// Ingredient name
	Name string `json:"name"`

	// ingredient ID for a similar (likely a different spelling)
	Parent *string `json:"parent,omitempty"`
}

// An Ingredient
type IngredientDetail struct {
	// Ingredients that are equivalent
	Children *[]IngredientDetail `json:"children,omitempty"`

	// A top level food
	Food *Food `json:"food,omitempty"`

	// An Ingredient
	Ingredient Ingredient `json:"ingredient"`

	// Ingredient name
	Name string `json:"name"`

	// Recipes referencing this ingredient
	Recipes []RecipeDetail `json:"recipes"`

	// mappings of equivalent units
	UnitMappings []UnitMapping `json:"unit_mappings"`
}

// IngredientKind defines model for IngredientKind.
type IngredientKind string

// A generic list (for pagination use)
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

// A meal, which bridges recipes to photos
type Meal struct {
	// when it was taken
	AteAt time.Time `json:"ate_at"`

	// id
	Id string `json:"id"`

	// public image
	Name    string        `json:"name"`
	Photos  []Photo       `json:"photos"`
	Recipes *[]MealRecipe `json:"recipes,omitempty"`
}

// A recipe that's part of a meal (a recipe at a specific amount)
type MealRecipe struct {
	// when it was taken
	Multiplier float64 `json:"multiplier"`

	// A revision of a recipe
	Recipe RecipeDetail `json:"recipe"`
}

// an update to the recipes on a mea
type MealRecipeUpdate struct {
	// todo
	Action MealRecipeUpdateAction `json:"action"`

	// multiplier
	Multiplier float64 `json:"multiplier"`

	// Recipe Id
	RecipeId string `json:"recipe_id"`
}

// todo
type MealRecipeUpdateAction string

// todo
type Nutrient struct {
	// todo
	Id int `json:"id"`

	// todo
	Name     string           `json:"name"`
	UnitName FoodNutrientUnit `json:"unit_name"`
}

// pages of Food
type PaginatedFoods struct {
	Foods *[]Food `json:"foods,omitempty"`

	// A generic list (for pagination use)
	Meta *Items `json:"meta,omitempty"`
}

// pages of IngredientDetail
type PaginatedIngredients struct {
	Ingredients *[]IngredientDetail `json:"ingredients,omitempty"`

	// A generic list (for pagination use)
	Meta *Items `json:"meta,omitempty"`
}

// pages of Meal
type PaginatedMeals struct {
	Meals *[]Meal `json:"meals,omitempty"`

	// A generic list (for pagination use)
	Meta *Items `json:"meta,omitempty"`
}

// pages of Photos
type PaginatedPhotos struct {
	// A generic list (for pagination use)
	Meta   *Items   `json:"meta,omitempty"`
	Photos *[]Photo `json:"photos,omitempty"`
}

// pages of Recipe
type PaginatedRecipeWrappers struct {
	// A generic list (for pagination use)
	Meta    *Items           `json:"meta,omitempty"`
	Recipes *[]RecipeWrapper `json:"recipes,omitempty"`
}

// A photo
type Photo struct {
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

	// where the photo came from
	Source PhotoSource `json:"source"`

	// width px
	Width int64 `json:"width"`
}

// where the photo came from
type PhotoSource string

// node?
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

// A revision of a recipe
type RecipeDetail struct {
	// when the version was created
	CreatedAt time.Time `json:"created_at"`

	// id
	Id string `json:"id"`

	// whether or not it is the most recent version
	IsLatestVersion bool `json:"is_latest_version"`

	// recipe name
	Name string `json:"name"`

	// Other versions
	OtherVersions *[]RecipeDetail `json:"other_versions,omitempty"`

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
	Version int64 `json:"version"`
}

// A revision of a recipe
type RecipeDetailInput struct {
	// when it created / updated
	Date *time.Time `json:"date,omitempty"`

	// recipe name
	Name string `json:"name"`

	// serving quantity
	Quantity int64 `json:"quantity"`

	// sections of the recipe
	Sections []RecipeSectionInput `json:"sections"`

	// num servings
	Servings *int64 `json:"servings,omitempty"`

	// book or websites
	Sources *[]RecipeSource `json:"sources,omitempty"`

	// serving unit
	Unit string `json:"unit"`
}

// A step in the recipe
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

// A step in the recipe
type RecipeSectionInput struct {
	// A range of time or a specific duration of time (in seconds)
	Duration *TimeRange `json:"duration,omitempty"`

	// x
	Ingredients []SectionIngredientInput `json:"ingredients"`

	// x
	Instructions []SectionInstructionInput `json:"instructions"`
}

// where the recipe came from (i.e. book/website)
type RecipeSource struct {
	// page number/section (if book)
	Page *string `json:"page,omitempty"`

	// title (if book)
	Title *string `json:"title,omitempty"`

	// url
	Url *string `json:"url,omitempty"`
}

// A recipe with subcomponents
type RecipeWrapper struct {
	// A revision of a recipe
	Detail RecipeDetail `json:"detail"`

	// id
	Id           string   `json:"id"`
	LinkedMeals  *[]Meal  `json:"linked_meals,omitempty"`
	LinkedPhotos *[]Photo `json:"linked_photos,omitempty"`
}

// A recipe with subcomponents
type RecipeWrapperInput struct {
	// A revision of a recipe
	Detail RecipeDetailInput `json:"detail"`

	// id
	Id *string `json:"id,omitempty"`
}

// A search result wrapper, which contains ingredients and recipes
type SearchResult struct {
	// The ingredients
	Ingredients *[]Ingredient `json:"ingredients,omitempty"`

	// A generic list (for pagination use)
	Meta *Items `json:"meta,omitempty"`

	// The recipes
	Recipes *[]RecipeWrapper `json:"recipes,omitempty"`
}

// Ingredients in a single section
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

// Ingredients in a single section
type SectionIngredientInput struct {
	// adjective
	Adjective *string `json:"adjective,omitempty"`

	// the various measures
	Amounts []Amount       `json:"amounts"`
	Kind    IngredientKind `json:"kind"`

	// recipe/ingredient name
	Name *string `json:"name,omitempty"`

	// optional
	Optional *bool `json:"optional,omitempty"`

	// raw line item (pre-import/scrape)
	Original *string `json:"original,omitempty"`

	// x
	Substitutes *[]SectionIngredientInput `json:"substitutes,omitempty"`

	// recipe/ingredient id
	TargetId *string `json:"target_id,omitempty"`
}

// Instructions in a single section
type SectionInstruction struct {
	// id
	Id string `json:"id"`

	// instruction
	Instruction string `json:"instruction"`
}

// Instructions in a single section
type SectionInstructionInput struct {
	// instruction
	Instruction string `json:"instruction"`
}

// A range of time or a specific duration of time (in seconds)
type TimeRange struct {
	// The maximum amount of seconds (if a range)
	Max int `json:"max"`

	// The minimum amount of seconds (or the total, if not a range)
	Min int `json:"min"`
}

// UnitConversionRequest defines model for UnitConversionRequest.
type UnitConversionRequest struct {
	// multiple amounts to try
	Input  []Amount                     `json:"input"`
	Target *UnitConversionRequestTarget `json:"target,omitempty"`

	// mappings of equivalent units
	UnitMappings []UnitMapping `json:"unit_mappings"`
}

// UnitConversionRequestTarget defines model for UnitConversionRequest.Target.
type UnitConversionRequestTarget string

// mappings
type UnitMapping struct {
	// amount and unit
	A Amount `json:"a"`

	// amount and unit
	B Amount `json:"b"`

	// source of the mapping
	Source *string `json:"source,omitempty"`
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
type CreateRecipesJSONBody RecipeWrapperInput

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

// CreateIngredientsJSONRequestBody defines body for CreateIngredients for application/json ContentType.
type CreateIngredientsJSONRequestBody CreateIngredientsJSONBody

// MergeIngredientsJSONRequestBody defines body for MergeIngredients for application/json ContentType.
type MergeIngredientsJSONRequestBody MergeIngredientsJSONBody

// UpdateRecipesForMealJSONRequestBody defines body for UpdateRecipesForMeal for application/json ContentType.
type UpdateRecipesForMealJSONRequestBody UpdateRecipesForMealJSONBody

// CreateRecipesJSONRequestBody defines body for CreateRecipes for application/json ContentType.
type CreateRecipesJSONRequestBody CreateRecipesJSONBody
