import { gql } from "apollo-boost";

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const _ = [
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
        same {
          uuid
          name
        }
        usdaFood {
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
