mod index;

use anyhow::{Context, Result};
use axum::{extract, response::IntoResponse, Json};
use indicatif::ProgressIterator;
use meilisearch_sdk::{client::Client, task_info::TaskInfo};
use openapi::models::RecipeDetail;
use serde::{de::DeserializeOwned, Serialize};
use strum::IntoEnumIterator;
use tracing::info;

pub use self::index::Index;

fn get_client() -> Client {
    let meilisearch_url = option_env!("MEILISEARCH_URL").unwrap_or("http://localhost:7700");
    let meilisearch_api_key = option_env!("MEILISEARCH_API_KEY").unwrap_or("FOO");
    Client::new(meilisearch_url, Some(meilisearch_api_key)).expect("failed to create client")
}

#[derive(Debug)]

pub struct Searcher {
    client: Client,
}
impl Searcher {
    pub fn new() -> Self {
        Self {
            client: get_client(),
        }
    }

    #[tracing::instrument]
    pub async fn search<T: Document>(
        &self,
        index: Index,
        limit: usize,
        name: &str,
    ) -> Result<Vec<T>> {
        let results = self
            .client
            .index(index)
            .search()
            .with_query(name)
            .with_limit(limit)
            .execute::<T>()
            .await
            .context(format!("search {index} for {name}"))?;
        info!("searched in {}ms", results.processing_time_ms);
        Ok(results.hits.into_iter().map(|x| x.result).collect())
    }

    pub async fn get_document<T: Document>(&self, index: Index, name: &str) -> Result<Option<T>> {
        match self.client.index(index).get_document::<T>(name).await {
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

    #[tracing::instrument(skip(data))]
    pub async fn load<T: Document>(self, data: &Vec<T>, index: Index) {
        let chunks: Vec<Vec<T>> = data.chunks(2000).map(|x| x.to_vec()).collect();

        let tasks: Vec<_> = chunks
            .iter()
            .progress()
            .map(|v| {
                let client = self.client.clone();
                let x = v.clone();
                let i = index;
                tokio::spawn(async move { client.index(i).add_documents(&x, None).await.unwrap() })
            })
            .collect();
        let _res = futures::future::join_all(tasks).await;

        info!(
            "loaded {} items in {} chunks into index {}",
            data.len(),
            chunks.len(),
            index,
        );
    }

    #[tracing::instrument]
    pub async fn init_indexes(self) {
        for x in Index::iter().progress() {
            let mut touched = false;
            if let Some(attr) = x.get_searchable_attributes() {
                let task: TaskInfo = self
                    .client
                    .index(x)
                    .set_searchable_attributes(attr)
                    .await
                    .unwrap();
                info!("configured searchable attributes for {}: {:?}", x, task);
                touched = true;
            }
            if let Some(attr) = x.get_filterable_attributes() {
                let task: TaskInfo = self
                    .client
                    .index(x)
                    .set_filterable_attributes(attr)
                    .await
                    .unwrap();
                info!("configured filterable attributes for {}: {:?}", x, task);
                touched = true;
            }
            if !touched {
                let task: TaskInfo = self.client.create_index(x.to_string(), None).await.unwrap();
                info!("created index without config {}: {:?}", x, task);
            }
        }
    }
}

pub trait Document: Clone + Serialize + DeserializeOwned + Send + Sync + 'static {}
impl<T: Clone + Serialize + DeserializeOwned + Send + Sync + 'static> Document for T {}

#[tracing::instrument(name = "route::index_recipe_detail", skip(cr))]
pub async fn index_recipe_detail(
    extract::Json(cr): extract::Json<Vec<RecipeDetail>>,
) -> impl IntoResponse {
    Searcher::new().load(&cr, Index::RecipeDetails).await;
    Json(())
}
