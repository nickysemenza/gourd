/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// RecipeDetailMeta : metadata about recipe detail



#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct RecipeDetailMeta {
    /// version of the recipe
    #[serde(rename = "version")]
    pub version: i32,
    /// whether or not it is the most recent version
    #[serde(rename = "is_latest_version")]
    pub is_latest_version: bool,
}

impl RecipeDetailMeta {
    /// metadata about recipe detail
    pub fn new(version: i32, is_latest_version: bool) -> RecipeDetailMeta {
        RecipeDetailMeta {
            version,
            is_latest_version,
        }
    }
}

