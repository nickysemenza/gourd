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
pub struct CompactRecipe {
    #[serde(rename = "meta")]
    pub meta: Box<crate::models::CompactRecipeMeta>,
    #[serde(rename = "sections")]
    pub sections: Vec<crate::models::CompactRecipeSection>,
}

impl CompactRecipe {
    pub fn new(meta: crate::models::CompactRecipeMeta, sections: Vec<crate::models::CompactRecipeSection>) -> CompactRecipe {
        CompactRecipe {
            meta: Box::new(meta),
            sections,
        }
    }
}

