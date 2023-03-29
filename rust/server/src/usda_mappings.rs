use actix_web::{web, HttpResponse};
use gourd_common::make_unit_mappings;
use openapi::models::FoodWrapper;

#[tracing::instrument(name = "route::unit_mappings_from_food")]
pub async fn unit_mappings_from_food(cr: web::Json<FoodWrapper>) -> HttpResponse {
    HttpResponse::Ok().json(make_unit_mappings(cr.0))
}
