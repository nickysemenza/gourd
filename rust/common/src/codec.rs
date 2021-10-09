use openapi::models::{
    RecipeDetail, RecipeDetailInput, RecipeSection, RecipeSectionInput, SectionIngredient,
    SectionIngredientInput, SectionInstruction, SectionInstructionInput,
};
use tracing::trace;

use crate::{parse_ingredient, si_to_ingredient};

#[derive(Serialize, Deserialize)]
pub enum CompactRecipeLine {
    Ing(ingredient::Ingredient),
    Ins(String),
}

#[derive(Serialize, Deserialize)]
pub struct CompactRecipe(Vec<Vec<CompactRecipeLine>>);

pub fn compact_recipe(r: RecipeDetailInput) -> CompactRecipe {
    let mut sections = Vec::new();
    for s in r.sections.iter() {
        let mut sec = Vec::new();
        for ing in s.ingredients.clone().into_iter() {
            let mut ing2 = ing.clone();
            ing2.amounts
                .retain(|a| a.source.as_ref().unwrap_or(&"".to_string()) != "calculated");
            sec.push(CompactRecipeLine::Ing(si_to_ingredient(ing2)));
        }
        for ins in s.instructions.iter() {
            sec.push(CompactRecipeLine::Ins(ins.instruction.clone()));
        }
        sections.push(sec);
    }
    return CompactRecipe(sections);
}

pub fn encode_recipe(r: RecipeDetailInput) -> String {
    let mut res = String::new();
    res.push_str(&format!("name: {}\n", r.name));
    res.push_str("---\n");
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
pub fn decode_recipe(r: String) -> RecipeDetailInput {
    trace!("decoding {}", r);
    let mut name = String::new();
    let parts: Vec<&str> = r.split("---\n").collect();
    if parts.len() == 2 {
        parts[0]
            .split("\n")
            .collect::<Vec<&str>>()
            .into_iter()
            .for_each(|s| match s.strip_prefix("name: ") {
                Some(n) => name = dbg!(n).to_string(),
                None => {}
            })
    }
    let raw_sections: Vec<&str> = dbg!(parts).last().unwrap().split("\n\n").collect();

    let sections = dbg!(raw_sections)
        .into_iter()
        .map(|s| {
            let mut instructions = vec![];
            let mut ingredients = vec![];
            let lines: Vec<&str> = s.split("\n").collect();
            for line in lines.into_iter() {
                match dbg!(line).strip_prefix(";") {
                    Some(i) => instructions.push(SectionInstructionInput::new(i.to_string())),

                    None => ingredients.push(parse_ingredient(line).unwrap()),
                };
            }
            RecipeSectionInput::new(instructions, ingredients)
        })
        .collect();

    RecipeDetailInput::new(sections, name, 0, "".to_string())
}
pub fn section_to_input(s: &RecipeSection) -> RecipeSectionInput {
    RecipeSectionInput::new(
        s.instructions
            .iter()
            .map(|i| section_instruction_to_input(i))
            .collect(),
        s.ingredients
            .iter()
            .map(|i| section_ingredient_to_input(i))
            .collect(),
    )
}
pub fn section_ingredient_to_input(s: &SectionIngredient) -> SectionIngredientInput {
    SectionIngredientInput {
        name: match s.kind {
            openapi::models::IngredientKind::Ingredient => {
                Some(s.ingredient.clone().unwrap().ingredient.name.clone())
            }

            openapi::models::IngredientKind::Recipe => Some(s.recipe.clone().unwrap().name.clone()),
        },
        original: s.original.clone(),
        ..SectionIngredientInput::new(s.id.clone(), s.kind, s.amounts.clone())
    }
}
pub fn section_instruction_to_input(s: &SectionInstruction) -> SectionInstructionInput {
    SectionInstructionInput::new(s.instruction.clone())
}
pub fn recipe_to_input(r: RecipeDetail) -> RecipeDetailInput {
    RecipeDetailInput::new(
        r.sections
            .into_iter()
            .map(|s| section_to_input(&s))
            .collect(),
        r.name,
        r.quantity,
        r.unit,
    )
}
// pub fn recipe_from_input(r: RecipeDetailInput) -> RecipeDetail {
//     RecipeDetail::new(
//         r.id,
//         // r.sections,
//         vec![], // TODO!
//         r.name,
//         r.quantity,
//         r.unit,
//         0,
//         false,
//         "".to_string(),
//     )
// }

#[cfg(test)]
mod tests {

    use openapi::models::{
        Amount, IngredientKind, RecipeDetail, RecipeSection, SectionIngredient, SectionInstruction,
    };
    use pretty_assertions::assert_eq;

    use crate::{
        bare_detail,
        codec::{decode_recipe, encode_recipe, recipe_to_input},
    };

    #[test]
    fn test_encode() {
        let si_1 = SectionIngredient {
            ingredient: Some(Box::new(bare_detail("foo".to_string()))),
            original: Some("12 g foo".to_string()),
            ..SectionIngredient::new(
                "".to_string(),
                IngredientKind::Ingredient,
                vec![Amount::new("g".to_string(), 12.0)],
            )
        };
        let si_2 = SectionIngredient {
            ingredient: Some(Box::new(bare_detail("bar".to_string()))),
            original: Some("14 g / 1.5 cups bar".to_string()),
            ..SectionIngredient::new(
                "".to_string(),
                IngredientKind::Ingredient,
                vec![
                    Amount::new("g".to_string(), 14.0),
                    Amount::new("cups".to_string(), 1.5),
                ],
            )
        };
        let si_3 = SectionIngredient {
            ingredient: Some(Box::new(bare_detail("bar".to_string()))),
            original: Some("2 g bar".to_string()),
            ..SectionIngredient::new(
                "".to_string(),
                IngredientKind::Ingredient,
                vec![Amount::new("g".to_string(), 2.0)],
            )
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
            "cake".to_string(),
            0,
            "".to_string(),
            0,
            false,
            "".to_string(),
        );
        let raw = "name: cake
---
12 g foo
14 g / 1.5 cups bar
2 g bar
;inst1

2 g bar
;inst2
;inst3";
        assert_eq!(encode_recipe(recipe_to_input(r.clone())), raw);
        let decoded = decode_recipe(raw.to_string());
        assert_eq!(dbg!(decoded), dbg!(recipe_to_input(r)));
    }
    #[test]
    fn test_encode_decode() {
        let r = "name: cookies
---
113 g / 1 stick butter
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
}
