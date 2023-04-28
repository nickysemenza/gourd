/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// FoodSearchResult : A meal, which bridges recipes to photos



#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct FoodSearchResult {
    #[serde(rename = "foods")]
    pub foods: Vec<crate::models::TempFood>,
}

impl FoodSearchResult {
    /// A meal, which bridges recipes to photos
    pub fn new(foods: Vec<crate::models::TempFood>) -> FoodSearchResult {
        FoodSearchResult {
            foods,
        }
    }
}


