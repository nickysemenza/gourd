use anyhow::bail;
use ingredient::{rich_text::Rich, unit::Measure};
use openapi::models::{
    Amount, CompactRecipe, CompactRecipeSection, RecipeDetail, RecipeDetailInput, RecipeSection,
    RecipeSectionInput, RecipeSource, SectionIngredient, SectionIngredientInput,
    SectionInstruction, SectionInstructionInput,
};
use tracing::trace;
use uuid::Uuid;

use crate::{
    converter::{measure_to_amount, si_to_ingredient},
    parser::{new_ingredient_parser, section_ingredient_from_parsed},
};

pub fn to_string(cr: CompactRecipe) -> Result<String, anyhow::Error> {
    let mut res = String::new();

    let section1 = format!("name: {}\n", cr.name);
    res.push_str(&section1);
    res.push_str(SEP);
    dbg!(res.clone());
    for s in cr.sections.into_iter() {
        for i in s.ingredients.into_iter() {
            res.push_str(&i);
            res.push('\n');
        }
        for i in s.instructions.into_iter() {
            res.push_str(format!(";{}", i).as_str());
            res.push('\n');
        }
        res.push('\n');
    }
    Ok(res.trim_end().to_string())
}
pub fn from_string(r: String) -> Result<CompactRecipe, anyhow::Error> {
    trace!("decoding {}", r);
    let parts: Vec<&str> = r.trim_start_matches(SEP).split(SEP).collect();
    if parts.len() != 2 {
        bail!("expected 2 parts: {:#?}", parts);
    }
    let name = parts[0]
        .to_string()
        .strip_prefix("name: ")
        .unwrap()
        .trim_end()
        .to_string();
    // let mut compact: CompactRecipe = serde_yaml::from_str(parts[0])?;
    let compact = CompactRecipe {
        name,
        url: None,
        image: None,
        id: Uuid::new_v4().to_string(),
        sections: parts[1]
            .split("\n\n")
            .collect::<Vec<&str>>()
            .iter()
            .map(|section_text_chunk| {
                let mut section = CompactRecipeSection {
                    ingredients: vec![],
                    instructions: vec![],
                };
                section_text_chunk.split("\n").into_iter().for_each(|l| {
                    match l.strip_prefix(";") {
                        Some(i) => section.instructions.push(i.to_string()),
                        None => section.ingredients.push(l.to_string()),
                    }
                });
                section
            })
            .collect(),
    };
    Ok(compact)
}
// condense down recipe detail input into a compact recipe
pub fn compact_recipe(r: RecipeDetailInput) -> CompactRecipe {
    let mut sections = Vec::new();
    for s in r.sections.iter() {
        let mut sec = CompactRecipeSection {
            ingredients: vec![],
            instructions: vec![],
        };
        for ing in s.ingredients.clone().into_iter() {
            let mut ing2 = ing.clone();
            ing2.amounts
                .retain(|a| a.source.as_ref().unwrap_or(&"".to_string()) != "calculated");
            sec.ingredients.push(si_to_ingredient(ing2).to_string());
        }
        for ins in s.instructions.iter() {
            sec.instructions.push(ins.instruction.clone());
        }
        sections.push(sec);
    }
    return CompactRecipe {
        id: Uuid::new_v4().to_string(),
        name: r.name,
        url: None,
        image: None,
        sections,
    };
}

