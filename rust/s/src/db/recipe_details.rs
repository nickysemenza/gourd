use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Deserialize, Serialize)]
#[sea_orm(table_name = "recipe_details")]
pub struct Model {
    #[sea_orm(primary_key)]
    id: String,
    #[sea_orm(column_type = "Text")]
    recipe: String,
    #[sea_orm(column_type = "Text")]
    name: String,
    // equipment: String,
    #[sea_orm(column_type = "Json")]
    source: Json,
    // servings: i32,
    // quantity: i32,
    // unit: String,
    // version: i32,
    // is_latest_version: bool,
    // #[sea_orm(column_type = "DateTime")]
    // created_at: DateTime,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(
        belongs_to = "super::recipes::Entity",
        from = "Column::Recipe",
        to = "super::recipes::Column::Id"
    )]
    Recipes,
}

impl Related<super::recipes::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::Recipes.def()
    }
}

impl ActiveModelBehavior for ActiveModel {}
