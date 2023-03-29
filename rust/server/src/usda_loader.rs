use actix_web::{web, HttpResponse};
use futures::future::join4;
use gourd_common::food_info_from_branded_food_item;
use meilisearch_sdk::task_info::TaskInfo;
use openapi::models::FoodResultByItem;
use openapi::models::{BrandedFoodItem, FoundationFoodItem, SrLegacyFoodItem, SurveyFoodItem};
use serde::Deserialize;
use serde_json::Value;
use std::{
    fs::File,
    io::{BufReader, Read},
    path::Path,
    time::Instant,
};
use tracing::info;

use crate::search::{get_client, load, Document, Index};

#[tracing::instrument]
pub fn read_file(filepath: &str) -> String {
    let file = File::open(filepath).expect("could not open file");
    let mut buffered_reader = BufReader::new(file);
    let mut contents = String::new();

    buffered_reader
        .read_to_string(&mut contents)
        .expect("could not read file into the string");
    contents
}

#[tracing::instrument]
async fn read_and_load<T: Document>(
    filename: &str,
    index: Index,
) -> Result<Vec<T>, Box<dyn std::error::Error + Send + Sync>> {
    let data: Vec<T> = read_from_file(filename, &index.get_top_level()).unwrap();
    load(&data, index).await;
    Ok(data)
}
#[tracing::instrument]
fn read_from_file<T: Document>(
    filename: &str,
    toplevel: &str,
) -> Result<Vec<T>, Box<dyn std::error::Error + Send + Sync>> {
    // Open the file
    info!("loading {}", filename);
    let start = Instant::now();
    // let root: Value = serde_json::from_str(data.as_str()).unwrap();
    let file = File::open(Path::new("/Users/nicky/dev/gourd/data/").join(filename))
        .expect("could not open file");
    let reader = BufReader::new(file);
    let root: Value = serde_json::from_reader(reader)?;

    let items = if let Some(items) = root.get(toplevel) {
        serde_json::from_value::<Vec<T>>(items.clone())?
    } else {
        panic!("No items found in JSON file");
    };
    info!(
        "loaded {} from {} in {:?}",
        items.len(),
        filename,
        start.elapsed()
    );
    Ok(items)
}

pub async fn load_json_into_search() {
    init_indexes().await;

    let res = join4(
        tokio::spawn(read_and_load::<SrLegacyFoodItem>(
            "FoodData_Central_sr_legacy_food_json_2021-10-28.json",
            Index::SRLegacyFoods,
        )),
        tokio::spawn(read_and_load::<BrandedFoodItem>(
            "FoodData_Central_branded_food_json_2022-10-28.json",
            Index::BrandedFoods,
        )),
        tokio::spawn(read_and_load::<SurveyFoodItem>(
            "FoodData_Central_survey_food_json_2022-10-28.json",
            Index::SurveyFoods,
        )),
        tokio::spawn(read_and_load::<FoundationFoodItem>(
            "FoodData_Central_foundation_food_json_2022-10-28.json",
            Index::FoundationFoods,
        )),
    );
    let res = res.await;
    res.0.unwrap().unwrap();
    res.1.unwrap().unwrap();
    res.2.unwrap().unwrap();
    res.3.unwrap().unwrap();
}

#[derive(Deserialize, Debug)]
pub struct URLInput {
    name: String,
}

#[tracing::instrument]
async fn search<T: Document>(
    client: &meilisearch_sdk::Client,
    index: Index,
    limit: usize,
    name: &str,
) -> Vec<T> {
    let results = client
        .index(index)
        .search()
        .with_query(name)
        .with_limit(limit)
        .execute::<T>()
        .await
        .unwrap();
    info!("searched in {}ms", results.processing_time_ms);
    results.hits.into_iter().map(|x| x.result).collect()
}

#[tracing::instrument(name = "route::search_usda")]
pub async fn search_usda(info: web::Query<URLInput>) -> HttpResponse {
    let name = info.name.as_str();
    let branded_food: Vec<BrandedFoodItem> =
        search(&get_client(), Index::BrandedFoods, 5, name).await;

    let wrapped_into_parent = branded_food
        .clone()
        .into_iter()
        .map(food_info_from_branded_food_item)
        .collect();
    let res = FoodResultByItem {
        branded_food,
        foundation_food: vec![],
        survey_food: vec![],
        legacy_food: vec![],
        info: wrapped_into_parent,
    };

    HttpResponse::Ok().json(res)
}

pub async fn init_indexes() {
    let client = get_client();

    for x in [Index::BrandedFoods] {
        if let Some(attr) = x.get_searchable_attributes() {
            let task: TaskInfo = client
                .index(x)
                .set_searchable_attributes(attr)
                .await
                .unwrap();
            info!("task {:?}", task);
        }
        if let Some(attr) = x.get_filterable_attributes() {
            let task: TaskInfo = client
                .index(x)
                .set_filterable_attributes(attr)
                .await
                .unwrap();
            info!("task {:?}", task);
        }
    }
}
