use actix_web::{web, HttpResponse};
use gourd_common::{convert_to, pan, unit::Measure};
use openapi::models::{
    Amount, IngredientKind, RecipeDetailInput, RecipeSectionInput, RecipeWrapperInput,
    SectionIngredientInput, SectionInstructionInput,
};
use serde::Deserialize;
use tracing::{debug, span};

use crate::scraper;

#[derive(Deserialize, Debug)]
pub struct Info {
    text: String,
}

#[derive(Deserialize, Debug)]
pub struct URLInput {
    url: String,
}

pub async fn parser(info: web::Query<Info>) -> HttpResponse {
    let root = span!(
        tracing::Level::TRACE,
        "parser",
        ingredient = info.text.to_string().as_str()
    );
    let _enter = root.enter();

    let i = gourd_common::parse_ingredient(&info.text);
    if i.is_err() {
        return HttpResponse::BadRequest().finish();
    }
    let foo = web::Json(i.unwrap());

    HttpResponse::Ok().json(actix_web::web::Json(foo.0)) // <- send response
}
pub async fn decode_recipe(info: web::Query<Info>) -> HttpResponse {
    let root = span!(tracing::Level::TRACE, "decode_recipe",);
    let _enter = root.enter();

    let detail = gourd_common::codec::decode_recipe(info.text.to_string());

    let foo = web::Json(detail);

    HttpResponse::Ok().json(web::Json(foo.0)) // <- send response
}
pub async fn amount_parser(info: web::Query<Info>) -> HttpResponse {
    let root = span!(
        tracing::Level::TRACE,
        "amount_parser",
        amount = info.text.to_string().as_str()
    );
    let _enter = root.enter();

    let i = gourd_common::parse_amount(&info.text);

    let foo = web::Json(i);

    HttpResponse::Ok().json(web::Json(foo.0)) // <- send response
}
pub async fn convert(r: web::Json<openapi::models::UnitConversionRequest>) -> HttpResponse {
    let root = span!(
        tracing::Level::TRACE,
        "convert",
        item = format!("{:#?}", r.0).to_string().as_str()
    );
    let _enter = root.enter();
    HttpResponse::Ok().json(convert_to(r.0))
}

pub async fn pans() -> HttpResponse {
    let p = pan::inventory();

    HttpResponse::Ok().json(actix_web::web::Json(p)) // <- send response
}

fn foo(instructions: String) -> Measure {
    let rich_text_tokens =
        gourd_common::ingredient::rich_text::parse(instructions.as_str()).unwrap_or_default();

    let mut amts = Vec::<gourd_common::unit::Measure>::new();
    let mut total_time =
        gourd_common::unit::Measure::parse(gourd_common::ingredient::Amount::new("second", 0.0));

    for token in rich_text_tokens.clone().into_iter() {
        match token {
            gourd_common::ingredient::rich_text::Chunk::Amount(amt) => {
                for a in amt.into_iter() {
                    let m = gourd_common::amount_to_measure2(a);
                    amts.push(m.clone());
                    if m.kind().unwrap() == gourd_common::ingredient::unit::kind::MeasureKind::Time
                    {
                        total_time = total_time.add(m).unwrap();
                    }
                }
            }
            _ => {}
        }
    }
    total_time
}
#[tracing::instrument(name = "route::debug_scrape")]
pub async fn debug_scrape(info: web::Query<URLInput>) -> HttpResponse {
    let url = info.url.as_str();

    let sc_result = scraper::scrape_recipe(url);
    let res = scrape_result_to_recipe(sc_result.clone());
    let rich_text_tokens =
        gourd_common::ingredient::rich_text::parse(sc_result.instructions.clone().as_str())
            .unwrap_or_default();

    // let mut amts = Vec::<gourd_common::unit::Measure>::new();
    // let mut total_time =
    //     gourd_common::unit::Measure::parse(gourd_common::ingredient::Amount::new("second", 0.0));

    // for token in rich_text_tokens.clone().into_iter() {
    //     match token {
    //         gourd_common::ingredient::rich_text::Chunk::Amount(amt) => {
    //             for a in amt.into_iter() {
    //                 let m = gourd_common::amount_to_measure2(a);
    //                 amts.push(m.clone());
    //                 if m.kind().unwrap() == gourd_common::unit::MeasureKind::Time {
    //                     total_time = total_time.add(m).unwrap();
    //                 }
    //             }
    //         }
    //         _ => {}
    //     }
    // }
    let total_time = foo(sc_result.clone().instructions);

    HttpResponse::Ok().json((
        sc_result,
        res,
        rich_text_tokens,
        // amts,
        total_time.as_bare().unwrap(),
    ))
}
pub fn scrape_result_to_recipe(sc_result: scraper::ScrapeResult) -> RecipeWrapperInput {
    let total_time_seconds = foo(sc_result.clone().instructions).as_raw();
    let sections = vec![RecipeSectionInput {
        duration: Some(Box::new(Amount::new(
            total_time_seconds.unit,
            total_time_seconds.value.into(),
        ))),
        instructions: sc_result
            .instructions
            .split('\n')
            .map(|x| SectionInstructionInput::new(x.to_string()))
            .collect(),
        ingredients: sc_result
            .ingredients
            .iter()
            .map(|x| {
                let _x = 1;
                gourd_common::parse_ingredient(&x).unwrap_or(SectionIngredientInput::new(
                    IngredientKind::Ingredient,
                    vec![],
                ))
            })
            .collect(),
    }];

    let detail = RecipeDetailInput::new(sections, sc_result.title, 0, "".to_string());
    return RecipeWrapperInput::new(detail);
}
#[tracing::instrument(name = "route::scrape")]
pub async fn scrape(info: web::Query<Info>) -> HttpResponse {
    let url = info.text.as_str();
    let sc_result = scraper::scrape_recipe(url);
    let res = scrape_result_to_recipe(sc_result);

    debug!("scraped {}", url.clone());
    HttpResponse::Ok().json(actix_web::web::Json(res)) // <- send response
}

#[cfg(test)]
mod tests {
    use super::*;
    use actix_web::Error;
    use actix_web::{test, web, App};

    #[actix_web::test]
    async fn test_parse() -> Result<(), Error> {
        let mut app = test::init_service(
            App::new().service(web::resource("/parse").route(web::get().to(parser))),
        )
        .await;

        let req = test::TestRequest::get()
            .uri("/parse?text=1%20cup%20(120%20grams)%20flour,%20lightly%20sifted")
            .param("text", "1 cup flour")
            .to_request();
        // let resp = app.call(req).await.unwrap();

        let resp = test::call_service(&mut app, req).await;
        assert!(resp.status().is_success());

        // assert_eq!(resp.status(), http::StatusCode::OK);

        let response_body: String = match resp.into_body() {
            actix_web::body::AnyBody::Bytes(bytes) => {
                std::str::from_utf8(&bytes).unwrap().to_string()
            }
            _ => panic!("Response error"),
        };

        assert_eq!(
            response_body,
            r##"{"name":"flour","kind":"ingredient","amounts":[{"unit":"cup","value":1.0},{"unit":"g","value":120.0}],"adjective":"lightly sifted","original":"1 cup (120 grams) flour, lightly sifted"}"##
        );

        Ok(())
    }
}
