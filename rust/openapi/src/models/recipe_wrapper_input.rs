/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// RecipeWrapperInput : A recipe with subcomponents



#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct RecipeWrapperInput {
    /// id
    #[serde(rename = "id", skip_serializing_if = "Option::is_none")]
    pub id: Option<String>,
    #[serde(rename = "detail")]
    pub detail: Box<crate::models::RecipeDetailInput>,
}

impl RecipeWrapperInput {
    /// A recipe with subcomponents
    pub fn new(detail: crate::models::RecipeDetailInput) -> RecipeWrapperInput {
        RecipeWrapperInput {
            id: None,
            detail: Box::new(detail),
        }
    }
}

