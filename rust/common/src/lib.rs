use std::collections::{hash_map::Entry, HashMap};

use openapi::models::{
    unit_conversion_request::Target, Amount, Ingredient, IngredientDetail, IngredientKind,
    RecipeDetail, SectionIngredient, SectionIngredientInput, UnitConversionRequest, UnitMapping,
};
use openapi::models::{BrandedFood, BrandedFoodItem, FoodInfo, FoodWrapper};
use tracing::info;

use ingredient::unit::kind::MeasureKind;
use ingredient::IngredientParser;

#[macro_use]
extern crate serde;

pub mod codec;
pub mod pan;
pub mod unit;
pub use ingredient;

fn section_ingredient_from_parsed(
    i: ingredient::Ingredient,
    original: &str,
) -> SectionIngredientInput {
    let mut grams = 0.0;
    let mut oz = 0.0;
    let mut ml = 0.0;
    let mut unit: Option<String> = None;
    let mut amount: Option<f64> = None;
    let mut kind = IngredientKind::Ingredient;

    let mut amounts: Vec<Amount> = Vec::new();
    for x in i.amounts.iter() {
        if is_gram(&x.unit) {
            grams = x.value.into();
            amount = Some(grams);
            unit = Some("g".to_string());
        } else if is_oz(x) {
            oz = x.value.into();
        } else if is_ml(x) {
            ml = x.value.into();
        } else {
            unit = Some(x.unit.clone());
            amount = Some(x.value as f64);
        }
        if amount.is_some() && unit.is_some() {
            let unitstr = unit.clone().unwrap_or("unknown".to_string());
            if unitstr == IngredientKind::Recipe.to_string() {
                // if the ingredient amount has a unit of "recipe", then it's likely a recipe
                kind = IngredientKind::Recipe;
            }
            let mut a = Amount::new(unitstr, amount.unwrap_or(0.0));
            if x.upper_value.is_some() {
                a.upper_value = Some(x.upper_value.unwrap().into());
            }
            amounts.push(a);
        }
    }
    if grams == 0.0 {
        // if grams are absent, try converting oz
        if grams == 0.0 && oz != 0.0 {
            grams = ((oz * 28.35) as f64).round();
        }

        // // if grams are absent, try using mL
        if grams == 0.0 && ml != 0.0 {
            grams = ml
        }
        if grams != 0.0 {
            amounts.push(Amount::new("g".to_string(), grams));
        }
    }

    return SectionIngredientInput {
        adjective: i.modifier,
        name: Some(i.name),
        original: Some(original.to_string()),
        ..SectionIngredientInput::new(kind, amounts)
    };
}
//todo: put this in ing parser project
pub fn parse_ingredient(s: &str) -> Result<SectionIngredientInput, String> {
    let mut s2 = s.to_string();
    if s2.contains("((") {
        // for woksoflife.com
        s2 = s2.replace("((", "(");
        s2 = s2.replace("))", ")");
    }
    let i = dbg!(ingredient::from_str(s2.as_str()));
    Ok(section_ingredient_from_parsed(i, s))
}
pub fn parse_amount(s: &str) -> Vec<ingredient::Amount> {
    let i = new_ingredient_parser(false).parse_amount(s);
    info!("parsed {} into {:?}", s, i.clone());
    return i;
}
fn get_grams_si(si: SectionIngredient) -> f64 {
    for x in si.amounts.iter() {
        if is_gram(&x.unit) {
            return x.value;
        }
    }
    return 0.0;
}
pub fn get_grams(r: RecipeDetail) -> f64 {
    r.sections.iter().fold(0.0, |acc, s| {
        acc + s
            .ingredients
            .iter()
            .fold(0.0, |acc, i| acc + get_grams_si(i.clone()))
    })
}

pub fn sum_ingredients(
    r: RecipeDetail,
) -> (
    HashMap<String, Vec<SectionIngredient>>,
    HashMap<String, Vec<SectionIngredient>>,
) {
    let mut recipes = HashMap::new();
    let mut ing = HashMap::new();

    let flat_ingredients: Vec<SectionIngredient> = r
        .sections
        .iter()
        .flat_map(|section| section.ingredients.clone())
        .collect();

    flat_ingredients.iter().for_each(|i| {
        let (k, m) = match i.kind {
            IngredientKind::Recipe => (i.recipe.as_ref().unwrap().id.clone(), &mut recipes),
            IngredientKind::Ingredient => (
                i.ingredient.as_ref().unwrap().ingredient.id.clone(),
                &mut ing,
            ),
        };
        match m.entry(k) {
            Entry::Vacant(e) => {
                e.insert(vec![i.clone()]);
            }
            Entry::Occupied(mut e) => {
                e.get_mut().push(i.clone());
            }
        }
    });
    (recipes, ing)
}

