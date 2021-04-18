import React from "react";
import { Link } from "react-router-dom";
import { RecipeDetail, UnitMapping } from "../api/openapi-hooks/api";

export interface Props {
  recipe: RecipeDetail;
  multiplier?: number;
}
export const RecipeLink: React.FC<Props> = ({
  recipe: { name, version, is_latest_version, id },
  multiplier,
}) => (
  <div className="flex space-x-0.5">
    <Link to={`/recipe/${id}?multiplier=${multiplier || 1}`} className="link">
      <div
        className={`${is_latest_version ? "text-blue-800" : "text-blue-200"}`}
      >
        {name}
      </div>
    </Link>
    <div className="flex font-mono">v{version}</div>
    {multiplier && <div className="font-mono">@{multiplier}x</div>}
  </div>
);

export const UnitMappingList: React.FC<{ unit_mappings: UnitMapping[] }> = ({
  unit_mappings,
}) => (
  <ul>
    {unit_mappings &&
      unit_mappings.map((m) => (
        <li>
          {m.a.value} {m.a.unit} = {m.b.value} {m.b.unit} ({m.source})
        </li>
      ))}
  </ul>
);
