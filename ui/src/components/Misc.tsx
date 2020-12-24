import React from "react";
import { Link } from "react-router-dom";
import { RecipeDetail } from "../api/openapi-hooks/api";

export interface Props {
  recipe: RecipeDetail;
}
export const RecipeLink: React.FC<Props> = ({ recipe }) => (
  <div className="flex space-x-0.5">
    <Link to={`recipe/${recipe.id}`} className="link">
      <div
        className={`${
          recipe.is_latest_version ? "text-blue-800" : "text-blue-200"
        }`}
      >
        {recipe.name}
      </div>
    </Link>
    <div>
      <div className="flex font-mono">v{recipe.version}</div>
    </div>
  </div>
);
