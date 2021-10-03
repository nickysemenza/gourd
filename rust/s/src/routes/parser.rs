use actix_web::{web, HttpResponse};
use gourd_common::{convert_to, pan};
use openapi::models::{
    Amount, IngredientKind, RecipeDetail, RecipeSection, RecipeWrapper, SectionIngredient,
    SectionInstruction,
};
use pyo3::{types::PyModule, PyAny, Python};
use serde::{Deserialize, Serialize};
use sqlx::PgPool;
use tracing::{debug, error, span};

fn si_to_api(r: SI) -> SectionIngredient {
    SectionIngredient {
        kind: if r.ingredient.is_some() {
            IngredientKind::Ingredient
        } else {
            IngredientKind::Recipe
        },
        // section: r.section,
        id: r.id,
        // sort: r.sort,
        ingredient: None,
        recipe: None,
        amounts: r
            .amounts
            .iter()
            .map(|a| Amount {
                unit: a.unit.clone(),
                value: a.value,
                source: Some("todo".to_string()),
            })
            .collect(),
        adjective: r.adjective,
        optional: r.optional,
        original: r.original,
        substitutes: None,
        // substitutes_for: r.substitutes_for,
    }
}
/// This handler uses json extractor
pub async fn index(pool: web::Data<PgPool>) -> HttpResponse {
    let rows = get_test(&pool).await.unwrap();
    let data: Vec<SectionIngredient> = rows.into_iter().map(|r| si_to_api(r)).collect();
    // dbg!(a);

    HttpResponse::Ok().json(actix_web::web::Json(data)) // <- send response
}

#[derive(Deserialize, Debug)]
pub struct Info {
    text: String,
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
    if i.is_err() {
        let msg = format!("{:?}", i.err().unwrap());
        error!("parse error {}", msg);
        return HttpResponse::BadRequest().body(msg);
    }
    let foo = web::Json(i.unwrap());

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

#[derive(Debug, Clone, Serialize)]
pub struct SI {
    section: String,
    id: String,
    sort: Option<i32>,
    ingredient: Option<String>,
    recipe: Option<String>,
    amounts: sqlx::types::Json<Vec<Amount>>,
    adjective: Option<String>,
    optional: Option<bool>,
    original: Option<String>,
    substitutes_for: Option<String>,
}

#[tracing::instrument]
pub async fn get_test(pool: &PgPool) -> Result<Vec<SI>, sqlx::Error> {
    let res = sqlx::query_as!(
        SI,
        r#"
    select section, id, sort, ingredient, recipe, amounts as "amounts: sqlx::types::Json<Vec<Amount>>",
     adjective, optional, original, substitutes_for from recipe_section_ingredients;
            "#,
    )
    .fetch_all(pool)
    .await?;

    // dbg!(res);
    // let res2 = res.unwrap();
    Ok(res)
}
pub async fn pans() -> HttpResponse {
    let p = pan::inventory();

    HttpResponse::Ok().json(actix_web::web::Json(p)) // <- send response
}
#[tracing::instrument(name = "route::scrape")]
pub async fn scrape(info: web::Query<Info>) -> HttpResponse {
    let mut sc_result: (Vec<String>, String, String) = (vec![], "".to_string(), "".to_string());
    Python::with_gil(|py| {
        let syspath: &PyAny = py.import("sys").unwrap().get("path").unwrap();

        dbg!(syspath);
        let activators = PyModule::from_code(
            py,
            r#"
from recipe_scrapers import scrape_me            
def sc(x,y):
    res = scrape_me(x,wild_mode=y)
    return res.ingredients(), res.instructions(), res.title()
            "#,
            "recipe_scrape.py",
            "recipe_scrape",
        )
        .unwrap();

        dbg!(activators);
        sc_result = activators
            .getattr("sc")
            .unwrap()
            .call((info.text.clone(), true), None)
            .unwrap()
            .extract()
            .unwrap();
    });
    let sections = vec![RecipeSection::new(
        "".to_string(),
        sc_result
            .1
            .split('\n')
            .map(|x| SectionInstruction::new("".to_string(), x.to_string()))
            .collect(),
        sc_result
            .0
            .iter()
            .map(|x| {
                let _x = 1;
                gourd_common::parse_ingredient(&x).unwrap_or(SectionIngredient::new(
                    "".to_string(),
                    IngredientKind::Ingredient,
                    vec![],
                ))
            })
            .collect(),
    )];
    let detail = RecipeDetail::new("".to_string(), sections, sc_result.2, 0, "".to_string());
    let res = RecipeWrapper::new("".to_string(), detail);

    debug!("scraped {}", info.text.clone());
    HttpResponse::Ok().json(actix_web::web::Json(res)) // <- send response
}

#[cfg(test)]
mod tests {
    use super::*;
    use actix_web::Error;
    use actix_web::{test, web, App};

    #[actix_rt::test]
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
            r##"{"id":"","kind":"ingredient","ingredient":{"name":"","ingredient":{"id":"","name":"flour"},"recipes":[],"children":[],"unit_mappings":[]},"amounts":[{"unit":"cup","value":1.0},{"unit":"g","value":120.0}],"adjective":"lightly sifted","original":"1 cup (120 grams) flour, lightly sifted"}"##
        );

        Ok(())
    }
}
