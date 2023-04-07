use openapi::models::{Amount, IngredientKind, SectionIngredientInput, UnitMapping};
use tracing::info;

use ingredient::IngredientParser;

use crate::converter::amount_to_measure;
use crate::unit;

pub(crate) fn section_ingredient_from_parsed(
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

pub fn new_ingredient_parser(is_rich_text: bool) -> IngredientParser {
    IngredientParser::new(is_rich_text)
}

// fn get_grams_si(si: SectionIngredient) -> f64 {
//     for x in si.amounts.iter() {
//         if is_gram(&x.unit) {
//             return x.value;
//         }
//     }
//     return 0.0;
// }
// pub fn get_grams(r: RecipeDetail) -> f64 {
//     r.sections.iter().fold(0.0, |acc, s| {
//         acc + s
//             .ingredients
//             .iter()
//             .fold(0.0, |acc, i| acc + get_grams_si(i.clone()))
//     })
// }

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

#[cfg(test)]
mod tests {
    use std::collections::HashMap;

    use openapi::models::{
        Amount, Ingredient, IngredientDetail, IngredientKind, RecipeDetail, RecipeSection,
        SectionIngredient, SectionIngredientInput,
    };
    use pretty_assertions::assert_eq;

    use crate::{converter::sum_ingredients, parser::parse_ingredient};

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