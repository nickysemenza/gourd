/* Generated by restful-react */

import React from "react";
import {
  Get,
  GetProps,
  useGet,
  UseGetProps,
  Mutate,
  MutateProps,
  useMutate,
  UseMutateProps,
} from "restful-react";
export const SPEC_VERSION = "1.0.0";
/**
 * Ingredients in a single section
 */
export interface SectionIngredient {
  /**
   * id
   */
  id: string;
  kind: IngredientKind;
  recipe?: RecipeDetail;
  ingredient?: IngredientDetail;
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
}

/**
 * Ingredients in a single section
 */
export interface SectionIngredientInput {
  /**
   * recipe/ingredient id
   */
  target_id: string;
  /**
   * recipe/ingredient name
   */
  name?: string;
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
}

/**
 * Instructions in a single section
 */
export interface SectionInstruction {
  /**
   * id
   */
  id: string;
  /**
   * instruction
   */
  instruction: string;
}

/**
 * Instructions in a single section
 */
export interface SectionInstructionInput {
  /**
   * instruction
   */
  instruction: string;
}

/**
 * A step in the recipe
 */
export interface RecipeSection {
  /**
   * id
   */
  id: string;
  duration?: TimeRange;
  /**
   * x
   */
  instructions: SectionInstruction[];
  /**
   * x
   */
  ingredients: SectionIngredient[];
}

/**
 * A step in the recipe
 */
export interface RecipeSectionInput {
  duration?: TimeRange;
  /**
   * x
   */
  instructions: SectionInstructionInput[];
  /**
   * x
   */
  ingredients: SectionIngredientInput[];
}

/**
 * A recipe with subcomponents
 */
export interface RecipeWrapper {
  /**
   * id
   */
  id: string;
  detail: RecipeDetail;
}

/**
 * A recipe with subcomponents
 */
export interface RecipeWrapperInput {
  /**
   * id
   */
  id: string;
  detail: RecipeDetailInput;
}

/**
 * A recipe with subcomponents
 */
export interface Recipe {
  /**
   * id
   */
  id: string;
  /**
   * all the versions of the recipe
   */
  versions: RecipeDetail[];
}

/**
 * A revision of a recipe
 */
export interface RecipeDetail {
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
  sources?: RecipeSource[];
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
  /**
   * version of the recipe
   */
  version: number;
  /**
   * whether or not it is the most recent version
   */
  is_latest_version: boolean;
  /**
   * when the version was created
   */
  created_at: string;
  /**
   * Other versions
   */
  other_versions?: RecipeDetail[];
}

/**
 * A revision of a recipe
 */
export interface RecipeDetailInput {
  /**
   * id
   */
  id: string;
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
}

/**
 * An Ingredient
 */
export interface Ingredient {
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
}

/**
 * An Ingredient
 */
export interface IngredientDetail {
  /**
   * Ingredient name
   */
  name: string;
  ingredient: Ingredient;
  /**
   * Recipes referencing this ingredient
   */
  recipes: RecipeDetail[];
  /**
   * Ingredients that are equivalent
   */
  children: IngredientDetail[];
  food?: Food;
  /**
   * mappings of equivalent units
   */
  unit_mappings: UnitMapping[];
}

/**
 * mappings
 */
export interface UnitMapping {
  a: Amount;
  b: Amount;
  /**
   * source of the mapping
   */
  source?: string;
}

/**
 * amount and unit
 */
export interface Amount {
  /**
   * unit
   */
  unit: string;
  /**
   * value
   */
  value: number;
  /**
   * if it was explicit, inferred, etc
   */
  source?: string;
}

/**
 * where the recipe came from (i.e. book/website)
 */
export interface RecipeSource {
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
}

/**
 * an album containing `GooglePhoto`
 */
export interface GooglePhotosAlbum {
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
}

/**
 * A google photo
 */
export interface GooglePhoto {
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
   */
  created: string;
  /**
   * width px
   */
  width: number;
  /**
   * height px
   */
  height: number;
}

/**
 * A meal, which bridges recipes to photos
 */
export interface Meal {
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
   */
  ate_at: string;
  photos: GooglePhoto[];
  recipes?: MealRecipe[];
}

/**
 * A recipe that's part of a meal (a recipe at a specific amount)
 */
