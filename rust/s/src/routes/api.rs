use actix_web::{web, HttpResponse};
use openapi::models::{Amount, IngredientKind, SectionIngredient};
use sea_orm::entity::*;
use serde::Serialize;
use sqlx::PgPool;

use crate::db;

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

fn si_to_api(r: SI) -> SectionIngredient {
    SectionIngredient {
        kind: if r.ingredient_id.is_some() {
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
                upper_value: None,
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
    dbg!(data);

    let conn = sea_orm::Database::connect("postgresql://gourd:gourd@localhost:5555/food")
        .await
        .unwrap();

    let demo: Vec<(db::recipes::Model, Vec<db::recipe_details::Model>)> =
        db::recipes::Entity::find()
            .find_with_related(db::recipe_details::Entity)
            .all(&conn)
            .await
            .unwrap();

    HttpResponse::Ok().json(actix_web::web::Json(demo)) // <- send response
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
