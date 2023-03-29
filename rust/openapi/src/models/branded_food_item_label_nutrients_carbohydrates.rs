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
pub struct BrandedFoodItemLabelNutrientsCarbohydrates {
    #[serde(rename = "value", skip_serializing_if = "Option::is_none")]
    pub value: Option<f64>,
}

impl BrandedFoodItemLabelNutrientsCarbohydrates {
    pub fn new() -> BrandedFoodItemLabelNutrientsCarbohydrates {
        BrandedFoodItemLabelNutrientsCarbohydrates {
            value: None,
        }
    }
}


