use anyhow::{Context, Result};
use axum::extract::Query;
use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use axum::Json;
use gourd_common::usda::IntoFoodWrapper;
use indicatif::{ProgressBar, ProgressState, ProgressStyle};
use itertools::Itertools;
use openapi::models::TempFood;
use openapi::models::{BrandedFoodItem, FoundationFoodItem, SrLegacyFoodItem, SurveyFoodItem};
use serde::Deserialize;
use serde_json::{Deserializer, Value};
use std::fmt::Write;
use std::{fs::File, io::BufReader, path::Path, time::Instant};
use tracing::info;

use crate::error::AppError;
use crate::search::{Document, Index, Searcher};

#[tracing::instrument]
async fn read_and_load_stream<T: Document>(filename: &str, index: Index) -> Result<()> {
    let file = File::open(Path::new("/Users/nicky/dev/gourd/tmp/usda_json/").join(filename))
        .expect(format!("could not open file {}", filename).as_str());
    let pb = ProgressBar::new(file.metadata()?.len());
    pb.set_style(ProgressStyle::with_template("{spinner:.green} [{elapsed_precise}] [{wide_bar:.cyan/blue}] {bytes}/{total_bytes} ({eta})")
    ?
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
        let s = Searcher::new();
        s.load(&res, index).await;
        // let v = value?;
        // data.push(v);

        // println!("{:?}", value?);

        info!(
            "[stream] loaded {} from {} in {:?}",
            res.len(),
            filename,
            start.elapsed()
        );
    }

    // let data: Vec<T> = read_from_file(filename, &index.get_top_level())?;
    // load(&data, index).await;
    Ok(())
}
#[tracing::instrument]
fn read_from_file<T: Document>(filename: &str, toplevel: &str) -> Result<Vec<T>> {
    // Open the file
    info!("loading {}", filename);
    let start = Instant::now();
    // let root: Value = serde_json::from_str(data.as_str())?;
    let file = File::open(Path::new("/Users/nicky/dev/gourd/tmp/usda_json/").join(filename))
        .with_context(|| format!("could not open file {}", filename))?;
    let pb = ProgressBar::new(file.metadata()?.len());
    pb.set_style(ProgressStyle::with_template("{spinner:.green} [{elapsed_precise}] [{wide_bar:.cyan/blue}] {bytes}/{total_bytes} ({eta})")
    ?
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
    let data: Vec<T> = read_from_file(filename, &index.get_top_level_json_key())?;
    Searcher::new().load(&data, index).await;
    Ok(data)
}

pub async fn load_json_into_search() -> Result<()> {
    Searcher::new().init_indexes().await;
    // return;
    if false {
        read_and_load_stream::<SrLegacyFoodItem>("srlegacyfoods.ndjson", Index::SRLegacyFoods)
            .await?;
        read_and_load_stream::<BrandedFoodItem>("brandedfoods.ndjson", Index::BrandedFoods).await?;
        read_and_load_stream::<SurveyFoodItem>("surveyfoods.ndjson", Index::SurveyFoods).await?;
        read_and_load_stream::<FoundationFoodItem>(
            "foundationfoods.ndjson",
            Index::FoundationFoods,
        )
        .await?;
    } else {
        read_and_load::<SrLegacyFoodItem>(
            "FoodData_Central_sr_legacy_food_json_2021-10-28.json",
            Index::SRLegacyFoods,
        )
        .await?;
        read_and_load::<BrandedFoodItem>(
            "FoodData_Central_branded_food_json_2022-04-28.json",
            Index::BrandedFoods,
        )
        .await?;
        read_and_load::<SurveyFoodItem>(
            "FoodData_Central_survey_food_json_2022-10-28.json",
            Index::SurveyFoods,
        )
        .await?;
        read_and_load::<FoundationFoodItem>(
            "FoodData_Central_foundation_food_json_2022-04-28.json",
            Index::FoundationFoods,
        )
        .await?;
    }
    Ok(())
}

#[derive(Deserialize, Debug)]
pub struct URLInput {
    name: String,
}
async fn get_usda_by_id(id: &str) -> Result<Option<Result<TempFood>>> {
    let s = Searcher::new();

    let item = s
        .get_document::<BrandedFoodItem>(Index::BrandedFoods, id)
        .await?
        .map(|x| x.into_wrapper());
    if item.is_some() {
        return Ok(item);
    }
    let item = s
        .get_document::<SrLegacyFoodItem>(Index::SRLegacyFoods, id)
        .await?
        .map(|x| x.into_wrapper());
    if item.is_some() {
        return Ok(item);
    }
    let item = s
        .get_document::<SurveyFoodItem>(Index::SurveyFoods, id)
        .await?
        .map(|x| x.into_wrapper());
    if item.is_some() {
        return Ok(item);
    }
    let item = s
        .get_document::<FoundationFoodItem>(Index::FoundationFoods, id)
        .await?
        .map(|x| x.into_wrapper());
    if item.is_some() {
        return Ok(item);
    }
    Ok(None)
}

pub async fn get_usda(info: Query<URLInput>) -> Response {
    let id = info.name.as_str();

    match get_usda_by_id(id).await.unwrap() {
        Some(item) => match item {
            Ok(item) => Json(item).into_response(),
            Err(e) => AppError::from(e).into_response(),
        },
        None => (StatusCode::NOT_FOUND, "Not found").into_response(),
    }
}

async fn search_usda_by_name(name: &str) -> Result<Vec<TempFood>> {
    let s = Searcher::new();

    let mut a: Vec<TempFood> = s
        .search::<BrandedFoodItem>(Index::BrandedFoods, 5, name)
        .await
        .context("branded food")?
        .into_iter()
        .map(|x| x.into_wrapper())
        .collect::<Result<Vec<TempFood>>>()
        .context("branded food")?;
    let mut b = s
        .search::<SrLegacyFoodItem>(Index::SRLegacyFoods, 5, name)
        .await?
        .into_iter()
        .map(|x| x.into_wrapper())
        .collect::<Result<Vec<TempFood>>>()?;
    let mut c = s
        .search::<SurveyFoodItem>(Index::SurveyFoods, 5, name)
        .await?
        .into_iter()
        .map(|x| x.into_wrapper())
        .collect::<Result<Vec<TempFood>>>()?;
    let mut d = s
        .search::<FoundationFoodItem>(Index::FoundationFoods, 5, name)
        .await?
        .into_iter()
        .map(|x| x.into_wrapper())
        .collect::<Result<Vec<TempFood>>>()?;
    a.append(&mut b);
    a.append(&mut c);
    a.append(&mut d);
    Ok(a)
}
#[tracing::instrument(name = "route::search_usda")]
pub async fn search_usda(info: Query<URLInput>) -> Result<Json<Vec<TempFood>>, AppError> {
    Ok(Json(search_usda_by_name(info.name.as_str()).await?))
}
