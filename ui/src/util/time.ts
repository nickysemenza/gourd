import { Amount, RecipeSection } from "../api/openapi-hooks/api";
import { wasm } from "./wasmContext";

export const sumTimeRanges = (
  w: wasm,
  ranges: (Amount | undefined)[]
): Amount => {
  const ranges2 = ranges.filter((r) => r !== undefined) as Amount[];
  return w.sum_time_amounts(ranges2);
};

export const getTotalDuration = (w: wasm, sections: RecipeSection[]) =>
  sumTimeRanges(
    w,
    (sections || []).map((s) => s.duration).filter((t) => t !== undefined)
  );

export const formatTimeRange = (w?: wasm, range?: Amount) => {
  // sum_time_amounts converts the units to the best scaled
  return w && range ? w.format_amount(w.sum_time_amounts([range])) : "";
};
