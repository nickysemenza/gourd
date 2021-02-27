use std::collections::{hash_map::Entry, HashMap};

use openapi::models::{section_ingredient::Kind, Ingredient, RecipeDetail, SectionIngredient};

pub fn parse_ingredient(s: &str) -> Result<SectionIngredient, String> {
    let i = ingredient::from_str(s, true)?;

    let (unit, amount) = match i.amounts.iter().find(|&x| !is_gram(x)) {
        Some(i) => (Some(i.unit.clone()), Some(i.value as f64)),
        None => (None, None),
    };

    Ok(SectionIngredient {
        unit,
        amount,
        adjective: i.modifier,
        ingredient: Some(Ingredient::new("".to_string(), i.name)),
        ..SectionIngredient::new(
            "".to_string(),
            Kind::Ingredient,
            match i.amounts.iter().find(|&x| is_gram(x)) {
                Some(g) => g.value.into(),
                None => 0.0,
            },
        )
    })
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
    a.unit == "g" || a.unit == "grams"
}

#[cfg(test)]
mod tests {
    use super::*;
    use openapi::models::{
        section_ingredient::Kind, RecipeDetail, RecipeSection, SectionIngredient,
    };

    #[test]
    fn it_works() {
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
