// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package api

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
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

// BrandedFoodItem defines model for BrandedFoodItem.
type BrandedFoodItem struct {
	AvailableDate            *string          `json:"availableDate,omitempty"`
	BrandOwner               *string          `json:"brandOwner,omitempty"`
	BrandedFoodCategory      *string          `json:"brandedFoodCategory,omitempty"`
	DataSource               *string          `json:"dataSource,omitempty"`
	DataType                 string           `json:"dataType"`
	Description              string           `json:"description"`
	FdcId                    int              `json:"fdcId"`
	FoodClass                *string          `json:"foodClass,omitempty"`
	FoodNutrients            *[]FoodNutrient  `json:"foodNutrients,omitempty"`
	FoodUpdateLog            *[]FoodUpdateLog `json:"foodUpdateLog,omitempty"`
	GpcClassCode             *int             `json:"gpcClassCode,omitempty"`
	GtinUpc                  *string          `json:"gtinUpc,omitempty"`
	HouseholdServingFullText *string          `json:"householdServingFullText,omitempty"`
	Ingredients              *string          `json:"ingredients,omitempty"`
	LabelNutrients           *struct {
		Calcium *struct {
			Value *float64 `json:"value,omitempty"`
		} `json:"calcium,omitempty"`
		Calories *struct {
			Value *float64 `json:"value,omitempty"`
		} `json:"calories,omitempty"`
		Carbohydrates *struct {
			Value *float64 `json:"value,omitempty"`
		} `json:"carbohydrates,omitempty"`
		Cholesterol *struct {
			Value *float64 `json:"value,omitempty"`
		} `json:"cholesterol,omitempty"`
		Fat *struct {
			Value *float64 `json:"value,omitempty"`
		} `json:"fat,omitempty"`
		Fiber *struct {
			Value *float64 `json:"value,omitempty"`
		} `json:"fiber,omitempty"`
		Iron *struct {
			Value *float64 `json:"value,omitempty"`
		} `json:"iron,omitempty"`
		Potassium *struct {
			Value *float64 `json:"value,omitempty"`
		} `json:"potassium,omitempty"`
		Protein *struct {
			Value *float64 `json:"value,omitempty"`
		} `json:"protein,omitempty"`
		SaturatedFat *struct {
			Value *float64 `json:"value,omitempty"`
		} `json:"saturatedFat,omitempty"`
		Sodium *struct {
			Value *float64 `json:"value,omitempty"`
		} `json:"sodium,omitempty"`
		Sugars *struct {
			Value *float64 `json:"value,omitempty"`
		} `json:"sugars,omitempty"`
		TransFat *struct {
			Value *float64 `json:"value,omitempty"`
		} `json:"transFat,omitempty"`
	} `json:"labelNutrients,omitempty"`
	ModifiedDate         *string   `json:"modifiedDate,omitempty"`
	PreparationStateCode *string   `json:"preparationStateCode,omitempty"`
	PublicationDate      *string   `json:"publicationDate,omitempty"`
	ServingSize          *float64  `json:"servingSize,omitempty"`
	ServingSizeUnit      *string   `json:"servingSizeUnit,omitempty"`
	TradeChannel         *[]string `json:"tradeChannel,omitempty"`
}

// CompactRecipe defines model for CompactRecipe.
type CompactRecipe struct {
	Id       string                 `json:"id"`
	Image    *string                `json:"image,omitempty"`
	Name     string                 `json:"name"`
	Sections []CompactRecipeSection `json:"sections"`
	Url      *string                `json:"url,omitempty"`
}

// CompactRecipeSection defines model for CompactRecipeSection.
type CompactRecipeSection struct {
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
}

// config data
type ConfigData struct {
	GoogleClientId string `json:"google_client_id"`
	GoogleScopes   string `json:"google_scopes"`
}

