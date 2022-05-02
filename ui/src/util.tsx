import {
  Ingredient,
  SectionIngredient,
  IngredientKind,
  Amount,
  RecipeSection,
} from "./api/openapi-hooks/api";
import { RecipeWrapperInput } from "./api/openapi-fetch/models/RecipeWrapperInput";
import { wasm, RichItem } from "./wasmContext";

export const getIngredient = (si: Partial<SectionIngredient>) => {
  let name = "";
  let kind: IngredientKind = "ingredient";
  if (si.recipe) {
    name = si.recipe.name;
    kind = "recipe";
  } else if (si.ingredient) {
    name = si.ingredient.ingredient.name;
  }
  return { name, kind };
};

export const formatRichText = (w: wasm, text: RichItem[]) => {
  return text.map((t, x) => {
    if (t.kind === "Text") {
      return t.value;
    } else if (t.kind === "Ing") {
      return (
        <div
          className="inline text-orange-800 m-0 underline decoration-grey decoration-solid"
          key={x + "a"}
        >
          {t.value}
        </div>
      );
    } else if (t.kind === "Amount") {
      let val = t.value.pop();
      if (!val) {
        return null;
      }
      if (val.unit === "whole") {
        val.unit = "";
      }
      if (val.value === null) {
        val.value = 0;
      }
      return (
        <div
          className="inline text-green-800 m-0 underline decoration-grey decoration-solid"
          key={x}
        >
          {/* // sum_time_amounts converts the units to the best scaled */}
          {w.format_amount(val)}
        </div>
      );
    } else {
      return null;
    }
  });
};
export const formatText = (text: React.ReactText) => {
  // regexr.com/5mt55
  const re = /[\d]+ ?(F|C|°|°F)/g;

  if (typeof text === "number") {
    return text;
  }

  let pairs = [];
  const matches = [...text.matchAll(re)];
  if (matches.length === 0) {
    return text;
  }

  let lastProcessed = 0;
  for (const match of matches) {
    const matchStart = match.index || 0;
    const matchEnd = matchStart + match[0].length;
    pairs.push(text.substring(lastProcessed, matchStart));
    pairs.push(
      <code
        className="text-red-800 m-0 underline decoration-grey decoration-solid"
        key={matchStart}
      >
        {text.substring(matchStart, matchEnd)}
      </code>
    );
    lastProcessed = matchEnd;
  }
  pairs.push(text.substring(lastProcessed));
  return pairs;
};

export const formatTimeRange = (w?: wasm, range?: Amount) => {
  // sum_time_amounts converts the units to the best scaled
  return w && range ? w.format_amount(w.sum_time_amounts([range])) : "";
};

export const sumTimeRanges = (
  w: wasm,
  ranges: (Amount | undefined)[]
): Amount => {
  let ranges2 = ranges.filter((r) => r !== undefined) as Amount[];
  return w.sum_time_amounts(ranges2);
};

export const Code: React.FC<{
  children?: React.ReactNode;
}> = ({ children }) => (
  <code className="rounded-sm px-2 py-0.5 bg-gray-100 text-red-500 h-6 flex text-sm font-bold w-min">
    {children}
  </code>
);

export const blankRecipeWrapperInput = (
  name: string = "",
  id: string = ""
): RecipeWrapperInput => ({
  detail: {
    name,
    quantity: 0,
    unit: "",
    sections: [],
    tags: [],
  },
  id,
});
export const blankIngredient = (name: string): Ingredient => ({ name, id: "" });

export const scaledRound = (x: number) => x.toFixed(x < 10 ? 2 : 0);

export const getTotalDuration = (w: wasm, sections: RecipeSection[]) =>
  sumTimeRanges(
    w,
    (sections || []).map((s) => s.duration).filter((t) => t !== undefined)
  );