fn is_gram(unit: &String) -> bool {
    unit == "g" || unit == "gram" || unit == "grams"
}
fn is_oz(a: &ingredient::Amount) -> bool {
    a.unit == "oz" || a.unit == "ounce" || a.unit == "ounces"
}
fn is_ml(a: &ingredient::Amount) -> bool {
    a.unit == "ml" || a.unit == "mL"
}

pub fn amount_from_ingredient(e: &ingredient::Amount) -> Amount {
    Amount {
        unit: e.unit.clone(),
        value: e.value,
        upper_value: e.upper_value,
        source: None,
    }
}

pub fn parse_unit_mappings(um: Vec<UnitMapping>) -> Vec<(unit::Measure, unit::Measure)> {
    um.iter()
        .map(|u| {
            return (
                amount_to_measure(*u.a.clone()),
                amount_to_measure(*u.b.clone()),
            );
        })
        .collect()
}
pub fn convert_to(req: UnitConversionRequest) -> Option<Amount> {
    let equivalencies = parse_unit_mappings(req.unit_mappings);
    let target = match req.target.unwrap_or(Target::Other) {
        Target::Weight => MeasureKind::Weight,
        Target::Volume => MeasureKind::Volume,
        Target::Money => MeasureKind::Money,
        Target::Calories => MeasureKind::Calories,
        Target::Other => MeasureKind::Other,
    };
    if req.input.len() == 0 {
        return None;
    }
    return match unit::convert(
        amount_to_measure(req.input[0].clone()),
        target.clone(),
        equivalencies.clone(),
    ) {
        Some(a) => Some(measure_to_amount(a).unwrap()),
        None => {
            if target == MeasureKind::Weight {
                // try again to convert to ml, and then use that as grams
                return match unit::convert(
                    amount_to_measure(req.input[0].clone()),
                    MeasureKind::Volume,
                    equivalencies,
                ) {
                    Some(a) => {
                        let mut a = measure_to_amount(a).unwrap();
                        a.unit = "gram".to_string();
                        info!("no grams for {:#?} using volume with density 1", req.input);
                        return Some(a);
                    }
                    None => None,
                };
            }
            return None;
        }
    };
}
pub fn tmp_normalize(a: Amount) -> Amount {
    let m = amount_to_measure(a.clone()).as_raw();
    return Amount {
        unit: m.unit,
        value: m.value,
        upper_value: m.upper_value,
        source: a.source.clone(),
    };
}
pub fn amount_to_measure(a: Amount) -> unit::Measure {
    unit::Measure::parse(ingredient::Amount {
        unit: a.unit,
        value: a.value,
        upper_value: a.upper_value,
    })
}
pub fn amount_to_measure2(a: ingredient::Amount) -> unit::Measure {
    unit::Measure::parse(ingredient::Amount {
        unit: a.unit,
        value: a.value,
        upper_value: a.upper_value,
    })
}
pub fn measure_to_amount(m: unit::Measure) -> anyhow::Result<Amount> {
    let m1 = m.as_bare()?;
    Ok(Amount::new(m1.unit, m1.value.into()))
}
pub fn si_to_ingredient(s: SectionIngredientInput) -> ingredient::Ingredient {
    let mut amounts = vec![];
    for a in s.amounts.iter() {
        amounts.push(ingredient::Amount {
            unit: a.unit.clone(),
            value: a.value,
            upper_value: a.upper_value,
        });
    }

    return ingredient::Ingredient {
        name: s.name.unwrap_or_default(),
        modifier: s.adjective,
        amounts,
    };
}
#[allow(dead_code)]
fn bare_detail(name: String) -> IngredientDetail {
    IngredientDetail::new(
        Ingredient::new("".to_string(), name.to_string()),
        vec![],
        vec![],
    )
}

#[cfg(test)]
mod tests {
    use std::collections::HashMap;

    use openapi::models::{
        Amount, Ingredient, IngredientDetail, IngredientKind, RecipeDetail, RecipeSection,
        SectionIngredient, SectionIngredientInput,
    };
    use pretty_assertions::assert_eq;