// turn the recipe into a text block
pub fn encode_recipe(r: RecipeDetailInput) -> Result<String, anyhow::Error> {
    let compact = compact_recipe(r);
    to_string(compact)
}
pub fn decode_recipe(r: String) -> Result<(RecipeDetailInput, Vec<Rich>), anyhow::Error> {
    let compact = from_string(r)?;
    expand_recipe(compact)
}
// handle some cases that the parser can't yet
pub fn fixup_ingredient_line_item(ing: String) -> String {
    ing.replace("Vegetable, canola or coconut oil", "Vegetable oil,")
        .replace(
            "Diamond Crystal kosher salt or 1 teaspoon table salt (4g)",
            "Diamond Crystal kosher salt",
        )
        .replace("unrefined or light brown sugar", "brown sugar")
        .replace(
            "fresh Thai basil leaves or regular basil leaves",
            "thai basil",
        )
        .replace("Egg, (s)", "egg")
}
const SEP: &str = "---\n";
pub fn expand_recipe(r: CompactRecipe) -> Result<(RecipeDetailInput, Vec<Rich>), anyhow::Error> {
    let mut rtt: Vec<Rich> = vec![];
    let sections = r
        .sections
        .into_iter()
        .map(|s| {
            let mut instructions = vec![];
            let mut ingredients = vec![];
            let mut total_time = Measure::parse_new("second", 0.0);

            for ing in s.ingredients.into_iter() {
                ingredients.push(section_ingredient_from_parsed(
                    ingredient::from_str(&fixup_ingredient_line_item(ing.clone())),
                    &ing,
                ))
            }
            let rtp = ingredient::rich_text::RichParser {
                ingredient_names: ingredients
                    .iter()
                    .map(|i| i.name.clone())
                    .filter_map(|x| x)
                    .collect(),
                ip: new_ingredient_parser(true),
            };
            for i in s.instructions.into_iter() {
                let rich_text_tokens = dbg!(rtp.clone().parse(&i).unwrap_or_default());
                rtt.push(rich_text_tokens.clone());

                for token in rich_text_tokens.into_iter() {
                    match token {
                        ingredient::rich_text::Chunk::Measure(amt) => {
                            for m in amt.into_iter() {
                                if m.kind().unwrap() == ingredient::unit::kind::MeasureKind::Time {
                                    total_time = total_time.add(m).unwrap();
                                }
                            }
                        }
                        _ => {}
                    }
                }
                instructions.push(SectionInstructionInput::new(
                    i.strip_prefix(" ")
                        .unwrap_or(&i) // trim leading space if exsists to support `;` or `; `
                        .to_string(),
                ))
            }
            RecipeSectionInput {
                duration: match total_time.values().0 == 0.0 {
                    true => None,
                    false => Some(Box::new(Amount {
                        source: Some("parsed sum".to_string()),

                        ..measure_to_amount(total_time)
                    })),
                },
                ..RecipeSectionInput::new(instructions, ingredients)
            }
        })
        .collect();

    Ok((
        RecipeDetailInput {
            sources: match r.url {
                Some(u) => Some(vec![RecipeSource {
                    url: Some(u),
                    image_url: r.image,
                    ..RecipeSource::new()
                }]),
                None => None,
            },
            ..RecipeDetailInput::new(sections, r.name, 0, "".to_string(), vec![])
        },
        rtt,
    ))
}
// turn the text block back into a recipe

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
        target_id: if let "" = s.id.as_str() {
            None
        } else {
            Some(s.id.clone())
        },
        ..SectionIngredientInput::new(s.kind, s.amounts.clone())
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
        r.tags,
    )
}

#[cfg(test)]
mod tests {

    use openapi::models::{
        Amount, IngredientKind, RecipeDetail, RecipeSection, SectionIngredient, SectionInstruction,
    };
    use pretty_assertions::assert_eq;

    use crate::{
        codec::{decode_recipe, encode_recipe, recipe_to_input},
        converter::bare_detail,
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
            vec![],
            0,
            "".to_string(),
            0,
            false,
            "".to_string(),
            vec![],
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
        assert_eq!(
            dbg!(encode_recipe(recipe_to_input(r.clone()))).unwrap(),
            raw
        );
        let decoded = decode_recipe(raw.to_string()).unwrap().0;
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
173 g / 1.3 cups flour
6 g / 2 tsp salt
6 g / 1 tsp baking soda
;food processor until combined
;add to mixer

100 g / 0.5 recipe CS Pecan Brittle
100 g / 1 cup oats
;add to mixer";
        let recipe = decode_recipe(r.to_string()).unwrap().0;
        assert_eq!(r, encode_recipe(recipe).unwrap());
    }
}
