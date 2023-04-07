/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// InputFoodFoundation : applies to Foundation foods. Not all inputFoods will have all fields.



#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct InputFoodFoundation {
    #[serde(rename = "id", skip_serializing_if = "Option::is_none")]
    pub id: Option<i32>,
    #[serde(rename = "foodDescription", skip_serializing_if = "Option::is_none")]
    pub food_description: Option<String>,
    #[serde(rename = "inputFood", skip_serializing_if = "Option::is_none")]
    pub input_food: Option<Box<crate::models::SampleFoodItem>>,
}

impl InputFoodFoundation {
    /// applies to Foundation foods. Not all inputFoods will have all fields.
    pub fn new() -> InputFoodFoundation {
        InputFoodFoundation {
            id: None,
            food_description: None,
            input_food: None,
        }
    }
}

