import React from "react";
import { RecipeDetail } from "../api/openapi-fetch";
import Debug from "./Debug";
import { sumIngredients } from "./RecipeEditorUtils";

const RecipeDiff: React.FC<{ details: RecipeDetail[] }> = ({ details }) => {
  const uniqueIngredientIDs = [
    ...new Set(
      details
        .map((detail) =>
          Object.keys(sumIngredients(detail.sections).ingredients)
        )
        .flat()
    ),
  ];
  return <Debug data={{ uniqueIngredientIDs, details }} />;
};

export default RecipeDiff;
