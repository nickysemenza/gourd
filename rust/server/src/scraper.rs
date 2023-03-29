use anyhow::Result;
use openapi::models::{CompactRecipe, CompactRecipeMeta, CompactRecipeSection};

use crate::search::get_client;

#[tracing::instrument(name = "route::scrape_recipe")]
pub async fn scrape_recipe(url: &str) -> Result<CompactRecipe> {
    let s = recipe_scraper_fetcher::Fetcher::new();
    let res = s.scrape_url(url).await?;

    let client = get_client();
    let mut res2 = res.clone();
    res2.url = res2.url.replace(|c: char| !c.is_alphanumeric(), "");
    client
        .index("scraped_recipes")
        .add_documents(&vec![res2], Some("url"))
        .await
        .unwrap();

    let meta = CompactRecipeMeta {
        name: res.name,
        image: res.image,
        url: Some(url.to_string()),
    };

    let compact = CompactRecipe {
        sections: vec![CompactRecipeSection {
            ingredients: res.ingredients,
            instructions: res.instructions,
        }],
        meta: Box::new(meta),
    };

    // client
    //     .index("recipes")
    //     .add_documents(&vec![compact.clone()], Some("url"))
    //     .await
    //     .unwrap();

    Ok(compact)
}
