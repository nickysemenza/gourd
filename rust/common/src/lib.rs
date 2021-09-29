use std::collections::{hash_map::Entry, HashMap};

use openapi::models::{
    unit_conversion_request::Target, Amount, Ingredient, IngredientDetail, IngredientKind,
    RecipeDetail, SectionIngredient, UnitConversionRequest, UnitMapping,
};
use unit::MeasureKind;

#[macro_use]
extern crate serde;

pub mod codec;
pub mod pan;
pub mod unit;

fn section_ingredient_from_parsed(i: ingredient::Ingredient, original: &str) -> SectionIngredient {
    let mut grams = 0.0;
    let mut oz = 0.0;
    let mut ml = 0.0;
    let mut unit: Option<String> = None;
    let mut amount: Option<f64> = None;
    let mut kind = IngredientKind::Ingredient;

    let mut amounts: Vec<Amount> = Vec::new();
    for x in i.amounts.iter() {
        if is_gram(x.unit.clone()) {
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
            if unitstr == "recipe" {
                // if the ingredient amount has a unit of "recipe", then it's likely a recipe
                kind = IngredientKind::Recipe;
            }
            amounts.push(Amount::new(unitstr, amount.unwrap_or(0.0)));
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

    return SectionIngredient {
        adjective: i.modifier,
        ingredient: if kind == IngredientKind::Ingredient {
            Some(Box::new(IngredientDetail::new(
                "".to_string(),
                Ingredient::new("".to_string(), i.name.clone()),
                vec![],
                vec![],
                vec![],
            )))
        } else {
            None
        },
        recipe: if kind == IngredientKind::Recipe {
            Some(Box::new(RecipeDetail::new(
                "".to_string(),
                vec![],
                i.name,
                0,
                "".to_string(),
            )))
        } else {
            None
        },
        original: Some(original.to_string()),
        ..SectionIngredient::new("".to_string(), kind, amounts)
    };
}
pub fn parse_ingredient(s: &str) -> Result<SectionIngredient, String> {
    let i = dbg!(ingredient::from_str(s, true))?;
    Ok(section_ingredient_from_parsed(i, s))
}
pub fn parse_amount(s: &str) -> Result<Vec<ingredient::Amount>, String> {
    let i = ingredient::parse_amount(s);
    println!("parsed {} into {:?}", s, i.clone().unwrap_or(vec![]));
    return i;
}
fn get_grams_si(si: SectionIngredient) -> f64 {
    for x in si.amounts.iter() {
        if is_gram(x.unit.clone()) {
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

fn is_gram(unit: String) -> bool {
    unit == "g" || unit == "gram" || unit == "grams"
}
fn is_oz(a: &ingredient::Amount) -> bool {
    a.unit == "oz" || a.unit == "ounce" || a.unit == "ounces"
}
fn is_ml(a: &ingredient::Amount) -> bool {
    a.unit == "ml" || a.unit == "mL"
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
    return match amount_to_measure(req.input[0].clone())
        .convert(target.clone(), equivalencies.clone())
    {
        Some(a) => Some(measure_to_amount(a)),
        None => {
            if target == MeasureKind::Weight {
                // try again to convert to ml, and then use that as grams
                return match amount_to_measure(req.input[0].clone())
                    .convert(MeasureKind::Volume, equivalencies)
                {
                    Some(a) => {
                        let mut a = measure_to_amount(a);
                        a.unit = "gram".to_string();
                        return Some(a);
                    }
                    None => None,
                };
            }
            return None;
        }
    };
}
pub fn amount_to_measure(a: Amount) -> unit::Measure {
    unit::Measure::parse(unit::BareMeasurement::new(a.unit, a.value as f32))
}
pub fn measure_to_amount(m: unit::Measure) -> Amount {
    let m1 = m.as_bare();
    Amount::new(m1.unit, m1.value.into())
}
pub fn si_to_ingredient(s: SectionIngredient) -> ingredient::Ingredient {
    let mut amounts = vec![];
    for a in s.amounts.iter() {
        amounts.push(ingredient::Amount::new(&a.unit, a.value as f32));
    }

    let mut name = String::new();
    if s.ingredient.is_some() {
        name = s.ingredient.unwrap().ingredient.name;
    }
    if s.recipe.is_some() {
        name = s.recipe.unwrap().name;
    }
    return ingredient::Ingredient {
        name,
        modifier: s.adjective,
        amounts,
    };
}
#[allow(dead_code)]
fn bare_detail(name: String) -> IngredientDetail {
    IngredientDetail::new(
        "".to_string(),
        Ingredient::new("".to_string(), name.to_string()),
        vec![],
        vec![],
        vec![],
    )
}

#[cfg(test)]
mod tests {
    use std::collections::HashMap;

    use openapi::models::{
        Amount, Ingredient, IngredientDetail, IngredientKind, RecipeDetail, RecipeSection,
        SectionIngredient,
    };
    use pretty_assertions::assert_eq;

    use crate::{bare_detail, parse_ingredient, sum_ingredients};

    #[test]
    fn test_from_parsed() {
        assert_eq!(
            parse_ingredient("118 grams / 0.5 cups water").unwrap(),
            SectionIngredient {
                original: Some("118 grams / 0.5 cups water".to_string()),
                ingredient: Some(Box::new(bare_detail("water".to_string()))),
                ..SectionIngredient::new(
                    "".to_string(),
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
            SectionIngredient {
                original: Some("118 ml / 0.5 cups water".to_string()),
                ingredient: Some(Box::new(bare_detail("water".to_string()))),
                ..SectionIngredient::new(
                    "".to_string(),
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
            SectionIngredient {
                original: Some("4 oz / 1/2 cup water".to_string()),
                ingredient: Some(Box::new(bare_detail("water".to_string()))),
                ..SectionIngredient::new(
                    "".to_string(),
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
                "foo".to_string(),
                Ingredient::new("a".to_string(), "foo".to_string()),
                vec![],
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
                "bar".to_string(),
                Ingredient::new("b".to_string(), "bar".to_string()),
                vec![],
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
                "bar".to_string(),
                Ingredient::new("b".to_string(), "bar".to_string()),
                vec![],
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
            0,
            "".to_string(),
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