// holds name/id and multiplier for a Kind of entity
type EntitySummary struct {
	// recipe_detail or ingredient id
	Id   string         `json:"id"`
	Kind IngredientKind `json:"kind"`

	// multiplier
	Multiplier float64 `json:"multiplier"`

	// recipe or ingredient name
	Name string `json:"name"`
}

// A generic error message
type Error struct {
	Message string  `json:"message"`
	TraceId *string `json:"trace_id,omitempty"`
}

// FoodAttribute defines model for FoodAttribute.
type FoodAttribute struct {
	FoodAttributeType *struct {
		Description *string `json:"description,omitempty"`
		Id          *int    `json:"id,omitempty"`
		Name        *string `json:"name,omitempty"`
	} `json:"FoodAttributeType,omitempty"`
	Id             *int    `json:"id,omitempty"`
	SequenceNumber *int    `json:"sequenceNumber,omitempty"`
	Value          *string `json:"value,omitempty"`
}

// FoodCategory defines model for FoodCategory.
type FoodCategory struct {
	Code        *string `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
	Id          *int32  `json:"id,omitempty"`
}

// FoodComponent defines model for FoodComponent.
type FoodComponent struct {
	DataPoints      *int     `json:"dataPoints,omitempty"`
	GramWeight      *float32 `json:"gramWeight,omitempty"`
	Id              *int32   `json:"id,omitempty"`
	IsRefuse        *bool    `json:"isRefuse,omitempty"`
	MinYearAcquired *int     `json:"minYearAcquired,omitempty"`
	Name            *string  `json:"name,omitempty"`
	PercentWeight   *float32 `json:"percentWeight,omitempty"`
}

// FoodNutrient defines model for FoodNutrient.
type FoodNutrient struct {
	Amount                 *float64                `json:"amount,omitempty"`
	DataPoints             *int64                  `json:"dataPoints,omitempty"`
	FoodNutrientDerivation *FoodNutrientDerivation `json:"foodNutrientDerivation,omitempty"`
	Id                     int                     `json:"id"`
	Max                    *float64                `json:"max,omitempty"`
	Median                 *float64                `json:"median,omitempty"`
	Min                    *float64                `json:"min,omitempty"`

	// a food nutrient
	Nutrient                *Nutrient                `json:"nutrient,omitempty"`
	NutrientAnalysisDetails *NutrientAnalysisDetails `json:"nutrientAnalysisDetails,omitempty"`
	Type                    *string                  `json:"type,omitempty"`
}

// FoodNutrientDerivation defines model for FoodNutrientDerivation.
type FoodNutrientDerivation struct {
	Code               *string             `json:"code,omitempty"`
	Description        *string             `json:"description,omitempty"`
	FoodNutrientSource *FoodNutrientSource `json:"foodNutrientSource,omitempty"`
	Id                 *int32              `json:"id,omitempty"`
}

// FoodNutrientSource defines model for FoodNutrientSource.
type FoodNutrientSource struct {
	Code        *string `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
	Id          *int32  `json:"id,omitempty"`
}

// FoodPortion defines model for FoodPortion.
type FoodPortion struct {
	Amount             *float64     `json:"amount,omitempty"`
	DataPoints         *int32       `json:"dataPoints,omitempty"`
	GramWeight         *float64     `json:"gramWeight,omitempty"`
	Id                 *int32       `json:"id,omitempty"`
	MeasureUnit        *MeasureUnit `json:"measureUnit,omitempty"`
	MinYearAcquired    *int         `json:"minYearAcquired,omitempty"`
	Modifier           *string      `json:"modifier,omitempty"`
	PortionDescription *string      `json:"portionDescription,omitempty"`
	SequenceNumber     *int         `json:"sequenceNumber,omitempty"`
}

// A meal, which bridges recipes to photos
type FoodSearchResult struct {
	Foods []TempFood `json:"foods"`
}