export interface MealRecipe {
  /**
   * when it was taken
   */
  multiplier: number;
  recipe: RecipeDetail;
}

/**
 * A search result wrapper, which contains ingredients and recipes
 */
export interface SearchResult {
  /**
   * The ingredients
   */
  ingredients?: Ingredient[];
  /**
   * The recipes
   */
  recipes?: RecipeWrapper[];
  meta?: Items;
}

/**
 * A generic error message
 */
export interface Error {
  message: string;
}

/**
 * todo
 */
export interface AuthResp {
  user: { [key: string]: any };
  jwt: string;
}

/**
 * A generic list (for pagination use)
 */
export interface Items {
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
}

/**
 * A range of time or a specific duration of time (in seconds)
 */
export interface TimeRange {
  /**
   * The minimum amount of seconds (or the total, if not a range)
   */
  min: number;
  /**
   * The maximum amount of seconds (if a range)
   */
  max: number;
}

/**
 * pages of Recipe
 */
export interface PaginatedRecipes {
  recipes?: Recipe[];
  meta?: Items;
}

/**
 * pages of Recipe
 */
export interface PaginatedRecipeWrappers {
  recipes?: RecipeWrapper[];
  meta?: Items;
}

/**
 * pages of IngredientDetail
 */
export interface PaginatedIngredients {
  ingredients?: IngredientDetail[];
  meta?: Items;
}

/**
 * pages of GooglePhoto
 */
export interface PaginatedPhotos {
  photos?: GooglePhoto[];
  meta?: Items;
}

/**
 * pages of Meal
 */
export interface PaginatedMeals {
  meals?: Meal[];
  meta?: Items;
}

/**
 * pages of Food
 */
export interface PaginatedFoods {
  foods?: Food[];
  meta?: Items;
}

/**
 * todo
 */
export interface Nutrient {
  /**
   * todo
   */
  id: number;
  /**
   * todo
   */
  name: string;
  unit_name: FoodNutrientUnit;
}

/**
 * todo
 */
export interface FoodNutrient {
  nutrient: Nutrient;
  amount: number;
  data_points: number;
}

/**
 * food category, set for some
 */
export interface FoodCategory {
  /**
   * Food description
   */
  code: string;
  /**
   * Food description
   */
  description: string;
}

/**
 * branded_food
 */
export interface BrandedFood {
  brand_owner?: string;
  ingredients?: string;
  serving_size: number;
  serving_size_unit: string;
  household_serving?: string;
  branded_food_category?: string;
}

/**
 * food_portion
 */
export interface FoodPortion {
  id: number;
  amount: number;
  portion_description: string;
  modifier: string;
  gram_weight: number;
}

/**
 * A top level food
 */
export interface Food {
  /**
   * FDC Id
   */
  fdc_id: number;
  /**
   * Food description
   */
  description: string;
  data_type: FoodDataType;
  category?: FoodCategory;
  /**
   * todo
   */
  nutrients: FoodNutrient[];
  /**
   * portion datapoints
   */
  portions?: FoodPortion[];
  branded_info?: BrandedFood;
  /**
   * mappings of equivalent units
   */
  unit_mappings: UnitMapping[];
}

/**
 * an update to the recipes on a mea
 */
export interface MealRecipeUpdate {
  /**
   * Recipe Id
   */
  recipe_id: string;
  /**
   * multiplier
   */
  multiplier: number;
  /**
   * todo
   */
  action: "add" | "remove";
}

export type FoodDataType =
  | "foundation_food"
  | "sample_food"
  | "market_acquisition"
  | "survey_fndds_food"
  | "sub_sample_food"
  | "agricultural_acquisition"
  | "sr_legacy_food"
  | "branded_food";

export type FoodNutrientUnit =
  | "UG"
  | "G"
  | "IU"
  | "kJ"
  | "KCAL"
  | "MG"
  | "MG_ATE"
  | "SP_GR";

export type IngredientKind = "ingredient" | "recipe";

/**
 * node?
 */
export interface RecipeDependency {
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
  ingredient_kind: IngredientKind;
}

export interface UnitConversionRequest {
  target?: "weight" | "volume" | "money" | "calories" | "other";
  /**
   * multiple amounts to try
   */
  input: Amount[];
  /**
   * mappings of equivalent units
   */
  unit_mappings: UnitMapping[];
}

