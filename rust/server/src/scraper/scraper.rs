use anyhow::{Context, Result};
use openapi::models::{CompactRecipe, CompactRecipeMeta, CompactRecipeSection};
use pyo3::{types::PyModule, PyAny, Python};

#[tracing::instrument(name = "route::scrape_recipe")]
pub fn scrape_recipe(url: &str) -> Result<CompactRecipe> {
    let mut sc_result: (Vec<String>, String, String, String) =
        (vec![], "".to_string(), "".to_string(), "".to_string());
    Python::with_gil(|py| -> Result<()> {
        let syspath: &PyAny = py.import("sys").unwrap().getattr("path").unwrap();

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

    Ok(CompactRecipe {
        sections: vec![CompactRecipeSection {
            ingredients: sc_result.0.into_iter().map(|i| i.to_string()).collect(),
            instructions: sc_result.1.split('\n').map(|i| i.to_string()).collect(),
        }],
        meta: Box::new(CompactRecipeMeta {
            name: sc_result.2,
            image: match sc_result.3 == "" {
                true => None,
                false => Some(sc_result.3),
            },
            url: Some(url.to_string()),
        }),
    })
}
