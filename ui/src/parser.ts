import { RecipeWrapper, SectionIngredient } from "./api/openapi-hooks/api";
import { getIngredient } from "./util";

export const parseLine = (line: string): [number, string, string] => {
  if (line.includes(",")) {
    const re = /(\d*)(g (.*?), )(\S*)/;
    const matches = re.exec(line);
    if (!matches) throw new Error("invalid match for: " + line);
    return [parseFloat(matches[1]), matches[3], matches[4].toLowerCase()];
  } else {
    const re = /(\d*)(g (.*?))(\S*)/;
    const matches = re.exec(line);
    if (!matches) throw new Error("invalid match for: " + line);
    return [parseFloat(matches[1]), matches[4].toLowerCase(), ""];
  }
};

export const parseRecipe = (text: string): RecipeWrapper | undefined => {
  const lines = text.split(/\r?\n/);
  const ingredients: Array<SectionIngredient> = lines.map((line) => {
    const [grams, ingredient, adjective] = parseLine(line);
    const si: SectionIngredient = {
      id: "",
      ingredient: { id: "", name: ingredient },
      kind: "ingredient",
      grams,
      adjective,
      amount: 0,
      unit: "",
      optional: false,
    };
    return si;
  });
  const recipe: RecipeWrapper = {
    detail: {
      name: "",
      id: "",
      quantity: 0,
      unit: "",
      sections: [
        { id: "", duration: { min: 0, max: 0 }, ingredients, instructions: [] },
      ],
    },

    id: "",
  };
  return recipe;
};
export const encodeIngredient = (ingredient: SectionIngredient): string => {
  const { grams, adjective, amount, unit } = ingredient;
  let res = grams + "g";
  if (amount !== 0 && unit !== "") res += ` (${amount} ${unit})`;
  res += ` ${getIngredient(ingredient).name}`;
  if (adjective !== "") res += `, ${adjective}`;
  return res;
};

export const encodeRecipe = (recipe: RecipeWrapper): string =>
  recipe && recipe.detail.sections
    ? recipe.detail.sections
        .map((section) => section.ingredients.map(encodeIngredient).join("\n"))
        .join("\n")
    : "";
