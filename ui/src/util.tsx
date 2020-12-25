import {
  Ingredient,
  SectionIngredient,
  RecipeDetail,
} from "./api/openapi-hooks/api";
import { TimeRange } from "./api/openapi-fetch";

export const getIngredient = (
  si: Partial<SectionIngredient>
): { name: "" } | RecipeDetail | Ingredient => {
  if (si.recipe) {
    return si.recipe;
  } else if (si.ingredient) {
    return si.ingredient;
  }
  return { name: "" };
};

export const formatText = (text: React.ReactText) => {
  const re = /[\d]* ?F/g;

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
      <code className="text-red-800 mx-1">
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

  return (
    (h > 0 ? `${h}h` : "") + (m > 0 ? `${m}m` : "") + (s > 0 ? `${s}s` : "")
  );
};
export const formatTimeRange = (range: TimeRange) => {
  const { min, max } = range;
  let items = [formatSeconds(min)];
  if (max > 0) {
    items.push(" - ", formatSeconds(max));
  }
  return <div className="flex">{items}</div>;
};
