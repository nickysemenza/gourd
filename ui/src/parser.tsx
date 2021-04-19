import { Ingredient } from "gourd_rs";
import { RecipeWrapper, SectionIngredient } from "./api/openapi-hooks/api";
import { getIngredient } from "./util";
import { wasm } from "./wasm";

export const encodeIngredient = (
  ingredient: SectionIngredient,
  format: (i: Ingredient) => string | undefined
): string => {
  const { grams, adjective, amount, unit } = ingredient;
  let i: Ingredient = {
    name: getIngredient(ingredient).name,
    modifier: adjective,
    amounts: [
      ...(grams !== 0 ? [{ value: grams, unit: "g" }] : []),
      ...(amount && amount !== 0 && unit ? [{ value: amount, unit }] : []),
    ],
  };

  return format(i) || "";
};

export const encodeRecipe = (recipe: RecipeWrapper, w: wasm): string =>
  recipe && recipe.detail.sections
    ? recipe.detail.sections
        .map(
          (section) =>
            `${section.ingredients
              .map((i) => encodeIngredient(i, w.format_ingredient))
              .join("\n")}\n${section.instructions
              .map((i) => ";" + i.instruction)
              .join("\n")}`
        )
        .join("\n\n")
    : "";
