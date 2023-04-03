use actix_web::{web, HttpResponse};
use anyhow::Result;
use futures::future::join4;
use gourd_common::usda::food_info_from_branded_food_item;
use itertools::Itertools;
use openapi::models::FoodResultByItem;
use openapi::models::{BrandedFoodItem, FoundationFoodItem, SrLegacyFoodItem, SurveyFoodItem};
use serde::Deserialize;
use serde_json::{Deserializer, Value};
use std::{fs::File, io::BufReader, path::Path, time::Instant};

use tracing::info;

use crate::search::{get_client, init_indexes, load, Document, Index};

#[tracing::instrument]
async fn read_and_load_stream<T: Document>(filename: &str, index: Index) -> Result<()> {
    let file = File::open(Path::new("/Users/nicky/dev/gourd/tmp/usda_json/").join(filename))
        .expect(format!("could not open file {}", filename).as_str());
    let reader = BufReader::new(file);
    let stream = Deserializer::from_reader(reader).into_iter::<T>();
    // info!("stream {} is {}", index, stream.sum());
    // let mut data: Vec<T> = vec![];
    for value in &stream.chunks(2000) {
        // info!("got {} chunk of {}", value.sum(), index);
        // let res = value.collect::<Vec<Result<T, serde_json::Error>>>();
        // res.
        let start = Instant::now();
        let res: Vec<T> = value.map(|x| x.unwrap()).collect();
        // load(&res, index).await;
        // let v = value.unwrap();
        // data.push(v);

        // println!("{:?}", value.unwrap());

        info!(
            "loaded {} from {} in {:?}",
            res.len(),
            filename,
            start.elapsed()
        );
    }

    // let data: Vec<T> = read_from_file(filename, &index.get_top_level()).unwrap();
    // load(&data, index).await;
    Ok(())
}
#[tracing::instrument]
fn read_from_file<T: Document>(filename: &str, toplevel: &str) -> Result<Vec<T>> {
    // Open the file
    info!("loading {}", filename);
    let start = Instant::now();
    // let root: Value = serde_json::from_str(data.as_str()).unwrap();
    let file = File::open(Path::new("/Users/nicky/dev/gourd/tmp/usda_json/").join(filename))
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

#[tracing::instrument]
async fn read_and_load<T: Document>(filename: &str, index: Index) -> Result<Vec<T>> {
    let data: Vec<T> = read_from_file(filename, &index.get_top_level()).unwrap();
    load(&data, index).await;
    Ok(data)
}

pub async fn load_json_into_search() {
    // let root: Value = serde_json::from_reader(reader)?;

    init_indexes().await;

    // return;
    if true {
        read_and_load_stream::<SrLegacyFoodItem>(
            // "FoodData_Central_sr_legacy_food_json_2021-10-28.json",
            "srlegacyfoods.ndjson",
            Index::SRLegacyFoods,
        )
        .await
        .unwrap();
        read_and_load_stream::<BrandedFoodItem>(
            // "FoodData_Central_branded_food_json_2022-10-28.json",
            "brandedfoods.ndjson",
            Index::BrandedFoods,
        )
        .await
        .unwrap();
        read_and_load_stream::<SurveyFoodItem>(
            // "FoodData_Central_survey_food_json_2022-10-28.json",
            "surveyfoods.ndjson",
            Index::SurveyFoods,
        )
        .await
        .unwrap();
        read_and_load_stream::<FoundationFoodItem>(
            // "FoodData_Central_foundation_food_json_2022-10-28.json",
            "foundationfoods.ndjson",
            Index::FoundationFoods,
        )
        .await
        .unwrap();
    } else {
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
