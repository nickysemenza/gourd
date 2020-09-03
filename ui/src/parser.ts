import {
  Recipe,
  SectionIngredient,
  SectionIngredientKind,
  GetRecipeByUuidQuery,
} from "./generated/graphql";

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

export const parseRecipe = (text: string): Recipe | undefined => {
  const lines = text.split(/\r?\n/);
  const ingredients: Array<SectionIngredient> = lines.map((line) => {
    const [grams, ingredient, adjective] = parseLine(line);
    const si: SectionIngredient = {
      uuid: "",
      info: { uuid: "", name: ingredient },
      kind: SectionIngredientKind.Ingredient,
      grams,
      adjective,
      amount: 0,
      unit: "",
      optional: false,
    };
    return si;
  });
  const recipe: Recipe = {
    name: "",
    meals: [],
    uuid: "",
    totalMinutes: 0,
    unit: "",
    sections: [{ uuid: "", minutes: 0, ingredients, instructions: [] }],
    notes: [],
  };
  return recipe;
};
export const encodeIngredient = (
  ingredient: Pick<SectionIngredient, "grams" | "info" | "adjective"> | any
): string => {
  const { info, grams, adjective, amount, unit } = ingredient;
  let res = grams + "g";
  if (amount !== 0 && unit !== "") res += ` (${amount} ${unit})`;
  res += ` ${info.name}`;
  if (adjective !== "") res += `, ${adjective}`;
  return res;
};

export const encodeRecipe = (recipe: GetRecipeByUuidQuery["recipe"]): string =>
  recipe && recipe.sections
    ? recipe.sections
        .map((section) => section.ingredients.map(encodeIngredient).join("\n"))
        .join("\n")
    : "";