// FoodUpdateLog defines model for FoodUpdateLog.
type FoodUpdateLog struct {
	AvailableDate            *string          `json:"availableDate,omitempty"`
	BrandOwner               *string          `json:"brandOwner,omitempty"`
	BrandedFoodCategory      *string          `json:"brandedFoodCategory,omitempty"`
	Changes                  *string          `json:"changes,omitempty"`
	DataSource               *string          `json:"dataSource,omitempty"`
	DataType                 *string          `json:"dataType,omitempty"`
	Description              *string          `json:"description,omitempty"`
	FdcId                    *int             `json:"fdcId,omitempty"`
	FoodAttributes           *[]FoodAttribute `json:"foodAttributes,omitempty"`
	FoodClass                *string          `json:"foodClass,omitempty"`
	GtinUpc                  *string          `json:"gtinUpc,omitempty"`
	HouseholdServingFullText *string          `json:"householdServingFullText,omitempty"`
	Ingredients              *string          `json:"ingredients,omitempty"`
	ModifiedDate             *string          `json:"modifiedDate,omitempty"`
	PublicationDate          *string          `json:"publicationDate,omitempty"`
	ServingSize              *float64         `json:"servingSize,omitempty"`
	ServingSizeUnit          *string          `json:"servingSizeUnit,omitempty"`
}

// FoundationFoodItem defines model for FoundationFoodItem.
type FoundationFoodItem struct {
	DataType                  string                       `json:"dataType"`
	Description               string                       `json:"description"`
	FdcId                     int                          `json:"fdcId"`
	FoodCategory              *FoodCategory                `json:"foodCategory,omitempty"`
	FoodClass                 *string                      `json:"foodClass,omitempty"`
	FoodComponents            *[]FoodComponent             `json:"foodComponents,omitempty"`
	FoodNutrients             *[]FoodNutrient              `json:"foodNutrients,omitempty"`
	FoodPortions              *[]FoodPortion               `json:"foodPortions,omitempty"`
	FootNote                  *string                      `json:"footNote,omitempty"`
	InputFoods                *[]InputFoodFoundation       `json:"inputFoods,omitempty"`
	IsHistoricalReference     *bool                        `json:"isHistoricalReference,omitempty"`
	NdbNumber                 *int                         `json:"ndbNumber,omitempty"`
	NutrientConversionFactors *[]NutrientConversionFactors `json:"nutrientConversionFactors,omitempty"`
	PublicationDate           *string                      `json:"publicationDate,omitempty"`
	ScientificName            *string                      `json:"scientificName,omitempty"`
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
	FdcId *int `json:"fdc_id,omitempty"`

	// id
	Id string `json:"id"`

	// Ingredient name
	Name string `json:"name"`

	// ingredient ID for a similar (likely a different spelling)
	Parent *string `json:"parent,omitempty"`
}

// IngredientKind defines model for IngredientKind.
type IngredientKind string

// details about ingredients
type IngredientMapping struct {
	Aliases []string `json:"aliases"`
	FdcId   *int     `json:"fdc_id,omitempty"`
	Name    string   `json:"name"`

	// mappings of equivalent units
	UnitMappings []UnitMapping `json:"unit_mappings"`
}

// list of IngredientMapping
type IngredientMappingsPayload struct {
	// mappings of equivalent units
	IngredientMappings []IngredientMapping `json:"ingredient_mappings"`
}

// todo
type IngredientUsage struct {
	// multiple amounts to try
	Amounts []Amount `json:"amounts"`

	// multiplier
	Multiplier float64 `json:"multiplier"`

	// mappings of equivalent units
	RequiredBy []EntitySummary `json:"required_by"`
}

// An Ingredient
type IngredientWrapper struct {
	// Ingredients that are equivalent
	Children *[]IngredientWrapper `json:"children,omitempty"`
	Food     *TempFood            `json:"food,omitempty"`

	// An Ingredient
	Ingredient Ingredient `json:"ingredient"`

	// Recipes referencing this ingredient
	Recipes []RecipeDetail `json:"recipes"`

	// mappings of equivalent units
	UnitMappings []UnitMapping `json:"unit_mappings"`
}

