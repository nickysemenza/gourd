import { Ingredient, SectionIngredient, Recipe } from "./api/openapi-hooks/api";

export const getIngredient = (
  si: Partial<SectionIngredient>
): { name: "" } | Recipe | Ingredient => {
  if (si.recipe) {
    return si.recipe;
  } else if (si.ingredient) {
    return si.ingredient;
  }
  return { name: "" };
};
