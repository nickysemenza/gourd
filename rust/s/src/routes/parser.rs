use actix_web::{web, HttpResponse};
use openapi::models::{IngredientKind, SectionIngredient};
use opentelemetry::{global, trace::Tracer};
use opentelemetry::{trace::get_active_span, KeyValue};
use serde::{Deserialize, Serialize};
use sqlx::{types::BigDecimal, PgPool};

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
        amount: match r.amount {
            Some(f) => bigdecimal::ToPrimitive::to_f64(&f),
            None => None,
        },
        unit: r.unit,
        grams: match r.grams {
            Some(f) => bigdecimal::ToPrimitive::to_f64(&f).unwrap_or_default(),
            None => 0.0,
        },
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

#[derive(Deserialize)]
pub struct Info {
    text: String,
}

pub async fn parser(info: web::Query<Info>) -> HttpResponse {
    global::tracer("my-component").start("parser");

    get_active_span(|span| {
        span.add_event(
            "parse".to_string(),
            vec![KeyValue::new("ingredient", info.text.to_string())],
        );
    });

    let i = gourd_common::parse_ingredient(&info.text);
    if i.is_err() {
        return HttpResponse::BadRequest().finish();
    }
    let foo = web::Json(i.unwrap());

    HttpResponse::Ok().json(actix_web::web::Json(foo.0)) // <- send response
}
pub async fn amount_parser(info: web::Query<Info>) -> HttpResponse {
    global::tracer("my-component").start("amount_parser");

    get_active_span(|span| {
        span.add_event(
            "amount_parse".to_string(),
            vec![KeyValue::new("line", info.text.to_string())],
        );
    });

    let i = gourd_common::parse_amount(&info.text);
    if i.is_err() {
        return HttpResponse::BadRequest().finish();
    }
    let foo = web::Json(i.unwrap());

    HttpResponse::Ok().json(actix_web::web::Json(foo.0)) // <- send response
}

#[derive(Debug, Clone, Serialize)]
pub struct SI {
    section: String,
    id: String,
    sort: Option<i32>,
    ingredient: Option<String>,
    recipe: Option<String>,
    amount: Option<BigDecimal>,
    unit: Option<String>,
    grams: Option<BigDecimal>,
    adjective: Option<String>,
    optional: Option<bool>,
    original: Option<String>,
    substitutes_for: Option<String>,
}

pub async fn get_test(pool: &PgPool) -> Result<Vec<SI>, sqlx::Error> {
    let res = sqlx::query_as!(
        SI,
        r#"
    select * from recipe_section_ingredients;
            "#,
    )
    .fetch_all(pool)
    .await?;

    // dbg!(res);
    // let res2 = res.unwrap();
    Ok(res)
}

#[cfg(test)]
mod tests {
    use super::*;
    use actix_web::{dev::Service, Error};
    use actix_web::{http, test, web, App};

    #[actix_rt::test]
    async fn test_parse() -> Result<(), Error> {
        let app = test::init_service(
            App::new().service(web::resource("/parse").route(web::get().to(parser))),
        )
        .await;

        let req = test::TestRequest::get()
            .uri("/parse?text=1%20cup%20(120%20grams)%20flour,%20lighty%20sifted")
            .param("text", "1 cup flour")
            .to_request();
        let resp = app.call(req).await.unwrap();

        assert_eq!(resp.status(), http::StatusCode::OK);

        let response_body = match resp.response().body().as_ref() {
            Some(actix_web::body::Body::Bytes(bytes)) => bytes,
            _ => panic!("Response error"),
        };

        assert_eq!(
            response_body,
            r##"{"id":"","kind":"ingredient","ingredient":{"id":"","name":"flour"},"grams":120.0,"amount":1.0,"unit":"cup","adjective":"lighty sifted"}"##
        );

        Ok(())
    }
}
