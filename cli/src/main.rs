use serde::{Deserialize, Serialize};

mod ingredient;

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
                        let i = ingredient::Ingredient::from_str(i_str).unwrap();
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
                    ingredient::Ingredient {
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

use graphql_client::{GraphQLQuery, Response};
use std::io::prelude::*;
use std::{
    fs::{self, File},
    path::Path,
};

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

use clap::{App, Arg, SubCommand};

fn main() {
    let matches = App::new("food CLI")
        .version("1.0")
        .arg(
            Arg::with_name("url")
                .short("u")
                .long("url")
                .value_name("URL")
                .help("Sets a custom URL")
                .takes_value(true),
        )
        .arg(
            Arg::with_name("v")
                .short("v")
                .multiple(true)
                .help("Sets the level of verbosity"),
        )
        .subcommand(
            SubCommand::with_name("get")
                .about("downloads recipes by uuid")
                .arg(
                    Arg::with_name("uuid")
                        .short("id")
                        .long("uuid")
                        .help("get by uuid")
                        .takes_value(true),
                ),
        )
        .subcommand(
            SubCommand::with_name("upload")
                .about("upload recipe by file")
                .arg(
                    Arg::with_name("file")
                        .short("f")
                        .long("file")
                        .value_name("FILE")
                        .help("")
                        .takes_value(true),
                ),
        )
        .get_matches();

    // Gets a value for config if supplied by user, or defaults to "default.conf"
    let url = matches
        .value_of("url")
        .unwrap_or("http://localhost:4242/query");

    let c = APIClient {
        url: url.to_string(),
    };

    if let Some(matches) = matches.subcommand_matches("get") {
        // let uuid = .unwrap();

        let uuid = match matches.value_of("uuid") {
            None => panic!("uuid is required"),
            Some(u) => u,
        };

        let r = c.get_recipe(uuid.to_string());
        write_yaml_to_file(r);
    }
    if let Some(matches) = matches.subcommand_matches("upload") {
        let file = match matches.value_of("file") {
            None => panic!("file is required"),
            Some(u) => u,
        };
        let r = read_yaml_file(file.to_string());
        let res = c.set_recipe(r);
        println!("uploaded: {:#?}", res);
    }
}
pub struct APIClient {
    url: String,
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

pub fn write_yaml_to_file(r: Recipe) {
    let y = serde_yaml::to_string(&r).unwrap();
    println!("{}", y);
    let file_name = format!("{}.yaml", r.uuid);
    let path = Path::new(&file_name);
    let display = path.display();

    let mut file = match File::create(&path) {
        Err(why) => panic!("couldn't create {}: {}", display, why),
        Ok(file) => file,
    };

    match file.write_all(y.as_bytes()) {
        Err(why) => panic!("couldn't write to {}: {}", display, why),
        Ok(_) => println!("successfully wrote to {}", display),
    }
}

pub fn read_yaml_file(filename: String) -> Recipe {
    let contents = fs::read_to_string(filename).expect("Something went wrong reading the file");

    let recipe: Recipe = serde_yaml::from_str(&contents).unwrap();
    return recipe;
}
