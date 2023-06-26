/**
 * Generated by @openapi-codegen
 *
 * @version 1.0.0
 */
/**
 * Ingredients in a single section
 */
export type SectionIngredient = {
  /**
   * id
   */
  id: string;
  /**
   * what kind of ingredient
   */
  kind: IngredientKind;
  recipe?: RecipeDetail;
  ingredient?: IngredientWrapper;
  /**
   * the various measures
   */
  amounts: Amount[];
  /**
   * adjective
   */
  adjective?: string;
  /**
   * optional
   */
  optional?: boolean;
  /**
   * raw line item (pre-import/scrape)
   */
  original?: string;
  /**
   * x
   */
  substitutes?: SectionIngredient[];
};

/**
 * Ingredients in a single section
 */
export type SectionIngredientInput = {
  /**
   * recipe/ingredient id
   */
  target_id?: string;
  /**
   * recipe/ingredient name
   */
  name?: string;
  /**
   * what kind of ingredient, for target_id or name
   */
  kind: IngredientKind;
  /**
   * the various measures
   */
  amounts: Amount[];
  /**
   * adjective
   */
  adjective?: string;
  /**
   * optional
   */
  optional?: boolean;
  /**
   * raw line item (pre-import/scrape)
   */
  original?: string;
  /**
   * x
   */
  substitutes?: SectionIngredientInput[];
};

/**
 * Instructions in a single section
 */
export type SectionInstruction = {
  /**
   * id
   */
  id: string;
  /**
   * instruction
   */
  instruction: string;
};

/**
 * Instructions in a single section
 */
export type SectionInstructionInput = {
  /**
   * instruction
   */
  instruction: string;
};

/**
 * A step in the recipe
 */
export type RecipeSection = {
  /**
   * id
   */
  id: string;
  duration?: Amount;
  /**
   * x
   */
  instructions: SectionInstruction[];
  /**
   * x
   */
  ingredients: SectionIngredient[];
};

/**
 * A step in the recipe
 */
export type RecipeSectionInput = {
  duration?: Amount;
  /**
   * x
   */
  instructions: SectionInstructionInput[];
  /**
   * x
   */
  ingredients: SectionIngredientInput[];
};

/**
 * A recipe with subcomponents, including some "generated" fields to enhance data
 */
export type RecipeWrapper = {
  /**
   * id
   */
  id: string;
  detail: RecipeDetail;
  linked_meals?: Meal[];
  linked_photos?: Photo[];
  /**
   * Other versions
   */
  other_versions?: RecipeDetail[];
};

/**
 * A recipe with subcomponents
 */
export type RecipeWrapperInput = {
  /**
   * id
   */
  id?: string;
  detail: RecipeDetailInput;
};

/**
 * metadata about recipe detail
 */
export type RecipeDetailMeta = {
  /**
   * version of the recipe
   */
  version: number;
  /**
   * whether or not it is the most recent version
   */
  is_latest_version: boolean;
};

/**
 * A revision of a recipe. does not include any "generated" fields. everything directly from db
 */
export type RecipeDetail = {
  /**
   * id
   */
  id: string;
  /**
   * sections of the recipe
   */
  sections: RecipeSection[];
  /**
   * recipe name
   */
  name: string;
  /**
   * book or websites
   */
  sources: RecipeSource[];
  serving_info: RecipeServingInfo;
  meta: RecipeDetailMeta;
  /**
   * when the version was created
   *
   * @format date-time
   */
  created_at: string;
  /**
   * tags
   */
  tags: string[];
};

/**
 * A revision of a recipe
 */
export type RecipeDetailInput = {
  /**
   * sections of the recipe
   */
  sections: RecipeSectionInput[];
  /**
   * recipe name
   */
  name: string;
  /**
   * book or websites
   */
  sources?: RecipeSource[];
  serving_info: RecipeServingInfo;
  /**
   * when it created / updated
   *
   * @format date-time
   */
  date?: string;
  /**
   * tags
   */
  tags: string[];
};

