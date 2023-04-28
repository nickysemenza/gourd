use anyhow::Result;
use ingredient::{
    unit::{self, Measure, Unit},
    IngredientParser,
};
use openapi::models::{Amount, IngredientKind, SectionIngredientInput, UnitMapping};

pub use crate::converter::amount_to_measure;

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
    for y in i.amounts.iter() {
        let measure_unit = y.unit();
        let (value, upper_value, foo1) = y.values();
        if measure_unit == Unit::Gram {
            grams = value;
            amount = Some(grams);
            unit = Some(foo1);
        } else if measure_unit == Unit::Ounce {
            oz = value;
        } else if measure_unit == Unit::Milliliter {
            ml = value;
        } else {
            unit = Some(foo1);
            amount = Some(value);
        }
        if amount.is_some() && unit.is_some() {
            let unitstr = unit.clone().unwrap_or("unknown".to_string());
            if unitstr == IngredientKind::Recipe.to_string() {
                // if the ingredient amount has a unit of "recipe", then it's likely a recipe
                kind = IngredientKind::Recipe;
            }
            let mut a = Amount::new(unitstr, amount.unwrap_or(0.0));
            a.upper_value = upper_value;
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
pub fn parse_amount(s: &str) -> Result<Vec<unit::Measure>> {
    new_ingredient_parser(false).parse_amount(s)
}

pub fn new_ingredient_parser(is_rich_text: bool) -> IngredientParser {
    IngredientParser::new(is_rich_text)
}

pub fn amount_from_ingredient(e1: &unit::Measure) -> Amount {
    let (value, upper_value, unit) = e1.normalize().values();
    Amount {
        unit,
        value,
        upper_value,
        source: None,
    }
}

pub fn parse_unit_mappings(um: Vec<UnitMapping>) -> Vec<(Measure, Measure)> {
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
                        Amount::new("cup".to_string(), 0.5),
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
                        Amount::new("cup".to_string(), 0.5),
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