export interface AuthLoginQueryParams {
  /**
   * Google code
   */
  code: string;
}

export type AuthLoginProps = Omit<
  MutateProps<AuthResp, Error, AuthLoginQueryParams, void, void>,
  "path" | "verb"
>;

/**
 * Google Login callback
 *
 * Second step of https://developers.google.com/identity/sign-in/web/backend-auth#send-the-id-token-to-your-server
 */
export const AuthLogin = (props: AuthLoginProps) => (
  <Mutate<AuthResp, Error, AuthLoginQueryParams, void, void>
    verb="POST"
    path={`/auth`}
    {...props}
  />
);

export type UseAuthLoginProps = Omit<
  UseMutateProps<AuthResp, Error, AuthLoginQueryParams, void, void>,
  "path" | "verb"
>;

/**
 * Google Login callback
 *
 * Second step of https://developers.google.com/identity/sign-in/web/backend-auth#send-the-id-token-to-your-server
 */
export const useAuthLogin = (props: UseAuthLoginProps) =>
  useMutate<AuthResp, Error, AuthLoginQueryParams, void, void>(
    "POST",
    `/auth`,
    props
  );

export interface ListPhotosQueryParams {
  /**
   * The number of items to skip before starting to collect the result set.
   */
  offset?: number;
  /**
   * The numbers of items to return.
   */
  limit?: number;
}

export type ListPhotosProps = Omit<
  GetProps<PaginatedPhotos, Error, ListPhotosQueryParams, void>,
  "path"
>;

/**
 * List all photos
 *
 * todo
 */
export const ListPhotos = (props: ListPhotosProps) => (
  <Get<PaginatedPhotos, Error, ListPhotosQueryParams, void>
    path={`/photos`}
    {...props}
  />
);

export type UseListPhotosProps = Omit<
  UseGetProps<PaginatedPhotos, Error, ListPhotosQueryParams, void>,
  "path"
>;

/**
 * List all photos
 *
 * todo
 */
export const useListPhotos = (props: UseListPhotosProps) =>
  useGet<PaginatedPhotos, Error, ListPhotosQueryParams, void>(`/photos`, props);

export interface ListAllAlbumsResponse {
  /**
   * The list of albums
   */
  albums?: GooglePhotosAlbum[];
}

export type ListAllAlbumsProps = Omit<
  GetProps<ListAllAlbumsResponse, Error, void, void>,
  "path"
>;

/**
 * List all albums
 *
 * todo
 */
export const ListAllAlbums = (props: ListAllAlbumsProps) => (
  <Get<ListAllAlbumsResponse, Error, void, void> path={`/albums`} {...props} />
);

export type UseListAllAlbumsProps = Omit<
  UseGetProps<ListAllAlbumsResponse, Error, void, void>,
  "path"
>;

/**
 * List all albums
 *
 * todo
 */
export const useListAllAlbums = (props: UseListAllAlbumsProps) =>
  useGet<ListAllAlbumsResponse, Error, void, void>(`/albums`, props);

export interface SearchQueryParams {
  /**
   * The number of items to skip before starting to collect the result set.
   */
  offset?: number;
  /**
   * The numbers of items to return.
   */
  limit?: number;
  /**
   * The search query (name).
   */
  name: string;
}

export type SearchProps = Omit<
  GetProps<SearchResult, Error, SearchQueryParams, void>,
  "path"
>;

/**
 * Search recipes and ingredients
 *
 * todo
 */
export const Search = (props: SearchProps) => (
  <Get<SearchResult, Error, SearchQueryParams, void>
    path={`/search`}
    {...props}
  />
);

export type UseSearchProps = Omit<
  UseGetProps<SearchResult, Error, SearchQueryParams, void>,
  "path"
>;

/**
 * Search recipes and ingredients
 *
 * todo
 */
export const useSearch = (props: UseSearchProps) =>
  useGet<SearchResult, Error, SearchQueryParams, void>(`/search`, props);

export interface ListMealsQueryParams {
  /**
   * The number of items to skip before starting to collect the result set.
   */
  offset?: number;
  /**
   * The numbers of items to return.
   */
  limit?: number;
}

