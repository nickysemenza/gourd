use openapi::models::{Amount, BrandedFood, BrandedFoodItem, FoodInfo, FoodWrapper, UnitMapping};
use tracing::info;

use crate::{amount_from_ingredient, parse_amount};

pub fn food_info_from_branded_food_item(x: BrandedFoodItem) -> FoodInfo {
    let x = FoodWrapper {
        fdc_id: x.fdc_id,
        description: x.description,
        data_type: openapi::models::FoodDataType::FoundationFood, // todo
        category: None,
        nutrients: x.food_nutrients.unwrap_or_default(),
        portions: None,
        branded_info: Some(Box::new(BrandedFood {
            serving_size: x.serving_size.unwrap(),
            serving_size_unit: x.serving_size_unit.unwrap(),
            brand_owner: x.brand_owner,
            ingredients: x.ingredients,
            // fdcId:  2317162 as None as the text which messes up the wasm
            household_serving: match x.household_serving_full_text.clone().unwrap_or_default()
                == "None"
            {
                true => None,
                false => x.household_serving_full_text,
            },
            branded_food_category: x.branded_food_category,
        })),
    };
    FoodInfo::new(x.clone(), make_unit_mappings(x))
}

#[tracing::instrument]
pub fn make_unit_mappings(food: FoodWrapper) -> Vec<UnitMapping> {
    let mut mappings = Vec::new();

    if let Some(branded_food) = food.branded_info {
        if let Some(serving) = branded_food.household_serving {
            info!("going to parse {} for {}", serving, food.fdc_id);
            let res = parse_amount(&serving);
            if let Some(b) = res.first() {
                info!("found {} servings: {:?}", res.len(), res);
                let mapping = UnitMapping {
                    a: Box::new(Amount {
                        unit: branded_food.serving_size_unit,
                        value: branded_food.serving_size,
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

    if let Some(portions) = food.portions {
        info!("found {} portions", portions.len());
        for p in portions {
            let a = Box::new(Amount::new("grams".to_string(), p.gram_weight));
            if p.portion_description != "" {
                let res = parse_amount(&p.portion_description);
                if let Some(b) = res.first() {
                    let mapping = UnitMapping {
                        a,
                        b: Box::new(amount_from_ingredient(b)),
                        source: Some("fdc p1".to_string()),
                    };
                    mappings.push(mapping);
                }
            } else {
                let mapping = UnitMapping {
                    a: Box::new(Amount {
                        unit: "grams".to_string(),
                        value: p.gram_weight,
                        upper_value: None,
                        source: None,
                    }),
                    b: Box::new(Amount::new(p.modifier, p.amount)),
                    source: Some("fdc p2".to_string()),
                };
                mappings.push(mapping);
            }
        }
    }

    for n in food.nutrients {
        if let Some(nutrient) = n.nutrient {
            if let Some(unit) = nutrient.unit_name {
                if unit.to_ascii_lowercase() == "kcal" {
                    mappings.push(UnitMapping::new_with_source(
                        Amount::new("kcal".to_string(), n.amount.unwrap()),
                        Amount::new("grams".to_string(), 100.0),
                        "fdc p",
                    ));
                }
            }
        }
    }

    info!("found {} mappings", mappings.len());

    mappings
}

#[cfg(test)]
mod tests {
    use openapi::models::BrandedFoodItem;

    use super::food_info_from_branded_food_item;

    #[test]
    fn test_branded_item_conversion() {
        let res: BrandedFoodItem =
            serde_json::from_str(include_str!("sample_branded_food.json")).unwrap();
        assert_eq!(res.fdc_id, 2082103);

        let converted = food_info_from_branded_food_item(res.clone());
        assert_eq!(res.fdc_id, converted.wrapper.fdc_id);
        assert_eq!(1, converted.unit_mappings.len());
        assert_eq!(500.0, converted.unit_mappings[0].a.value);
    }
}
