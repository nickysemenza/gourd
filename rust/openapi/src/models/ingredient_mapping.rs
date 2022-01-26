/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// IngredientMapping : details about ingredients



#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct IngredientMapping {
    #[serde(rename = "name")]
    pub name: String,
    #[serde(rename = "fdc_id", skip_serializing_if = "Option::is_none")]
    pub fdc_id: Option<i32>,
    #[serde(rename = "aliases")]
    pub aliases: Vec<String>,
    /// mappings of equivalent units
    #[serde(rename = "unit_mappings")]
    pub unit_mappings: Vec<crate::models::UnitMapping>,
}

impl IngredientMapping {
    /// details about ingredients
    pub fn new(name: String, aliases: Vec<String>, unit_mappings: Vec<crate::models::UnitMapping>) -> IngredientMapping {
        IngredientMapping {
            name,
            fdc_id: None,
            aliases,
            unit_mappings,
        }
    }
}


