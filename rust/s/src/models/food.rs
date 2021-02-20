/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// Food : A top level food



#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct Food {
    /// FDC Id
    #[serde(rename = "fdc_id")]
    pub fdc_id: i32,
    /// Food description
    #[serde(rename = "description")]
    pub description: String,
    #[serde(rename = "data_type")]
    pub data_type: crate::models::FoodDataType,
    #[serde(rename = "category", skip_serializing_if = "Option::is_none")]
    pub category: Option<crate::models::FoodCategory>,
    /// todo
    #[serde(rename = "nutrients")]
    pub nutrients: Vec<crate::models::FoodNutrient>,
    /// portion datapoints
    #[serde(rename = "portions", skip_serializing_if = "Option::is_none")]
    pub portions: Option<Vec<crate::models::FoodPortion>>,
    #[serde(rename = "branded_info", skip_serializing_if = "Option::is_none")]
    pub branded_info: Option<crate::models::BrandedFood>,
}

impl Food {
    /// A top level food
    pub fn new(fdc_id: i32, description: String, data_type: crate::models::FoodDataType, nutrients: Vec<crate::models::FoodNutrient>) -> Food {
        Food {
            fdc_id,
            description,
            data_type,
            category: None,
            nutrients,
            portions: None,
            branded_info: None,
        }
    }
}


