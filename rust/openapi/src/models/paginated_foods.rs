/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// PaginatedFoods : pages of Food



#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct PaginatedFoods {
    #[serde(rename = "foods", skip_serializing_if = "Option::is_none")]
    pub foods: Option<Vec<crate::models::Food>>,
    #[serde(rename = "meta")]
    pub meta: Box<crate::models::Items>,
}

impl PaginatedFoods {
    /// pages of Food
    pub fn new(meta: crate::models::Items) -> PaginatedFoods {
        PaginatedFoods {
            foods: None,
            meta: Box::new(meta),
        }
    }
}


