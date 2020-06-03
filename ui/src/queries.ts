import { gql } from "apollo-boost";

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const _ = [
  gql`
    query getRecipeByUUID($uuid: String!) {
      recipe(uuid: $uuid) {
        uuid
        name
        totalMinutes
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
        totalMinutes
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
  `,
];
