import gql from "graphql-tag";
import * as React from "react";
import * as ApolloReactCommon from "@apollo/react-common";
import * as ApolloReactComponents from "@apollo/react-components";
import * as ApolloReactHoc from "@apollo/react-hoc";
import * as ApolloReactHooks from "@apollo/react-hooks";
export type Maybe<T> = T | null;
export type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type Ingredient = {
  __typename?: "Ingredient";
  uuid: Scalars["String"];
  name: Scalars["String"];
  recipes?: Maybe<Array<Recipe>>;
};

export type SectionInstruction = {
  __typename?: "SectionInstruction";
  uuid: Scalars["String"];
  instruction: Scalars["String"];
};

export type IngredientInfo = Ingredient | Recipe;

export enum SectionIngredientKind {
  Recipe = "recipe",
  Ingredient = "ingredient",
}

export type SectionIngredient = {
  __typename?: "SectionIngredient";
  uuid: Scalars["String"];
  info: IngredientInfo;
  kind: SectionIngredientKind;
  grams: Scalars["Float"];
  amount: Scalars["Float"];
  unit: Scalars["String"];
  adjective: Scalars["String"];
  optional: Scalars["Boolean"];
};

export type Section = {
  __typename?: "Section";
  uuid: Scalars["String"];
  minutes: Scalars["Int"];
  instructions: Array<SectionInstruction>;
  ingredients: Array<SectionIngredient>;
};

export type Recipe = {
  __typename?: "Recipe";
  uuid: Scalars["String"];
  name: Scalars["String"];
  totalMinutes: Scalars["Int"];
  unit: Scalars["String"];
  sections: Array<Section>;
  meals: Array<Meal>;
};

export type Meal = {
  __typename?: "Meal";
  uuid: Scalars["String"];
  name: Scalars["String"];
  notionURL?: Maybe<Scalars["String"]>;
  imageURLs: Array<Scalars["String"]>;
};

export type RecipeInput = {
  uuid: Scalars["String"];
  name: Scalars["String"];
  totalMinutes?: Maybe<Scalars["Int"]>;
  unit?: Maybe<Scalars["String"]>;
  sections?: Maybe<Array<SectionInput>>;
};

export type SectionInstructionInput = {
  instruction: Scalars["String"];
};

export type SectionIngredientInput = {
  infoUUID: Scalars["String"];
  kind: SectionIngredientKind;
  grams: Scalars["Float"];
  amount?: Maybe<Scalars["Float"]>;
  unit?: Maybe<Scalars["String"]>;
  adjective?: Maybe<Scalars["String"]>;
  optional?: Maybe<Scalars["Boolean"]>;
};

export type SectionInput = {
  minutes: Scalars["Int"];
  instructions: Array<SectionInstructionInput>;
  ingredients: Array<SectionIngredientInput>;
};

export type NewRecipe = {
  name: Scalars["String"];
};

export type Nutrient = {
  __typename?: "Nutrient";
  id: Scalars["Int"];
  name: Scalars["String"];
  unitName: Scalars["String"];
};

export type FoodNutrientDerivation = {
  __typename?: "FoodNutrientDerivation";
  code: Scalars["String"];
  description: Scalars["String"];
};

export type FoodNutrient = {
  __typename?: "FoodNutrient";
  nutrient: Nutrient;
  amount: Scalars["Float"];
  dataPoints: Scalars["Int"];
};

export type FoodCategory = {
  __typename?: "FoodCategory";
  code: Scalars["String"];
  description: Scalars["String"];
};

export type Food = {
  __typename?: "Food";
  fdcId: Scalars["Int"];
  description: Scalars["String"];
  dataType: FoodDataType;
  category?: Maybe<FoodCategory>;
  nutrients: Array<FoodNutrient>;
  brandedInfo?: Maybe<BrandedFood>;
};

export type BrandedFood = {
  __typename?: "BrandedFood";
  brandOwner?: Maybe<Scalars["String"]>;
  ingredients?: Maybe<Scalars["String"]>;
  servingSize: Scalars["Float"];
  servingSizeUnit: Scalars["String"];
  householdServing?: Maybe<Scalars["String"]>;
  brandedFoodCategory?: Maybe<Scalars["String"]>;
};

