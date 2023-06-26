/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// PaginatedIngredients : pages of IngredientWrapper



#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct PaginatedIngredients {
    #[serde(rename = "ingredients", skip_serializing_if = "Option::is_none")]
    pub ingredients: Option<Vec<crate::models::IngredientWrapper>>,
    #[serde(rename = "meta")]
    pub meta: Box<crate::models::Items>,
}

impl PaginatedIngredients {
    /// pages of IngredientWrapper
    pub fn new(meta: crate::models::Items) -> PaginatedIngredients {
        PaginatedIngredients {
            ingredients: None,
            meta: Box::new(meta),
        }
    }
}


