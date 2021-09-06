import {
  Food,
  Ingredient,
  IngredientDetail,
  Meal,
  RecipeDetail,
  RecipeSection,
  RecipeSource,
  RecipeWrapper,
  SectionIngredient,
  TimeRange,
} from "../api/openapi-hooks/api";
import update from "immutability-helper";
import { wasm } from "../wasm";
import { try_convert } from "./UnitConvertDemo";
import { Amount, UnitConversionRequestTargetEnum } from "../api/openapi-fetch";
export type Override = {
  sectionID: number;
  ingredientID: number;
  subIndex?: number;
  value: number;
  attr: IngredientAttr;
};

export type RecipeTweaks = {
  override?: Override;
  multiplier: number;
  edit: boolean;
};

export type IngredientAttr = "grams" | "amount";
export type IngredientKind = SectionIngredient["kind"];

export type FoodsById = {
  [key: number]: Food;
};
export type IngDetailsById = {
  [key: string]: IngredientDetail;
};

export const updateIngredientInfo = (
  recipe: RecipeWrapper,
  sectionID: number,
  ingredientID: number,
  ingredient: Ingredient,
  kind: IngredientKind
) => {
  const { id, name, fdc_id } = ingredient;
  return update(recipe, {
    detail: {
      sections: {
        [sectionID]: {
          ingredients: {
            [ingredientID]: {
              recipe: {
                $set:
                  kind === "recipe"
                    ? { id, name, quantity: 0, unit: "", sections: [] }
                    : undefined,
              },
              ingredient: {
                $set:
                  kind === "ingredient"
                    ? {
                        ingredient: { id, name, fdc_id },
                        recipes: [],
                        children: [],
                        unit_mappings: [],
                      }
                    : undefined,
              },
              kind: { $set: kind },
            },
          },
        },
      },
    },
  });
};

export const isOverride = (
  tweaks: RecipeTweaks,
  sectionID: number,
  ingredientID: number,
  subIndex: number | undefined,
  attr: IngredientAttr
) =>
  tweaks.override?.ingredientID === ingredientID &&
  tweaks.override.sectionID === sectionID &&
  tweaks.override.subIndex === subIndex &&
  tweaks.override.attr === attr;
export const adjustIngredientValue = (
  tweaks: RecipeTweaks,
  sectionID: number,
  ingredientID: number,
  subIndex: number | undefined,
  value: number,
  attr: IngredientAttr
) =>
  (isOverride(tweaks, sectionID, ingredientID, subIndex, attr) &&
    tweaks.override?.value) ||
  value * tweaks.multiplier;

export const updateInstruction = (
  recipe: RecipeWrapper,
  tweaks: RecipeTweaks,
  sectionID: number,
  instructionID: number,
  value: string
) =>
  tweaks.edit
    ? update(recipe, {
        detail: {
          sections: {
            [sectionID]: {
              instructions: {
                [instructionID]: { instruction: { $set: value } },
              },
            },
          },
        },
      })
    : recipe;

export const updateRecipeName = (recipe: RecipeWrapper, value: string) =>
  update(recipe, {
    detail: { name: { $set: value } },
  });

export const updateRecipeSource = (
  recipe: RecipeWrapper,
  index: number,
  value: string,
  field: keyof RecipeSource
) =>
  update(recipe, {
    detail: { sources: { [index]: { [field]: { $set: value } } } },
  });

export const addInstruction = (recipe: RecipeWrapper, sectionID: number) =>
  update(recipe, {
    detail: {
      sections: {
        [sectionID]: {
          instructions: {
            $push: [{ id: "", instruction: "" }],
          },
        },
      },
    },
  });

export const addIngredient = (recipe: RecipeWrapper, sectionID: number) =>
  update(recipe, {
    detail: {
      sections: {
        [sectionID]: {
          ingredients: {
            $push: [
              {
                id: "",
                kind: "ingredient",
                // info: { name: "", id: "", __typename: "Ingredient" },
                amounts: [
                  { unit: "grams", value: 0 },
                  { unit: "", value: 0 },
                ],
                adjective: "",
                optional: false,
              },
            ],
          },
        },
      },
    },
  });