export enum FoodDataType {
  FoundationFood = "foundation_food",
  SampleFood = "sample_food",
  MarketAcquisition = "market_acquisition",
  SurveyFnddsFood = "survey_fndds_food",
  SubSampleFood = "sub_sample_food",
  AgriculturalAcquisition = "agricultural_acquisition",
  SrLegacyFood = "sr_legacy_food",
  BrandedFood = "branded_food",
}

export type Mutation = {
  __typename?: "Mutation";
  createRecipe: Recipe;
  updateRecipe: Recipe;
  createIngredient: Ingredient;
  upsertIngredient: Scalars["String"];
};

export type MutationCreateRecipeArgs = {
  recipe?: Maybe<NewRecipe>;
};

export type MutationUpdateRecipeArgs = {
  recipe?: Maybe<RecipeInput>;
};

export type MutationCreateIngredientArgs = {
  name: Scalars["String"];
};

export type MutationUpsertIngredientArgs = {
  name: Scalars["String"];
  kind: SectionIngredientKind;
};

export type Query = {
  __typename?: "Query";
  recipes: Array<Recipe>;
  recipe?: Maybe<Recipe>;
  ingredients: Array<Ingredient>;
  ingredient?: Maybe<Ingredient>;
  food?: Maybe<Food>;
  foods?: Maybe<Array<Food>>;
};

export type QueryRecipesArgs = {
  searchQuery?: Scalars["String"];
};

export type QueryRecipeArgs = {
  uuid: Scalars["String"];
};

export type QueryIngredientsArgs = {
  searchQuery?: Scalars["String"];
};

export type QueryIngredientArgs = {
  uuid: Scalars["String"];
};

export type QueryFoodArgs = {
  fdcId: Scalars["Int"];
};

export type QueryFoodsArgs = {
  searchQuery: Scalars["String"];
  dataType?: Maybe<FoodDataType>;
  foodCategoryID?: Maybe<Scalars["Int"]>;
};

export type GetRecipeByUuidQueryVariables = {
  uuid: Scalars["String"];
};

export type GetRecipeByUuidQuery = { __typename?: "Query" } & {
  recipe?: Maybe<
    { __typename?: "Recipe" } & Pick<
      Recipe,
      "uuid" | "name" | "totalMinutes" | "unit"
    > & {
        meals: Array<
          { __typename?: "Meal" } & Pick<Meal, "uuid" | "name" | "imageURLs">
        >;
        sections: Array<
          { __typename?: "Section" } & Pick<Section, "minutes" | "uuid"> & {
              ingredients: Array<
                { __typename?: "SectionIngredient" } & Pick<
                  SectionIngredient,
                  | "uuid"
                  | "kind"
                  | "grams"
                  | "amount"
                  | "unit"
                  | "adjective"
                  | "optional"
                > & {
                    info:
                      | ({ __typename: "Ingredient" } & Pick<
                          Ingredient,
                          "uuid" | "name"
                        >)
                      | ({ __typename: "Recipe" } & Pick<
                          Recipe,
                          "uuid" | "name"
                        >);
                  }
              >;
              instructions: Array<
                { __typename?: "SectionInstruction" } & Pick<
                  SectionInstruction,
                  "instruction" | "uuid"
                >
              >;
            }
        >;
      }
  >;
};

export type GetRecipesQueryVariables = {};

export type GetRecipesQuery = { __typename?: "Query" } & {
  recipes: Array<
    { __typename?: "Recipe" } & Pick<
      Recipe,
      "uuid" | "name" | "totalMinutes" | "unit"
    >
  >;
};

export type UpdateRecipeMutationVariables = {
  recipe: RecipeInput;
};

export type UpdateRecipeMutation = { __typename?: "Mutation" } & {
  updateRecipe: { __typename?: "Recipe" } & Pick<Recipe, "uuid" | "name">;
};

export type CreateRecipeMutationVariables = {
  recipe: NewRecipe;
};

export type CreateRecipeMutation = { __typename?: "Mutation" } & {
  createRecipe: { __typename?: "Recipe" } & Pick<Recipe, "uuid" | "name">;
};

export type GetIngredientsQueryVariables = {};

export type GetIngredientsQuery = { __typename?: "Query" } & {
  ingredients: Array<
    { __typename?: "Ingredient" } & Pick<Ingredient, "uuid" | "name"> & {
        recipes?: Maybe<
          Array<{ __typename?: "Recipe" } & Pick<Recipe, "uuid" | "name">>
        >;
      }
  >;
};

