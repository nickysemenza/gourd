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
        meals {
          uuid
          name
          imageURLs
        }
        source {
          meta
          name
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
  gql`
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
  `,
  gql`
    mutation createIngredient($name: String!, $kind: SectionIngredientKind!) {
      upsertIngredient(name: $name, kind: $kind)
    }
  `,
];