/**
 * An Ingredient
 */
export type Ingredient = {
  /**
   * id
   */
  id: string;
  /**
   * Ingredient name
   */
  name: string;
  /**
   * ingredient ID for a similar (likely a different spelling)
   */
  parent?: string;
  /**
   * FDC id equivalent to this ingredient
   */
  fdc_id?: number;
};

/**
 * An Ingredient
 */
export type IngredientWrapper = {
  ingredient: Ingredient;
  /**
   * Recipes referencing this ingredient
   */
  recipes: RecipeDetail[];
  /**
   * Ingredients that are equivalent
   */
  children?: IngredientWrapper[];
  food?: TempFood;
  /**
   * mappings of equivalent units
   */
  unit_mappings: UnitMapping[];
};

/**
 * mappings
 */
export type UnitMapping = {
  a: Amount;
  b: Amount;
  /**
   * source of the mapping
   */
  source?: string;
};

/**
 * list of IngredientMapping
 */
export type IngredientMappingsPayload = {
  /**
   * mappings of equivalent units
   */
  ingredient_mappings: IngredientMapping[];
};

/**
 * details about ingredients
 */
export type IngredientMapping = {
  name: string;
  fdc_id?: number;
  aliases: string[];
  /**
   * mappings of equivalent units
   */
  unit_mappings: UnitMapping[];
};

/**
 * amount and unit
 */
export type Amount = {
  /**
   * unit
   */
  unit: string;
  /**
   * value
   *
   * @format double
   */
  value: number;
  /**
   * value
   *
   * @format double
   */
  upper_value?: number;
  /**
   * if it was explicit, inferred, etc
   */
  source?: string;
};

/**
 * recipe servings info
 */
export type RecipeServingInfo = {
  /**
   * num servings
   */
  servings?: number;
  /**
   * serving quantity
   */
  quantity: number;
  /**
   * serving unit
   */
  unit: string;
};

/**
 * where the recipe came from (i.e. book/website)
 */
export type RecipeSource = {
  /**
   * url
   */
  url?: string;
  /**
   * title (if book)
   */
  title?: string;
  /**
   * page number/section (if book)
   */
  page?: string;
  /**
   * image url
   */
  image_url?: string;
};

/**
 * an album containing `Photo`
 */
export type GooglePhotosAlbum = {
  /**
   * id
   */
  id: string;
  /**
   * title
   */
  title: string;
  /**
   * product_url
   */
  product_url: string;
  /**
   * usecase
   */
  usecase: string;
};

/**
 * A photo
 */
export type Photo = {
  /**
   * id
   */
  id: string;
  /**
   * public image
   */
  base_url: string;
  /**
   * blur hash
   */
  blur_hash?: string;
  /**
   * when it was taken
   *
   * @format date-time
   */
  taken_at?: string;
  /**
   * width px
   *
   * @format int64
   */
  width: number;
  /**
   * height px
   *
   * @format int64
   */
  height: number;
  /**
   * where the photo came from
   */
  source: "google" | "notion";
};

/**
 * A meal, which bridges recipes to photos
 */
export type Meal = {
  /**
   * id
   */
  id: string;
  /**
   * public image
   */
  name: string;
  /**
   * when it was taken
   *
   * @format date-time
   */
  ate_at: string;
  photos: Photo[];
  recipes?: MealRecipe[];
};

/**
 * A recipe that's part of a meal (a recipe at a specific amount)
 */
export type MealRecipe = {
  /**
   * when it was taken
   *
   * @format double
   */
  multiplier: number;
  recipe: RecipeDetail;
};

/**
 * A search result wrapper, which contains ingredients and recipes
 */
export type SearchResult = {
  /**
   * The ingredients
   */
  ingredients?: IngredientWrapper[];
  /**
   * The recipes
   */
  recipes?: RecipeWrapper[];
  meta?: Items;
};

