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
pub struct TempFoodWrapper {
    #[serde(rename = "fdcId")]
    pub fdc_id: i32,
    #[serde(rename = "availableDate", skip_serializing_if = "Option::is_none")]
    pub available_date: Option<String>,
    #[serde(rename = "brandOwner", skip_serializing_if = "Option::is_none")]
    pub brand_owner: Option<String>,
    #[serde(rename = "dataSource", skip_serializing_if = "Option::is_none")]
    pub data_source: Option<String>,
    #[serde(rename = "dataType")]
    pub data_type: String,
    #[serde(rename = "description")]
    pub description: String,
    #[serde(rename = "foodClass", skip_serializing_if = "Option::is_none")]
    pub food_class: Option<String>,
    #[serde(rename = "gtinUpc", skip_serializing_if = "Option::is_none")]
    pub gtin_upc: Option<String>,
    #[serde(rename = "householdServingFullText", skip_serializing_if = "Option::is_none")]
    pub household_serving_full_text: Option<String>,
    #[serde(rename = "ingredients", skip_serializing_if = "Option::is_none")]
    pub ingredients: Option<String>,
    #[serde(rename = "modifiedDate", skip_serializing_if = "Option::is_none")]
    pub modified_date: Option<String>,
    #[serde(rename = "publicationDate", skip_serializing_if = "Option::is_none")]
    pub publication_date: Option<String>,
    #[serde(rename = "servingSize", skip_serializing_if = "Option::is_none")]
    pub serving_size: Option<f64>,
    #[serde(rename = "servingSizeUnit", skip_serializing_if = "Option::is_none")]
    pub serving_size_unit: Option<String>,
    #[serde(rename = "preparationStateCode", skip_serializing_if = "Option::is_none")]
    pub preparation_state_code: Option<String>,
    #[serde(rename = "brandedFoodCategory", skip_serializing_if = "Option::is_none")]
    pub branded_food_category: Option<String>,
    #[serde(rename = "tradeChannel", skip_serializing_if = "Option::is_none")]
    pub trade_channel: Option<Vec<String>>,
    #[serde(rename = "gpcClassCode", skip_serializing_if = "Option::is_none")]
    pub gpc_class_code: Option<i32>,
    #[serde(rename = "foodNutrients", skip_serializing_if = "Option::is_none")]
    pub food_nutrients: Option<Vec<crate::models::FoodNutrient>>,
    #[serde(rename = "foodUpdateLog", skip_serializing_if = "Option::is_none")]
    pub food_update_log: Option<Vec<crate::models::FoodUpdateLog>>,
    #[serde(rename = "labelNutrients", skip_serializing_if = "Option::is_none")]
    pub label_nutrients: Option<Box<crate::models::BrandedFoodItemLabelNutrients>>,
    #[serde(rename = "footNote", skip_serializing_if = "Option::is_none")]
    pub foot_note: Option<String>,
    #[serde(rename = "isHistoricalReference", skip_serializing_if = "Option::is_none")]
    pub is_historical_reference: Option<bool>,
    #[serde(rename = "ndbNumber", skip_serializing_if = "Option::is_none")]
    pub ndb_number: Option<i32>,
    #[serde(rename = "scientificName", skip_serializing_if = "Option::is_none")]
    pub scientific_name: Option<String>,
    #[serde(rename = "foodCategory", skip_serializing_if = "Option::is_none")]
    pub food_category: Option<Box<crate::models::FoodCategory>>,
    #[serde(rename = "foodComponents", skip_serializing_if = "Option::is_none")]
    pub food_components: Option<Vec<crate::models::FoodComponent>>,
    #[serde(rename = "foodPortions", skip_serializing_if = "Option::is_none")]
    pub food_portions: Option<Vec<crate::models::FoodPortion>>,
    #[serde(rename = "inputFoods", skip_serializing_if = "Option::is_none")]
    pub input_foods: Option<Vec<crate::models::InputFoodSurvey>>,
    #[serde(rename = "nutrientConversionFactors", skip_serializing_if = "Option::is_none")]
    pub nutrient_conversion_factors: Option<Vec<crate::models::NutrientConversionFactors>>,
    #[serde(rename = "endDate", skip_serializing_if = "Option::is_none")]
    pub end_date: Option<String>,
    #[serde(rename = "foodCode", skip_serializing_if = "Option::is_none")]
    pub food_code: Option<String>,
    #[serde(rename = "startDate", skip_serializing_if = "Option::is_none")]
    pub start_date: Option<String>,
    #[serde(rename = "foodAttributes", skip_serializing_if = "Option::is_none")]
    pub food_attributes: Option<Vec<crate::models::FoodAttribute>>,
    #[serde(rename = "wweiaFoodCategory", skip_serializing_if = "Option::is_none")]
    pub wweia_food_category: Option<Box<crate::models::WweiaFoodCategory>>,
}

impl TempFoodWrapper {
    pub fn new(fdc_id: i32, data_type: String, description: String) -> TempFoodWrapper {
        TempFoodWrapper {
            fdc_id,
            available_date: None,
            brand_owner: None,
            data_source: None,
            data_type,
            description,
            food_class: None,
            gtin_upc: None,
            household_serving_full_text: None,
            ingredients: None,
            modified_date: None,
            publication_date: None,
            serving_size: None,
            serving_size_unit: None,
            preparation_state_code: None,
            branded_food_category: None,
            trade_channel: None,
            gpc_class_code: None,
            food_nutrients: None,
            food_update_log: None,
            label_nutrients: None,
            foot_note: None,
            is_historical_reference: None,
            ndb_number: None,
            scientific_name: None,
            food_category: None,
            food_components: None,
            food_portions: None,
            input_foods: None,
            nutrient_conversion_factors: None,
            end_date: None,
            food_code: None,
            start_date: None,
            food_attributes: None,
            wweia_food_category: None,
        }
    }
}


