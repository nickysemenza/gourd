use actix_web::{web, HttpResponse};
use anyhow::{Context, Result};
use gourd_common::usda::{
    branded_food_into_wrapper, foundation_food_into_wrapper, sr_legacy_food_into_wrapper,
    survey_food_into_wrapper,
};
use indicatif::{ProgressBar, ProgressState, ProgressStyle};
use itertools::Itertools;
use openapi::models::TempFood;
use openapi::models::{BrandedFoodItem, FoundationFoodItem, SrLegacyFoodItem, SurveyFoodItem};
use serde::Deserialize;
use serde_json::{Deserializer, Value};
use std::fmt::Write;
use std::{fs::File, io::BufReader, path::Path, time::Instant};
use tracing::info;

use crate::search::{get_client, init_indexes, load, Document, Index};

#[tracing::instrument]
async fn read_and_load_stream<T: Document>(filename: &str, index: Index) -> Result<()> {
    let file = File::open(Path::new("/Users/nicky/dev/gourd/tmp/usda_json/").join(filename))
        .expect(format!("could not open file {}", filename).as_str());
    let pb = ProgressBar::new(file.metadata()?.len());
    pb.set_style(ProgressStyle::with_template("{spinner:.green} [{elapsed_precise}] [{wide_bar:.cyan/blue}] {bytes}/{total_bytes} ({eta})")
    .unwrap()
    .with_key("eta", |state: &ProgressState, w: &mut dyn Write| write!(w, "{:.1}s", state.eta().as_secs_f64()).unwrap())
    .progress_chars("#>-"));

    let reader = BufReader::new(pb.wrap_read(file));
    let stream = Deserializer::from_reader(reader).into_iter::<T>();
    // info!("stream {} is {}", index, stream.sum());
    // let mut data: Vec<T> = vec![];
    for value in &stream.chunks(10000) {
        // info!("got {} chunk of {}", value.sum(), index);
        // let res = value.collect::<Vec<Result<T, serde_json::Error>>>();
        // res.
        let start = Instant::now();
        let res: Vec<T> = value.map(|x| x.unwrap()).collect();
        load(&res, index).await;
        // let v = value.unwrap();
        // data.push(v);

        // println!("{:?}", value.unwrap());

        info!(
            "[stream] loaded {} from {} in {:?}",
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
        .with_context(|| format!("could not open file {}", filename))?;
    let pb = ProgressBar::new(file.metadata()?.len());
    pb.set_style(ProgressStyle::with_template("{spinner:.green} [{elapsed_precise}] [{wide_bar:.cyan/blue}] {bytes}/{total_bytes} ({eta})")
    .unwrap()
    .with_key("eta", |state: &ProgressState, w: &mut dyn Write| write!(w, "{:.1}s", state.eta().as_secs_f64()).unwrap())
    .progress_chars("#>-"));

    let reader = BufReader::new(pb.wrap_read(file));
    let root: Value = serde_json::from_reader(reader)?;

    let items = if let Some(items) = root.get(toplevel) {
        serde_json::from_value::<Vec<T>>(items.clone())?
    } else {
        panic!("No items found in JSON file");
    };
    info!(
        "[full] loaded {} from {} in {:?}",
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
    if false {
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
        read_and_load::<SrLegacyFoodItem>(
            "FoodData_Central_sr_legacy_food_json_2021-10-28.json",
            Index::SRLegacyFoods,
        )
        .await
        .unwrap();
        read_and_load::<BrandedFoodItem>(
            "FoodData_Central_branded_food_json_2022-04-28.json",
            Index::BrandedFoods,
        )
        .await
        .unwrap();
        read_and_load::<SurveyFoodItem>(
            "FoodData_Central_survey_food_json_2022-10-28.json",
            Index::SurveyFoods,
        )
        .await
        .unwrap();
        read_and_load::<FoundationFoodItem>(
            "FoodData_Central_foundation_food_json_2022-04-28.json",
            Index::FoundationFoods,
        )
        .await
        .unwrap();
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

async fn get_document<T: Document>(
    client: &meilisearch_sdk::Client,
    index: Index,
    name: &str,
) -> Result<Option<T>> {
    match client.index(index).get_document::<T>(name).await {
        Ok(d) => Ok(Some(d)),
        Err(e) => match e {
            meilisearch_sdk::errors::Error::Meilisearch(m) => match m.error_code {
                meilisearch_sdk::errors::ErrorCode::DocumentNotFound => Ok(None),
                _ => Err(m.into()),
            },

            _ => Err(e.into()),
        },
    }
}

#[tracing::instrument(name = "route::search_usda")]
pub async fn get_usda(info: web::Query<URLInput>) -> HttpResponse {
    let id = info.name.as_str();
    let item = get_document::<BrandedFoodItem>(&get_client(), Index::BrandedFoods, id)
        .await
        .unwrap();
    if let Some(item) = item {
        return HttpResponse::Ok().json(branded_food_into_wrapper(item));
    };

    let item = get_document::<SrLegacyFoodItem>(&get_client(), Index::SRLegacyFoods, id)
        .await
        .unwrap();
    if let Some(item) = item {
        return HttpResponse::Ok().json(sr_legacy_food_into_wrapper(item));
    };

    let item = get_document::<SurveyFoodItem>(&get_client(), Index::SurveyFoods, id)
        .await
        .unwrap();
    if let Some(item) = item {
        return HttpResponse::Ok().json(survey_food_into_wrapper(item));
    };

    let item = get_document::<FoundationFoodItem>(&get_client(), Index::FoundationFoods, id)
        .await
        .unwrap();
    if let Some(item) = item {
        return HttpResponse::Ok().json(foundation_food_into_wrapper(item));
    };

    HttpResponse::NotFound().json(format!("{} not found", id).to_string())
}

#[tracing::instrument(name = "route::search_usda")]
pub async fn search_usda(info: web::Query<URLInput>) -> HttpResponse {
    let name = info.name.as_str();
    let branded_food: Vec<BrandedFoodItem> =
        search(&get_client(), Index::BrandedFoods, 5, name).await;
    let legacy_food: Vec<SrLegacyFoodItem> =
        search(&get_client(), Index::SRLegacyFoods, 5, name).await;

    let mut a: Vec<TempFood> = branded_food
        .into_iter()
        .map(branded_food_into_wrapper)
        .collect();
    let mut b = legacy_food
        .into_iter()
        .map(sr_legacy_food_into_wrapper)
        .collect();
    a.append(&mut b);

    HttpResponse::Ok().json(a)
}

// #[tracing::instrument(name = "route::search_usda")]
// pub async fn search_usda(info: web::Query<URLInput>) -> HttpResponse {