/**
 * A generic error message
 *
 * @example {"message":"Something went wrong"}
 */
export type Error = {
  message: string;
  trace_id?: string;
};

/**
 * config data
 */
export type ConfigData = {
  google_scopes: string;
  google_client_id: string;
};

/**
 * todo
 */
export type AuthResp = {
  user: Record<string, any>;
  jwt: string;
};

/**
 * A generic list (for pagination use)
 */
export type Items = {
  /**
   * What number page this is
   */
  page_number: number;
  /**
   * How many items were requested for this page
   */
  limit: number;
  /**
   * todo
   */
  offset: number;
  /**
   * Total number of items across all pages
   */
  total_count: number;
  /**
   * Total number of pages available
   */
  page_count: number;
};

/**
 * pages of Recipe
 */
export type PaginatedRecipeWrappers = {
  recipes?: RecipeWrapper[];
  meta: Items;
};

/**
 * pages of IngredientWrapper
 */
export type PaginatedIngredients = {
  ingredients?: IngredientWrapper[];
  meta: Items;
};

/**
 * pages of Photos
 */
export type PaginatedPhotos = {
  photos?: Photo[];
  meta: Items;
};

/**
 * pages of Meal
 */
export type PaginatedMeals = {
  meals?: Meal[];
  meta: Items;
};

/**
 * pages of Food
 */
export type PaginatedFoods = {
  foods?: TempFood[];
  meta: Items;
};

/**
 * an update to the recipes on a mea
 */
export type MealRecipeUpdate = {
  /**
   * Recipe Id
   */
  recipe_id: string;
  /**
   * multiplier
   *
   * @format double
   * @minimum 0
   */
  multiplier: number;
  /**
   * todo
   */
  action: "add" | "remove";
};

export type IngredientKind = "ingredient" | "recipe";

/**
 * represents a relationship between recipe and ingredient, the latter of which can also be a recipe.
 */
export type RecipeDependency = {
  /**
   * recipe_id
   */
  recipe_id: string;
  /**
   * id
   */
  recipe_name: string;
  /**
   * id
   */
  ingredient_id: string;
  /**
   * id
   */
  ingredient_name: string;
  /**
   * what kind of ingredient
   */
  ingredient_kind: IngredientKind;
};

/**
 * holds name/id and multiplier for a Kind of entity
 */
export type EntitySummary = {
  /**
   * recipe_detail or ingredient id
   */
  id: string;
  /**
   * recipe or ingredient name
   */
  name: string;
  /**
   * multiplier
   *
   * @format double
   */
  multiplier: number;
  /**
   * what kind of entity
   */
  kind: IngredientKind;
};

/**
 * todo
 */
export type IngredientUsage = {
  /**
   * multiplier
   *
   * @format double
   */
  multiplier: number;
  /**
   * multiple amounts to try
   */
  amounts: Amount[];
  /**
   * mappings of equivalent units
   */
  required_by: EntitySummary[];
};

/**
 * holds information
 */
export type UsageValue = {
  /**
   * multiplier
   */
  ings: IngredientUsage[];
  /**
   * amounts
   */
  sum: Amount[];
  /**
   * detail about it
   */
  meta: EntitySummary;
};

export type UnitConversionRequest = {
  target?: "weight" | "volume" | "money" | "calories" | "other";
  /**
   * multiple amounts to try
   */
  input: Amount[];
  /**
   * mappings of equivalent units
   */
  unit_mappings: UnitMapping[];
};

export type SumsResponse = {
  /**
   * mappings of equivalent units
   */
  sums: UsageValue[];
  by_recipe: {
    [key: string]: UsageValue[];
  };
};

export type CompactRecipeSection = {
  ingredients: string[];
  instructions: string[];
};

export type Unused = {
  _compact?: CompactRecipe;
  _convert?: UnitConversionRequest;
};