export type ListMealsProps = Omit<
  GetProps<PaginatedMeals, Error, ListMealsQueryParams, void>,
  "path"
>;

/**
 * List all meals
 *
 * todo
 */
export const ListMeals = (props: ListMealsProps) => (
  <Get<PaginatedMeals, Error, ListMealsQueryParams, void>
    path={`/meals`}
    {...props}
  />
);

export type UseListMealsProps = Omit<
  UseGetProps<PaginatedMeals, Error, ListMealsQueryParams, void>,
  "path"
>;

/**
 * List all meals
 *
 * todo
 */
export const useListMeals = (props: UseListMealsProps) =>
  useGet<PaginatedMeals, Error, ListMealsQueryParams, void>(`/meals`, props);

export interface GetMealByIdPathParams {
  /**
   * The id of the meal to retrieve
   */
  meal_id: string;
}

export type GetMealByIdProps = Omit<
  GetProps<Meal, Error, void, GetMealByIdPathParams>,
  "path"
> &
  GetMealByIdPathParams;

/**
 * Info for a specific meal
 *
 * todo
 */
export const GetMealById = ({ meal_id, ...props }: GetMealByIdProps) => (
  <Get<Meal, Error, void, GetMealByIdPathParams>
    path={`/meals/${meal_id}`}
    {...props}
  />
);

export type UseGetMealByIdProps = Omit<
  UseGetProps<Meal, Error, void, GetMealByIdPathParams>,
  "path"
> &
  GetMealByIdPathParams;

/**
 * Info for a specific meal
 *
 * todo
 */
export const useGetMealById = ({ meal_id, ...props }: UseGetMealByIdProps) =>
  useGet<Meal, Error, void, GetMealByIdPathParams>(
    (paramsInPath: GetMealByIdPathParams) => `/meals/${paramsInPath.meal_id}`,
    { pathParams: { meal_id }, ...props }
  );

export interface UpdateRecipesForMealPathParams {
  /**
   * The id of the meal to retrieve
   */
  meal_id: string;
}

export type UpdateRecipesForMealProps = Omit<
  MutateProps<
    Meal,
    Error,
    void,
    MealRecipeUpdate,
    UpdateRecipesForMealPathParams
  >,
  "path" | "verb"
> &
  UpdateRecipesForMealPathParams;

/**
 * Update the recipes associated with a given meal
 *
 * todo
 */
export const UpdateRecipesForMeal = ({
  meal_id,
  ...props
}: UpdateRecipesForMealProps) => (
  <Mutate<Meal, Error, void, MealRecipeUpdate, UpdateRecipesForMealPathParams>
    verb="PATCH"
    path={`/meals/${meal_id}/recipes`}
    {...props}
  />
);

export type UseUpdateRecipesForMealProps = Omit<
  UseMutateProps<
    Meal,
    Error,
    void,
    MealRecipeUpdate,
    UpdateRecipesForMealPathParams
  >,
  "path" | "verb"
> &
  UpdateRecipesForMealPathParams;

/**
 * Update the recipes associated with a given meal
 *
 * todo
 */
export const useUpdateRecipesForMeal = ({
  meal_id,
  ...props
}: UseUpdateRecipesForMealProps) =>
  useMutate<
    Meal,
    Error,
    void,
    MealRecipeUpdate,
    UpdateRecipesForMealPathParams
  >(
    "PATCH",
    (paramsInPath: UpdateRecipesForMealPathParams) =>
      `/meals/${paramsInPath.meal_id}/recipes`,
    { pathParams: { meal_id }, ...props }
  );

export interface ListIngredientsQueryParams {
  /**
   * The number of items to skip before starting to collect the result set.
   */
  offset?: number;
  /**
   * The numbers of items to return.
   */
  limit?: number;
  /**
   * ids
   */
  ingredient_id?: string[];
}

export type ListIngredientsProps = Omit<
  GetProps<PaginatedIngredients, Error, ListIngredientsQueryParams, void>,
  "path"
>;

/**
 * List all ingredients
 *
 * todo
 */
export const ListIngredients = (props: ListIngredientsProps) => (
  <Get<PaginatedIngredients, Error, ListIngredientsQueryParams, void>
    path={`/ingredients`}
    {...props}
  />
);

