/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// PaginatedRecipeWrappers : pages of Recipe



#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct PaginatedRecipeWrappers {
    #[serde(rename = "recipes", skip_serializing_if = "Option::is_none")]
    pub recipes: Option<Vec<crate::models::RecipeWrapper>>,
    #[serde(rename = "meta", skip_serializing_if = "Option::is_none")]
    pub meta: Option<Box<crate::models::Items>>,
}

impl PaginatedRecipeWrappers {
    /// pages of Recipe
    pub fn new() -> PaginatedRecipeWrappers {
        PaginatedRecipeWrappers {
            recipes: None,
            meta: None,
        }
    }
}

