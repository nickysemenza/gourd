/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */


/// 
#[derive(Clone, Copy, Debug, Eq, PartialEq, Ord, PartialOrd, Hash, Serialize, Deserialize)]
pub enum FoodNutrientUnit {
    #[serde(rename = "UG")]
    Ug,
    #[serde(rename = "G")]
    G,
    #[serde(rename = "IU")]
    Iu,
    #[serde(rename = "kJ")]
    KJ,
    #[serde(rename = "KCAL")]
    Kcal,
    #[serde(rename = "MG")]
    Mg,
    #[serde(rename = "MG_ATE")]
    MgAte,
    #[serde(rename = "SP_GR")]
    SpGr,

}

impl ToString for FoodNutrientUnit {
    fn to_string(&self) -> String {
        match self {
            Self::Ug => String::from("UG"),
            Self::G => String::from("G"),
            Self::Iu => String::from("IU"),
            Self::KJ => String::from("kJ"),
            Self::Kcal => String::from("KCAL"),
            Self::Mg => String::from("MG"),
            Self::MgAte => String::from("MG_ATE"),
            Self::SpGr => String::from("SP_GR"),
        }
    }
}

impl Default for FoodNutrientUnit {
    fn default() -> FoodNutrientUnit {
        Self::Ug
    }
}