export type CompactRecipe = {
  id: string;
  name: string;
  url?: string;
  image?: string;
  sections: CompactRecipeSection[];
};

/**
 * A meal, which bridges recipes to photos
 */
export type FoodSearchResult = {
  foods: TempFood[];
};

export type TempFood = {
  wrapper:
    | BrandedFoodItem
    | FoundationFoodItem
    | SRLegacyFoodItem
    | SurveyFoodItem;
  branded_food?: BrandedFoodItem;
  foundation_food?: FoundationFoodItem;
  legacy_food?: SRLegacyFoodItem;
  survey_food?: SurveyFoodItem;
  /**
   * mappings of equivalent units
   */
  unit_mappings: UnitMapping[];
  foodNutrients?: FoodNutrient[];
};

/**
 * a food nutrient
 */
export type Nutrient = {
  /**
   * @format int
   * @example 1005
   */
  id?: number;
  /**
   * @example 305
   */
  number?: string;
  /**
   * @example Carbohydrate, by difference
   */
  name?: string;
  /**
   * @format int
   * @example 1110
   */
  rank?: number;
  /**
   * @example g
   */
  unitName?: string;
};

export type FoodNutrientSource = {
  /**
   * @format int32
   * @example 9
   */
  id?: number;
  /**
   * @example 12
   */
  code?: string;
  /**
   * @example Manufacturer's analytical; partial documentation
   */
  description?: string;
};

export type FoodNutrientDerivation = {
  /**
   * @format int32
   * @example 75
   */
  id?: number;
  /**
   * @example LCCD
   */
  code?: string;
  /**
   * @example Calculated from a daily value percentage per serving size measure
   */
  description?: string;
  foodNutrientSource?: FoodNutrientSource;
};

export type NutrientAcquisitionDetails = {
  /**
   * @example 321632
   */
  sampleUnitId?: number;
  /**
   * @example 12/2/2005
   */
  purchaseDate?: string;
  /**
   * @example TRUSSVILLE
   */
  storeCity?: string;
  /**
   * @example AL
   */
  storeState?: string;
};

export type NutrientAnalysisDetails = {
  /**
   * @example 343866
   */
  subSampleId?: number;
  /**
   * @format double
   * @example 0
   */
  amount?: number;
  /**
   * @example 1005
   */
  nutrientId?: number;
  /**
   * @example 10.2135/cropsci2017.04.0244
   */
  labMethodDescription?: string;
  labMethodOriginalDescription?: string;
  /**
   * @format url
   * @example https://doi.org/10.2135/cropsci2017.04.0244
   */
  labMethodLink?: string;
  /**
   * @example DOI for Beans
   */
  labMethodTechnique?: string;
  nutrientAcquisitionDetails?: NutrientAcquisitionDetails[];
};

export type FoodNutrient = {
  /**
   * @format int
   * @example 167514
   */
  id: number;
  /**
   * @format double
   * @example 0
   */
  amount?: number;
  /**
   * @format int64
   * @example 49
   */
  dataPoints?: number;
  /**
   * @format double
   * @example 73.73
   */
  min?: number;
  /**
   * @format double
   * @example 91.8
   */
  max?: number;
  /**
   * @format double
   * @example 90.3
   */
  median?: number;
  /**
   * @example FoodNutrient
   */
  type?: string;
  nutrient?: Nutrient;
  foodNutrientDerivation?: FoodNutrientDerivation;
  nutrientAnalysisDetails?: NutrientAnalysisDetails;
};

export type FoodAttribute = {
  /**
   * @example 25117
   */
  id?: number;
  /**
   * @example 1
   */
  sequenceNumber?: number;
  /**
   * @example Moisture change: -5.0%
   */
  value?: string;
  FoodAttributeType?: {
    /**
     * @example 1002
     */
    id?: number;
    /**
     * @example Adjustments
     */
    name?: string;
    /**
     * @example Adjustments made to foods, including moisture and fat changes.
     */
    description?: string;
  };
};