export type UseListIngredientsProps = Omit<
  UseGetProps<PaginatedIngredients, Error, ListIngredientsQueryParams, void>,
  "path"
>;

/**
 * List all ingredients
 *
 * todo
 */
export const useListIngredients = (props: UseListIngredientsProps) =>
  useGet<PaginatedIngredients, Error, ListIngredientsQueryParams, void>(
    `/ingredients`,
    props
  );

export type CreateIngredientsProps = Omit<
  MutateProps<Ingredient, Error, void, Ingredient, void>,
  "path" | "verb"
>;

/**
 * Create a ingredient
 *
 * todo
 */
export const CreateIngredients = (props: CreateIngredientsProps) => (
  <Mutate<Ingredient, Error, void, Ingredient, void>
    verb="POST"
    path={`/ingredients`}
    {...props}
  />
);

export type UseCreateIngredientsProps = Omit<
  UseMutateProps<Ingredient, Error, void, Ingredient, void>,
  "path" | "verb"
>;

/**
 * Create a ingredient
 *
 * todo
 */
export const useCreateIngredients = (props: UseCreateIngredientsProps) =>
  useMutate<Ingredient, Error, void, Ingredient, void>(
    "POST",
    `/ingredients`,
    props
  );

export interface ListRecipesQueryParams {
  /**
   * The number of items to skip before starting to collect the result set.
   */
  offset?: number;
  /**
   * The numbers of items to return.
   */
  limit?: number;
}

export type ListRecipesProps = Omit<
  GetProps<PaginatedRecipes, Error, ListRecipesQueryParams, void>,
  "path"
>;

/**
 * List all recipes
 *
 * todo
 */
export const ListRecipes = (props: ListRecipesProps) => (
  <Get<PaginatedRecipes, Error, ListRecipesQueryParams, void>
    path={`/recipes`}
    {...props}
  />
);

export type UseListRecipesProps = Omit<
  UseGetProps<PaginatedRecipes, Error, ListRecipesQueryParams, void>,
  "path"
>;

/**
 * List all recipes
 *
 * todo
 */
export const useListRecipes = (props: UseListRecipesProps) =>
  useGet<PaginatedRecipes, Error, ListRecipesQueryParams, void>(
    `/recipes`,
    props
  );

export type CreateRecipesProps = Omit<
  MutateProps<RecipeWrapper, Error, void, RecipeWrapperInput, void>,
  "path" | "verb"
>;

/**
 * Create a recipe
 *
 * todo
 */
export const CreateRecipes = (props: CreateRecipesProps) => (
  <Mutate<RecipeWrapper, Error, void, RecipeWrapperInput, void>
    verb="POST"
    path={`/recipes`}
    {...props}
  />
);

export type UseCreateRecipesProps = Omit<
  UseMutateProps<RecipeWrapper, Error, void, RecipeWrapperInput, void>,
  "path" | "verb"
>;

/**
 * Create a recipe
 *
 * todo
 */
export const useCreateRecipes = (props: UseCreateRecipesProps) =>
  useMutate<RecipeWrapper, Error, void, RecipeWrapperInput, void>(
    "POST",
    `/recipes`,
    props
  );

export interface GetRecipeByIdPathParams {
  /**
   * The id of the recipe to retrieve
   */
  recipe_id: string;
}

export type GetRecipeByIdProps = Omit<
  GetProps<RecipeWrapper, Error, void, GetRecipeByIdPathParams>,
  "path"
> &
  GetRecipeByIdPathParams;

/**
 * Info for a specific recipe
 *
 * todo
 */
export const GetRecipeById = ({ recipe_id, ...props }: GetRecipeByIdProps) => (
  <Get<RecipeWrapper, Error, void, GetRecipeByIdPathParams>
    path={`/recipes/${recipe_id}`}
    {...props}
  />
);

export type UseGetRecipeByIdProps = Omit<
  UseGetProps<RecipeWrapper, Error, void, GetRecipeByIdPathParams>,
  "path"
> &
  GetRecipeByIdPathParams;

/**
 * Info for a specific recipe
 *
 * todo
 */
