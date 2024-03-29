/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// PaginatedMeals : pages of Meal



#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct PaginatedMeals {
    #[serde(rename = "meals", skip_serializing_if = "Option::is_none")]
    pub meals: Option<Vec<crate::models::Meal>>,
    #[serde(rename = "meta")]
    pub meta: Box<crate::models::Items>,
}

impl PaginatedMeals {
    /// pages of Meal
    pub fn new(meta: crate::models::Items) -> PaginatedMeals {
        PaginatedMeals {
            meals: None,
            meta: Box::new(meta),
        }
    }
}


