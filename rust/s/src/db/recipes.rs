use sea_orm::entity::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, PartialEq, DeriveEntityModel, Deserialize, Serialize)]
#[sea_orm(table_name = "recipes")]
pub struct Model {
    #[sea_orm(primary_key)]
    id: String,
}

#[derive(Copy, Clone, Debug, EnumIter, DeriveRelation)]
pub enum Relation {
    #[sea_orm(has_many = "super::recipe_details::Entity")]
    RecipeDetails,
}

impl Related<super::recipe_details::Entity> for Entity {
    fn to() -> RelationDef {
        Relation::RecipeDetails.def()
    }
}

impl ActiveModelBehavior for ActiveModel {}