export type GetFoodQueryVariables = {
  fdc_id: Scalars["Int"];
};

export type GetFoodQuery = { __typename?: "Query" } & {
  food?: Maybe<
    { __typename?: "Food" } & Pick<Food, "description" | "dataType"> & {
        category?: Maybe<
          { __typename?: "FoodCategory" } & Pick<
            FoodCategory,
            "code" | "description"
          >
        >;
        nutrients: Array<
          { __typename?: "FoodNutrient" } & Pick<FoodNutrient, "amount"> & {
              nutrient: { __typename?: "Nutrient" } & Pick<
                Nutrient,
                "name" | "unitName"
              >;
            }
        >;
      }
  >;
};

export type SearchIngredientsAndRecipesQueryVariables = {
  searchQuery: Scalars["String"];
};

export type SearchIngredientsAndRecipesQuery = { __typename?: "Query" } & {
  ingredients: Array<
    { __typename?: "Ingredient" } & Pick<Ingredient, "uuid" | "name">
  >;
  recipes: Array<{ __typename?: "Recipe" } & Pick<Recipe, "uuid" | "name">>;
};

export type CreateIngredientMutationVariables = {
  name: Scalars["String"];
};

export type CreateIngredientMutation = { __typename?: "Mutation" } & {
  createIngredient: { __typename?: "Ingredient" } & Pick<
    Ingredient,
    "uuid" | "name"
  >;
};

export const GetRecipeByUuidDocument = gql`
  query getRecipeByUUID($uuid: String!) {
    recipe(uuid: $uuid) {
      uuid
      name
      totalMinutes
      unit
      meals {
        uuid
        name
        imageURLs
      }
      sections {
        minutes
        uuid
        ingredients {
          uuid
          kind
          info {
            __typename
            ... on Ingredient {
              uuid
              name
            }
            ... on Recipe {
              uuid
              name
            }
          }
          grams
          amount
          unit
          adjective
          optional
        }
        instructions {
          instruction
          uuid
        }
      }
    }
  }
`;
export type GetRecipeByUuidComponentProps = Omit<
  ApolloReactComponents.QueryComponentOptions<
    GetRecipeByUuidQuery,
    GetRecipeByUuidQueryVariables
  >,
  "query"
> &
  (
    | { variables: GetRecipeByUuidQueryVariables; skip?: boolean }
    | { skip: boolean }
  );

export const GetRecipeByUuidComponent = (
  props: GetRecipeByUuidComponentProps
) => (
  <ApolloReactComponents.Query<
    GetRecipeByUuidQuery,
    GetRecipeByUuidQueryVariables
  >
    query={GetRecipeByUuidDocument}
    {...props}
  />
);

export type GetRecipeByUuidProps<
  TChildProps = {},
  TDataName extends string = "data"
> = {
  [key in TDataName]: ApolloReactHoc.DataValue<
    GetRecipeByUuidQuery,
    GetRecipeByUuidQueryVariables
  >;
} &
  TChildProps;
export function withGetRecipeByUuid<
  TProps,
  TChildProps = {},
  TDataName extends string = "data"
>(
  operationOptions?: ApolloReactHoc.OperationOption<
    TProps,
    GetRecipeByUuidQuery,
    GetRecipeByUuidQueryVariables,
    GetRecipeByUuidProps<TChildProps, TDataName>
  >
) {
  return ApolloReactHoc.withQuery<
    TProps,
    GetRecipeByUuidQuery,
    GetRecipeByUuidQueryVariables,
    GetRecipeByUuidProps<TChildProps, TDataName>
  >(GetRecipeByUuidDocument, {
    alias: "getRecipeByUuid",
    ...operationOptions,
  });
}

/**
 * __useGetRecipeByUuidQuery__
 *
 * To run a query within a React component, call `useGetRecipeByUuidQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetRecipeByUuidQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetRecipeByUuidQuery({
 *   variables: {
 *      uuid: // value for 'uuid'
 *   },
 * });
 */
