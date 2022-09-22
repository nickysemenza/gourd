/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// RecipeDependency : represents a relationship between recipe and ingredient, the latter of which can also be a recipe.



#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct RecipeDependency {
    /// recipe_id
    #[serde(rename = "recipe_id")]
    pub recipe_id: String,
    /// id
    #[serde(rename = "recipe_name")]
    pub recipe_name: String,
    /// id
    #[serde(rename = "ingredient_id")]
    pub ingredient_id: String,
    /// id
    #[serde(rename = "ingredient_name")]
    pub ingredient_name: String,
    #[serde(rename = "ingredient_kind")]
    pub ingredient_kind: crate::models::IngredientKind,
}

impl RecipeDependency {
    /// represents a relationship between recipe and ingredient, the latter of which can also be a recipe.
    pub fn new(recipe_id: String, recipe_name: String, ingredient_id: String, ingredient_name: String, ingredient_kind: crate::models::IngredientKind) -> RecipeDependency {
        RecipeDependency {
            recipe_id,
            recipe_name,
            ingredient_id,
            ingredient_name,
            ingredient_kind,
        }
    }
}