export type FoodUpdateLog = {
  /**
   * @example 534358
   */
  fdcId?: number;
  /**
   * @example 8/18/2018
   */
  availableDate?: string;
  /**
   * @example Kar Nut Products Company
   */
  brandOwner?: string;
  /**
   * @example LI
   */
  dataSource?: string;
  /**
   * @example Branded
   */
  dataType?: string;
  /**
   * @example NUT 'N BERRY MIX
   */
  description?: string;
  /**
   * @example Branded
   */
  foodClass?: string;
  /**
   * @example 077034085228
   */
  gtinUpc?: string;
  /**
   * @example 1 ONZ
   */
  householdServingFullText?: string;
  /**
   * @example PEANUTS (PEANUTS, PEANUT AND/OR SUNFLOWER OIL). RAISINS. DRIED CRANBERRIES (CRANBERRIES, SUGAR, SUNFLOWER OIL). SUNFLOWER KERNELS AND ALMONDS (SUNFLOWER KERNELS AND ALMONDS, PEANUT AND/OR SUNFLOWER OIL).
   */
  ingredients?: string;
  /**
   * @example 8/18/2018
   */
  modifiedDate?: string;
  /**
   * @example 4/1/2019
   */
  publicationDate?: string;
  /**
   * @format double
   * @example 28
   */
  servingSize?: number;
  /**
   * @example g
   */
  servingSizeUnit?: string;
  /**
   * @example Popcorn, Peanuts, Seeds & Related Snacks
   */
  brandedFoodCategory?: string;
  /**
   * @example Nutrient Added, Nutrient Updated
   */
  changes?: string;
  foodAttributes?: FoodAttribute[];
};

export type BrandedFoodItem = {
  /**
   * @example 534358
   */
  fdcId: number;
  /**
   * @example 8/18/2018
   */
  availableDate?: string;
  /**
   * @example Kar Nut Products Company
   */
  brandOwner?: string;
  /**
   * @example LI
   */
  dataSource?: string;
  /**
   * @example Branded
   */
  dataType: string;
  /**
   * @example NUT 'N BERRY MIX
   */
  description: string;
  /**
   * @example Branded
   */
  foodClass?: string;
  /**
   * @example 077034085228
   */
  gtinUpc?: string;
  /**
   * @example 1 ONZ
   */
  householdServingFullText?: string;
  /**
   * @example PEANUTS (PEANUTS, PEANUT AND/OR SUNFLOWER OIL). RAISINS. DRIED CRANBERRIES (CRANBERRIES, SUGAR, SUNFLOWER OIL). SUNFLOWER KERNELS AND ALMONDS (SUNFLOWER KERNELS AND ALMONDS, PEANUT AND/OR SUNFLOWER OIL).
   */
  ingredients?: string;
  /**
   * @example 8/18/2018
   */
  modifiedDate?: string;
  /**
   * @example 4/1/2019
   */
  publicationDate?: string;
  /**
   * @format double
   * @example 28
   */
  servingSize?: number;
  /**
   * @example g
   */
  servingSizeUnit?: string;
  /**
   * @example UNPREPARED
   */
  preparationStateCode?: string;
  /**
   * @example Popcorn, Peanuts, Seeds & Related Snacks
   */
  brandedFoodCategory?: string;
  /**
   * @example CHILD_NUTRITION_FOOD_PROGRAMS
   * @example GROCERY
   */
  tradeChannel?: string[];
  /**
   * @example 50161800
   */
  gpcClassCode?: number;
  foodNutrients?: FoodNutrient[];
  foodUpdateLog?: FoodUpdateLog[];
  labelNutrients?: {
    fat?: {
      /**
       * @format double
       * @example 8.9992
       */
      value?: number;
    };
    saturatedFat?: {
      /**
       * @format double
       * @example 0.9996
       */
      value?: number;
    };
    transFat?: {
      /**
       * @format double
       * @example 0
       */
      value?: number;
    };
    cholesterol?: {
      /**
       * @format double
       * @example 0
       */
      value?: number;
    };
    sodium?: {
      /**
       * @format double
       * @example 0
       */
      value?: number;
    };
    carbohydrates?: {
      /**
       * @format double
       * @example 12.0008
       */
      value?: number;
    };
    fiber?: {
      /**
       * @format double
       * @example 1.988
       */
      value?: number;
    };
    sugars?: {
      /**
       * @format double
       * @example 7.9996
       */
      value?: number;
    };
    protein?: {
      /**
       * @format double
       * @example 4.0012
       */
      value?: number;
    };
    calcium?: {
      /**
       * @format double
       * @example 19.88
       */
      value?: number;
    };
    iron?: {
      /**
       * @format double
       * @example 0.7196
       */
      value?: number;
    };
    potassium?: {
      /**
       * @format double
       * @example 159.88
       */
      value?: number;
    };
    calories?: {
      /**
       * @format double
       * @example 140
       */
      value?: number;
    };
  };
};