export function useGetRecipeByUuidQuery(
  baseOptions?: ApolloReactHooks.QueryHookOptions<
    GetRecipeByUuidQuery,
    GetRecipeByUuidQueryVariables
  >
) {
  return ApolloReactHooks.useQuery<
    GetRecipeByUuidQuery,
    GetRecipeByUuidQueryVariables
  >(GetRecipeByUuidDocument, baseOptions);
}
export function useGetRecipeByUuidLazyQuery(
  baseOptions?: ApolloReactHooks.LazyQueryHookOptions<
    GetRecipeByUuidQuery,
    GetRecipeByUuidQueryVariables
  >
) {
  return ApolloReactHooks.useLazyQuery<
    GetRecipeByUuidQuery,
    GetRecipeByUuidQueryVariables
  >(GetRecipeByUuidDocument, baseOptions);
}
export type GetRecipeByUuidQueryHookResult = ReturnType<
  typeof useGetRecipeByUuidQuery
>;
export type GetRecipeByUuidLazyQueryHookResult = ReturnType<
  typeof useGetRecipeByUuidLazyQuery
>;
export type GetRecipeByUuidQueryResult = ApolloReactCommon.QueryResult<
  GetRecipeByUuidQuery,
  GetRecipeByUuidQueryVariables
>;
export const GetRecipesDocument = gql`
  query getRecipes {
    recipes {
      uuid
      name
      totalMinutes
      unit
    }
  }
`;
export type GetRecipesComponentProps = Omit<
  ApolloReactComponents.QueryComponentOptions<
    GetRecipesQuery,
    GetRecipesQueryVariables
  >,
  "query"
>;

export const GetRecipesComponent = (props: GetRecipesComponentProps) => (
  <ApolloReactComponents.Query<GetRecipesQuery, GetRecipesQueryVariables>
    query={GetRecipesDocument}
    {...props}
  />
);

export type GetRecipesProps<
  TChildProps = {},
  TDataName extends string = "data"
> = {
  [key in TDataName]: ApolloReactHoc.DataValue<
    GetRecipesQuery,
    GetRecipesQueryVariables
  >;
} &
  TChildProps;
export function withGetRecipes<
  TProps,
  TChildProps = {},
  TDataName extends string = "data"
>(
  operationOptions?: ApolloReactHoc.OperationOption<
    TProps,
    GetRecipesQuery,
    GetRecipesQueryVariables,
    GetRecipesProps<TChildProps, TDataName>
  >
) {
  return ApolloReactHoc.withQuery<
    TProps,
    GetRecipesQuery,
    GetRecipesQueryVariables,
    GetRecipesProps<TChildProps, TDataName>
  >(GetRecipesDocument, {
    alias: "getRecipes",
    ...operationOptions,
  });
}

/**
 * __useGetRecipesQuery__
 *
 * To run a query within a React component, call `useGetRecipesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetRecipesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetRecipesQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetRecipesQuery(
  baseOptions?: ApolloReactHooks.QueryHookOptions<
    GetRecipesQuery,
    GetRecipesQueryVariables
  >
) {
  return ApolloReactHooks.useQuery<GetRecipesQuery, GetRecipesQueryVariables>(
    GetRecipesDocument,
    baseOptions
  );
}
export function useGetRecipesLazyQuery(
  baseOptions?: ApolloReactHooks.LazyQueryHookOptions<
    GetRecipesQuery,
    GetRecipesQueryVariables
  >
) {
  return ApolloReactHooks.useLazyQuery<
    GetRecipesQuery,
    GetRecipesQueryVariables
  >(GetRecipesDocument, baseOptions);
}
export type GetRecipesQueryHookResult = ReturnType<typeof useGetRecipesQuery>;
export type GetRecipesLazyQueryHookResult = ReturnType<
  typeof useGetRecipesLazyQuery
>;
export type GetRecipesQueryResult = ApolloReactCommon.QueryResult<
  GetRecipesQuery,
  GetRecipesQueryVariables
>;
export const UpdateRecipeDocument = gql`
  mutation updateRecipe($recipe: RecipeInput!) {
    updateRecipe(recipe: $recipe) {
      uuid
      name
    }
  }
`;
export type UpdateRecipeMutationFn = ApolloReactCommon.MutationFunction<
  UpdateRecipeMutation,
  UpdateRecipeMutationVariables
>;
export type UpdateRecipeComponentProps = Omit<
  ApolloReactComponents.MutationComponentOptions<
    UpdateRecipeMutation,
    UpdateRecipeMutationVariables
  >,
  "mutation"
>;

export const UpdateRecipeComponent = (props: UpdateRecipeComponentProps) => (
  <ApolloReactComponents.Mutation<
    UpdateRecipeMutation,
    UpdateRecipeMutationVariables
  >
    mutation={UpdateRecipeDocument}
    {...props}
  />
);

export type UpdateRecipeProps<
  TChildProps = {},
  TDataName extends string = "mutate"
> = {
  [key in TDataName]: ApolloReactCommon.MutationFunction<
    UpdateRecipeMutation,
    UpdateRecipeMutationVariables
  >;
} &
  TChildProps;
export function withUpdateRecipe<
  TProps,
  TChildProps = {},
  TDataName extends string = "mutate"
>(
  operationOptions?: ApolloReactHoc.OperationOption<
    TProps,
    UpdateRecipeMutation,
    UpdateRecipeMutationVariables,
    UpdateRecipeProps<TChildProps, TDataName>
  >
) {
  return ApolloReactHoc.withMutation<
    TProps,
    UpdateRecipeMutation,
    UpdateRecipeMutationVariables,
    UpdateRecipeProps<TChildProps, TDataName>
  >(UpdateRecipeDocument, {
    alias: "updateRecipe",
    ...operationOptions,
  });
}

/**
 * __useUpdateRecipeMutation__
 *
 * To run a mutation, you first call `useUpdateRecipeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateRecipeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateRecipeMutation, { data, loading, error }] = useUpdateRecipeMutation({
 *   variables: {
 *      recipe: // value for 'recipe'
 *   },
 * });
 */
