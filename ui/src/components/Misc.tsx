import Graphviz from "graphviz-react";
import React, { useContext } from "react";
import { Link } from "react-router-dom";
import {
  UnitConversionRequest,
  UnitConversionRequestTargetEnum,
} from "../api/openapi-fetch";
import { RecipeDetail, UnitMapping } from "../api/openapi-hooks/api";
import { scaledRound } from "../util";
import { WasmContext } from "../wasm";

export interface Props {
  recipe: RecipeDetail;
  multiplier?: number;
}
export const RecipeLink: React.FC<Props> = ({
  recipe: { name, version, is_latest_version, id },
  multiplier,
}) => (
  <div className="inline-flex">
    <Link
      to={`/recipe/${id}?multiplier=${multiplier || 1}`}
      className={`font-bold pr-0.5 underline ${
        is_latest_version
          ? "decoration-blue-300 text-blue-800"
          : "decoration-red-300 text-purple-400"
      }`}
    >
      {name}
    </Link>
    <div className="flex font-mono small">v{version}</div>
    {multiplier && <div className="font-mono">@{multiplier}x</div>}
  </div>
);

export const UnitMappingList: React.FC<{
  unit_mappings: UnitMapping[];
  includeDot?: boolean;
}> = ({ unit_mappings, includeDot = false }) => {
  const w = useContext(WasmContext);
  let dot = "";
  if (unit_mappings.length > 0 && w) {
    let foo: UnitConversionRequest = {
      target: UnitConversionRequestTargetEnum.MONEY,
      unit_mappings,
      input: [{ unit: "grams", value: 100 }],
    };

    try {
      dot = w.make_dag(foo);
    } catch (e) {
      console.error({ e });
    }
  }
  if (!unit_mappings || !w) return null;
  return (
    <div>
      {dot !== "" && includeDot && (
        <Graphviz
          dot={dot}
          options={{ width: 300, height: null }}
          className="w-full"
        />
      )}
      <div className="w-60">
        {unit_mappings.map((m, x) => (
          <div
            key={x}
            className="flex text-sm text-gray-700 dark:text-gray-400"
          >
            <p>
              {scaledRound(m.a.value)} {m.a.unit}
            </p>
            <p className="text-center px-1">=</p>
            <p>
              {scaledRound(m.b.value)} {m.b.unit}
            </p>
            <p className="text-xs pl-1">{m.source}</p>
          </div>
        ))}
      </div>
    </div>
  );
};
