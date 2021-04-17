use std::collections::{hash_map::Entry, HashMap};

use openapi::models::{
    unit_conversion_request::Target, Amount, Ingredient, IngredientKind, RecipeDetail,
    SectionIngredient, UnitConversionRequest,
};
use unit::MeasureKind;

pub mod unit;
fn section_ingredient_from_parsed(i: ingredient::Ingredient) -> SectionIngredient {
    let mut grams = 0.0;
    let mut oz = 0.0;
    let mut ml = 0.0;
    let mut unit: Option<String> = None;
    let mut amount: Option<f64> = None;

    for x in i.amounts.iter() {
        if is_gram(x) {
            grams = x.value.into();
        } else if is_oz(x) {
            oz = x.value.into();
        } else if is_ml(x) {
            ml = x.value.into();
        } else {
            unit = Some(x.unit.clone());
            amount = Some(x.value as f64);
        }
    }

    // if grams are absent, try converting oz
    if grams == 0.0 && oz != 0.0 {
        grams = ((oz * 28.35) as f64).round();
    } else if unit.is_none() && oz != 0.0 {
        unit = Some("oz".to_string());
        amount = Some(oz);
    }

    // if grams are absent, try using mL
    if grams == 0.0 && ml != 0.0 {
        grams = ml
    } else if unit.is_none() && ml != 0.0 {
        unit = Some("ml".to_string());
        amount = Some(ml);
    }

    return SectionIngredient {
        unit,
        amount,
        adjective: i.modifier,
        ingredient: Some(Ingredient::new("".to_string(), i.name)),
        ..SectionIngredient::new("".to_string(), IngredientKind::Ingredient, grams)
    };
}
pub fn parse_ingredient(s: &str) -> Result<SectionIngredient, String> {
    let i = dbg!(ingredient::from_str(s, true))?;
    Ok(section_ingredient_from_parsed(i))
}
pub fn parse_amount(s: &str) -> Result<Vec<ingredient::Amount>, String> {
    let i = dbg!(ingredient::parse_amount(s))?;
    Ok(i)
}
pub fn get_grams(r: RecipeDetail) -> f64 {
    r.sections.iter().fold(0.0, |acc, s| {
        acc + s.ingredients.iter().fold(0.0, |acc, i| acc + i.grams)
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
            IngredientKind::Ingredient => (i.ingredient.as_ref().unwrap().id.clone(), &mut ing),
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

fn is_gram(a: &ingredient::Amount) -> bool {
    a.unit == "g" || a.unit == "gram" || a.unit == "grams"
}
fn is_oz(a: &ingredient::Amount) -> bool {
    a.unit == "oz" || a.unit == "ounce" || a.unit == "ounces"
}
fn is_ml(a: &ingredient::Amount) -> bool {
    a.unit == "ml" || a.unit == "mL"
}

pub fn convert_to(req: UnitConversionRequest) -> Option<Amount> {
    let equivalencies: Vec<(unit::Measure, unit::Measure)> = req
        .unit_mappings
        .iter()
        .map(|u| {
            return (
                amount_to_measure(u.a.clone()),
                amount_to_measure(u.b.clone()),
            );
        })
        .collect();
    let target = match req.target.unwrap_or(Target::Other) {
        Target::Weight => MeasureKind::Weight,
        Target::Volume => MeasureKind::Volume,
        Target::Money => MeasureKind::Money,
        Target::Other => MeasureKind::Other,
    };
    let foo = amount_to_measure(req.input[0].clone()).convert(target, equivalencies);
    // return foo.and_then(measure_to_amount);
    // return None;
    return match foo {
        Some(a) => Some(measure_to_amount(a)),
        None => None,
    };
}
pub fn amount_to_measure(a: Amount) -> unit::Measure {
    unit::Measure::parse(unit::BareMeasurement::new(a.unit, a.value as f32))
}
pub fn measure_to_amount(m: unit::Measure) -> Amount {
    let m1 = dbg!(m.as_bare());
    Amount::new(m1.unit, m1.value.into())
}

#[cfg(test)]
mod tests {
    use super::*;
    use openapi::models::{IngredientKind, RecipeDetail, RecipeSection, SectionIngredient};

    #[test]
    fn test_from_parsed() {
        assert_eq!(
            parse_ingredient("118 grams / 0.5 cups water").unwrap(),
            SectionIngredient {
                amount: Some(0.5),
                unit: Some("cups".to_string()),
                ingredient: Some(Ingredient::new("".to_string(), "water".to_string())),
                ..SectionIngredient::new("".to_string(), IngredientKind::Ingredient, 118.0)
            }
        );
        // ml is parsed as grams if not present
        assert_eq!(
            parse_ingredient("118 ml / 0.5 cups water").unwrap(),
            SectionIngredient {
                amount: Some(0.5),
                unit: Some("cups".to_string()),
                ingredient: Some(Ingredient::new("".to_string(), "water".to_string())),
                ..SectionIngredient::new("".to_string(), IngredientKind::Ingredient, 118.0)
            }
        );
        assert_eq!(
            parse_ingredient("4 oz / 1/2 cup water").unwrap(),
            SectionIngredient {
                amount: Some(0.5),
                unit: Some("cup".to_string()),
                ingredient: Some(Ingredient::new("".to_string(), "water".to_string())),
                ..SectionIngredient::new("".to_string(), IngredientKind::Ingredient, 113.0)
            }
        );
    }

    #[test]
    fn test_sum_ingredients() {
        let si_1 = SectionIngredient {
            ingredient: Some(Ingredient::new("a".to_string(), "foo".to_string())),
            ..SectionIngredient::new("".to_string(), IngredientKind::Ingredient, 12.0)
        };
        let si_2 = SectionIngredient {
            ingredient: Some(Ingredient::new("b".to_string(), "bar".to_string())),
            ..SectionIngredient::new("".to_string(), IngredientKind::Ingredient, 14.0)
        };
        let si_3 = SectionIngredient {
            ingredient: Some(Ingredient::new("b".to_string(), "bar".to_string())),
            ..SectionIngredient::new("".to_string(), IngredientKind::Ingredient, 2.0)
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
