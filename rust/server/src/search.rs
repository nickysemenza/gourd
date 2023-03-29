use std::fmt;

use serde::{de::DeserializeOwned, Serialize};
use tracing::info;

pub fn get_client() -> meilisearch_sdk::Client {
    let meilisearch_url = option_env!("MEILISEARCH_URL").unwrap_or("http://localhost:7700");
    let meilisearch_api_key = option_env!("MEILISEARCH_API_KEY").unwrap_or("FOO");
    meilisearch_sdk::Client::new(meilisearch_url, meilisearch_api_key)
}

#[derive(Debug, Clone, Copy)]
pub enum Index {
    BrandedFoods,
    FoundationFoods,
    SurveyFoods,
    SRLegacyFoods,
}
impl Index {
    pub fn get_top_level(&self) -> String {
        match self {
            Index::BrandedFoods
            | Index::FoundationFoods
            | Index::SurveyFoods
            | Index::SRLegacyFoods => self.to_string(),
        }
    }
    pub fn get_searchable_attributes(&self) -> Option<Vec<&str>> {
        match self {
            Index::BrandedFoods => Some(vec![
                "description",
                "ingredients",
                "brandOwner",
                "fdcId",
                "brandedFoodCategory",
            ]),
            Index::FoundationFoods => None,
            Index::SurveyFoods => None,
            Index::SRLegacyFoods => None,
        }
    }
    pub fn get_filterable_attributes(&self) -> Option<Vec<&str>> {
        match self {
            Index::BrandedFoods => Some(vec![
                "brandOwner",
                "brandedFoodCategory",
                "ingredients",
                "description",
            ]),
            Index::FoundationFoods => None,
            Index::SurveyFoods => None,
            Index::SRLegacyFoods => None,
        }
    }
}

impl Into<String> for Index {
    fn into(self) -> String {
        self.to_string()
    }
}
impl fmt::Display for Index {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Index::BrandedFoods => write!(f, "BrandedFoods"),
            Index::FoundationFoods => write!(f, "FoundationFoods"),
            Index::SurveyFoods => write!(f, "SurveyFoods"),
            Index::SRLegacyFoods => write!(f, "SRLegacyFoods"),
        }
    }
}

pub trait Document: Clone + Serialize + DeserializeOwned + Send + Sync + 'static {}
impl<T: Clone + Serialize + DeserializeOwned + Send + Sync + 'static> Document for T {}

#[tracing::instrument(skip(data))]
pub async fn load<T: Document>(data: &Vec<T>, index: Index) {
    // let chunks: Vec<&[FoundationFoodItem]> = legacy.chunks(10).collect();
    let chunks: Vec<Vec<T>> = data.chunks(2000).map(|x| x.to_vec()).collect();
    info!("going to load {} chunks into index {}", chunks.len(), index);

    let tasks: Vec<_> = chunks
        .iter()
        .map(|v| {
            let client = get_client();
            let x = v.clone();
            let i = index.clone();
            tokio::spawn(async move { client.index(i).add_documents(&x, None).await.unwrap() })
        })
        .collect();
    futures::future::join_all(tasks).await;
    info!("finished loading")
}
