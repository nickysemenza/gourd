use super::ingredient::Ingredient;
use graphql_client::{GraphQLQuery, Response};
use serde::{Deserialize, Serialize};

#[derive(Debug, PartialEq, Default, Deserialize, Serialize)]
#[serde(default)]
pub struct Recipe {
    pub uuid: String,
    pub name: String,
    pub total_minutes: i64,
    pub unit: String,
    pub sections: Vec<RecipeSection>,
}
#[derive(Debug, PartialEq, Default, Deserialize, Serialize, Clone)]
pub struct RecipeSection {
    pub minutes: i64,
    pub ingredients: Vec<String>,
    pub instructions: Vec<String>,
}

impl Recipe {
    fn to_api_recipe(&self, c: &APIClient) -> update_recipe::RecipeInput {
        let sections: Vec<update_recipe::SectionInput> = self
            .sections
            .iter()
            .map(|section| update_recipe::SectionInput {
                minutes: section.minutes,
                ingredients: section
                    .ingredients
                    .iter()
                    .map(|i_str| {
                        let i = Ingredient::from_str(i_str).unwrap();
                        return update_recipe::SectionIngredientInput {
                            kind: if i.is_recipe {
                                update_recipe::SectionIngredientKind::recipe
                            } else {
                                update_recipe::SectionIngredientKind::ingredient
                            },
                            info_uuid: c.get_uuid(i.name, i.is_recipe),
                            grams: i.grams,
                            amount: Some(i.amount),
                            unit: Some(i.unit.to_string()),
                            adjective: Some(i.adjective.to_string()),
                            optional: Some(i.optional),
                        };
                    })
                    .collect(),
                instructions: section
                    .instructions
                    .iter()
                    .map(|i| update_recipe::SectionInstructionInput {
                        instruction: i.to_string(),
                    })
                    .collect(),
            })
            .collect();
        let ri = update_recipe::RecipeInput {
            name: self.name.to_string(),
            uuid: self.uuid.to_string(),
            unit: Some(self.unit.to_string()),
            total_minutes: Some(self.total_minutes),
            sections: Some(sections),
        };
        return ri;
    }
    fn from_api_recipe(
        recipe: get_recipe_by_uuid::GetRecipeByUuidRecipe,
    ) -> Result<Self, &'static str> {
        let rs: &Vec<RecipeSection> = &recipe
        .sections
        .iter()
        .map(|section| RecipeSection {
            minutes: section.minutes,
            instructions: section.instructions.iter().map(|i| i.instruction.to_string()).collect(),
            ingredients: section
                .ingredients
                .iter()
                .map(|i| {
                    // let mut name = String::new();
                    let (name, is_recipe) = match &i.info {
                            get_recipe_by_uuid::GetRecipeByUuidRecipeSectionsIngredientsInfo::Recipe(val) => {
                               (val.name.to_string(), true)
                            }
                            get_recipe_by_uuid::GetRecipeByUuidRecipeSectionsIngredientsInfo::Ingredient(val) => {
                               (val.name.to_string(), false)
                            }
                        };
                    Ingredient {
                        name,
                        is_recipe,
                        grams: i.grams,
                        amount: i.amount,
                        unit: i.unit.to_string(),
                        adjective: i.adjective.to_string(),
                        optional: i.optional,
                    }.to_string()
                })
                .collect(),
        })
        .collect();

        let r = Recipe {
            sections: rs.to_vec(),
            name: recipe.name,
            uuid: recipe.uuid,
            unit: recipe.unit,
            total_minutes: recipe.total_minutes,
        };
        Ok(r)
    }
}

#[derive(GraphQLQuery)]
#[graphql(
    schema_path = "../graph/schema.graphql",
    query_path = "query.graphql",
    response_derives = "Debug"
)]
struct GetRecipeByUUID;

#[derive(GraphQLQuery)]
#[graphql(
    schema_path = "../graph/schema.graphql",
    query_path = "query_set.graphql",
    response_derives = "Debug"
)]
struct UpdateRecipe;

#[derive(GraphQLQuery)]
#[graphql(
    schema_path = "../graph/schema.graphql",
    query_path = "query_upsert_ingredient.graphql",
    response_derives = "Debug"
)]
struct UpsertIngredient;

pub struct APIClient {
    pub url: String,
}
impl APIClient {
    pub fn get_uuid(&self, name: String, is_recipe: bool) -> String {
        let q = UpsertIngredient::build_query(upsert_ingredient::Variables {
            name,
            kind: if is_recipe {
                upsert_ingredient::SectionIngredientKind::recipe
            } else {
                upsert_ingredient::SectionIngredientKind::ingredient
            },
        });
        let client = reqwest::Client::new();
        let mut res = client.post(&self.url).json(&q).send().unwrap();
        let response_body: Response<upsert_ingredient::ResponseData> = res.json().unwrap();
        return response_body.data.unwrap().upsert_ingredient;
    }
    pub fn get_recipe(&self, uuid: String) -> Recipe {
        let q = GetRecipeByUUID::build_query(get_recipe_by_uuid::Variables { uuid });
        let client = reqwest::Client::new();
        let mut res = client.post(&self.url).json(&q).send().unwrap();
        let response_body: Response<get_recipe_by_uuid::ResponseData> = res.json().unwrap();
        let recipe = response_body.data.unwrap().recipe.unwrap();

        return Recipe::from_api_recipe(recipe).unwrap();
    }
    pub fn set_recipe(&self, input_recipe: Recipe) {
        let input = input_recipe.to_api_recipe(self);
        let q = UpdateRecipe::build_query(update_recipe::Variables { recipe: input });

        let client = reqwest::Client::new();
        let mut res = client.post(&self.url).json(&q).send().unwrap();
        let response_body: Response<update_recipe::ResponseData> = res.json().unwrap();
        let recipe: update_recipe::UpdateRecipeUpdateRecipe =
            response_body.data.unwrap().update_recipe;

        println!("{:#?}", recipe);
        // return Recipe::from_api_recipe(recipe).unwrap();
    }
}
