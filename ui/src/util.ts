import { Recipe, RecipeInput } from "./generated/graphql";

export const recipeToRecipeInput = (recipe: Partial<Recipe>): RecipeInput => {
  const { name, totalMinutes, unit, uuid, sections } = recipe;
  return {
    name: name + "a" || "",
    uuid: uuid || "",
    totalMinutes,
    unit,
    sections: (sections || []).map(
      ({ minutes, instructions, ingredients }) => ({
        minutes,
        instructions: instructions.map(({ instruction }) => ({ instruction })),
        ingredients: ingredients.map(
          ({ kind, info, grams, amount, unit, adjective, optional }) => ({
            infoUUID: info.uuid,
            kind,
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
