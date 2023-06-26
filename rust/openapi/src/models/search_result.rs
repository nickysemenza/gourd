/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// SearchResult : A search result wrapper, which contains ingredients and recipes



#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct SearchResult {
    /// The ingredients
    #[serde(rename = "ingredients", skip_serializing_if = "Option::is_none")]
    pub ingredients: Option<Vec<crate::models::IngredientWrapper>>,
    /// The recipes
    #[serde(rename = "recipes", skip_serializing_if = "Option::is_none")]
    pub recipes: Option<Vec<crate::models::RecipeWrapper>>,
    #[serde(rename = "meta", skip_serializing_if = "Option::is_none")]
    pub meta: Option<Box<crate::models::Items>>,
}

impl SearchResult {
    /// A search result wrapper, which contains ingredients and recipes
    pub fn new() -> SearchResult {
        SearchResult {
            ingredients: None,
            recipes: None,
            meta: None,
        }
    }
}


