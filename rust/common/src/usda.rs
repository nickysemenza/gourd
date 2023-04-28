use crate::{amount_from_ingredient, parse_amount};
use anyhow::{Context, Result};
use openapi::models::{
    Amount, BrandedFoodItem, FoundationFoodItem, SrLegacyFoodItem, SurveyFoodItem, TempFood,
    TempFoodWrapper, UnitMapping,
};
use tracing::info;

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

pub fn get_cal_mappings(nutrients: Option<Vec<openapi::models::FoodNutrient>>) -> Vec<UnitMapping> {
    if let Some(m) = kcal_from_nutrients(nutrients.unwrap_or_default()) {
        return vec![m];
    }
    return vec![];
}

pub trait IntoFoodWrapper {
    fn into_wrapper(self) -> Result<TempFood>;
}
impl IntoFoodWrapper for BrandedFoodItem {
    fn into_wrapper(self) -> Result<TempFood> {
        let mut mappings = get_cal_mappings(self.food_nutrients.clone());

        if let Some(serving) = self.household_serving_full_text.clone() {
            if serving != "" {
                // catch 614264
                info!("going to parse {} for {}", serving, self.fdc_id);
                let res = parse_amount(&serving).context("failed to parse serving size")?;
                if let Some(b) = res.first() {
                    info!("found {} servings: {:?}", res.len(), res);
                    let mapping = UnitMapping {
                        a: Box::new(Amount {
                            unit: self.serving_size_unit.clone().unwrap(),
                            value: self.serving_size.unwrap(),
                            upper_value: None,
                            source: None,
                        }),
                        b: Box::new(amount_from_ingredient(b)),
                        source: Some("fdc hs".to_string()),
                    };
                    mappings.push(mapping);
                }
            }
        }

        Ok(TempFood {
            branded_food: Some(Box::new(self.clone())),
            food_nutrients: self.food_nutrients,
            ..TempFood::new(
                TempFoodWrapper::new(self.fdc_id, self.data_type, self.description),
                mappings,
            )
        })
    }
}
impl IntoFoodWrapper for SrLegacyFoodItem {
    fn into_wrapper(self) -> Result<TempFood> {
        Ok(TempFood {
            legacy_food: Some(Box::new(self.clone())),
            food_nutrients: self.food_nutrients.clone(),
            ..TempFood::new(
                TempFoodWrapper::new(self.fdc_id, self.data_type, self.description),
                get_cal_mappings(self.food_nutrients),
            )
        })
    }
}
impl IntoFoodWrapper for FoundationFoodItem {
    fn into_wrapper(self) -> Result<TempFood> {
        Ok(TempFood {
            foundation_food: Some(Box::new(self.clone())),
            food_nutrients: self.food_nutrients.clone(),
            ..TempFood::new(
                TempFoodWrapper::new(self.fdc_id, self.data_type, self.description),
                get_cal_mappings(self.food_nutrients),
            )
        })
    }
}
impl IntoFoodWrapper for SurveyFoodItem {
    fn into_wrapper(self) -> Result<TempFood> {
        Ok(TempFood {
            survey_food: Some(Box::new(self.clone())),
            ..TempFood::new(
                TempFoodWrapper::new(self.fdc_id, self.data_type, self.description),
                vec![],
            )
        })
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