export function useUpdateRecipeMutation(
  baseOptions?: ApolloReactHooks.MutationHookOptions<
    UpdateRecipeMutation,
    UpdateRecipeMutationVariables
  >
) {
  return ApolloReactHooks.useMutation<
    UpdateRecipeMutation,
    UpdateRecipeMutationVariables
  >(UpdateRecipeDocument, baseOptions);
}
export type UpdateRecipeMutationHookResult = ReturnType<
  typeof useUpdateRecipeMutation
>;
export type UpdateRecipeMutationResult = ApolloReactCommon.MutationResult<
  UpdateRecipeMutation
>;
export type UpdateRecipeMutationOptions = ApolloReactCommon.BaseMutationOptions<
  UpdateRecipeMutation,
  UpdateRecipeMutationVariables
>;
export const CreateRecipeDocument = gql`
  mutation createRecipe($recipe: NewRecipe!) {
    createRecipe(recipe: $recipe) {
      uuid
      name
    }
  }
`;
export type CreateRecipeMutationFn = ApolloReactCommon.MutationFunction<
  CreateRecipeMutation,
  CreateRecipeMutationVariables
>;
export type CreateRecipeComponentProps = Omit<
  ApolloReactComponents.MutationComponentOptions<
    CreateRecipeMutation,
    CreateRecipeMutationVariables
  >,
  "mutation"
>;

export const CreateRecipeComponent = (props: CreateRecipeComponentProps) => (
  <ApolloReactComponents.Mutation<
    CreateRecipeMutation,
    CreateRecipeMutationVariables
  >
    mutation={CreateRecipeDocument}
    {...props}
  />
);

export type CreateRecipeProps<
  TChildProps = {},
  TDataName extends string = "mutate"
> = {
  [key in TDataName]: ApolloReactCommon.MutationFunction<
    CreateRecipeMutation,
    CreateRecipeMutationVariables
  >;
} &
  TChildProps;
export function withCreateRecipe<
  TProps,
  TChildProps = {},
  TDataName extends string = "mutate"
>(
  operationOptions?: ApolloReactHoc.OperationOption<
    TProps,
    CreateRecipeMutation,
    CreateRecipeMutationVariables,
    CreateRecipeProps<TChildProps, TDataName>
  >
) {
  return ApolloReactHoc.withMutation<
    TProps,
    CreateRecipeMutation,
    CreateRecipeMutationVariables,
    CreateRecipeProps<TChildProps, TDataName>
  >(CreateRecipeDocument, {
    alias: "createRecipe",
    ...operationOptions,
  });
}

/**
 * __useCreateRecipeMutation__
 *
 * To run a mutation, you first call `useCreateRecipeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateRecipeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createRecipeMutation, { data, loading, error }] = useCreateRecipeMutation({
 *   variables: {
 *      recipe: // value for 'recipe'
 *   },
 * });
 */
export function useCreateRecipeMutation(
  baseOptions?: ApolloReactHooks.MutationHookOptions<
    CreateRecipeMutation,
    CreateRecipeMutationVariables
  >
) {
  return ApolloReactHooks.useMutation<
    CreateRecipeMutation,
    CreateRecipeMutationVariables
  >(CreateRecipeDocument, baseOptions);
}
export type CreateRecipeMutationHookResult = ReturnType<
  typeof useCreateRecipeMutation