export const useGetRecipeById = ({
  recipe_id,
  ...props
}: UseGetRecipeByIdProps) =>
  useGet<RecipeWrapper, Error, void, GetRecipeByIdPathParams>(
    (paramsInPath: GetRecipeByIdPathParams) =>
      `/recipes/${paramsInPath.recipe_id}`,
    { pathParams: { recipe_id }, ...props }
  );

export interface GetRecipesByIdsQueryParams {
  /**
   * detail ids
   */
  recipe_id: string[];
}

export type GetRecipesByIdsProps = Omit<
  GetProps<PaginatedRecipeWrappers, unknown, GetRecipesByIdsQueryParams, void>,
  "path"
>;

/**
 * Get recipes
 *
 * get recipes by ids
 */
export const GetRecipesByIds = (props: GetRecipesByIdsProps) => (
  <Get<PaginatedRecipeWrappers, unknown, GetRecipesByIdsQueryParams, void>
    path={`/recipes/bulk`}
    {...props}
  />
);

export type UseGetRecipesByIdsProps = Omit<
  UseGetProps<
    PaginatedRecipeWrappers,
    unknown,
    GetRecipesByIdsQueryParams,
    void
  >,
  "path"
>;

/**
 * Get recipes
 *
 * get recipes by ids
 */
export const useGetRecipesByIds = (props: UseGetRecipesByIdsProps) =>
  useGet<PaginatedRecipeWrappers, unknown, GetRecipesByIdsQueryParams, void>(
    `/recipes/bulk`,
    props
  );

export interface ConvertIngredientToRecipePathParams {
  /**
   * The id of the ingredient
   */
  ingredient_id: string;
}

export type ConvertIngredientToRecipeProps = Omit<
  MutateProps<
    RecipeDetail,
    Error,
    void,
    void,
    ConvertIngredientToRecipePathParams
  >,
  "path" | "verb"
> &
  ConvertIngredientToRecipePathParams;

/**
 * Converts an ingredient to a recipe, updating all recipes depending on it.
 *
 * todo
 */
export const ConvertIngredientToRecipe = ({
  ingredient_id,
  ...props
}: ConvertIngredientToRecipeProps) => (
  <Mutate<RecipeDetail, Error, void, void, ConvertIngredientToRecipePathParams>
    verb="POST"
    path={`/ingredients/${ingredient_id}/convert_to_recipe`}
    {...props}
  />
);

export type UseConvertIngredientToRecipeProps = Omit<
  UseMutateProps<
    RecipeDetail,
    Error,
    void,
    void,
    ConvertIngredientToRecipePathParams
  >,
  "path" | "verb"
> &
  ConvertIngredientToRecipePathParams;

/**
 * Converts an ingredient to a recipe, updating all recipes depending on it.
 *
 * todo
 */
export const useConvertIngredientToRecipe = ({
  ingredient_id,
  ...props
}: UseConvertIngredientToRecipeProps) =>
  useMutate<
    RecipeDetail,
    Error,
    void,
    void,
    ConvertIngredientToRecipePathParams
  >(
    "POST",
    (paramsInPath: ConvertIngredientToRecipePathParams) =>
      `/ingredients/${paramsInPath.ingredient_id}/convert_to_recipe`,
    { pathParams: { ingredient_id }, ...props }
  );

export interface AssociateFoodWithIngredientQueryParams {
  /**
   * The FDC id of the food to link to the ingredient
   */
  fdc_id: number;
}

export interface AssociateFoodWithIngredientPathParams {
  /**
   * The id of the ingredient
   */
  ingredient_id: string;
}

export type AssociateFoodWithIngredientProps = Omit<
  MutateProps<
    RecipeDetail,
    Error,
    AssociateFoodWithIngredientQueryParams,
    void,
    AssociateFoodWithIngredientPathParams
  >,
  "path" | "verb"
> &
  AssociateFoodWithIngredientPathParams;

/**
 * Assosiates a food with a given ingredient
 *
 * todo
 */
export const AssociateFoodWithIngredient = ({
  ingredient_id,
  ...props
}: AssociateFoodWithIngredientProps) => (
  <Mutate<
    RecipeDetail,
    Error,
    AssociateFoodWithIngredientQueryParams,
    void,
    AssociateFoodWithIngredientPathParams
  >
    verb="POST"
    path={`/ingredients/${ingredient_id}/associate_food`}
    {...props}
  />
);

