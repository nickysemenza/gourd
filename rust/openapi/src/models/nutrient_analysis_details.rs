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
pub struct NutrientAnalysisDetails {
    #[serde(rename = "subSampleId", skip_serializing_if = "Option::is_none")]
    pub sub_sample_id: Option<i32>,
    #[serde(rename = "amount", skip_serializing_if = "Option::is_none")]
    pub amount: Option<f64>,
    #[serde(rename = "nutrientId", skip_serializing_if = "Option::is_none")]
    pub nutrient_id: Option<i32>,
    #[serde(rename = "labMethodDescription", skip_serializing_if = "Option::is_none")]
    pub lab_method_description: Option<String>,
    #[serde(rename = "labMethodOriginalDescription", skip_serializing_if = "Option::is_none")]
    pub lab_method_original_description: Option<String>,
    #[serde(rename = "labMethodLink", skip_serializing_if = "Option::is_none")]
    pub lab_method_link: Option<String>,
    #[serde(rename = "labMethodTechnique", skip_serializing_if = "Option::is_none")]
    pub lab_method_technique: Option<String>,
    #[serde(rename = "nutrientAcquisitionDetails", skip_serializing_if = "Option::is_none")]
    pub nutrient_acquisition_details: Option<Vec<crate::models::NutrientAcquisitionDetails>>,
}

impl NutrientAnalysisDetails {
    pub fn new() -> NutrientAnalysisDetails {
        NutrientAnalysisDetails {
            sub_sample_id: None,
            amount: None,
            nutrient_id: None,
            lab_method_description: None,
            lab_method_original_description: None,
            lab_method_link: None,
            lab_method_technique: None,
            nutrient_acquisition_details: None,
        }
    }
}