export type FoodCategory = {
  /**
   * @format int32
   * @example 11
   */
  id?: number;
  /**
   * @example 1100
   */
  code?: string;
  /**
   * @example Vegetables and Vegetable Products
   */
  description?: string;
};

export type FoodComponent = {
  /**
   * @format int32
   * @example 59929
   */
  id?: number;
  /**
   * @example External fat
   */
  name?: string;
  /**
   * @example 24
   */
  dataPoints?: number;
  /**
   * @example 2.1
   */
  gramWeight?: number;
  /**
   * @example true
   */
  isRefuse?: boolean;
  /**
   * @example 2011
   */
  minYearAcquired?: number;
  /**
   * @example 0.5
   */
  percentWeight?: number;
};

export type MeasureUnit = {
  /**
   * @format int32
   * @example 999
   */
  id?: number;
  /**
   * @example undetermined
   */
  abbreviation?: string;
  /**
   * @example undetermined
   */
  name?: string;
};

export type FoodPortion = {
  /**
   * @format int32
   * @example 135806
   */
  id?: number;
  /**
   * @format double
   * @example 1
   */
  amount?: number;
  /**
   * @format int32
   * @example 9
   */
  dataPoints?: number;
  /**
   * @format double
   * @example 91
   */
  gramWeight?: number;
  /**
   * @example 2011
   */
  minYearAcquired?: number;
  /**
   * @example 10205
   */
  modifier?: string;
  /**
   * @example 1 cup
   */
  portionDescription?: string;
  /**
   * @example 1
   */
  sequenceNumber?: number;
  measureUnit?: MeasureUnit;
};

export type SampleFoodItem = {
  /**
   * @example 45551
   */
  fdcId: number;
  /**
   * @example Sample
   */
  datatype?: string;
  /**
   * @example Beef, Tenderloin Roast, select, roasted, comp5, lean (34BLTR)
   */
  description: string;
  /**
   * @example Composite
   */
  foodClass?: string;
  /**
   * @example 4/1/2019
   */
  publicationDate?: string;
  foodAttributes?: FoodCategory[];
};

/**
 * applies to Foundation foods. Not all inputFoods will have all fields.
 */
export type InputFoodFoundation = {
  /**
   * @example 45551
   */
  id?: number;
  /**
   * @example Beef, Tenderloin Roast, select, roasted, comp5, lean (34BLTR)
   */
  foodDescription?: string;
  inputFood?: SampleFoodItem;
};

export type NutrientConversionFactors = {
  /**
   * @example .ProteinConversionFactor
   */
  type?: string;
  /**
   * @format double
   * @example 6.25
   */
  value?: number;
};