export type UseAssociateFoodWithIngredientProps = Omit<
  UseMutateProps<
    RecipeDetail,
    Error,
    AssociateFoodWithIngredientQueryParams,
    void,
    AssociateFoodWithIngredientPathParams
  >,
  "path" | "verb"
> &
  AssociateFoodWithIngredientPathParams;

/**
 * Assosiates a food with a given ingredient
 *
 * todo
 */
export const useAssociateFoodWithIngredient = ({
  ingredient_id,
  ...props
}: UseAssociateFoodWithIngredientProps) =>
  useMutate<
    RecipeDetail,
    Error,
    AssociateFoodWithIngredientQueryParams,
    void,
    AssociateFoodWithIngredientPathParams
  >(
    "POST",
    (paramsInPath: AssociateFoodWithIngredientPathParams) =>
      `/ingredients/${paramsInPath.ingredient_id}/associate_food`,
    { pathParams: { ingredient_id }, ...props }
  );

export interface MergeIngredientsPathParams {
  /**
   * The id of the ingredient to merge into
   */
  ingredient_id: string;
}

export interface MergeIngredientsRequestBody {
  ingredient_ids: string[];
}

export type MergeIngredientsProps = Omit<
  MutateProps<
    Ingredient,
    Error,
    void,
    MergeIngredientsRequestBody,
    MergeIngredientsPathParams
  >,
  "path" | "verb"
> &
  MergeIngredientsPathParams;

/**
 * Merges the provide ingredients in the body into the param
 *
 * todo
 */
export const MergeIngredients = ({
  ingredient_id,
  ...props
}: MergeIngredientsProps) => (
  <Mutate<
    Ingredient,
    Error,
    void,
    MergeIngredientsRequestBody,
    MergeIngredientsPathParams
  >
    verb="POST"
    path={`/ingredients/${ingredient_id}/merge`}
    {...props}
  />
);

export type UseMergeIngredientsProps = Omit<
  UseMutateProps<
    Ingredient,
    Error,
    void,
    MergeIngredientsRequestBody,
    MergeIngredientsPathParams
  >,
  "path" | "verb"
> &
  MergeIngredientsPathParams;

/**
 * Merges the provide ingredients in the body into the param
 *
 * todo
 */
export const useMergeIngredients = ({
  ingredient_id,
  ...props
}: UseMergeIngredientsProps) =>
  useMutate<
    Ingredient,
    Error,
    void,
    MergeIngredientsRequestBody,
    MergeIngredientsPathParams
  >(
    "POST",
    (paramsInPath: MergeIngredientsPathParams) =>
      `/ingredients/${paramsInPath.ingredient_id}/merge`,
    { pathParams: { ingredient_id }, ...props }
  );

export interface GetFoodByIdPathParams {
  /**
   * The fdc id
   */
  fdc_id: number;
}

export type GetFoodByIdProps = Omit<
  GetProps<Food, Error, void, GetFoodByIdPathParams>,
  "path"
> &
  GetFoodByIdPathParams;

/**
 * get a FDC entry by id
 *
 * todo
 */
export const GetFoodById = ({ fdc_id, ...props }: GetFoodByIdProps) => (
  <Get<Food, Error, void, GetFoodByIdPathParams>
    path={`/foods/${fdc_id}`}
    {...props}
  />
);

export type UseGetFoodByIdProps = Omit<
  UseGetProps<Food, Error, void, GetFoodByIdPathParams>,
  "path"
> &
  GetFoodByIdPathParams;

/**
 * get a FDC entry by id
 *
 * todo
 */
export const useGetFoodById = ({ fdc_id, ...props }: UseGetFoodByIdProps) =>
  useGet<Food, Error, void, GetFoodByIdPathParams>(
    (paramsInPath: GetFoodByIdPathParams) => `/foods/${paramsInPath.fdc_id}`,
    { pathParams: { fdc_id }, ...props }
  );

export interface GetIngredientByIdPathParams {
  /**
   * The id of the ingredient to get into
   */
  ingredient_id: string;
}

export type GetIngredientByIdProps = Omit<
  GetProps<IngredientDetail, Error, void, GetIngredientByIdPathParams>,
  "path"
