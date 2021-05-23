use std::collections::{hash_map::Entry, HashMap};

use openapi::models::{
    unit_conversion_request::Target, Amount, Ingredient, IngredientKind, RecipeDetail,
    RecipeSection, SectionIngredient, SectionInstruction, UnitConversionRequest, UnitMapping,
};
use unit::MeasureKind;

#[macro_use]
extern crate serde;

pub mod unit;
fn section_ingredient_from_parsed(i: ingredient::Ingredient, original: &str) -> SectionIngredient {
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
        original: Some(original.to_string()),
        ..SectionIngredient::new("".to_string(), IngredientKind::Ingredient, grams)
    };
}
pub fn parse_ingredient(s: &str) -> Result<SectionIngredient, String> {
    let i = dbg!(ingredient::from_str(s, true))?;
    Ok(section_ingredient_from_parsed(i, s))
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
pub fn parse_unit_mappings(um: Vec<UnitMapping>) -> Vec<(unit::Measure, unit::Measure)> {
    um.iter()
        .map(|u| {
            return (
                amount_to_measure(u.a.clone()),
                amount_to_measure(u.b.clone()),
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
    return match amount_to_measure(req.input[0].clone()).convert(target, equivalencies) {
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
pub fn si_to_ingredient(s: SectionIngredient) -> ingredient::Ingredient {
    let mut amounts = vec![];
    if s.grams > 0.0 {
        amounts.push(ingredient::Amount {
            unit: "g".to_string(),
            value: s.grams as f32,
        })
    }
    if s.amount.is_some() && s.unit.is_some() {
        amounts.push(ingredient::Amount {
            unit: s.unit.unwrap(),
            value: s.amount.unwrap() as f32,
        })
    }
    let mut name = String::new();
    if s.ingredient.is_some() {
        name = s.ingredient.unwrap().name;
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
#[derive(Serialize, Deserialize)]
pub enum CompactRecipeLine {
    Ing(ingredient::Ingredient),
    Ins(String),
}
#[derive(Serialize, Deserialize)]
pub struct CompactRecipe(Vec<Vec<CompactRecipeLine>>);

pub fn compact_recipe(r: RecipeDetail) -> CompactRecipe {
    let mut sections = Vec::new();
    for s in r.sections.iter() {
        let mut sec = Vec::new();
        for ing in s.ingredients.clone().into_iter() {
            sec.push(CompactRecipeLine::Ing(si_to_ingredient(ing)));
        }
        for ins in s.instructions.iter() {
            sec.push(CompactRecipeLine::Ins(ins.instruction.clone()));
        }
        sections.push(sec);
    }
    return CompactRecipe(sections);
}

pub fn encode_recipe(r: RecipeDetail) -> String {
    let mut res = String::new();
    for s in compact_recipe(r).0.into_iter() {
        for i in s.into_iter() {
            res.push_str(
                match i {
                    CompactRecipeLine::Ing(ing) => ing.to_string(),
                    CompactRecipeLine::Ins(ins) => format!(";{}", ins),
                }
                .as_str(),
            );
            res.push('\n');
        }
        res.push('\n');
    }
    return res.trim_end().to_string();
}
pub fn decode_recipe(r: String) -> RecipeDetail {
    let raw_sections: Vec<&str> = r.split("\n\n").collect();

    let sections = dbg!(raw_sections)
        .into_iter()
        .map(|s| {
            let mut instructions = vec![];
            let mut ingredients = vec![];
            let lines: Vec<&str> = s.split("\n").collect();
            for line in lines.into_iter() {
                match dbg!(line).strip_prefix(";") {
                    Some(i) => {
                        instructions.push(SectionInstruction::new("".to_string(), i.to_string()))
                    }

                    None => ingredients.push(parse_ingredient(line).unwrap()),
                };
            }
            RecipeSection::new("".to_string(), instructions, ingredients)
        })
        .collect();

    RecipeDetail::new("".to_string(), sections, "".to_string(), 0, "".to_string())
}

#[cfg(test)]
mod tests {
    use super::*;
    use openapi::models::{
        IngredientKind, RecipeDetail, RecipeSection, SectionIngredient, SectionInstruction,
    };

    #[test]
    fn test_encode() {
        let si_1 = SectionIngredient {
            ingredient: Some(Ingredient::new("".to_string(), "foo".to_string())),
            ..SectionIngredient::new("".to_string(), IngredientKind::Ingredient, 12.0)
        };
        let si_2 = SectionIngredient {
            ingredient: Some(Ingredient::new("".to_string(), "bar".to_string())),
            amount: Some(1.5),
            unit: Some("cups".to_string()),
            ..SectionIngredient::new("".to_string(), IngredientKind::Ingredient, 14.0)
        };
        let si_3 = SectionIngredient {
            ingredient: Some(Ingredient::new("".to_string(), "bar".to_string())),
            ..SectionIngredient::new("".to_string(), IngredientKind::Ingredient, 2.0)
        };
        let r = RecipeDetail::new(
            "".to_string(),
            vec![
                RecipeSection::new(
                    "".to_string(),
                    vec![SectionInstruction::new("".to_string(), "inst1".to_string())],
                    vec![si_1.clone(), si_2.clone(), si_3.clone()],
                ),
                RecipeSection::new(
                    "".to_string(),
                    vec![
                        SectionInstruction::new("".to_string(), "inst2".to_string()),
                        SectionInstruction::new("".to_string(), "inst3".to_string()),
                    ],
                    vec![si_3.clone()],
                ),
            ],
            "".to_string(),
            0,
            "".to_string(),
        );
        let raw = "12 g foo
14 g / 1.5 cups bar
2 g bar
;inst1

2 g bar
;inst2
;inst3";
        assert_eq!(encode_recipe(r.clone()), raw);
        let decoded = decode_recipe(raw.to_string());
        assert_eq!(dbg!(decoded), dbg!(r));
    }
    #[test]
    fn test_encode_decode() {
        let r = "113 g / 1 stick butter
;brown, add to mixer

113 g / 1 stick butter, cold
;add  to mixer with hot brown butter, 
;let cool

150 g / 75 cups brown sugar
100 g / 5 cups sugar
;add to butters, cream

100 g / 2 large eggs, cold
13 g / 1 tbsp vanilla extract
;add to sugar/butter, mix

100 g / 0.5 recipe CS Pecan Brittle
100 g / 1 cup oats
173 g / 1.3 cup flour
6 g / 2 tsp salt
6 g / 1 tsp baking soda
;food processor until combined
;add to mixer

100 g / 0.5 recipe CS Pecan Brittle
100 g / 1 cup oats
;add to mixer";
        let recipe = decode_recipe(r.to_string());
        assert_eq!(r, encode_recipe(recipe));
    }

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
