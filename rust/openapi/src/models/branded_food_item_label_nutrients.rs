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
pub struct BrandedFoodItemLabelNutrients {
    #[serde(rename = "fat", skip_serializing_if = "Option::is_none")]
    pub fat: Option<Box<crate::models::BrandedFoodItemLabelNutrientsFat>>,
    #[serde(rename = "saturatedFat", skip_serializing_if = "Option::is_none")]
    pub saturated_fat: Option<Box<crate::models::BrandedFoodItemLabelNutrientsSaturatedFat>>,
    #[serde(rename = "transFat", skip_serializing_if = "Option::is_none")]
    pub trans_fat: Option<Box<crate::models::BrandedFoodItemLabelNutrientsTransFat>>,
    #[serde(rename = "cholesterol", skip_serializing_if = "Option::is_none")]
    pub cholesterol: Option<Box<crate::models::BrandedFoodItemLabelNutrientsTransFat>>,
    #[serde(rename = "sodium", skip_serializing_if = "Option::is_none")]
    pub sodium: Option<Box<crate::models::BrandedFoodItemLabelNutrientsTransFat>>,
    #[serde(rename = "carbohydrates", skip_serializing_if = "Option::is_none")]
    pub carbohydrates: Option<Box<crate::models::BrandedFoodItemLabelNutrientsCarbohydrates>>,
    #[serde(rename = "fiber", skip_serializing_if = "Option::is_none")]
    pub fiber: Option<Box<crate::models::BrandedFoodItemLabelNutrientsFiber>>,
    #[serde(rename = "sugars", skip_serializing_if = "Option::is_none")]
    pub sugars: Option<Box<crate::models::BrandedFoodItemLabelNutrientsSugars>>,
    #[serde(rename = "protein", skip_serializing_if = "Option::is_none")]
    pub protein: Option<Box<crate::models::BrandedFoodItemLabelNutrientsProtein>>,
    #[serde(rename = "calcium", skip_serializing_if = "Option::is_none")]
    pub calcium: Option<Box<crate::models::BrandedFoodItemLabelNutrientsCalcium>>,
    #[serde(rename = "iron", skip_serializing_if = "Option::is_none")]
    pub iron: Option<Box<crate::models::BrandedFoodItemLabelNutrientsIron>>,
    #[serde(rename = "potassium", skip_serializing_if = "Option::is_none")]
    pub potassium: Option<Box<crate::models::BrandedFoodItemLabelNutrientsPotassium>>,
    #[serde(rename = "calories", skip_serializing_if = "Option::is_none")]
    pub calories: Option<Box<crate::models::BrandedFoodItemLabelNutrientsCalories>>,
}

impl BrandedFoodItemLabelNutrients {
    pub fn new() -> BrandedFoodItemLabelNutrients {
        BrandedFoodItemLabelNutrients {
            fat: None,
            saturated_fat: None,
            trans_fat: None,
            cholesterol: None,
            sodium: None,
            carbohydrates: None,
            fiber: None,
            sugars: None,
            protein: None,
            calcium: None,
            iron: None,
            potassium: None,
            calories: None,
        }
    }
}