export const addSection = (recipe: RecipeWrapper) =>
  update(recipe, {
    detail: {
      sections: {
        $push: [
          {
            id: "",
            duration: { min: 0, max: 0 },
            ingredients: [],
            instructions: [],
          },
        ],
      },
    },
  });
export const setDetail = (recipe: RecipeWrapper, detail: RecipeDetail) =>
  update(recipe, { detail: { $set: detail } });
export const updateTimeRange = (
  recipe: RecipeWrapper,
  sectionID: number,
  value?: TimeRange
) =>
  !!value
    ? update(recipe, {
        detail: {
          sections: {
            [sectionID]: {
              duration: {
                $set: value,
              },
            },
          },
        },
      })
    : recipe;

export type I = keyof Pick<RecipeSection, "ingredients" | "instructions">;
const calculateMoveI = (
  recipe: RecipeWrapper,
  sectionIndex: number,
  index: number,
  movingUp: boolean,
  i: I
) => {
  const { sections } = recipe.detail;

  const numI = sections[sectionIndex][i].length;
  const numSections = sections.length;
  const firstInSection = index === 0;
  const lastInSection = index === numI - 1;

  let newSectionIndex = sectionIndex;
  let newInIndex: number;
  if (firstInSection && movingUp) {
    // needs to go to prior section
    newSectionIndex--;
    if (newSectionIndex < 0) {
      // out of bounds
      return null;
    }
    newInIndex = sections[newSectionIndex][i].length;
  } else if (!firstInSection && movingUp) {
    // prior row in same section
    newInIndex = index - 1;
  } else if (lastInSection && !movingUp) {
    // needs to go to next section
    newSectionIndex++;
    if (newSectionIndex > numSections - 1) {
      // out of bounds
      return null;
    }
    newInIndex = 0;
  } else {
    // next row in same section
    newInIndex = index + 1;
  }

  return { newSectionIndex, newInIndex };
};
export const canMoveI = (
  recipe: RecipeWrapper,
  sectionIndex: number,
  index: number,
  movingUp: boolean,
  i: I
) => !!calculateMoveI(recipe, sectionIndex, index, movingUp, i);
export const moveI = (
  recipe: RecipeWrapper,
  sectionIndex: number,
  index: number,
  movingUp: boolean,
  i: I
) => {
  const coords = calculateMoveI(recipe, sectionIndex, index, movingUp, i);
  if (!coords) return recipe;
  const { newSectionIndex, newInIndex } = coords;
  console.log("moving!", {
    sectionIndex,
    newSectionIndex,
    index,
    newInIndex,
  });
  const target = recipe.detail.sections[sectionIndex][i][index];
  return update(recipe, {
    detail: {
      sections:
        newSectionIndex === sectionIndex
          ? {
              [sectionIndex]: {
                [i]: {
                  $splice: [
                    [index, 1],
                    [newInIndex, 0, target],
                  ],
                },
              },
            }
          : {
              [sectionIndex]: {
                [i]: {
                  $splice: [[index, 1]],
                },
              },
              [newSectionIndex]: {
                [i]: {
                  $splice: [[newInIndex, 0, target]],
                },
              },
            },
    },
  });
};
export const delI = (
  recipe: RecipeWrapper,
  sectionIndex: number,
  index: number,
  i: I
) =>
  update(recipe, {
    detail: {
      sections: {
        [sectionIndex]: {
          [i]: {
            $splice: [[index, 1]],
          },
        },
      },
    },
  });

export const pushMealRecipe = (
  meals: Meal[],
  mealIndex: number,
  recipe: RecipeDetail
) => {
  // debugger
  return update(meals, {
    [mealIndex]: {
      recipes: {
        $push: [{ multiplier: 1, recipe }],
      },
    },
  });
};
// returns the 1-indexed count of the instruction, across all sections.
export const getGlobalInstructionNumber = (
  recipe: RecipeWrapper,
  sectionIndex: number,
  instructionIndex: number
) =>
  recipe.detail.sections
    .slice(0, sectionIndex)
    .map((x) => x.instructions.length)
    .reduce((a, b) => a + b, 0) +
  instructionIndex +
  1;

export const flatIngredients = (
  sections: RecipeSection[]
): SectionIngredient[] =>
  sections
    .map((section) =>
      section.ingredients
        .map((ingredient) => [
          ingredient,
          ...flatIngredients(ingredient.recipe?.sections || []),
        ])
        .flat()
    )
    .flat();