>;
export type CreateRecipeMutationResult = ApolloReactCommon.MutationResult<
  CreateRecipeMutation
>;
export type CreateRecipeMutationOptions = ApolloReactCommon.BaseMutationOptions<
  CreateRecipeMutation,
  CreateRecipeMutationVariables
>;
export const GetIngredientsDocument = gql`
  query getIngredients {
    ingredients {
      uuid
      name
      recipes {
        uuid
        name
      }
    }
  }
`;
export type GetIngredientsComponentProps = Omit<
  ApolloReactComponents.QueryComponentOptions<
    GetIngredientsQuery,
    GetIngredientsQueryVariables
  >,
  "query"
>;

export const GetIngredientsComponent = (
  props: GetIngredientsComponentProps
) => (
  <ApolloReactComponents.Query<
    GetIngredientsQuery,
    GetIngredientsQueryVariables
  >
    query={GetIngredientsDocument}
    {...props}
  />
);

export type GetIngredientsProps<
  TChildProps = {},
  TDataName extends string = "data"
> = {
  [key in TDataName]: ApolloReactHoc.DataValue<
    GetIngredientsQuery,
    GetIngredientsQueryVariables
  >;
} &
  TChildProps;
export function withGetIngredients<
  TProps,
  TChildProps = {},
  TDataName extends string = "data"
>(
  operationOptions?: ApolloReactHoc.OperationOption<
    TProps,
    GetIngredientsQuery,
    GetIngredientsQueryVariables,
    GetIngredientsProps<TChildProps, TDataName>
  >
) {
  return ApolloReactHoc.withQuery<
    TProps,
    GetIngredientsQuery,
    GetIngredientsQueryVariables,
    GetIngredientsProps<TChildProps, TDataName>
  >(GetIngredientsDocument, {
    alias: "getIngredients",
    ...operationOptions,
  });
}

/**
 * __useGetIngredientsQuery__
 *
 * To run a query within a React component, call `useGetIngredientsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetIngredientsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetIngredientsQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetIngredientsQuery(
  baseOptions?: ApolloReactHooks.QueryHookOptions<
    GetIngredientsQuery,
    GetIngredientsQueryVariables
  >
) {
  return ApolloReactHooks.useQuery<
    GetIngredientsQuery,
    GetIngredientsQueryVariables
  >(GetIngredientsDocument, baseOptions);
}
export function useGetIngredientsLazyQuery(
  baseOptions?: ApolloReactHooks.LazyQueryHookOptions<
    GetIngredientsQuery,
    GetIngredientsQueryVariables
  >
) {
  return ApolloReactHooks.useLazyQuery<
    GetIngredientsQuery,
    GetIngredientsQueryVariables
  >(GetIngredientsDocument, baseOptions);
}
export type GetIngredientsQueryHookResult = ReturnType<
  typeof useGetIngredientsQuery
>;
export type GetIngredientsLazyQueryHookResult = ReturnType<
  typeof useGetIngredientsLazyQuery
>;
export type GetIngredientsQueryResult = ApolloReactCommon.QueryResult<
  GetIngredientsQuery,
  GetIngredientsQueryVariables
>;
export const GetFoodDocument = gql`
  query getFood($fdc_id: Int!) {
    food(fdcId: $fdc_id) {
      description
      dataType
      category {
        code
        description
      }
      nutrients {
        nutrient {
          name
          unitName
        }
        amount
      }
    }
  }
`;
export type GetFoodComponentProps = Omit<
  ApolloReactComponents.QueryComponentOptions<
    GetFoodQuery,
    GetFoodQueryVariables
  >,
  "query"
> &
  ({ variables: GetFoodQueryVariables; skip?: boolean } | { skip: boolean });

export const GetFoodComponent = (props: GetFoodComponentProps) => (
  <ApolloReactComponents.Query<GetFoodQuery, GetFoodQueryVariables>
    query={GetFoodDocument}
    {...props}
  />
);

export type GetFoodProps<
  TChildProps = {},
  TDataName extends string = "data"
> = {
  [key in TDataName]: ApolloReactHoc.DataValue<
    GetFoodQuery,
    GetFoodQueryVariables
  >;
} &
  TChildProps;
export function withGetFood<
  TProps,
  TChildProps = {},
  TDataName extends string = "data"
>(
  operationOptions?: ApolloReactHoc.OperationOption<
    TProps,
    GetFoodQuery,
    GetFoodQueryVariables,
    GetFoodProps<TChildProps, TDataName>
  >
) {
  return ApolloReactHoc.withQuery<
    TProps,
    GetFoodQuery,
    GetFoodQueryVariables,
    GetFoodProps<TChildProps, TDataName>
  >(GetFoodDocument, {
    alias: "getFood",
    ...operationOptions,
  });
}

/**
 * __useGetFoodQuery__
 *
 * To run a query within a React component, call `useGetFoodQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetFoodQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetFoodQuery({
 *   variables: {
 *      fdc_id: // value for 'fdc_id'
 *   },
 * });
 */
