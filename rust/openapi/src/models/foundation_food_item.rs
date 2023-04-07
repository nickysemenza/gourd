/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */




#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct FoundationFoodItem {
    #[serde(rename = "fdcId")]
    pub fdc_id: i32,
    #[serde(rename = "dataType")]
    pub data_type: String,
    #[serde(rename = "description")]
    pub description: String,
    #[serde(rename = "foodClass", skip_serializing_if = "Option::is_none")]
    pub food_class: Option<String>,
    #[serde(rename = "footNote", skip_serializing_if = "Option::is_none")]
    pub foot_note: Option<String>,
    #[serde(rename = "isHistoricalReference", skip_serializing_if = "Option::is_none")]
    pub is_historical_reference: Option<bool>,
    #[serde(rename = "ndbNumber", skip_serializing_if = "Option::is_none")]
    pub ndb_number: Option<i32>,
    #[serde(rename = "publicationDate", skip_serializing_if = "Option::is_none")]
    pub publication_date: Option<String>,
    #[serde(rename = "scientificName", skip_serializing_if = "Option::is_none")]
    pub scientific_name: Option<String>,
    #[serde(rename = "foodCategory", skip_serializing_if = "Option::is_none")]
    pub food_category: Option<Box<crate::models::SchemasFoodCategory>>,
    #[serde(rename = "foodComponents", skip_serializing_if = "Option::is_none")]
    pub food_components: Option<Vec<crate::models::FoodComponent>>,
    #[serde(rename = "foodNutrients", skip_serializing_if = "Option::is_none")]
    pub food_nutrients: Option<Vec<crate::models::FoodNutrient>>,
    #[serde(rename = "foodPortions", skip_serializing_if = "Option::is_none")]
    pub food_portions: Option<Vec<crate::models::SchemasFoodPortion>>,
    #[serde(rename = "inputFoods", skip_serializing_if = "Option::is_none")]
    pub input_foods: Option<Vec<crate::models::InputFoodFoundation>>,
    #[serde(rename = "nutrientConversionFactors", skip_serializing_if = "Option::is_none")]
    pub nutrient_conversion_factors: Option<Vec<crate::models::NutrientConversionFactors>>,
}

impl FoundationFoodItem {
    pub fn new(fdc_id: i32, data_type: String, description: String) -> FoundationFoodItem {
        FoundationFoodItem {
            fdc_id,
            data_type,
            description,
            food_class: None,
            foot_note: None,
            is_historical_reference: None,
            ndb_number: None,
            publication_date: None,
            scientific_name: None,
            food_category: None,
            food_components: None,
            food_nutrients: None,
            food_portions: None,
            input_foods: None,
            nutrient_conversion_factors: None,
        }
    }
}

