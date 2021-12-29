use anyhow::{Context, Result};
use pyo3::{types::PyModule, PyAny, Python};
use serde::Serialize;

#[derive(Debug, Clone, Serialize)]
pub struct ScrapeResult {
    pub ingredients: Vec<String>,
    pub instructions: String,
    pub title: String,
    pub url: String,
    pub image: String,
}

#[tracing::instrument(name = "route::scrape_recipe")]
pub fn scrape_recipe(url: &str) -> Result<ScrapeResult> {
    let mut sc_result: (Vec<String>, String, String, String) =
        (vec![], "".to_string(), "".to_string(), "".to_string());
    Python::with_gil(|py| -> Result<()> {
        let syspath: &PyAny = py.import("sys").unwrap().get("path").unwrap();

        dbg!(syspath);
        let activators = PyModule::from_code(
            py,
            r#"
from recipe_scrapers import scrape_me            
def sc(x,y):
    res = scrape_me(x,wild_mode=y)
    return res.ingredients(), res.instructions(), res.title(), res.image()
            "#,
            "recipe_scrape.py",
            "recipe_scrape",
        )
        .context("failed to build py")?;

        dbg!(activators);
        sc_result = activators
            .getattr("sc")
            .context("failed to get attribute")?
            .call((url.clone(), true), None)
            .context("failed to call")?
            .extract()
            .context("failed to extract")?;
        Ok(())
    })
    .context("failed to parse")?;
    Ok(ScrapeResult {
        ingredients: sc_result.0,
        instructions: sc_result.1,
        title: sc_result.2,
        url: url.to_string(),
        image: sc_result.3,
    })
}
