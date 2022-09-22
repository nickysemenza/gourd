/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// IngredientDetail : An Ingredient



#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct IngredientDetail {
    #[serde(rename = "ingredient")]
    pub ingredient: Box<crate::models::Ingredient>,
    /// Recipes referencing this ingredient
    #[serde(rename = "recipes")]
    pub recipes: Vec<crate::models::RecipeDetail>,
    /// Ingredients that are equivalent
    #[serde(rename = "children", skip_serializing_if = "Option::is_none")]
    pub children: Option<Vec<crate::models::IngredientDetail>>,
    #[serde(rename = "food", skip_serializing_if = "Option::is_none")]
    pub food: Option<Box<crate::models::Food>>,
    /// mappings of equivalent units
    #[serde(rename = "unit_mappings")]
    pub unit_mappings: Vec<crate::models::UnitMapping>,
}

impl IngredientDetail {
    /// An Ingredient
    pub fn new(ingredient: crate::models::Ingredient, recipes: Vec<crate::models::RecipeDetail>, unit_mappings: Vec<crate::models::UnitMapping>) -> IngredientDetail {
        IngredientDetail {
            ingredient: Box::new(ingredient),
            recipes,
            children: None,
            food: None,
            unit_mappings,
        }
    }
}


