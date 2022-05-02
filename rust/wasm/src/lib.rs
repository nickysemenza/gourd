mod utils;

use gourd_common::{
    convert_to,
    ingredient::Amount,
    parse_unit_mappings, sum_ingredients,
    unit::{add_time_amounts, make_graph, print_graph},
};
use openapi::models::{RecipeDetail, RecipeDetailInput, UnitConversionRequest};
use tracing::{error, info};
use wasm_bindgen::prelude::*;

// When the `wee_alloc` feature is enabled, use `wee_alloc` as the global
// allocator.
#[cfg(feature = "wee_alloc")]
#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

#[wasm_bindgen(start)]
pub fn start() -> Result<(), JsValue> {
    // print pretty errors in wasm https://github.com/rustwasm/console_error_panic_hook
    // This is not needed for tracing_wasm to work, but it is a common tool for getting proper error line numbers for panics.
    console_error_panic_hook::set_once();

    // Add this line:
    tracing_wasm::set_as_global_default();

    info!("hello");
    Ok(())
}

#[wasm_bindgen]
pub fn parse(input: &str) -> String {
    utils::set_panic_hook();
    gourd_common::ingredient::from_str(input).to_string()
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
  | { kind: "Ing"; value: string }
  | { kind: "Amount"; value: Amount[] }
"#;

#[wasm_bindgen]
extern "C" {
    #[wasm_bindgen(typescript_type = "Ingredient")]
    pub type IIngredient;
    #[wasm_bindgen(typescript_type = "Amount")]
    #[derive(Debug)]
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
    let i = gourd_common::ingredient::from_str(input);
    Ok(JsValue::from_serde(&i).unwrap().into())
}

#[wasm_bindgen]
pub fn parse3(input: &str) -> JsValue {
    let i = gourd_common::ingredient::from_str(input);
    JsValue::from_serde(&i).unwrap()
}

#[wasm_bindgen]
pub fn parse4(input: &str) -> JsValue {
    let si = gourd_common::parse_ingredient(input).unwrap();
    JsValue::from_serde(&si).unwrap()
}

#[wasm_bindgen]
pub fn format_ingredient(val: &IIngredient) -> String {
    let i: gourd_common::ingredient::Ingredient = val.into_serde().unwrap();
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
    let ip = gourd_common::new_ingredient_parser();
    let i = ip.parse_amount(input);
    Ok(JsValue::from_serde(&i).unwrap().into())
}

// pub fn decode_recipe(input: &str)
#[wasm_bindgen]
pub fn encode_recipe_text(recipe_detail: &JsValue) -> String {
    utils::set_panic_hook();
    // info!("detail: {:?}", recipe_detail);
    let r: RecipeDetailInput = recipe_detail.into_serde().unwrap();
    gourd_common::codec::encode_recipe(r).unwrap()
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
    let detail = gourd_common::codec::decode_recipe(r).unwrap();
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
pub fn rich(r: String, ings: &JsValue) -> Result<RichItems, JsValue> {
    utils::set_panic_hook();
    let ings2: Vec<String> = ings.into_serde().unwrap();
    info!("rich2: {:?}", ings2);
    let rtp = gourd_common::ingredient::rich_text::RichParser {
        ingredient_names: ings2,
        ip: gourd_common::new_ingredient_parser(),
    };
    match rtp.parse(r.as_str()) {
        Ok(r) => Ok(JsValue::from_serde(&r).unwrap().into()),
        Err(e) => Err(JsValue::from_str(&e.to_string())),
    }
}

#[wasm_bindgen]
pub fn format_amount(amount: &IAmount) -> String {
    utils::set_panic_hook();
    let a1: Result<Amount, _> = amount.into_serde();
    match a1 {
        Ok(a) => format!("{}", a),
        Err(e) => {
            error!("failed to format {:#?}: {:?}", amount, e);
            format!("{}", e)
        }
    }
}

#[wasm_bindgen]
pub fn sum_time_amounts(amount: &IAmounts) -> IAmount {
    utils::set_panic_hook();
    let r: Vec<Amount> = amount.into_serde().unwrap();
    let sum = add_time_amounts(r);
    info!("sum {}", sum);
    JsValue::from_serde(&sum).unwrap().into()
}
