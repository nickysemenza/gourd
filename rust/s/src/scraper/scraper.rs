use pyo3::{types::PyModule, PyAny, Python};
use serde::Serialize;

#[derive(Debug, Clone, Serialize)]
pub struct ScrapeResult {
    pub ingredients: Vec<String>,
    pub instructions: String,
    pub title: String,
    pub url: String,
}

#[tracing::instrument(name = "route::scrape_recipe")]
pub fn scrape_recipe(url: &str) -> ScrapeResult {
    let mut sc_result: (Vec<String>, String, String) = (vec![], "".to_string(), "".to_string());
    Python::with_gil(|py| {
        let syspath: &PyAny = py.import("sys").unwrap().get("path").unwrap();

        dbg!(syspath);
        let activators = PyModule::from_code(
            py,
            r#"
from recipe_scrapers import scrape_me            
def sc(x,y):
    res = scrape_me(x,wild_mode=y)
    return res.ingredients(), res.instructions(), res.title()
            "#,
            "recipe_scrape.py",
            "recipe_scrape",
        )
        .unwrap();

        dbg!(activators);
        sc_result = activators
            .getattr("sc")
            .unwrap()
            .call((url.clone(), true), None)
            .unwrap()
            .extract()
            .unwrap();
    });
    return ScrapeResult {
        ingredients: sc_result.0,
        instructions: sc_result.1,
        title: sc_result.2,
        url: url.to_string(),
    };
}