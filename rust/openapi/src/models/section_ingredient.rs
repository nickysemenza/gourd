/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// SectionIngredient : Ingredients in a single section



#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct SectionIngredient {
    /// id
    #[serde(rename = "id")]
    pub id: String,
    #[serde(rename = "kind")]
    pub kind: crate::models::IngredientKind,
    #[serde(rename = "recipe", skip_serializing_if = "Option::is_none")]
    pub recipe: Option<Box<crate::models::RecipeDetail>>,
    #[serde(rename = "ingredient", skip_serializing_if = "Option::is_none")]
    pub ingredient: Option<Box<crate::models::IngredientWrapper>>,
    /// the various measures
    #[serde(rename = "amounts")]
    pub amounts: Vec<crate::models::Amount>,
    /// adjective
    #[serde(rename = "adjective", skip_serializing_if = "Option::is_none")]
    pub adjective: Option<String>,
    /// optional
    #[serde(rename = "optional", skip_serializing_if = "Option::is_none")]
    pub optional: Option<bool>,
    /// raw line item (pre-import/scrape)
    #[serde(rename = "original", skip_serializing_if = "Option::is_none")]
    pub original: Option<String>,
    /// x
    #[serde(rename = "substitutes", skip_serializing_if = "Option::is_none")]
    pub substitutes: Option<Vec<crate::models::SectionIngredient>>,
}

impl SectionIngredient {
    /// Ingredients in a single section
    pub fn new(id: String, kind: crate::models::IngredientKind, amounts: Vec<crate::models::Amount>) -> SectionIngredient {
        SectionIngredient {
            id,
            kind,
            recipe: None,
            ingredient: None,
            amounts,
            adjective: None,
            optional: None,
            original: None,
            substitutes: None,
        }
    }
}


