/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// EntitySummary : holds name/id and multiplier for a Kind of entity



#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct EntitySummary {
    /// recipe_detail or ingredient id
    #[serde(rename = "id")]
    pub id: String,
    /// recipe or ingredient name
    #[serde(rename = "name")]
    pub name: String,
    /// multiplier
    #[serde(rename = "multiplier")]
    pub multiplier: f64,
    #[serde(rename = "kind")]
    pub kind: crate::models::IngredientKind,
}

impl EntitySummary {
    /// holds name/id and multiplier for a Kind of entity
    pub fn new(id: String, name: String, multiplier: f64, kind: crate::models::IngredientKind) -> EntitySummary {
        EntitySummary {
            id,
            name,
            multiplier,
            kind,
        }
    }
}


