mod utils;

use wasm_bindgen::prelude::*;

// When the `wee_alloc` feature is enabled, use `wee_alloc` as the global
// allocator.
#[cfg(feature = "wee_alloc")]
#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

#[wasm_bindgen]
extern "C" {
    fn alert(s: &str);
}

#[wasm_bindgen]
pub fn greet() {
    alert("Hello, gourd!");
}

#[wasm_bindgen]
pub fn parse(input: &str) -> String {
    utils::set_panic_hook();
    ingredient::from_str(input, true).unwrap().to_string()
    // return "foo".to_string();
}

// #[wasm_bindgen]
// pub fn parse2(input: &str) -> IngredientA {
//     let i = ingredient::from_str(input, true).unwrap();
//     IngredientA {
//         name: i.name,
//         foo: true,
//         amounts: i.amounts,
//         modifier: i.modifier,
//     }
// }

#[wasm_bindgen(typescript_custom_section)]
const ITEXT_STYLE: &'static str = r#"
interface Ingredient {
    amounts: Amount[];
    modifier: string;
    name: string;
  }
  
  interface Amount {
    unit: string;
    value: number;
  }
"#;

#[wasm_bindgen]
extern "C" {
    #[wasm_bindgen()]
    pub type ITextStyle;
    #[wasm_bindgen(typescript_type = "-1 | 0 | 1")]
    type Foo;
}

#[wasm_bindgen]
pub fn parse3(input: &str) -> JsValue {
    let i = ingredient::from_str(input, true).unwrap();
    JsValue::from_serde(&i).unwrap()
}

#[wasm_bindgen]
pub fn parse4(input: &str) -> JsValue {
    let si = gourd_common::parse_ingredient(input).unwrap();
    JsValue::from_serde(&si).unwrap()
}


#[wasm_bindgen]
pub fn format_ingredient(val: &JsValue) -> String {
    let i: ingredient::Ingredient = val.into_serde().unwrap();
    return i.to_string();
}

// #[wasm_bindgen]
// pub fn parse2(input: &str) -> Ingredient {
//     ingredient(input).unwrap();
// }

// #[wasm_bindgen]
// #[derive(Default)]
// pub struct IngredientA {
//     name: String,
//     foo: bool,
//     amounts: Vec<Amount>,
//     modifier: Option<String>,
// }
