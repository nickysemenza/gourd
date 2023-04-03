/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// RecipeWrapper : A recipe with subcomponents, including some \"generated\" fields to enhance data



#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct RecipeWrapper {
    /// id
    #[serde(rename = "id")]
    pub id: String,
    #[serde(rename = "detail")]
    pub detail: Box<crate::models::RecipeDetail>,
    #[serde(rename = "linked_meals", skip_serializing_if = "Option::is_none")]
    pub linked_meals: Option<Vec<crate::models::Meal>>,
    #[serde(rename = "linked_photos", skip_serializing_if = "Option::is_none")]
    pub linked_photos: Option<Vec<crate::models::Photo>>,
    /// Other versions
    #[serde(rename = "other_versions", skip_serializing_if = "Option::is_none")]
    pub other_versions: Option<Vec<crate::models::RecipeDetail>>,
}

impl RecipeWrapper {
    /// A recipe with subcomponents, including some \"generated\" fields to enhance data
    pub fn new(id: String, detail: crate::models::RecipeDetail) -> RecipeWrapper {
        RecipeWrapper {
            id,
            detail: Box::new(detail),
            linked_meals: None,
            linked_photos: None,
            other_versions: None,
        }
    }
}


