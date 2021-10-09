import {
  Ingredient,
  SectionIngredient,
  RecipeDetail,
} from "./api/openapi-hooks/api";
import { RecipeWrapper, TimeRange } from "./api/openapi-fetch";
import parse from "parse-duration";
import { RecipeWrapperInput } from "./api/openapi-fetch/models/RecipeWrapperInput";

export const getIngredient = (
  si: Partial<SectionIngredient>
): { name: "" } | RecipeDetail | Ingredient => {
  if (si.recipe) {
    return si.recipe;
  } else if (si.ingredient) {
    return si.ingredient.ingredient;
  }
  return { name: "" };
};

export const formatText = (text: React.ReactText) => {
  // regexr.com/5mt55
  const re = /[\d]+ ?F/g;

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
      <code className="text-red-800 mx-1" key={matchStart}>
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
export const formatTimeRange = (range?: TimeRange) => {
  if (!range) return "";
  const { min, max } = range;
  let items = [formatSeconds(min)];
  if (max > 0 && max !== min) {
    items.push(" - ", formatSeconds(max));
  }
  return items.join("");
};

export const parseTimeRange = (input: string): TimeRange | null => {
  const parts = input.split(" - ");
  if (parts.length === 0 || parts.length > 2) return null;
  return {
    min: (parse(parts[0]) || 0) / 1000,
    max: ((parts.length === 2 && parse(parts[1])) || 0) / 1000,
  };
};

export const sumTimeRanges = (ranges: (TimeRange | undefined)[]): TimeRange => {
  let totalDuration: TimeRange = { min: 0, max: 0 };
  ranges.forEach((r) => {
    if (!!r) {
      totalDuration.min += r.min;
      totalDuration.max += r.max || r.min; // if max is 0, it means it's a finite and not a range, so use min
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
