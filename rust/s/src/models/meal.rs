/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// Meal : A meal, which bridges recipes to photos



#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct Meal {
    /// id
    #[serde(rename = "id")]
    pub id: String,
    /// public image
    #[serde(rename = "name")]
    pub name: String,
    /// when it was taken
    #[serde(rename = "ate_at")]
    pub ate_at: String,
    #[serde(rename = "photos")]
    pub photos: Vec<crate::models::GooglePhoto>,
    #[serde(rename = "recipes", skip_serializing_if = "Option::is_none")]
    pub recipes: Option<Vec<crate::models::MealRecipe>>,
}

impl Meal {
    /// A meal, which bridges recipes to photos
    pub fn new(id: String, name: String, ate_at: String, photos: Vec<crate::models::GooglePhoto>) -> Meal {
        Meal {
            id,
            name,
            ate_at,
            photos,
            recipes: None,
        }
    }
}


