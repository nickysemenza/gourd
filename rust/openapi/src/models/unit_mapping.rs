/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// UnitMapping : mappings



#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct UnitMapping {
    #[serde(rename = "a")]
    pub a: Box<crate::models::Amount>,
    #[serde(rename = "b")]
    pub b: Box<crate::models::Amount>,
    /// source of the mapping
    #[serde(rename = "source", skip_serializing_if = "Option::is_none")]
    pub source: Option<String>,
}

impl UnitMapping {
    /// mappings
    pub fn new(a: crate::models::Amount, b: crate::models::Amount) -> UnitMapping {
        UnitMapping {
            a: Box::new(a),
            b: Box::new(b),
            source: None,
        }
    }
}