export type Stats = {
  grams?: number;
  dollars?: number;
  kcal?: number;
};
const sumStats = (a: Stats, b: Stats) => {
  a.grams = (a.grams || 0) + (b.grams || 0);
  a.dollars = (a.dollars || 0) + (b.dollars || 0);
  a.kcal = (a.kcal || 0) + (b.kcal || 0);
  return a;
};
const scaleStats = (a: Stats, amount: number) => {
  a.grams = (a.grams || 0) * amount;
  a.dollars = (a.dollars || 0) * amount;
  a.kcal = (a.kcal || 0) * amount;
  return a;
};
export const getStats = (
  w: wasm,
  si: SectionIngredient,
  ing_hints: IngDetailsById,
  multiplier: number
): Stats => {
  const grams = getGrams(w, si, ing_hints);
  if (si.kind === "recipe") {
    console.log("foo", si);
    const total = countTotals(si.recipe?.sections || [], w, ing_hints);
    return scaleStats(
      total,
      total.grams && grams ? (grams / total.grams) * multiplier : 1
    );
  }
  return {
    grams,
    dollars: w_convert(
      w,
      si,
      ing_hints,
      multiplier,
      UnitConversionRequestTargetEnum.MONEY
    ),
    kcal: w_convert(
      w,
      si,
      ing_hints,
      multiplier,
      UnitConversionRequestTargetEnum.CALORIES
    ),
  };
};
export const countTotals = (
  sections: RecipeSection[],
  w: wasm,
  ing_hints: IngDetailsById
) =>
  flatIngredients(sections)
    .map((si) =>
      getStats(
        w,
        si,
        ing_hints,
        1 //todo
      )
    )
    .reduce((a, b) => sumStats(a, b), { grams: 0, dollars: 0 });

export const sumIngredients = (sections: RecipeSection[]) => {
  let recipes: Record<string, SectionIngredient[]> = {};
  let ingredients: Record<string, SectionIngredient[]> = {};

  flatIngredients(sections).forEach((i) => {
    switch (i.kind) {
      case "recipe":
        const r = i;
        if (!!r) {
          //todo: don't group by recipe/ingredient
          ingredients[r.id] = [...(ingredients[r.id] || []), r];
        }
        break;
      case "ingredient":
        const item = i;
        if (!!item) {
          ingredients[item.id] = [...(ingredients[item.id] || []), item];
        }
        break;
    }
  });

  return { recipes, ingredients };
};

export const getCalories = (food: Food) => {
  const first = food.nutrients.find((n) => n.nutrient.unit_name === "KCAL");
  return (!!first && first.amount) || 0;
};

export const calCalc = (
  sections: RecipeSection[],
  hints: FoodsById,
  multiplier: number
) => {
  console.group("nutrients");
  const ingredientsSum = sumIngredients(sections);
  const uniqIng = ingredientsSum.ingredients;
  let totalCal = 0;

  let ingredientsWithNutrients: Array<{
    ingredient: string;
    nutrients: Map<string, number>;
  }> = [];
  const totalNutrients = new Map<string, number>();
  // const foo = [];
  Object.keys(uniqIng).forEach((k) => {
    uniqIng[k].forEach((si) => {
      if (!!si.ingredient) {
        const fdc_id = si.ingredient.ingredient.fdc_id;
        if (fdc_id !== undefined) {
          const hint = hints[fdc_id];
          if (hint !== undefined) {
            const scalingFactor = (getGramsFromSI(si) / 100) * multiplier;
            const cal = getCalories(hint) * scalingFactor;
            const ingNutrients = new Map<string, number>();
            hint.nutrients.forEach((n) => {
              const { name, unit_name } = n.nutrient;
              const label = `${name} (${unit_name})`;
              if (n.amount <= 0) return;
              totalNutrients.set(
                label,
                n.amount * scalingFactor + (totalNutrients.get(label) || 0)
              );
              ingNutrients.set(label, n.amount * scalingFactor);
            });
            ingredientsWithNutrients.push({
              ingredient: si.ingredient.ingredient.name,
              nutrients: ingNutrients,
            });
            totalCal += cal;
            console.log(
              `${si.ingredient.ingredient.name}: ${getGramsFromSI(
                si
              )}g = ${scalingFactor}x of ${hint.description}`,
              cal
            );
          }
        }
      }
      if (!!si.recipe) {
        console.log("recursive calCalc");
        const {
          totalCal: totalCalSub,
          // ingredientsWithNutrients: ingredientsWithNutrientsSub,
          // totalNutrients: totalNutrientsSub,
          //todo: multiplier for recipes is wrong, needs to be based on amount of target recipe used
        } = calCalc(si.recipe.sections, hints, 0.1);
        totalCal += totalCalSub;
      }
    });
  });
  console.log("TOTAL", totalCal);
  console.log("foo", totalNutrients, ingredientsWithNutrients);
  console.groupEnd();
  return { totalCal, ingredientsWithNutrients, totalNutrients };
};

