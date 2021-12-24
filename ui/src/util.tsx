import {
  Ingredient,
  SectionIngredient,
  IngredientKind,
  Amount,
  RecipeSection,
} from "./api/openapi-hooks/api";
import parse from "parse-duration";
import { RecipeWrapperInput } from "./api/openapi-fetch/models/RecipeWrapperInput";
import { RichItem } from "gourd_rs";

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

export const formatRichText = (text: RichItem[]) => {
  return text.map((t, x) => {
    if (t.kind === "Text") {
      return t.value;
    } else if (t.kind === "Amount") {
      let val = t.value.pop();
      return val ? (
        <div
          className="inline text-green-800 m-0 underline decoration-grey decoration-solid"
          key={x}
        >
          {val.value}
          {val.upper_value && `-${val.upper_value}`}{" "}
          {val.unit !== "whole" && val.unit}
        </div>
      ) : null;
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

const formatSeconds = (seconds: number) => {
  // https://stackoverflow.com/a/6312999
  const secs = Math.round(seconds);
  const h = Math.floor(secs / (60 * 60));
  const divisor_for_minutes = secs % (60 * 60);
  const m = Math.floor(divisor_for_minutes / 60);
  const s = Math.ceil(divisor_for_minutes % 60);

  let vals = [];
  vals.push(h > 0 ? `${h} hr` : null);
  vals.push(m > 0 ? `${m} min` : null);
  vals.push(s > 0 ? `${s} sec` : null);
  return vals.join(" ");
};
export const formatTimeRange = (range?: Amount) => {
  if (!range) return "";
  const { value, upper_value } = range;
  let items = [formatSeconds(value)];
  if (upper_value) {
    items.push(" - ", formatSeconds(upper_value));
  }
  return items.join("");
};

export const parseTimeRange = (input: string): Amount | null => {
  // todo: wasm this
  const parts = input.split(" - ");
  if (parts.length === 0 || parts.length > 2) return null;
  return {
    unit: "seconds",
    value: (parse(parts[0]) || 0) / 1000,
    upper_value: parts.length === 2 ? (parse(parts[1]) || 0) / 1000 : undefined,
  };
};

export const sumTimeRanges = (ranges: (Amount | undefined)[]): Amount => {
  // todo: wasm this
  let totalDuration: Amount = { value: 0, unit: "seconds" };
  ranges.forEach((r) => {
    if (!!r) {
      totalDuration.value += r.value;

      if (r.upper_value !== undefined) {
        if (totalDuration.upper_value === undefined) {
          totalDuration.upper_value = 0;
        }
        totalDuration.upper_value += r.upper_value;
      }
    }
  });

  return totalDuration;
};

export const Code: React.FC = ({ children }) => (
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
  },
  id,
});
export const blankIngredient = (name: string): Ingredient => ({ name, id: "" });

export const scaledRound = (x: number) => x.toFixed(x < 10 ? 2 : 0);

export const getTotalDuration = (sections: RecipeSection[]) =>
  sumTimeRanges(sections.map((s) => s.duration).filter((t) => t !== undefined));
