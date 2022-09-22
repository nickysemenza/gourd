/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */




#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct MergeIngredientsRequest {
    #[serde(rename = "ingredient_ids")]
    pub ingredient_ids: Vec<String>,
}

impl MergeIngredientsRequest {
    pub fn new(ingredient_ids: Vec<String>) -> MergeIngredientsRequest {
        MergeIngredientsRequest {
            ingredient_ids,
        }
    }
}