export type FoundationFoodItem = {
  /**
   * @example 747448
   */
  fdcId: number;
  /**
   * @example Foundation
   */
  dataType: string;
  /**
   * @example Strawberries, raw
   */
  description: string;
  /**
   * @example FinalFood
   */
  foodClass?: string;
  /**
   * @example Source number reflects the actual number of samples analyzed for a nutrient. Repeat nutrient analyses may have been done on the same sample with the values shown.
   */
  footNote?: string;
  /**
   * @example false
   */
  isHistoricalReference?: boolean;
  /**
   * @example 9316
   */
  ndbNumber?: number;
  /**
   * @example 12/16/2019
   */
  publicationDate?: string;
  /**
   * @example Fragaria X ananassa
   */
  scientificName?: string;
  foodCategory?: FoodCategory;
  foodComponents?: FoodComponent[];
  foodNutrients?: FoodNutrient[];
  foodPortions?: FoodPortion[];
  inputFoods?: InputFoodFoundation[];
  nutrientConversionFactors?: NutrientConversionFactors[];
};

export type SRLegacyFoodItem = {
  /**
   * @example 170379
   */
  fdcId: number;
  /**
   * @example SR Legacy
   */
  dataType: string;
  /**
   * @example Broccoli, raw
   */
  description: string;
  /**
   * @example FinalFood
   */
  foodClass?: string;
  /**
   * @example true
   */
  isHistoricalReference?: boolean;
  /**
   * @example 11090
   */
  ndbNumber?: number;
  /**
   * @example 4/1/2019
   */
  publicationDate?: string;
  /**
   * @example Brassica oleracea var. italica
   */
  scientificName?: string;
  foodCategory?: FoodCategory;
  foodNutrients?: FoodNutrient[];
  nutrientConversionFactors?: NutrientConversionFactors[];
};

export type SurveyFoodItem = {
  /**
   * @example 337985
   */
  fdcId: number;
  /**
   * @example Survey (FNDDS)
   */
  dataType: string;
  /**
   * @example Beef curry
   */
  description: string;
  /**
   * @example 12/31/2014
   */
  endDate?: string;
  /**
   * @example Survey
   */
  foodClass?: string;
  /**
   * @example 27116100
   */
  foodCode?: string;
  /**
   * @example 4/1/2019
   */
  publicationDate?: string;
  /**
   * @example 1/1/2013
   */
  startDate?: string;
  foodAttributes?: FoodAttribute[];
  foodPortions?: FoodPortion[];
  inputFoods?: InputFoodSurvey[];
  wweiaFoodCategory?: WweiaFoodCategory;
};

export type RetentionFactor = {
  /**
   * @example 235
   */
  id?: number;
  /**
   * @example 3460
   */
  code?: number;
  /**
   * @example VEG, ROOTS, ETC, SAUTEED
   */
  description?: string;
};

/**
 * applies to Survey (FNDDS). Not all inputFoods will have all fields.
 */
export type InputFoodSurvey = {
  /**
   * @example 18146
   */
  id?: number;
  /**
   * @format double
   * @example 1.5
   */
  amount?: number;
  /**
   * @example Spices, curry powder
   */
  foodDescription?: string;
  /**
   * @example 2015
   */
  ingredientCode?: number;
  /**
   * @example Spices, curry powder
   */
  ingredientDescription?: string;
  /**
   * @format double
   * @example 9.45
   */
  ingredientWeight?: number;
  /**
   * @example 21000
   */
  portionCode?: string;
  /**
   * @example 1 tablespoon
   */
  portionDescription?: string;
  /**
   * @example 6
   */
  sequenceNumber?: number;
  /**
   * @example 0
   */
  surveyFlag?: number;
  /**
   * @example TB
   */
  unit?: string;
  inputFood?: SurveyFoodItem;
  retentionFactor?: RetentionFactor;
};

export type WweiaFoodCategory = {
  /**
   * @example 3002
   */
  wweiaFoodCategoryCode?: number;
  /**
   * @example Meat mixed dishes
   */
  wweiaFoodCategoryDescription?: string;
};
