use async_recursion::async_recursion;
use bigdecimal::BigDecimal;
use futures::StreamExt;
use openapi::models::Amount;
use rand::distributions::Alphanumeric;
use rand::{thread_rng, Rng};
use sea_query::{Expr, Iden, PostgresQueryBuilder, Query};
use serde::{Deserialize, Serialize};
use sqlx::PgPool;
use tracing::{error, info};
sea_query::sea_query_driver_postgres!();
use bigdecimal::FromPrimitive;
use ormx::Insert;
use sea_query_driver_postgres::bind_query;
use std::convert::TryFrom;

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
// pub struct Bar {}
// impl Bar {

#[tracing::instrument]
pub async fn foo(filename: &str, pool: &PgPool) -> anyhow::Result<()> {
    let f = std::fs::File::open(filename)?;
    let d: Vec<Root2> = serde_yaml::from_reader(f)?;
    info!("Read mappings: {}", d.len());

    futures::stream::iter(d.into_iter())
        .for_each(|x| async move {
            let ing = get_ingredient_by_name(pool, x.name).await.unwrap();
            let mut fdc_id = x.fdc_id;

            if fdc_id.is_none() && x.upc.is_some() {
                fdc_id = get_fdc_id_from_upc(pool, x.upc.unwrap())
                    .await
                    .expect("Failed to get fdc_id from upc");
            }

            if fdc_id.is_some() {
                sqlx::query!(
                    "UPDATE ingredients set fdc_id = $1 where id = $2",
                    i32::try_from(fdc_id.unwrap()).unwrap(),
                    ing.id,
                )
                .execute(pool)
                .await
                .expect("ok");
            }

            let child_ids: Vec<String> = futures::future::join_all(
                x.aliases
                    .into_iter()
                    .map(|a| async move { get_ingredient_by_name(pool, a).await.unwrap().id })
                    .collect::<Vec<_>>(),
            )
            .await;
            merge_ingredients(pool, ing.clone().id, child_ids.clone())
                .await
                .expect("merge");

            futures::future::join_all(x.unit_mappings.into_iter().map(|u| async {
                let ins = InsertIngredientUnitMapping {
                    ingredient: ing.clone().id,
                    unit_a: u.a.unit,
                    amount_a: BigDecimal::from_f64(u.a.value).unwrap(),
                    unit_b: u.b.unit,
                    amount_b: BigDecimal::from_f64(u.b.value).unwrap(),
                    source: u.source.unwrap_or_default(),
                };
                match ins
                    .insert(&mut *pool.acquire().await.expect("msfoog"))
                    .await
                {
                    Ok(_) => {}
                    Err(e) => {
                        error!("insert error: {}", e);
                        let mut ok = false;
                        if let sqlx::Error::Database(ref e2) = e {
                            if e2.constraint() == Some("ingredient_unit_attrs_unique") {
                                ok = true
                            }
                        }
                        if !ok {
                            anyhow::bail!("{:?}", e);
                        }
                    }
                };
                Ok(())
            }))
            .await;
        })
        .await;

    Ok(())
}
// }
//
// #[derive(Debug, Clone, Serialize)]
// pub struct Ingredient {
//     id: String,
//     name: String,
//     fdc_id: Option<i32>,
//     parent: Option<String>,
// }

#[derive(Debug, Clone, Serialize, ormx::Table)]
#[ormx(table = "ingredient_units", id = id, insertable, deletable)]
pub struct IngredientUnitMapping {
    #[ormx(default)]
    id: i32,
    ingredient: String,
    unit_a: String,
    amount_a: sqlx::types::BigDecimal,
    unit_b: String,
    amount_b: sqlx::types::BigDecimal,
    source: String,
}
#[derive(Debug, Clone, Serialize)]
// #[ormx(table = "ingredients", id = id, insertable, deletable)]
pub struct Ingredient {
    // #[ormx(get_one = get_by_ingredient_id)]
    id: String,
    // #[ormx(get_one = get_by_ingredient_name)]
    name: String,
    // #[ormx(set)]
    fdc_id: Option<i32>,
    parent: Option<String>,
}

fn id(prefix: &str) -> String {
    // let random_bytes = rand::thread_rng().gen::<[u8; 4]>();
    // let rstr: String = rand::thread_rng()
    //     .sample_iter(&Alphanumeric)
    //     .take(4)
    //     .collect();
    let s: String = thread_rng()
        .sample_iter(&Alphanumeric)
        .take(6)
        .map(char::from)
        .collect();

    return format!("{}_{}", prefix, s);
}

#[tracing::instrument]
#[async_recursion]
pub async fn get_ingredient_by_name(
    pool: &PgPool,
    name: String,
) -> Result<Ingredient, sqlx::Error> {
    let res = sqlx::query_as!(
        Ingredient,
        r#"
        SELECT * FROM ingredients
        WHERE lower(name) = lower($1) LIMIT 1
            "#,
        name
    )
    .fetch_optional(pool)
    .await?;

    if res.is_some() {
        return Ok(res.unwrap());
    }
    sqlx::query!(
        "INSERT INTO ingredients (id, name) VALUES ($1, $2)",
        id("i"),
        name
    )
    .execute(pool)
    .await?;

    return get_ingredient_by_name(pool, name).await;
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
#[derive(Iden)]
enum Ingredients {
    Table,
    Parent,
    Id,
}

// #[derive(Iden)]
// enum IngredientUnitMappings {
//     Table,
//     Ingredient,
//     UnitA,
//     AmountA,
//     UnitB,
//     AmountB,
//     Source,
// }

#[tracing::instrument]
pub async fn merge_ingredients(
    pool: &PgPool,
    parent_id: String,
    children: Vec<String>,
) -> Result<(), sqlx::Error> {
    let mut tx = pool.begin().await?;

    let (sql, values) = Query::update()
        .table(Ingredients::Table)
        .values(vec![(Ingredients::Parent, parent_id.into())])
        .and_where(Expr::col(Ingredients::Id).is_in(children))
        .build(PostgresQueryBuilder);
    bind_query(sqlx::query(&sql), &values)
        .execute(&mut tx)
        .await?;

    tx.commit().await?;
    Ok(())
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
