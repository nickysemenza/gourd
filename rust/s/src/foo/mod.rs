use openapi::models::Amount;
use serde::{Deserialize, Serialize};
use sqlx::PgPool;

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct Root2 {
    pub name: String,
    pub fdc_id: Option<i32>,
    #[serde(default)]
    pub unit_mappings: Vec<openapi::models::UnitMapping>,
    #[serde(default)]
    pub aliases: Vec<String>,
    pub upc: Option<String>,
}

#[derive(Debug, Clone, Serialize)]
pub struct IngredientUnitMapping {
    id: i32,
    ingredient_id: String,
    unit_a: String,
    amount_a: sqlx::types::BigDecimal,
    unit_b: String,
    amount_b: sqlx::types::BigDecimal,
    source: String,
}

#[tracing::instrument]
pub async fn get_fdc_id_from_upc(pool: &PgPool, upc: String) -> Result<Option<i32>, sqlx::Error> {
    let res = sqlx::query!(
        "select fdc_id from usda_branded_food where gtin_upc = $1 order by fdc_id desc limit 1",
        upc
    )
    .fetch_optional(pool)
    .await?;

    if res.is_some() {
        return Ok(Some(res.unwrap().fdc_id));
    }
    return Ok(None);
}

#[derive(Debug, Clone, Serialize)]
pub struct SI {
    section_id: String,
    id: String,
    sort: Option<i32>,
    ingredient_id: Option<String>,
    recipe_id: Option<String>,
    amounts: sqlx::types::Json<Vec<Amount>>,
    adjective: Option<String>,
    optional: Option<bool>,
    original: Option<String>,
    sub_for_ingredient_id: Option<String>,
}

#[tracing::instrument]
pub async fn get_test(pool: &PgPool) -> Result<Vec<SI>, sqlx::Error> {
    let res = sqlx::query_as!(
        SI,
        r#"
    select section_id, id, sort, ingredient_id, recipe_id, amounts as "amounts: sqlx::types::Json<Vec<Amount>>",
     adjective, optional, original, sub_for_ingredient_id from recipe_section_ingredients;
            "#,
    )
    .fetch_all(pool)
    .await?;

    // dbg!(res);
    // let res2 = res.unwrap();
    Ok(res)
}