// applies to Foundation foods. Not all inputFoods will have all fields.
type InputFoodFoundation struct {
	FoodDescription *string         `json:"foodDescription,omitempty"`
	Id              *int            `json:"id,omitempty"`
	InputFood       *SampleFoodItem `json:"inputFood,omitempty"`
}

// applies to Survey (FNDDS). Not all inputFoods will have all fields.
type InputFoodSurvey struct {
	Amount                *float64         `json:"amount,omitempty"`
	FoodDescription       *string          `json:"foodDescription,omitempty"`
	Id                    *int             `json:"id,omitempty"`
	IngredientCode        *int             `json:"ingredientCode,omitempty"`
	IngredientDescription *string          `json:"ingredientDescription,omitempty"`
	IngredientWeight      *float64         `json:"ingredientWeight,omitempty"`
	InputFood             *SurveyFoodItem  `json:"inputFood,omitempty"`
	PortionCode           *string          `json:"portionCode,omitempty"`
	PortionDescription    *string          `json:"portionDescription,omitempty"`
	RetentionFactor       *RetentionFactor `json:"retentionFactor,omitempty"`
	SequenceNumber        *int             `json:"sequenceNumber,omitempty"`
	SurveyFlag            *int             `json:"surveyFlag,omitempty"`
	Unit                  *string          `json:"unit,omitempty"`
}

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

	// A revision of a recipe. does not include any "generated" fields. everything directly from db
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

// MeasureUnit defines model for MeasureUnit.
type MeasureUnit struct {
	Abbreviation *string `json:"abbreviation,omitempty"`
	Id           *int32  `json:"id,omitempty"`
	Name         *string `json:"name,omitempty"`
}

