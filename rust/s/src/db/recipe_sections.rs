use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Deserialize, Serialize)]
#[sea_orm(table_name = "recipe_sections")]
pub struct Model {
    #[sea_orm(primary_key)]
    id: String,
    #[sea_orm(column_type = "Text")]
    recipe_detail: String,
    #[sea_orm(column_type = "Json")]
    duration_timerange: Json,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(
        belongs_to = "super::recipe_sections::Entity",
        from = "Column::RecipeDetail",
        to = "super::recipe_sections::Column::Id"
    )]
    RecipeSections,
}

impl Related<super::recipes::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::RecipeSections.def()
    }
}

impl ActiveModelBehavior for ActiveModel {}
