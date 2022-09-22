/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */




#[derive(Clone, Debug, PartialEq, Default, Serialize, Deserialize)]
pub struct ScrapeRecipeRequest {
    #[serde(rename = "url")]
    pub url: String,
}

impl ScrapeRecipeRequest {
    pub fn new(url: String) -> ScrapeRecipeRequest {
        ScrapeRecipeRequest {
            url,
        }
    }
}

