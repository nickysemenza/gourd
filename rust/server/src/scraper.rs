use anyhow::Result;
use openapi::models::{CompactRecipe, CompactRecipeMeta, CompactRecipeSection};

#[tracing::instrument(name = "route::scrape_recipe")]
pub async fn scrape_recipe(url: &str) -> Result<CompactRecipe> {
    let s = recipe_scraper_fetcher::Fetcher::new();
    let res = s.scrape_url(url).await?;

    let meta = CompactRecipeMeta {
        name: res.name,
        image: res.image,
        url: Some(url.to_string()),
    };
    Ok(CompactRecipe {
        sections: vec![CompactRecipeSection {
            ingredients: res.ingredients,
            instructions: res.instructions,
        }],
        meta: Box::new(meta),
    })
}