export function useGetFoodQuery(
  baseOptions?: ApolloReactHooks.QueryHookOptions<
    GetFoodQuery,
    GetFoodQueryVariables
  >
) {
  return ApolloReactHooks.useQuery<GetFoodQuery, GetFoodQueryVariables>(
    GetFoodDocument,
    baseOptions
  );
}
export function useGetFoodLazyQuery(
  baseOptions?: ApolloReactHooks.LazyQueryHookOptions<
    GetFoodQuery,
    GetFoodQueryVariables
  >
) {
  return ApolloReactHooks.useLazyQuery<GetFoodQuery, GetFoodQueryVariables>(
    GetFoodDocument,
    baseOptions
  );
}
export type GetFoodQueryHookResult = ReturnType<typeof useGetFoodQuery>;
export type GetFoodLazyQueryHookResult = ReturnType<typeof useGetFoodLazyQuery>;
export type GetFoodQueryResult = ApolloReactCommon.QueryResult<
  GetFoodQuery,
  GetFoodQueryVariables
>;
export const SearchIngredientsAndRecipesDocument = gql`
  query searchIngredientsAndRecipes($searchQuery: String!) {
    ingredients(searchQuery: $searchQuery) {
      uuid
      name
    }
    recipes(searchQuery: $searchQuery) {
      uuid
      name
    }
  }
`;
export type SearchIngredientsAndRecipesComponentProps = Omit<
  ApolloReactComponents.QueryComponentOptions<
    SearchIngredientsAndRecipesQuery,
    SearchIngredientsAndRecipesQueryVariables
  >,
  "query"
> &
  (
    | { variables: SearchIngredientsAndRecipesQueryVariables; skip?: boolean }
    | { skip: boolean }
  );

export const SearchIngredientsAndRecipesComponent = (
  props: SearchIngredientsAndRecipesComponentProps
) => (
  <ApolloReactComponents.Query<
    SearchIngredientsAndRecipesQuery,
    SearchIngredientsAndRecipesQueryVariables
  >
    query={SearchIngredientsAndRecipesDocument}
    {...props}
  />
);

export type SearchIngredientsAndRecipesProps<
  TChildProps = {},
  TDataName extends string = "data"
> = {
  [key in TDataName]: ApolloReactHoc.DataValue<
    SearchIngredientsAndRecipesQuery,
    SearchIngredientsAndRecipesQueryVariables
  >;
} &
  TChildProps;
export function withSearchIngredientsAndRecipes<
  TProps,
  TChildProps = {},
  TDataName extends string = "data"
>(
  operationOptions?: ApolloReactHoc.OperationOption<
    TProps,
    SearchIngredientsAndRecipesQuery,
    SearchIngredientsAndRecipesQueryVariables,
    SearchIngredientsAndRecipesProps<TChildProps, TDataName>
  >
) {
  return ApolloReactHoc.withQuery<
    TProps,
    SearchIngredientsAndRecipesQuery,
    SearchIngredientsAndRecipesQueryVariables,
    SearchIngredientsAndRecipesProps<TChildProps, TDataName>
  >(SearchIngredientsAndRecipesDocument, {
    alias: "searchIngredientsAndRecipes",
    ...operationOptions,
  });
}

/**
 * __useSearchIngredientsAndRecipesQuery__
 *
 * To run a query within a React component, call `useSearchIngredientsAndRecipesQuery` and pass it any options that fit your needs.
 * When your component renders, `useSearchIngredientsAndRecipesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useSearchIngredientsAndRecipesQuery({
 *   variables: {
 *      searchQuery: // value for 'searchQuery'
 *   },
 * });
 */
