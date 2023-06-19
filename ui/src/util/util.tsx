import {
  Ingredient,
  SectionIngredient,
  IngredientKind,
} from "../api/openapi-hooks/api";
import { RecipeWrapperInput } from "../api/openapi-fetch/models/RecipeWrapperInput";

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

export const blankRecipeWrapperInput = (
  name = "",
  id = ""
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
