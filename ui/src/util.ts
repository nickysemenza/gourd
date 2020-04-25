import { Recipe, RecipeInput } from "./generated/graphql";

export const recipeToRecipeInput = (recipe: Partial<Recipe>): RecipeInput => {
  const { name, total_minutes, unit, uuid } = recipe;
  return {
    name: name + "a" || "",
    uuid: uuid || "",
    total_minutes,
    unit,
    sections: [],
  };
};