// a food nutrient
type Nutrient struct {
	Id       *int    `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	Number   *string `json:"number,omitempty"`
	Rank     *int    `json:"rank,omitempty"`
	UnitName *string `json:"unitName,omitempty"`
}

// NutrientAcquisitionDetails defines model for NutrientAcquisitionDetails.
type NutrientAcquisitionDetails struct {
	PurchaseDate *string `json:"purchaseDate,omitempty"`
	SampleUnitId *int    `json:"sampleUnitId,omitempty"`
	StoreCity    *string `json:"storeCity,omitempty"`
	StoreState   *string `json:"storeState,omitempty"`
}

// NutrientAnalysisDetails defines model for NutrientAnalysisDetails.
type NutrientAnalysisDetails struct {
	Amount                       *float64                      `json:"amount,omitempty"`
	LabMethodDescription         *string                       `json:"labMethodDescription,omitempty"`
	LabMethodLink                *string                       `json:"labMethodLink,omitempty"`
	LabMethodOriginalDescription *string                       `json:"labMethodOriginalDescription,omitempty"`
	LabMethodTechnique           *string                       `json:"labMethodTechnique,omitempty"`
	NutrientAcquisitionDetails   *[]NutrientAcquisitionDetails `json:"nutrientAcquisitionDetails,omitempty"`
	NutrientId                   *int                          `json:"nutrientId,omitempty"`
	SubSampleId                  *int                          `json:"subSampleId,omitempty"`
}

// NutrientConversionFactors defines model for NutrientConversionFactors.
type NutrientConversionFactors struct {
	Type  *string  `json:"type,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// pages of Food
type PaginatedFoods struct {
	Foods *[]TempFood `json:"foods,omitempty"`

	// A generic list (for pagination use)
	Meta Items `json:"meta"`
}

// pages of IngredientWrapper
type PaginatedIngredients struct {
	Ingredients *[]IngredientWrapper `json:"ingredients,omitempty"`

	// A generic list (for pagination use)
	Meta Items `json:"meta"`
}

// pages of Meal
type PaginatedMeals struct {
	Meals *[]Meal `json:"meals,omitempty"`

	// A generic list (for pagination use)
	Meta Items `json:"meta"`
}

// pages of Photos
type PaginatedPhotos struct {
	// A generic list (for pagination use)
	Meta   Items    `json:"meta"`
	Photos *[]Photo `json:"photos,omitempty"`
}

// pages of Recipe
type PaginatedRecipeWrappers struct {
	// A generic list (for pagination use)
	Meta    Items            `json:"meta"`
	Recipes *[]RecipeWrapper `json:"recipes,omitempty"`
}

// A photo
type Photo struct {
	// public image
	BaseUrl string `json:"base_url"`

	// blur hash
	BlurHash *string `json:"blur_hash,omitempty"`

	// height px
	Height int64 `json:"height"`

	// id
	Id string `json:"id"`

	// where the photo came from
	Source PhotoSource `json:"source"`

	// when it was taken
	TakenAt *time.Time `json:"taken_at,omitempty"`

	// width px
	Width int64 `json:"width"`
}

// where the photo came from
type PhotoSource string

// represents a relationship between recipe and ingredient, the latter of which can also be a recipe.
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

// A revision of a recipe. does not include any "generated" fields. everything directly from db
type RecipeDetail struct {
	// when the version was created
	CreatedAt time.Time `json:"created_at"`

	// id
	Id string `json:"id"`

	// metadata about recipe detail
	Meta RecipeDetailMeta `json:"meta"`

	// recipe name
	Name string `json:"name"`

	// sections of the recipe
	Sections []RecipeSection `json:"sections"`

	// recipe servings info
	ServingInfo RecipeServingInfo `json:"serving_info"`

	// book or websites
	Sources []RecipeSource `json:"sources"`

	// tags
	Tags []string `json:"tags"`
}

// A revision of a recipe
type RecipeDetailInput struct {
	// when it created / updated
	Date *time.Time `json:"date,omitempty"`

	// recipe name
	Name string `json:"name"`

	// sections of the recipe
	Sections []RecipeSectionInput `json:"sections"`

	// recipe servings info
	ServingInfo RecipeServingInfo `json:"serving_info"`

	// book or websites
	Sources *[]RecipeSource `json:"sources,omitempty"`

	// tags
	Tags []string `json:"tags"`
}

// metadata about recipe detail
type RecipeDetailMeta struct {
	// whether or not it is the most recent version
	IsLatestVersion bool `json:"is_latest_version"`

	// version of the recipe
	Version int `json:"version"`
}

// A step in the recipe
type RecipeSection struct {
	// amount and unit
	Duration *Amount `json:"duration,omitempty"`

	// id
	Id string `json:"id"`

	// x
	Ingredients []SectionIngredient `json:"ingredients"`

	// x
	Instructions []SectionInstruction `json:"instructions"`
}

// A step in the recipe
type RecipeSectionInput struct {
	// amount and unit
	Duration *Amount `json:"duration,omitempty"`

	// x
	Ingredients []SectionIngredientInput `json:"ingredients"`

	// x
	Instructions []SectionInstructionInput `json:"instructions"`
}

// recipe servings info
type RecipeServingInfo struct {
	// serving quantity
	Quantity int `json:"quantity"`

	// num servings
	Servings *int `json:"servings,omitempty"`

	// serving unit
	Unit string `json:"unit"`
}

// where the recipe came from (i.e. book/website)
type RecipeSource struct {
	// image url
	ImageUrl *string `json:"image_url,omitempty"`

	// page number/section (if book)
	Page *string `json:"page,omitempty"`

	// title (if book)
	Title *string `json:"title,omitempty"`

	// url
	Url *string `json:"url,omitempty"`
}

// A recipe with subcomponents, including some "generated" fields to enhance data
type RecipeWrapper struct {
	// A revision of a recipe. does not include any "generated" fields. everything directly from db
	Detail RecipeDetail `json:"detail"`

	// id
	Id           string   `json:"id"`
	LinkedMeals  *[]Meal  `json:"linked_meals,omitempty"`
	LinkedPhotos *[]Photo `json:"linked_photos,omitempty"`

	// Other versions
	OtherVersions *[]RecipeDetail `json:"other_versions,omitempty"`
}

// A recipe with subcomponents
type RecipeWrapperInput struct {
	// A revision of a recipe
	Detail RecipeDetailInput `json:"detail"`

	// id
	Id *string `json:"id,omitempty"`
}

// RetentionFactor defines model for RetentionFactor.
type RetentionFactor struct {
	Code        *int    `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
	Id          *int    `json:"id,omitempty"`
}

// SRLegacyFoodItem defines model for SRLegacyFoodItem.
type SRLegacyFoodItem struct {
	DataType                  string                       `json:"dataType"`
	Description               string                       `json:"description"`
	FdcId                     int                          `json:"fdcId"`
	FoodCategory              *FoodCategory                `json:"foodCategory,omitempty"`
	FoodClass                 *string                      `json:"foodClass,omitempty"`
	FoodNutrients             *[]FoodNutrient              `json:"foodNutrients,omitempty"`
	IsHistoricalReference     *bool                        `json:"isHistoricalReference,omitempty"`
	NdbNumber                 *int                         `json:"ndbNumber,omitempty"`
	NutrientConversionFactors *[]NutrientConversionFactors `json:"nutrientConversionFactors,omitempty"`
	PublicationDate           *string                      `json:"publicationDate,omitempty"`
	ScientificName            *string                      `json:"scientificName,omitempty"`
}

// SampleFoodItem defines model for SampleFoodItem.
type SampleFoodItem struct {
	Datatype        *string         `json:"datatype,omitempty"`
	Description     string          `json:"description"`
	FdcId           int             `json:"fdcId"`
	FoodAttributes  *[]FoodCategory `json:"foodAttributes,omitempty"`
	FoodClass       *string         `json:"foodClass,omitempty"`
	PublicationDate *string         `json:"publicationDate,omitempty"`
}

// A search result wrapper, which contains ingredients and recipes
type SearchResult struct {
	// The ingredients
	Ingredients *[]IngredientWrapper `json:"ingredients,omitempty"`

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
	Ingredient *IngredientWrapper `json:"ingredient,omitempty"`
	Kind       IngredientKind     `json:"kind"`

	// optional
	Optional *bool `json:"optional,omitempty"`

	// raw line item (pre-import/scrape)
	Original *string `json:"original,omitempty"`

	// A revision of a recipe. does not include any "generated" fields. everything directly from db
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

// SumsResponse defines model for SumsResponse.
type SumsResponse struct {
	ByRecipe SumsResponse_ByRecipe `json:"by_recipe"`

	// mappings of equivalent units
	Sums []UsageValue `json:"sums"`
}

// SumsResponse_ByRecipe defines model for SumsResponse.ByRecipe.
type SumsResponse_ByRecipe struct {
	AdditionalProperties map[string][]UsageValue `json:"-"`
}

// SurveyFoodItem defines model for SurveyFoodItem.
type SurveyFoodItem struct {
	DataType          string             `json:"dataType"`
	Description       string             `json:"description"`
	EndDate           *string            `json:"endDate,omitempty"`
	FdcId             int                `json:"fdcId"`
	FoodAttributes    *[]FoodAttribute   `json:"foodAttributes,omitempty"`
	FoodClass         *string            `json:"foodClass,omitempty"`
	FoodCode          *string            `json:"foodCode,omitempty"`
	FoodPortions      *[]FoodPortion     `json:"foodPortions,omitempty"`
	InputFoods        *[]InputFoodSurvey `json:"inputFoods,omitempty"`
	PublicationDate   *string            `json:"publicationDate,omitempty"`
	StartDate         *string            `json:"startDate,omitempty"`
	WweiaFoodCategory *WweiaFoodCategory `json:"wweiaFoodCategory,omitempty"`
}

// TempFood defines model for TempFood.
type TempFood struct {
	BrandedFood    *BrandedFoodItem    `json:"branded_food,omitempty"`
	FoodNutrients  *[]FoodNutrient     `json:"foodNutrients,omitempty"`
	FoundationFood *FoundationFoodItem `json:"foundation_food,omitempty"`
	LegacyFood     *SRLegacyFoodItem   `json:"legacy_food,omitempty"`
	SurveyFood     *SurveyFoodItem     `json:"survey_food,omitempty"`

	// mappings of equivalent units
	UnitMappings []UnitMapping `json:"unit_mappings"`
	Wrapper      interface{}   `json:"wrapper"`
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

// holds information
type UsageValue struct {
	// multiplier
	Ings []IngredientUsage `json:"ings"`

	// holds name/id and multiplier for a Kind of entity
	Meta EntitySummary `json:"meta"`

	// amounts
	Sum []Amount `json:"sum"`
}

// WweiaFoodCategory defines model for WweiaFoodCategory.
type WweiaFoodCategory struct {
	WweiaFoodCategoryCode        *int    `json:"wweiaFoodCategoryCode,omitempty"`
	WweiaFoodCategoryDescription *string `json:"wweiaFoodCategoryDescription,omitempty"`
}

// Unused defines model for _unused.
type Unused struct {
	Compact *CompactRecipe         `json:"_compact,omitempty"`
	Convert *UnitConversionRequest `json:"_convert,omitempty"`
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

// LoadIngredientMappingsJSONBody defines parameters for LoadIngredientMappings.
type LoadIngredientMappingsJSONBody IngredientMappingsPayload

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

// ScrapeRecipeJSONBody defines parameters for ScrapeRecipe.
type ScrapeRecipeJSONBody struct {
	Url string `json:"url"`
}

// SumRecipesJSONBody defines parameters for SumRecipes.
type SumRecipesJSONBody struct {
	Inputs []EntitySummary `json:"inputs"`
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

// DoSyncParams defines parameters for DoSync.
type DoSyncParams struct {
	// how many days to lookback
	LookbackDays int `json:"lookback_days"`
}

// CreateIngredientsJSONRequestBody defines body for CreateIngredients for application/json ContentType.
type CreateIngredientsJSONRequestBody CreateIngredientsJSONBody

// MergeIngredientsJSONRequestBody defines body for MergeIngredients for application/json ContentType.
type MergeIngredientsJSONRequestBody MergeIngredientsJSONBody

// UpdateRecipesForMealJSONRequestBody defines body for UpdateRecipesForMeal for application/json ContentType.
type UpdateRecipesForMealJSONRequestBody UpdateRecipesForMealJSONBody

// LoadIngredientMappingsJSONRequestBody defines body for LoadIngredientMappings for application/json ContentType.
type LoadIngredientMappingsJSONRequestBody LoadIngredientMappingsJSONBody

// CreateRecipesJSONRequestBody defines body for CreateRecipes for application/json ContentType.
type CreateRecipesJSONRequestBody CreateRecipesJSONBody

// ScrapeRecipeJSONRequestBody defines body for ScrapeRecipe for application/json ContentType.
type ScrapeRecipeJSONRequestBody ScrapeRecipeJSONBody

// SumRecipesJSONRequestBody defines body for SumRecipes for application/json ContentType.
type SumRecipesJSONRequestBody SumRecipesJSONBody

// Getter for additional properties for SumsResponse_ByRecipe. Returns the specified
// element and whether it was found
func (a SumsResponse_ByRecipe) Get(fieldName string) (value []UsageValue, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for SumsResponse_ByRecipe
func (a *SumsResponse_ByRecipe) Set(fieldName string, value []UsageValue) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string][]UsageValue)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for SumsResponse_ByRecipe to handle AdditionalProperties
func (a *SumsResponse_ByRecipe) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string][]UsageValue)
		for fieldName, fieldBuf := range object {
			var fieldVal []UsageValue
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for SumsResponse_ByRecipe to handle AdditionalProperties
func (a SumsResponse_ByRecipe) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}
