use std::collections::{hash_map::Entry, HashMap};

use openapi::models::{section_ingredient::Kind, Ingredient, RecipeDetail, SectionIngredient};

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
        ..SectionIngredient::new("".to_string(), Kind::Ingredient, grams)
    };
}
pub fn parse_ingredient(s: &str) -> Result<SectionIngredient, String> {
    let i = dbg!(ingredient::from_str(s, true))?;
    Ok(section_ingredient_from_parsed(i))
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
            Kind::Recipe => (i.recipe.as_ref().unwrap().id.clone(), &mut recipes),
            Kind::Ingredient => (i.ingredient.as_ref().unwrap().id.clone(), &mut ing),
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

#[cfg(test)]
mod tests {
    use super::*;
    use openapi::models::{
        section_ingredient::Kind, RecipeDetail, RecipeSection, SectionIngredient,
    };

    #[test]
    fn test_from_parsed() {
        assert_eq!(
            parse_ingredient("118 grams / 0.5 cups water").unwrap(),
            SectionIngredient {
                amount: Some(0.5),
                unit: Some("cups".to_string()),
                ingredient: Some(Ingredient::new("".to_string(), "water".to_string())),
                ..SectionIngredient::new("".to_string(), Kind::Ingredient, 118.0)
            }
        );
        // ml is parsed as grams if not present
        assert_eq!(
            parse_ingredient("118 ml / 0.5 cups water").unwrap(),
            SectionIngredient {
                amount: Some(0.5),
                unit: Some("cups".to_string()),
                ingredient: Some(Ingredient::new("".to_string(), "water".to_string())),
                ..SectionIngredient::new("".to_string(), Kind::Ingredient, 118.0)
            }
        );
        assert_eq!(
            parse_ingredient("4 oz / 1/2 cup water").unwrap(),
            SectionIngredient {
                amount: Some(0.5),
                unit: Some("cup".to_string()),
                ingredient: Some(Ingredient::new("".to_string(), "water".to_string())),
                ..SectionIngredient::new("".to_string(), Kind::Ingredient, 113.0)
            }
        );
    }

    #[test]
    fn test_sum_ingredients() {
        let si_1 = SectionIngredient {
            ingredient: Some(Ingredient::new("a".to_string(), "foo".to_string())),
            ..SectionIngredient::new("".to_string(), Kind::Ingredient, 12.0)
        };
        let si_2 = SectionIngredient {
            ingredient: Some(Ingredient::new("b".to_string(), "bar".to_string())),
            ..SectionIngredient::new("".to_string(), Kind::Ingredient, 14.0)
        };
        let si_3 = SectionIngredient {
            ingredient: Some(Ingredient::new("b".to_string(), "bar".to_string())),
            ..SectionIngredient::new("".to_string(), Kind::Ingredient, 2.0)
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