export function useSearchIngredientsAndRecipesQuery(
  baseOptions?: ApolloReactHooks.QueryHookOptions<
    SearchIngredientsAndRecipesQuery,
    SearchIngredientsAndRecipesQueryVariables
  >
) {
  return ApolloReactHooks.useQuery<
    SearchIngredientsAndRecipesQuery,
    SearchIngredientsAndRecipesQueryVariables
  >(SearchIngredientsAndRecipesDocument, baseOptions);
}
export function useSearchIngredientsAndRecipesLazyQuery(
  baseOptions?: ApolloReactHooks.LazyQueryHookOptions<
    SearchIngredientsAndRecipesQuery,
    SearchIngredientsAndRecipesQueryVariables
  >
) {
  return ApolloReactHooks.useLazyQuery<
    SearchIngredientsAndRecipesQuery,
    SearchIngredientsAndRecipesQueryVariables
  >(SearchIngredientsAndRecipesDocument, baseOptions);
}
export type SearchIngredientsAndRecipesQueryHookResult = ReturnType<
  typeof useSearchIngredientsAndRecipesQuery
>;
export type SearchIngredientsAndRecipesLazyQueryHookResult = ReturnType<
  typeof useSearchIngredientsAndRecipesLazyQuery
>;
export type SearchIngredientsAndRecipesQueryResult = ApolloReactCommon.QueryResult<
  SearchIngredientsAndRecipesQuery,
  SearchIngredientsAndRecipesQueryVariables
>;
export const CreateIngredientDocument = gql`
  mutation createIngredient($name: String!) {
    createIngredient(name: $name) {
      uuid
      name
    }
  }
`;
export type CreateIngredientMutationFn = ApolloReactCommon.MutationFunction<
  CreateIngredientMutation,
  CreateIngredientMutationVariables
>;
export type CreateIngredientComponentProps = Omit<
  ApolloReactComponents.MutationComponentOptions<
    CreateIngredientMutation,
    CreateIngredientMutationVariables
  >,
  "mutation"
>;

export const CreateIngredientComponent = (
  props: CreateIngredientComponentProps
) => (
  <ApolloReactComponents.Mutation<
    CreateIngredientMutation,
    CreateIngredientMutationVariables
  >
    mutation={CreateIngredientDocument}
    {...props}
  />
);

export type CreateIngredientProps<
  TChildProps = {},
  TDataName extends string = "mutate"
> = {
  [key in TDataName]: ApolloReactCommon.MutationFunction<
    CreateIngredientMutation,
    CreateIngredientMutationVariables
  >;
} &
  TChildProps;
export function withCreateIngredient<
  TProps,
  TChildProps = {},
  TDataName extends string = "mutate"
>(
  operationOptions?: ApolloReactHoc.OperationOption<
    TProps,
    CreateIngredientMutation,
    CreateIngredientMutationVariables,
    CreateIngredientProps<TChildProps, TDataName>
  >
) {
  return ApolloReactHoc.withMutation<
    TProps,
    CreateIngredientMutation,
    CreateIngredientMutationVariables,
    CreateIngredientProps<TChildProps, TDataName>
  >(CreateIngredientDocument, {
    alias: "createIngredient",
    ...operationOptions,
  });
}

/**
 * __useCreateIngredientMutation__
 *
 * To run a mutation, you first call `useCreateIngredientMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateIngredientMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createIngredientMutation, { data, loading, error }] = useCreateIngredientMutation({
 *   variables: {
 *      name: // value for 'name'
 *   },
 * });
 */
export function useCreateIngredientMutation(
  baseOptions?: ApolloReactHooks.MutationHookOptions<
    CreateIngredientMutation,
    CreateIngredientMutationVariables
  >
) {
  return ApolloReactHooks.useMutation<
    CreateIngredientMutation,
    CreateIngredientMutationVariables
  >(CreateIngredientDocument, baseOptions);
}
export type CreateIngredientMutationHookResult = ReturnType<
  typeof useCreateIngredientMutation
>;
export type CreateIngredientMutationResult = ApolloReactCommon.MutationResult<
  CreateIngredientMutation
>;
export type CreateIngredientMutationOptions = ApolloReactCommon.BaseMutationOptions<
  CreateIngredientMutation,
  CreateIngredientMutationVariables
>;