    use crate::{parse_ingredient, sum_ingredients};

    #[test]
    fn test_from_parsed() {
        assert_eq!(
            parse_ingredient("118 grams / 0.5 cups water").unwrap(),
            SectionIngredientInput {
                original: Some("118 grams / 0.5 cups water".to_string()),
                name: Some("water".to_string()),
                ..SectionIngredientInput::new(
                    IngredientKind::Ingredient,
                    vec![
                        Amount::new("g".to_string(), 118.0),
                        Amount::new("cups".to_string(), 0.5),
                    ]
                )
            }
        );
    }
    #[test]
    fn test_from_parsed_grams_from_ml() {
        // ml is parsed as grams if not present
        assert_eq!(
            parse_ingredient("118 ml / 0.5 cups water").unwrap(),
            SectionIngredientInput {
                original: Some("118 ml / 0.5 cups water".to_string()),
                name: Some("water".to_string()),
                ..SectionIngredientInput::new(
                    IngredientKind::Ingredient,
                    vec![
                        Amount::new("cups".to_string(), 0.5),
                        Amount::new("g".to_string(), 118.0),
                    ]
                )
            }
        );
    }
    #[test]
    fn test_from_parsed_grams_from_oz() {
        assert_eq!(
            parse_ingredient("4 oz / 1/2 cup water").unwrap(),
            SectionIngredientInput {
                original: Some("4 oz / 1/2 cup water".to_string()),
                name: Some("water".to_string()),
                ..SectionIngredientInput::new(
                    IngredientKind::Ingredient,
                    vec![
                        Amount::new("cup".to_string(), 0.5),
                        Amount::new("g".to_string(), 113.0),
                    ]
                )
            }
        );
    }

    #[test]
    fn test_sum_ingredients() {
        let si_1 = SectionIngredient {
            ingredient: Some(Box::new(IngredientDetail::new(
                Ingredient::new("a".to_string(), "foo".to_string()),
                vec![],
                vec![],
            ))),
            ..SectionIngredient::new(
                "".to_string(),
                IngredientKind::Ingredient,
                vec![Amount::new("g".to_string(), 12.0)],
            )
        };
        let si_2 = SectionIngredient {
            ingredient: Some(Box::new(IngredientDetail::new(
                Ingredient::new("b".to_string(), "bar".to_string()),
                vec![],
                vec![],
            ))),
            ..SectionIngredient::new(
                "".to_string(),
                IngredientKind::Ingredient,
                vec![Amount::new("g".to_string(), 14.0)],
            )
        };
        let si_3 = SectionIngredient {
            ingredient: Some(Box::new(IngredientDetail::new(
                Ingredient::new("b".to_string(), "bar".to_string()),
                vec![],
                vec![],
            ))),
            ..SectionIngredient::new(
                "".to_string(),
                IngredientKind::Ingredient,
                vec![Amount::new("g".to_string(), 2.0)],
            )
        };
        let r = RecipeDetail::new(
            "".to_string(),
            vec![RecipeSection::new(
                "".to_string(),
                vec![],
                vec![si_1.clone(), si_2.clone(), si_3.clone()],
            )],
            "".to_string(),
            vec![],
            0,
            "".to_string(),
            0,
            false,
            "".to_string(),
            vec![],
        );
        let expected: HashMap<String, Vec<SectionIngredient>> = [
            ("a".to_string(), vec![si_1.clone()]),
            ("b".to_string(), vec![si_2.clone(), si_3.clone()]),
        ]
        .iter()
        .cloned()
        .collect();
        assert_eq!(sum_ingredients(r), (HashMap::new(), expected));
    }
}

pub fn new_ingredient_parser(is_rich_text: bool) -> IngredientParser {
    IngredientParser::new(is_rich_text)
}

pub fn food_info_from_branded_food_item(x: BrandedFoodItem) -> FoodInfo {
    let x = FoodWrapper {
        fdc_id: x.fdc_id,
        description: x.description,
        data_type: openapi::models::FoodDataType::FoundationFood,
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
                if unit == "KCAL" {
                    let mapping = UnitMapping {
                        a: Box::new(Amount::new("kcal".to_string(), n.amount.unwrap())),
                        b: Box::new(Amount::new("grams".to_string(), 100.0)),
                        source: Some("fdc p".to_string()),
                    };
                    mappings.push(mapping);
                }
            }
        }
    }

    info!("found {} mappings", mappings.len());

    mappings
}
