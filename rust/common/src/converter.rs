use std::collections::{hash_map::Entry, HashMap};

use openapi::models::{
    unit_conversion_request::Target, Amount, Ingredient, IngredientDetail, IngredientKind,
    RecipeDetail, SectionIngredient, SectionIngredientInput, UnitConversionRequest,
};

use tracing::info;

use ingredient::unit::{kind::MeasureKind, Measure};

use crate::{amount_from_ingredient, parser::parse_unit_mappings};

#[tracing::instrument]
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
        .convert_measure_via_mappings(target.clone(), equivalencies.clone())
    {
        Some(a) => Some(measure_to_amount(a)),
        None => {
            if target == MeasureKind::Weight {
                // try again to convert to ml, and then use that as grams
                return match amount_to_measure(req.input[0].clone())
                    .convert_measure_via_mappings(MeasureKind::Volume, equivalencies)
                {
                    Some(a) => {
                        let mut a = measure_to_amount(a);
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
    let m = measure_to_amount(amount_to_measure(a.clone()));
    return Amount {
        unit: m.unit,
        value: m.value,
        upper_value: m.upper_value,
        source: a.source.clone(),
    };
}
pub fn amount_to_measure(a: Amount) -> Measure {
    Measure::from_parts(a.unit.as_ref(), a.value, a.upper_value)
}

pub fn measure_to_amount(m: Measure) -> Amount {
    amount_from_ingredient(&m)
}
pub fn si_to_ingredient(s: SectionIngredientInput) -> ingredient::Ingredient {
    let mut amounts = vec![];
    for a in s.amounts.iter() {
        amounts.push(Measure::from_parts(a.unit.as_ref(), a.value, a.upper_value));
    }

    return ingredient::Ingredient {
        name: s.name.unwrap_or_default(),
        modifier: s.adjective,
        amounts,
    };
}
#[allow(dead_code)]
pub fn bare_detail(name: String) -> IngredientDetail {
    IngredientDetail::new(
        Ingredient::new("".to_string(), name.to_string()),
        vec![],
        vec![],
    )
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

#[cfg(test)]
mod tests {
    use openapi::models::{
        unit_conversion_request::Target, Amount, UnitConversionRequest, UnitMapping,
    };

    use crate::convert_to;

    #[test]
    fn test_conversion_request() {
        let unit_mappings = vec![
            UnitMapping::new(
                Amount::new("piece".to_string(), 1.0),
                Amount::new("grams".to_string(), 100.0),
            ),
            UnitMapping::new(
                Amount::new("cents".to_string(), 1.0),
                Amount::new("grams".to_string(), 1.0),
            ),
        ];
        assert_eq!(
            convert_to(UnitConversionRequest {
                target: Some(Target::Weight),
                input: vec![Amount::new("piece".to_string(), 0.5)],
                unit_mappings: unit_mappings.clone(),
            }),
            Some(Amount::new("g".to_string(), 50.0)),
        );
        assert_eq!(
            convert_to(UnitConversionRequest {
                target: Some(Target::Money),
                input: vec![Amount::new("piece".to_string(), 0.5)],
                unit_mappings: unit_mappings.clone(),
            }),
            Some(Amount::new("$".to_string(), 0.5)),
        );
        assert_eq!(
            convert_to(UnitConversionRequest {
                target: Some(Target::Calories),
                input: vec![Amount::new("piece".to_string(), 0.5)],
                unit_mappings: unit_mappings.clone(),
            }),
            None,
        );
    }
}
