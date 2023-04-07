use strum_macros::{Display, EnumIter};

#[derive(EnumIter, Debug, Clone, Copy, Display)]
pub enum Index {
    BrandedFoods,
    FoundationFoods,
    SurveyFoods,
    SRLegacyFoods,
    ScrapedRecipes,
    RecipeDetails,
}
impl Index {
    pub fn get_top_level(&self) -> String {
        match self {
            Index::BrandedFoods
            | Index::FoundationFoods
            | Index::SurveyFoods
            | Index::SRLegacyFoods => self.to_string(),
            Index::ScrapedRecipes | Index::RecipeDetails => {
                panic!("ScrapedRecipes is not a top level index")
            }
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
            Index::ScrapedRecipes => None,
            Index::RecipeDetails => None,
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
            Index::ScrapedRecipes => Some(vec!["name", "url", "sections"]),
            Index::RecipeDetails => Some(vec!["unit", "is_latest_version", "tags"]),
        }
    }
}

impl Into<String> for Index {
    fn into(self) -> String {
        self.to_string()
    }
}