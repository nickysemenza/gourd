use openapi::models::{
    Amount, BrandedFoodItem, FoundationFoodItem, SrLegacyFoodItem, SurveyFoodItem, TempFood,
    TempFoodWrapper, UnitMapping,
};
use tracing::info;

use crate::{amount_from_ingredient, parse_amount};

pub fn kcal_from_nutrients(nutrients: Vec<openapi::models::FoodNutrient>) -> Option<UnitMapping> {
    for n in nutrients {
        if n.nutrient.is_some() {
            let n = n.nutrient.as_ref().unwrap();
            if n.unit_name == Some("kcal".to_string()) {
                return Some(UnitMapping {
                    a: Box::new(Amount {
                        unit: "grams".to_owned(),
                        value: 100.0,
                        upper_value: None,
                        source: None,
                    }),
                    b: Box::new(Amount {
                        unit: "kcal".to_string(),
                        value: n.number.clone().unwrap().parse::<f64>().unwrap(),
                        upper_value: None,
                        source: None,
                    }),
                    source: Some("fdc hs".to_string()),
                });
            }
        }
    }
    None
}

pub fn foo(nutrients: Option<Vec<openapi::models::FoodNutrient>>) -> Vec<UnitMapping> {
    if let Some(m) = kcal_from_nutrients(nutrients.unwrap_or_default()) {
        return vec![m];
    }
    return vec![];
}

pub fn branded_food_into_wrapper(x: BrandedFoodItem) -> TempFood {
    let mut mappings = foo(x.food_nutrients.clone());

    if let Some(serving) = x.household_serving_full_text.clone() {
        info!("going to parse {} for {}", serving, x.fdc_id);
        let res = parse_amount(&serving);
        if let Some(b) = res.first() {
            info!("found {} servings: {:?}", res.len(), res);
            let mapping = UnitMapping {
                a: Box::new(Amount {
                    unit: x.serving_size_unit.clone().unwrap(),
                    value: x.serving_size.unwrap(),
                    upper_value: None,
                    source: None,
                }),
                b: Box::new(amount_from_ingredient(b)),
                source: Some("fdc hs".to_string()),
            };
            mappings.push(mapping);
        }
    }

    TempFood {
        branded_food: Some(Box::new(x.clone())),
        food_nutrients: x.food_nutrients,
        ..TempFood::new(
            TempFoodWrapper::new(x.fdc_id, x.data_type, x.description),
            mappings,
        )
    }
}
pub fn sr_legacy_food_into_wrapper(x: SrLegacyFoodItem) -> TempFood {
    TempFood {
        legacy_food: Some(Box::new(x.clone())),
        food_nutrients: x.food_nutrients.clone(),
        ..TempFood::new(
            TempFoodWrapper::new(x.fdc_id, x.data_type, x.description),
            foo(x.food_nutrients),
        )
    }
}

pub fn foundation_food_into_wrapper(x: FoundationFoodItem) -> TempFood {
    TempFood {
        foundation_food: Some(Box::new(x.clone())),
        food_nutrients: x.food_nutrients.clone(),
        ..TempFood::new(
            TempFoodWrapper::new(x.fdc_id, x.data_type, x.description),
            foo(x.food_nutrients),
        )
    }
}
pub fn survey_food_into_wrapper(x: SurveyFoodItem) -> TempFood {
    TempFood {
        survey_food: Some(Box::new(x.clone())),
        ..TempFood::new(
            TempFoodWrapper::new(x.fdc_id, x.data_type, x.description),
            vec![],
        )
    }
}

#[cfg(test)]
mod tests {
    use openapi::models::{Amount, BrandedFoodItem, UnitMapping};

    use crate::usda::kcal_from_nutrients;

    #[test]
    fn test_branded_item_conversion() {
        let res: BrandedFoodItem =
            serde_json::from_str(include_str!("sample_branded_food.json")).unwrap();
        assert_eq!(res.fdc_id, 2082103);

        assert_eq!(
            Some(UnitMapping {
                a: Box::new(Amount::new("grams".to_string(), 100.0)),
                b: Box::new(Amount::new("kcal".to_string(), 208.0)),
                source: Some("fdc hs".to_string()),
            }),
            kcal_from_nutrients(res.food_nutrients.unwrap())
        );
    }
}
