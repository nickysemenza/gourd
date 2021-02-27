

use openapi::models::{section_ingredient::Kind, Ingredient, SectionIngredient};

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

fn is_gram(a: &ingredient::Amount) -> bool {
    a.unit == "g" || a.unit == "grams"
}


#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}