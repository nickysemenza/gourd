/*
 * Gourd Recipe Database
 *
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 * Generated by: https://openapi-generator.tech
 */

/// Photo : A photo



#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct Photo {
    /// id
    #[serde(rename = "id")]
    pub id: String,
    /// public image
    #[serde(rename = "base_url")]
    pub base_url: String,
    /// blur hash
    #[serde(rename = "blur_hash", skip_serializing_if = "Option::is_none")]
    pub blur_hash: Option<String>,
    /// when it was taken
    #[serde(rename = "created")]
    pub created: String,
    /// width px
    #[serde(rename = "width")]
    pub width: i64,
    /// height px
    #[serde(rename = "height")]
    pub height: i64,
    /// where the photo came from
    #[serde(rename = "source")]
    pub source: Source,
}

impl Photo {
    /// A photo
    pub fn new(id: String, base_url: String, created: String, width: i64, height: i64, source: Source) -> Photo {
        Photo {
            id,
            base_url,
            blur_hash: None,
            created,
            width,
            height,
            source,
        }
    }
}

/// where the photo came from
#[derive(Clone, Copy, Debug, Eq, PartialEq, Ord, PartialOrd, Hash, Serialize, Deserialize)]
pub enum Source {
    #[serde(rename = "google")]
    Google,
    #[serde(rename = "notion")]
    Notion,
}
