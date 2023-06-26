import Graphviz from "graphviz-react";
import React, { useContext } from "react";
import { Link } from "react-router-dom";

import {
  RecipeDetail,
  UnitConversionRequest,
  UnitMapping,
} from "../../api/react-query/gourdApiSchemas";
import { scaledRound } from "../../util/util";
import { WasmContext } from "../../util/wasmContext";
import { cn } from "../ui/lib";

export interface Props {
  recipe: Pick<RecipeDetail, "id" | "name" | "meta">;
  multiplier?: number;
}
export const RecipeLink: React.FC<Props> = ({
  recipe: { name, meta, id },
  multiplier,
}) => (
  <span className="inline-block">
    <Link
      to={`/recipe/${id}?multiplier=${multiplier || 1}`}
      className={`font-bold pr-0.5 underline ${
        meta.is_latest_version
          ? "decoration-blue-300 text-blue-800 dark:text-blue-400"
          : "decoration-red-300 text-purple-400"
      }`}
    >
      {name}
    </Link>
    {/* -1 is used to signal version is unknown */}
    {meta.version >= 0 && (
      <div className="inline font-mono text-sm">v{meta.version}</div>
    )}
    {multiplier && (
      <div className="inline font-mono text-sm">@{multiplier}x</div>
    )}
  </span>
);

export const UnitMappingList: React.FC<{
  unit_mappings: UnitMapping[];
  includeDot?: boolean;
}> = ({ unit_mappings, includeDot = false }) => {
  const w = useContext(WasmContext);
  let dot = "";
  if (unit_mappings.length > 0 && w) {
    const foo: UnitConversionRequest = {
      target: "money",
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
      <div className="">
        {unit_mappings.map((m, x) => {
          const hasKCalGrams = m.a.unit === "grams" && m.b.unit === "kcal";
          return (
            <div
              key={x}
              className={cn(`flex text-sm text-gray-700 dark:text-gray-400`)}
            >
              <div>
                {scaledRound(m.a.value)} {m.a.unit}
              </div>
              <div className="text-center px-1">=</div>
              <div>
                {scaledRound(m.b.value)} {m.b.unit}
              </div>
              <div className="text-xs pl-1">{m.source}</div>
              {hasKCalGrams && "💪"}
            </div>
          );
        })}
      </div>
    </div>
  );
};
