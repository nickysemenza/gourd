/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */




#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct InlineObject2 {
    #[serde(rename = "ingredient_ids")]
    pub ingredient_ids: Vec<String>,
}

impl InlineObject2 {
    pub fn new(ingredient_ids: Vec<String>) -> InlineObject2 {
        InlineObject2 {
            ingredient_ids,
        }
    }
}

