import {
  Ingredient,
  SectionIngredient,
  RecipeDetail,
} from "./api/openapi-hooks/api";

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
    // pairs.push()
    lastProcessed = matchEnd;
    // pairs.push([, ]);
  }
  pairs.push(text.substring(lastProcessed));
  // let res = [];
  // for
  return pairs;

  // console.log(pairs);
};
