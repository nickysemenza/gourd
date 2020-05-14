import { GraphQLClient } from "graphql-request";
import { print } from "graphql";
import gql from "graphql-tag";
export type Maybe<T> = T | null;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type Ingredient = {
  uuid: Scalars["String"];
  name: Scalars["String"];
  recipes?: Maybe<Array<Recipe>>;
};

export type SectionInstruction = {
  uuid: Scalars["String"];
  instruction: Scalars["String"];
};

export type IngredientInfo = Ingredient | Recipe;

export type SectionIngredient = {
  uuid: Scalars["String"];
  info: IngredientInfo;
  grams: Scalars["Float"];
  amount: Scalars["Float"];
  unit: Scalars["String"];
  adjective: Scalars["String"];
  optional: Scalars["Boolean"];
};

export type Section = {
  uuid: Scalars["String"];
  minutes: Scalars["Int"];
  instructions: Array<SectionInstruction>;
  ingredients: Array<SectionIngredient>;
};

export type Recipe = {
  uuid: Scalars["String"];
  name: Scalars["String"];
  total_minutes: Scalars["Int"];
  unit: Scalars["String"];
  sections: Array<Section>;
};

export type RecipeInput = {
  uuid: Scalars["String"];
  name: Scalars["String"];
  total_minutes?: Maybe<Scalars["Int"]>;
  unit?: Maybe<Scalars["String"]>;
  sections?: Maybe<Array<SectionInput>>;
};

export type SectionInstructionInput = {
  instruction: Scalars["String"];
};

export type SectionIngredientInput = {
  name: Scalars["String"];
  grams: Scalars["Float"];
  amount: Scalars["Float"];
  unit: Scalars["String"];
  adjective: Scalars["String"];
  optional: Scalars["Boolean"];
};

export type SectionInput = {
  minutes: Scalars["Int"];
  instructions: Array<SectionInstructionInput>;
  ingredients: Array<SectionIngredientInput>;
};

export type NewRecipe = {
  name: Scalars["String"];
};

export type Mutation = {
  createRecipe: Recipe;
  updateRecipe: Recipe;
};

export type MutationCreateRecipeArgs = {
  recipe?: Maybe<NewRecipe>;
};

export type MutationUpdateRecipeArgs = {
  recipe?: Maybe<RecipeInput>;
};

export type Query = {
  recipes: Array<Recipe>;
  recipe?: Maybe<Recipe>;
  ingredients: Array<Ingredient>;
};

export type QueryRecipeArgs = {
  uuid: Scalars["String"];
};

export type GetRecipeByUuidQueryVariables = {
  uuid: Scalars["String"];
};

export type GetRecipeByUuidQuery = {
  recipe?: Maybe<
    { __typename: "Recipe" } & Pick<
      Recipe,
      "uuid" | "name" | "total_minutes" | "unit"
    > & {
        sections: Array<
          { __typename: "Section" } & Pick<Section, "minutes"> & {
              ingredients: Array<
                { __typename: "SectionIngredient" } & Pick<
                  SectionIngredient,
                  "uuid" | "grams" | "unit" | "amount"
                > & {
                    info:
                      | ({ __typename: "Ingredient" } & Pick<
                          Ingredient,
                          "name" | "uuid"
                        >)
                      | ({ __typename: "Recipe" } & Pick<
                          Recipe,
                          "name" | "uuid"
                        >);
                  }
              >;
              instructions: Array<
                { __typename: "SectionInstruction" } & Pick<
                  SectionInstruction,
                  "instruction" | "uuid"
                >
              >;
            }
        >;
      }
  >;
};

export type UpdateRecipeMutationVariables = {
  recipe: RecipeInput;
};

export type UpdateRecipeMutation = {
  updateRecipe: Pick<Recipe, "uuid" | "name">;
};

export const GetRecipeByUuidDocument = gql`
  query getRecipeByUUID($uuid: String!) {
    recipe(uuid: $uuid) {
      uuid
      name
      total_minutes
      unit
      sections {
        minutes
        ingredients {
          uuid
          info {
            __typename
            ... on Ingredient {
              name
              uuid
              __typename
            }
            ... on Recipe {
              name
              uuid
              __typename
            }
          }
          grams
          unit
          amount
          __typename
        }
        instructions {
          instruction
          uuid
          __typename
        }
        __typename
      }
      __typename
    }
  }
`;
export const UpdateRecipeDocument = gql`
  mutation updateRecipe($recipe: RecipeInput!) {
    updateRecipe(recipe: $recipe) {
      uuid
      name
    }
  }
`;

export type SdkFunctionWrapper = <T>(action: () => Promise<T>) => Promise<T>;

const defaultWrapper: SdkFunctionWrapper = (sdkFunction) => sdkFunction();
export function getSdk(
  client: GraphQLClient,
  withWrapper: SdkFunctionWrapper = defaultWrapper
) {
  return {
    getRecipeByUUID(
      variables: GetRecipeByUuidQueryVariables
    ): Promise<GetRecipeByUuidQuery> {
      return withWrapper(() =>
        client.request<GetRecipeByUuidQuery>(
          print(GetRecipeByUuidDocument),
          variables
        )
      );
    },
    updateRecipe(
      variables: UpdateRecipeMutationVariables
    ): Promise<UpdateRecipeMutation> {
      return withWrapper(() =>
        client.request<UpdateRecipeMutation>(
          print(UpdateRecipeDocument),
          variables
        )
      );
    },
  };
}
export type Sdk = ReturnType<typeof getSdk>;

export const GetRecipeByUuid = gql`
  query getRecipeByUUID($uuid: String!) {
    recipe(uuid: $uuid) {
      uuid
      name
      total_minutes
      unit
      sections {
        minutes
        ingredients {
          uuid
          info {
            __typename
            ... on Ingredient {
              name
              uuid
              __typename
            }
            ... on Recipe {
              name
              uuid
              __typename
            }
          }
          grams
          unit
          amount
          __typename
        }
        instructions {
          instruction
          uuid
          __typename
        }
        __typename
      }
      __typename
    }
  }
`;
export const UpdateRecipe = gql`
  mutation updateRecipe($recipe: RecipeInput!) {
    updateRecipe(recipe: $recipe) {
      uuid
      name
    }
  }
`;

export interface IntrospectionResultData {
  __schema: {
    types: {
      kind: string;
      name: string;
      possibleTypes: {
        name: string;
      }[];
    }[];
  };
}
const result: IntrospectionResultData = {
  __schema: {
    types: [
      {
        kind: "UNION",
        name: "IngredientInfo",
        possibleTypes: [
          {
            name: "Ingredient",
          },
          {
            name: "Recipe",
          },
        ],
      },
    ],
  },
};
export default result;
