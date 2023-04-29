use axum::{
    extract::{self, Query},
    response::IntoResponse,
    Json,
};
use gourd_common::{
    codec::expand_recipe,
    convert_to,
    ingredient::{self, unit::Measure},
    tmp_normalize,
};
use openapi::models::{Amount, CompactRecipe, RecipeWrapperInput};
use serde::{Deserialize, Serialize};
use tracing::{debug, span};

use crate::error::AppError;

#[derive(Deserialize, Debug)]
pub struct Info {
    text: String,
}

#[derive(Deserialize, Debug)]
pub struct URLInput {
    url: String,
}

pub async fn decode_recipe(info: Query<Info>) -> impl IntoResponse {
    let root = span!(tracing::Level::TRACE, "decode_recipe",);
    let _enter = root.enter();

    let detail = gourd_common::codec::decode_recipe(info.text.to_string())
        .unwrap()
        .0;

    Json(detail)

    // let foo = Json(detail);

    // HttpResponse::Ok().json(Json(foo.0)) // <- send response
}
pub async fn amount_parser(info: Query<Info>) -> Result<Json<Vec<Measure>>, AppError> {
    let root = span!(
        tracing::Level::TRACE,
        "amount_parser",
        amount = info.text.to_string().as_str()
    );
    let _enter = root.enter();

    let i = gourd_common::parse_amount(&info.text)?;

    Ok(Json(i))
}
pub async fn convert(r: Json<openapi::models::UnitConversionRequest>) -> impl IntoResponse {
    let root = span!(
        tracing::Level::TRACE,
        "convert",
        item = format!("{:#?}", r.0).to_string().as_str()
    );
    let _enter = root.enter();
    Json(convert_to(r.0))
}

#[derive(Clone, PartialEq, Debug, Serialize, Deserialize)]
pub struct DebugScrapeWrapper {
    compact_scrape_result: CompactRecipe,
    instructions_rich: Vec<ingredient::rich_text::Rich>,
    input: RecipeWrapperInput,
}
#[tracing::instrument(name = "route::debug_scrape")]
pub async fn debug_scrape(info: Query<URLInput>) -> Result<Json<DebugScrapeWrapper>, AppError> {
    let url = info.url.as_str();

    let compact_scrape_result = crate::scraper::scrape_recipe(url).await?;
    let a = expand_recipe(compact_scrape_result.clone()).unwrap();
    let res = RecipeWrapperInput::new(a.clone().0);

    Ok(Json(DebugScrapeWrapper {
        compact_scrape_result,
        instructions_rich: a.1,
        input: res,
    }))
}

#[tracing::instrument(name = "route::scrape")]
pub async fn scrape(info: Query<Info>) -> Result<Json<RecipeWrapperInput>, AppError> {
    let url = info.text.as_str();
    let sc_result = crate::scraper::scrape_recipe(url).await?;

    let res = RecipeWrapperInput::new(expand_recipe(sc_result.clone()).unwrap().0);

    debug!("scraped {}", url.clone());
    Ok(Json(res))
}

#[tracing::instrument(name = "route::expand_compact_to_input")]
pub async fn expand_compact_to_input(
    extract::Json(cr): extract::Json<CompactRecipe>,
) -> impl IntoResponse {
    let res = RecipeWrapperInput::new(expand_recipe(cr.clone()).unwrap().0);
    Json(res)
}
#[tracing::instrument(name = "route::normalize_amount")]
pub async fn normalize_amount(extract::Json(cr): extract::Json<Amount>) -> impl IntoResponse {
    Json(tmp_normalize(cr.clone()))
}
