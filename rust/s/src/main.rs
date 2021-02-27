use actix_web::{error, middleware, web, App, Error, HttpRequest, HttpResponse, HttpServer};
use actix_web_opentelemetry::RequestTracing;
use futures::StreamExt;
use openapi::models::{section_ingredient::Kind, Ingredient, SectionIngredient};
use opentelemetry::{
    global,
    sdk::{
        propagation::TraceContextPropagator,
        trace::{self, Sampler},
    },
    trace::Tracer,
};
use opentelemetry::{trace::get_active_span, KeyValue};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
struct MyObj {
    name: String,
    number: i32,
}

/// This handler uses json extractor
async fn index(item: web::Json<MyObj>) -> HttpResponse {
    println!("model: {:?}", &item);

    let foo = web::Json(SectionIngredient::new(
        "".to_string(),
        Kind::Ingredient,
        0.0,
    ));

    HttpResponse::Ok().json(foo.0) // <- send response
}

#[derive(Deserialize)]
struct Info {
    text: String,
}

async fn parser(info: web::Query<Info>) -> HttpResponse {
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

    HttpResponse::Ok().json(foo.0) // <- send response
}

/// This handler uses json extractor with limit
async fn extract_item(item: web::Json<MyObj>, req: HttpRequest) -> HttpResponse {
    println!("request: {:?}", req);
    println!("model: {:?}", item);

    HttpResponse::Ok().json(item.0) // <- send json response
}

const MAX_SIZE: usize = 262_144; // max payload size is 256k

/// This handler manually load request payload and parse json object
async fn index_manual(mut payload: web::Payload) -> Result<HttpResponse, Error> {
    // payload is a stream of Bytes objects
    let mut body = web::BytesMut::new();
    while let Some(chunk) = payload.next().await {
        let chunk = chunk?;
        // limit max size of in-memory payload
        if (body.len() + chunk.len()) > MAX_SIZE {
            return Err(error::ErrorBadRequest("overflow"));
        }
        body.extend_from_slice(&chunk);
    }

    // body is loaded, now we can deserialize serde-json
    let obj = serde_json::from_slice::<MyObj>(&body)?;
    Ok(HttpResponse::Ok().json(obj)) // <- send response
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    std::env::set_var("RUST_LOG", "actix_web=info");
    env_logger::init();

    // let tracer = opentelemetry::sdk::export::trace::stdout::new_pipeline()
    //     .with_trace_config(trace::config().with_default_sampler(Sampler::AlwaysOn))
    //     .install();
    // println!("tracer: {:?}", tracer);

    // global::set_text_map_propagator(opentelemetry_jaeger::Propagator::new());
    // not jaeger, this one: https://github.com/openebs/Mayastor/blob/master/control-plane/rest/service/src/main.rs#L64
    global::set_text_map_propagator(TraceContextPropagator::new());
    let (tracer, _uninstall) = opentelemetry_jaeger::new_pipeline()
        .with_service_name("actix_server")
        .with_trace_config(trace::config().with_default_sampler(Sampler::AlwaysOn))
        .install()
        .expect("pipeline install error");
    println!("tracer: {:?}", tracer);

    HttpServer::new(|| {
        App::new()
            .wrap(RequestTracing::new())
            // enable logger
            .wrap(middleware::Logger::default())
            .data(web::JsonConfig::default().limit(4096)) // <- limit size of the payload (global configuration)
            .service(web::resource("/extractor").route(web::post().to(index)))
            .service(web::resource("/parse").route(web::get().to(parser)))
            .service(
                web::resource("/extractor2")
                    .data(web::JsonConfig::default().limit(1024)) // <- limit size of the payload (resource level)
                    .route(web::post().to(extract_item)),
            )
            .service(web::resource("/manual").route(web::post().to(index_manual)))
            .service(web::resource("/").route(web::post().to(index)))
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}

// fn is_gram(a: &ingredient::Amount) -> bool {
//     a.unit == "g" || a.unit == "grams"
// }
// fn parse_ingredient(s: &str) -> Result<SectionIngredient, String> {
//     global::tracer("my-component").start("parse_ingredient");

//     get_active_span(|span| {
//         span.add_event(
//             "parse".to_string(),
//             vec![KeyValue::new("ingredient", s.to_string())],
//         );
//     });

//     let i = ingredient::from_str(s, true)?;

//     let (unit, amount) = match i.amounts.iter().find(|&x| !is_gram(x)) {
//         Some(i) => (Some(i.unit.clone()), Some(i.value as f64)),
//         None => (None, None),
//     };

//     Ok(SectionIngredient {
//         unit,
//         amount,
//         adjective: i.modifier,
//         ingredient: Some(Ingredient::new("".to_string(), i.name)),
//         ..SectionIngredient::new(
//             "".to_string(),
//             Kind::Ingredient,
//             match i.amounts.iter().find(|&x| is_gram(x)) {
//                 Some(g) => g.value.into(),
//                 None => 0.0,
//             },
//         )
//     })
// }

#[cfg(test)]
mod tests {
    use super::*;
    use actix_web::dev::Service;
    use actix_web::{http, test, web, App};

    #[actix_rt::test]
    async fn test_index() -> Result<(), Error> {
        let mut app =
            test::init_service(App::new().service(web::resource("/").route(web::post().to(index))))
                .await;

        let req = test::TestRequest::post()
            .uri("/")
            .set_json(&MyObj {
                name: "my-name".to_owned(),
                number: 43,
            })
            .to_request();
        let resp = app.call(req).await.unwrap();

        assert_eq!(resp.status(), http::StatusCode::OK);

        let response_body = match resp.response().body().as_ref() {
            Some(actix_web::body::Body::Bytes(bytes)) => bytes,
            _ => panic!("Response error"),
        };

        assert_eq!(
            response_body,
            r##"{"id":"","kind":"ingredient","grams":0.0}"##
        );

        Ok(())
    }
    #[actix_rt::test]
    async fn test_parse() -> Result<(), Error> {
        let mut app = test::init_service(
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
