import { Recipe, RecipeInput } from "./generated/graphql";

export const recipeToRecipeInput = (recipe: Partial<Recipe>): RecipeInput => {
  const { name, total_minutes, unit, uuid, sections } = recipe;
  return {
    name: name + "a" || "",
    uuid: uuid || "",
    total_minutes,
    unit,
    sections: (sections || []).map(
      ({ minutes, instructions, ingredients }) => ({
        minutes,
        instructions: instructions.map(({ instruction }) => ({ instruction })),
        ingredients: ingredients.map(
          ({ info, grams, amount, unit, adjective, optional }) => ({
            name: info.name,
            grams,
            amount,
            unit,
            adjective,
            optional,
          })
        ),
      })
    ),
  };
};
