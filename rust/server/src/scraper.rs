use anyhow::{Context, Result};
use openapi::models::{CompactRecipe, CompactRecipeSection};
use url::Url;

use crate::search::Searcher;

#[tracing::instrument(name = "route::scrape_recipe")]
pub async fn scrape_recipe(url: &str) -> Result<CompactRecipe> {
    let s = recipe_scraper_fetcher::Fetcher::new();
    let res = s
        .scrape_url(url)
        .await
        .with_context(|| format!("Failed to scrape {}", url))?;

    let parsed_url = Url::parse(url).unwrap();
    let compact = CompactRecipe {
        id: format!(
            "{}-{}",
            parsed_url.host_str().unwrap(),
            parsed_url.path().replace("/", "-")
        )
        .replace(|c: char| !c.is_alphanumeric() && c != '-' && c != '_', ""),

        name: res.name,
        image: res.image,
        url: Some(url.to_string()),
        sections: vec![CompactRecipeSection {
            ingredients: res.ingredients,
            instructions: res.instructions,
        }],
    };

    Searcher::new()
        .load(&vec![compact.clone()], crate::search::Index::ScrapedRecipes)
        .await;

    Ok(compact)
}
