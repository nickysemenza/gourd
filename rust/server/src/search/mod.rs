mod index;

use actix_web::{web, HttpResponse};
use meilisearch_sdk::task_info::TaskInfo;
use openapi::models::RecipeDetail;
use serde::{de::DeserializeOwned, Serialize};
use strum::IntoEnumIterator;
use tracing::info;

pub use self::index::Index;

pub fn get_client() -> meilisearch_sdk::Client {
    let meilisearch_url = option_env!("MEILISEARCH_URL").unwrap_or("http://localhost:7700");
    let meilisearch_api_key = option_env!("MEILISEARCH_API_KEY").unwrap_or("FOO");
    meilisearch_sdk::Client::new(meilisearch_url, meilisearch_api_key)
}

pub trait Document: Clone + Serialize + DeserializeOwned + Send + Sync + 'static {}
impl<T: Clone + Serialize + DeserializeOwned + Send + Sync + 'static> Document for T {}

#[tracing::instrument(skip(data))]
pub async fn load<T: Document>(data: &Vec<T>, index: Index) {
    // let chunks: Vec<&[FoundationFoodItem]> = legacy.chunks(10).collect();
    let chunks: Vec<Vec<T>> = data.chunks(2000).map(|x| x.to_vec()).collect();

    let tasks: Vec<_> = chunks
        .iter()
        .map(|v| {
            let client = get_client();
            let x = v.clone();
            let i = index.clone();
            tokio::spawn(async move { client.index(i).add_documents(&x, None).await.unwrap() })
        })
        .collect();
    let res = futures::future::join_all(tasks).await;
    // info!("finished loading {:?}", res);

    info!(
        "going load {} items in {} chunks into index {}: {:?}",
        data.len(),
        chunks.len(),
        index,
        res
    );
}

#[tracing::instrument(name = "route::index_recipe_detail", skip(cr))]
pub async fn index_recipe_detail(cr: web::Json<Vec<RecipeDetail>>) -> HttpResponse {
    load(&cr.0, Index::RecipeDetails).await;
    HttpResponse::Ok().json(())
}

#[tracing::instrument]
pub async fn init_indexes() {
    let client = get_client();

    for x in Index::iter() {
        if let Some(attr) = x.get_searchable_attributes() {
            let task: TaskInfo = client
                .index(x)
                .set_searchable_attributes(attr)
                .await
                .unwrap();
            info!("configured searchable attributes for {}: {:?}", x, task);
        }
        if let Some(attr) = x.get_filterable_attributes() {
            let task: TaskInfo = client
                .index(x)
                .set_filterable_attributes(attr)
                .await
                .unwrap();
            info!("configured filterable attributes for {}: {:?}", x, task);
        }
    }
}
