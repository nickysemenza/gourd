import { gql } from "apollo-boost";

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const _ = [
  gql`
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
              }
              ... on Recipe {
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
  `,
  gql`
    query getRecipes {
      recipes {
        uuid
        name
        total_minutes
        unit
      }
    }
  `,
  gql`
    mutation updateRecipe($recipe: RecipeInput!) {
      updateRecipe(recipe: $recipe) {
        uuid
        name
      }
    }
  `,
  gql`
    mutation createRecipe($recipe: NewRecipe!) {
      createRecipe(recipe: $recipe) {
        uuid
        name
      }
    }
  `,
  gql`
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
  `,
  gql`
    query getFood($fdc_id: Int!) {
      food(fdc_id: $fdc_id) {
        description
        data_type
        category {
          code
          description
        }
        nutrients {
          nutrient {
            name
            unit_name
          }
          amount
        }
      }
    }
  `,
];