> &
  GetIngredientByIdPathParams;

/**
 * Get a specific ingredient
 *
 * todo
 */
export const GetIngredientById = ({
  ingredient_id,
  ...props
}: GetIngredientByIdProps) => (
  <Get<IngredientDetail, Error, void, GetIngredientByIdPathParams>
    path={`/ingredients/${ingredient_id}`}
    {...props}
  />
);

export type UseGetIngredientByIdProps = Omit<
  UseGetProps<IngredientDetail, Error, void, GetIngredientByIdPathParams>,
  "path"
> &
  GetIngredientByIdPathParams;

/**
 * Get a specific ingredient
 *
 * todo
 */
export const useGetIngredientById = ({
  ingredient_id,
  ...props
}: UseGetIngredientByIdProps) =>
  useGet<IngredientDetail, Error, void, GetIngredientByIdPathParams>(
    (paramsInPath: GetIngredientByIdPathParams) =>
      `/ingredients/${paramsInPath.ingredient_id}`,
    { pathParams: { ingredient_id }, ...props }
  );

export interface SearchFoodsQueryParams {
  /**
   * The number of items to skip before starting to collect the result set.
   */
  offset?: number;
  /**
   * The numbers of items to return.
   */
  limit?: number;
  /**
   * The search query (name).
   */
  name: string;
  /**
   * The data types
   */
  data_types?: FoodDataType[];
}

export type SearchFoodsProps = Omit<
  GetProps<PaginatedFoods, Error, SearchFoodsQueryParams, void>,
  "path"
>;

/**
 * Search foods
 *
 * todo
 */
export const SearchFoods = (props: SearchFoodsProps) => (
  <Get<PaginatedFoods, Error, SearchFoodsQueryParams, void>
    path={`/foods/search`}
    {...props}
  />
);

export type UseSearchFoodsProps = Omit<
  UseGetProps<PaginatedFoods, Error, SearchFoodsQueryParams, void>,
  "path"
>;

/**
 * Search foods
 *
 * todo
 */
export const useSearchFoods = (props: UseSearchFoodsProps) =>
  useGet<PaginatedFoods, Error, SearchFoodsQueryParams, void>(
    `/foods/search`,
    props
  );

export interface GetFoodsByIdsQueryParams {
  /**
   * ids
   */
  fdc_id: number[];
}

export type GetFoodsByIdsProps = Omit<
  GetProps<PaginatedFoods, unknown, GetFoodsByIdsQueryParams, void>,
  "path"
>;

/**
 * Get foods
 *
 * get foods by ids
 */
export const GetFoodsByIds = (props: GetFoodsByIdsProps) => (
  <Get<PaginatedFoods, unknown, GetFoodsByIdsQueryParams, void>
    path={`/foods/bulk`}
    {...props}
  />
);

export type UseGetFoodsByIdsProps = Omit<
  UseGetProps<PaginatedFoods, unknown, GetFoodsByIdsQueryParams, void>,
  "path"
>;

/**
 * Get foods
 *
 * get foods by ids
 */
export const useGetFoodsByIds = (props: UseGetFoodsByIdsProps) =>
  useGet<PaginatedFoods, unknown, GetFoodsByIdsQueryParams, void>(
    `/foods/bulk`,
    props
  );

export interface RecipeDependenciesResponse {
  /**
   * all
   */
  items?: RecipeDependency[];
}

export type RecipeDependenciesProps = Omit<
  GetProps<RecipeDependenciesResponse, unknown, void, void>,
  "path"
>;

/**
 * Get foods
 *
 * recipe dependencies
 */
export const RecipeDependencies = (props: RecipeDependenciesProps) => (
  <Get<RecipeDependenciesResponse, unknown, void, void>
    path={`/data/recipe_dependencies`}
    {...props}
  />
);

export type UseRecipeDependenciesProps = Omit<
  UseGetProps<RecipeDependenciesResponse, unknown, void, void>,
  "path"
>;

/**
 * Get foods
 *
 * recipe dependencies
 */
export const useRecipeDependencies = (props: UseRecipeDependenciesProps) =>
  useGet<RecipeDependenciesResponse, unknown, void, void>(
    `/data/recipe_dependencies`,
    props
  );
