mod utils;

use gourd_common::{
    convert_to, ingredient, parse_unit_mappings, sum_ingredients,
    unit::{make_graph, print_graph},
};
use openapi::models::{RecipeDetail, RecipeDetailInput, UnitConversionRequest};
use wasm_bindgen::prelude::*;

// When the `wee_alloc` feature is enabled, use `wee_alloc` as the global
// allocator.
#[cfg(feature = "wee_alloc")]
#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

#[wasm_bindgen]
pub fn parse(input: &str) -> String {
    utils::set_panic_hook();
    ingredient::from_str(input).to_string()
    // return "foo".to_string();
}

#[wasm_bindgen(typescript_custom_section)]
const ITEXT_STYLE: &'static str = r#"
interface Ingredient {
    amounts: Amount[];
    modifier?: string;
    name: string;
  }
  
interface Amount {
  unit: string;
  value: number;
  upper_value?: number;
}

interface CompactR {
  Ing?: Ingredient;
  Ins?: string;
}
export type RichItem =
  | { kind: "Text"; value: string }
  | { kind: "Amount"; value: Amount[] }
"#;

#[wasm_bindgen]
extern "C" {
    #[wasm_bindgen(typescript_type = "Ingredient")]
    pub type IIngredient;
    #[wasm_bindgen(typescript_type = "Amount")]
    pub type IAmount;
    #[wasm_bindgen(typescript_type = "Amount[]")]
    pub type IAmounts;
    #[wasm_bindgen(typescript_type = "CompactR[][]")]
    pub type ICompactR;
    #[wasm_bindgen(typescript_type = "RichItem[]")]
    pub type RichItems;
}

#[wasm_bindgen]
pub fn parse2(input: &str) -> Result<IIngredient, JsValue> {
    let i = ingredient::from_str(input);
    Ok(JsValue::from_serde(&i).unwrap().into())
}

#[wasm_bindgen]
pub fn parse3(input: &str) -> JsValue {
    let i = ingredient::from_str(input);
    JsValue::from_serde(&i).unwrap()
}

#[wasm_bindgen]
pub fn parse4(input: &str) -> JsValue {
    let si = gourd_common::parse_ingredient(input).unwrap();
    JsValue::from_serde(&si).unwrap()
}

#[wasm_bindgen]
pub fn format_ingredient(val: &IIngredient) -> String {
    let i: ingredient::Ingredient = val.into_serde().unwrap();
    return i.to_string();
}

#[wasm_bindgen]
pub fn sum_ingr(recipe_detail: &JsValue) -> JsValue {
    let r: RecipeDetail = recipe_detail.into_serde().unwrap();
    let res = sum_ingredients(r);
    JsValue::from_serde(&res).unwrap()
}

#[wasm_bindgen]
pub fn dolla(conversion_request: &JsValue) -> Result<IAmount, JsValue> {
    utils::set_panic_hook();
    let req: UnitConversionRequest = conversion_request.into_serde().unwrap();
    return match convert_to(req) {
        Some(a) => Ok(JsValue::from_serde(&a).unwrap().into()),
        None => Err(JsValue::from_str("no parse result")),
    };
}

#[wasm_bindgen]
pub fn parse_amount(input: &str) -> Result<IAmounts, JsValue> {
    utils::set_panic_hook();
    let i = ingredient::parse_amount(input);
    Ok(JsValue::from_serde(&i).unwrap().into())
}

// pub fn decode_recipe(input: &str)
#[wasm_bindgen]
pub fn encode_recipe_text(recipe_detail: &JsValue) -> String {
    utils::set_panic_hook();
    let r: RecipeDetailInput = recipe_detail.into_serde().unwrap();
    gourd_common::codec::encode_recipe(r)
}
#[wasm_bindgen]
pub fn encode_recipe_to_compact_json(recipe_detail: &JsValue) -> ICompactR {
    utils::set_panic_hook();
    let r: RecipeDetailInput = recipe_detail.into_serde().unwrap();
    let c = gourd_common::codec::compact_recipe(r);
    JsValue::from_serde(&c).unwrap().into()
}

#[wasm_bindgen]
pub fn decode_recipe_text(r: String) -> JsValue {
    utils::set_panic_hook();
    let detail = gourd_common::codec::decode_recipe(r);
    JsValue::from_serde(&detail).unwrap()
}

#[wasm_bindgen]
pub fn make_dag(conversion_request: &JsValue) -> String {
    utils::set_panic_hook();
    let req: UnitConversionRequest = conversion_request.into_serde().unwrap();

    let equivalencies = parse_unit_mappings(req.unit_mappings);
    let g = make_graph(equivalencies);
    return print_graph(g);
}

#[wasm_bindgen]
pub fn rich(r: String) -> Result<RichItems, JsValue> {
    utils::set_panic_hook();
    match ingredient::rich_text::parse(r.as_str()) {
        Ok(r) => Ok(JsValue::from_serde(&r).unwrap().into()),
        Err(e) => Err(JsValue::from_str(&e.to_string())),
    }
}