export const getFDCIds = (sections: RecipeSection[]): number[] =>
  sections
    .map((section) =>
      section.ingredients
        .map((ingredient) => {
          if (!!ingredient.ingredient) {
            return [ingredient.ingredient.ingredient.fdc_id || 0];
          } else if (!!ingredient.recipe) {
            return getFDCIds(ingredient.recipe.sections);
          } else {
            return [0];
          }
        })
        .flat()
        .filter((id) => id !== 0)
    )
    .flat();
export const getHint = (
  ingredient: SectionIngredient,
  ing_hints: IngDetailsById
): IngredientDetail | undefined => {
  const ingredientId =
    ingredient?.ingredient?.ingredient.same_as ||
    ingredient?.ingredient?.ingredient.id ||
    "";
  const hint = ing_hints[ingredientId];
  return !!hint ? hint : undefined;
};
// TODO: move this to backend
export const inferGrams = (
  w: wasm,
  ingredient: SectionIngredient,
  ing_hints: IngDetailsById
): number | undefined => {
  const hint = getHint(ingredient, ing_hints);
  if (!hint) return undefined;
  return (
    ingredient.amounts
      .filter((a) => !isGram(a))
      .map(({ unit, value }) => {
        if (!unit || !value || unit === "" || value === 0) return undefined;

        const targets = [
          UnitConversionRequestTargetEnum.WEIGHT,
          UnitConversionRequestTargetEnum.VOLUME,
        ];
        let res = targets.map(
          (t) =>
            try_convert(
              w,
              hint.unit_mappings,
              [{ unit: unit || "", value: value || 0 }],
              t,
              ingredient.ingredient?.ingredient.name
            )?.value
        );
        return res.find((x) => x !== undefined && x !== 0);
      })
      .filter((x) => x !== undefined)
      .pop() || 0
  );
};
export const isGram = (a: Amount) =>
  a.unit === "grams" || a.unit === "gram" || a.unit === "g";
export const getGramsFromSI = (si: SectionIngredient) =>
  (si.amounts.filter((a) => isGram(a)).pop() || { value: 0 }).value;
const getGrams = (
  w: wasm,
  ingredient: SectionIngredient,
  ing_hints: IngDetailsById
) => {
  const grams = getGramsFromSI(ingredient);
  return grams || inferGrams(w, ingredient, ing_hints);
};
const w_convert = (
  w: wasm,
  ingredient: SectionIngredient,
  ing_hints: IngDetailsById,
  multiplier: number,
  to:
    | UnitConversionRequestTargetEnum.CALORIES
    | UnitConversionRequestTargetEnum.MONEY
): number | undefined => {
  // console.info("w_convert", ingredient);
  const hint = getHint(ingredient, ing_hints);
  if (!hint) return undefined;
  const grams = getGrams(w, ingredient, ing_hints);
  if (grams === 0) return undefined;
  const val = try_convert(
    w,
    hint.unit_mappings,
    [{ unit: "grams", value: grams || 0 }],
    to,
    ingredient.ingredient?.ingredient.name
  )?.value;
  return val ? val * multiplier : undefined;
};

// for baker's percentage cauclation we need the total mass of all flours (which together are '100%')
export const totalFlourMass = (sections: RecipeSection[]) =>
  (sections || []).reduce(
    (acc, section) =>
      acc +
      section.ingredients
        .filter((item) =>
          item.ingredient?.ingredient.name.toLowerCase().includes("flour")
        )
        .reduce((acc, ingredient) => acc + getGramsFromSI(ingredient), 0),
    0
  );
