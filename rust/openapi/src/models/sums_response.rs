/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */




#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct SumsResponse {
    /// mappings of equivalent units
    #[serde(rename = "sums")]
    pub sums: Vec<crate::models::UsageValue>,
    #[serde(rename = "by_recipe")]
    pub by_recipe: ::std::collections::HashMap<String, Vec<crate::models::UsageValue>>,
}

impl SumsResponse {
    pub fn new(sums: Vec<crate::models::UsageValue>, by_recipe: ::std::collections::HashMap<String, Vec<crate::models::UsageValue>>) -> SumsResponse {
        SumsResponse {
            sums,
            by_recipe,
        }
    }
}


