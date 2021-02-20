/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// FoodNutrient : todo



#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct FoodNutrient {
    #[serde(rename = "nutrient")]
    pub nutrient: crate::models::Nutrient,
    #[serde(rename = "amount")]
    pub amount: f64,
    #[serde(rename = "data_points")]
    pub data_points: i32,
}

impl FoodNutrient {
    /// todo
    pub fn new(nutrient: crate::models::Nutrient, amount: f64, data_points: i32) -> FoodNutrient {
        FoodNutrient {
            nutrient,
            amount,
            data_points,
        }
    }
}


