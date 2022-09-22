/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// RecipeSection : A step in the recipe



#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct RecipeSection {
    /// id
    #[serde(rename = "id")]
    pub id: String,
    #[serde(rename = "duration", skip_serializing_if = "Option::is_none")]
    pub duration: Option<Box<crate::models::Amount>>,
    /// x
    #[serde(rename = "instructions")]
    pub instructions: Vec<crate::models::SectionInstruction>,
    /// x
    #[serde(rename = "ingredients")]
    pub ingredients: Vec<crate::models::SectionIngredient>,
}

impl RecipeSection {
    /// A step in the recipe
    pub fn new(id: String, instructions: Vec<crate::models::SectionInstruction>, ingredients: Vec<crate::models::SectionIngredient>) -> RecipeSection {
        RecipeSection {
            id,
            duration: None,
            instructions,
            ingredients,
        }
    }
}


