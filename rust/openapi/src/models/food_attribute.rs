/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */




#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct FoodAttribute {
    #[serde(rename = "id", skip_serializing_if = "Option::is_none")]
    pub id: Option<i32>,
    #[serde(rename = "sequenceNumber", skip_serializing_if = "Option::is_none")]
    pub sequence_number: Option<i32>,
    #[serde(rename = "value", skip_serializing_if = "Option::is_none")]
    pub value: Option<String>,
    #[serde(rename = "FoodAttributeType", skip_serializing_if = "Option::is_none")]
    pub food_attribute_type: Option<Box<crate::models::FoodAttributeFoodAttributeType>>,
}

impl FoodAttribute {
    pub fn new() -> FoodAttribute {
        FoodAttribute {
            id: None,
            sequence_number: None,
            value: None,
            food_attribute_type: None,
        }
    }
}


