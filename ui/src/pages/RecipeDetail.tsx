import React, { useState, useEffect, useContext } from "react";

import { useHistory, useParams } from "react-router-dom";
import RecipeDetailTable, {
  UpdateIngredientProps,
} from "../components/RecipeDetailTable";
import update, { Spec } from "immutability-helper";
import { useHotkeys } from "react-hotkeys-hook";
import { encodeRecipe } from "../parser";
import {
  useGetRecipeById,
  useCreateRecipes,
  RecipeSource,
  SectionIngredient,
} from "../api/openapi-hooks/api";
import { formatTimeRange, sumTimeRanges } from "../util";
import {
  getCalories,
  Override,
  RecipeTweaks,
  replaceIngredients,
  sumIngredients,
  updateRecipeName,
  updateRecipeSource,
} from "../components/RecipeEditorUtils";
import { ButtonGroup } from "../components/Button";
import { Edit, Eye, Save, X } from "react-feather";
import { singular } from "pluralize";
import Nutrition from "../components/Nutrition";
import { WasmContext } from "../wasm";
import InstructionsListParser from "../components/InstructionsListParser";

const RecipeDetail: React.FC = () => {
  let { id } = useParams() as { id?: string };
  let history = useHistory();
  const { error, data } = useGetRecipeById({
    recipe_id: id || "",
  });

  const w = useContext(WasmContext);

  const [multiplier, setMultiplier] = useState(1.0);
  const [override, setOverride] = useState<Override>();
  const [edit, setEdit] = useState(false);
  const [recipe, setRecipe] = useState(data);

  const tweaks: RecipeTweaks = { override, multiplier, edit };
  const { mutate: post } = useCreateRecipes({
    onMutate: (_) => {
      // setRecipe(data);
    },
  });

  const resetMultiplier = () => setMultiplier(1);
  const toggleEdit = () => {
    resetMultiplier();
    setEdit(!edit);
  };
  const saveUpdate = async () => {
    if (recipe) {
      const updated = await post(recipe);
      setEdit(false);
      history.push(`/recipe/${updated.detail.id}`);
    }
  };

  useHotkeys("e", () => {
    toggleEdit();
  });
  useHotkeys("r", () => {
    resetMultiplier();
  });
  useHotkeys("s", () => {
    saveUpdate();
  });

  useEffect(() => {
    if (data?.detail) {
      setRecipe(data);
    }
  }, [data]);

  const e = error; // || saveError;
  if (e) {
    console.error({ e });
    // todo: extract to error component

    return (
      <div role="alert">
        <div className="bg-red-500 text-white font-bold rounded-t px-4 py-2">
          oops
        </div>
        <div className="border border-t-0 border-red-400 rounded-b bg-red-100 px-4 py-3 text-red-700">
          <p>{e.message}</p>
        </div>
      </div>
    );
  }

  if (!recipe) return null;

  const { detail } = recipe;
  const { quantity, unit } = detail;

  const updateIngredient = ({
    sectionID,
    ingredientID,
    subIndex,
    value,
    attr,
  }: UpdateIngredientProps) => {
    const foo: Spec<SectionIngredient> = {
      grams: {
        $apply: (v) => (attr === "grams" ? parseFloat(value) : v),
      },
      amount: {
        $apply: (v) => (attr === "amount" ? parseFloat(value) : v),
      },
      unit: {
        $apply: (v) => (attr === "unit" ? value : v),
      },
      adjective: {
        $apply: (v) => (attr === "adjective" ? value : v),
      },
      optional: {
        $apply: (v) => (attr === "optional" ? value === "true" : v),
      },
    };
    if (edit) {
      setRecipe(
        update(recipe, {
          detail: {
            sections: {
              [sectionID]: {
                ingredients: {
                  [ingredientID]:
                    subIndex === undefined
                      ? foo
                      : { substitutes: { [subIndex]: foo } },
                },
              },
            },
          },
        })
      );
    } else {
      const newValue = parseFloat(value.endsWith(".") ? value + "0" : value);
      const { grams, amount } =
        subIndex === undefined
          ? detail.sections[sectionID]!.ingredients[ingredientID]
          : (detail.sections[sectionID]!.ingredients[ingredientID]
              .substitutes || [])[subIndex];

      if (attr === "grams" || attr === "amount") {
        setOverride({
          sectionID,
          ingredientID,
          subIndex,
          value: newValue,
          attr,
        });

        setMultiplier(
          Math.round(
            (newValue / (attr === "grams" ? grams : amount || 0) +
              Number.EPSILON) *
              100
          ) / 100
        );
      }
    }
  };

  const totalDuration = sumTimeRanges(
    detail.sections.map((s) => s.duration).filter((t) => t !== undefined)
  );

  const sourceTypes: (keyof RecipeSource)[] = ["url", "title", "page"];

  console.group("nutrients");
  const ingredientsSum = sumIngredients(recipe.detail.sections);
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
      if (si.ingredient === undefined) return;
      const fdc_id = si.ingredient.fdc_id;
      if (fdc_id !== undefined) {
        const hint = recipe.food_hints && recipe.food_hints[fdc_id];
        if (hint !== undefined) {
          const scalingFactor = si.grams / 100;
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
            ingredient: si.ingredient.name,
            nutrients: ingNutrients,
          });
          totalCal += cal;
          console.log(
            `${si.ingredient.name}: ${si.grams}g = ${scalingFactor}x of ${hint.description}`,
            cal
          );
        }
      }
    });
  });
  console.log("TOTAL", totalCal);
  console.log("foo", totalNutrients, ingredientsWithNutrients);
  console.groupEnd();

  const totalGrams = recipe.detail.sections
    .map((section) => section.ingredients.map((ingredient) => ingredient))
    .flat()
    .map((si) => si.grams)
    .reduce((a, b) => a + b, 0);

  return (
    <div>
      <div className="lg:flex lg:items-center lg:justify-between mb-2 ">
        <div>
          {edit ? (
            <input
              className="border-2 w-96"
              value={detail.name}
              onChange={(e) =>
                setRecipe(updateRecipeName(recipe, e.target.value))
              }
            ></input>
          ) : (
            <div className="text-gray-900 flex">
              <h2 className="text-2xl font-bold leading-7 ">{detail.name}</h2>
              {!!detail.version && (
                <h4 className="text-small self-end pl-1">
                  version {detail.version}
                </h4>
              )}
            </div>
          )}

          <div className="flex flex-col">
            {detail.unit !== "" && (
              <div className="text-sm text-gray-600">
                Makes {detail.quantity} {detail.unit}
              </div>
            )}
            <div className="text-sm">
              Takes {formatTimeRange(totalDuration)}
            </div>
          </div>
          <div>
            {(detail.sources || []).map((source, i) => (
              <div className="flex text-gray-600 space-x-1" key={i}>
                <div className="text-xs font-bold uppercase self-center">
                  from:
                </div>
                {edit ? (
                  <div className="flex">
                    {sourceTypes.map((key) => (
                      <input
                        className="border-2 w-96"
                        value={source[key]}
                        onChange={(e) =>
                          setRecipe(
                            updateRecipeSource(recipe, i, e.target.value, key)
                          )
                        }
                      />
                    ))}
                  </div>
                ) : (
                  <div className="flex space-x-1">
                    {!!source.url && (
                      <a
                        href={source.url}
                        target="_blank"
                        rel="noreferrer"
                        className="text-indigo-600 underline"
                      >
                        {source.url}
                      </a>
                    )}
                    {!!source.title && <div>{source.title}</div>}
                    {!!source.page && <div>(pg. {source.page})</div>}
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
        <div className="self-start">
          <ButtonGroup
            buttons={[
              {
                onClick: resetMultiplier,
                disabled: multiplier === 1,
                text: "reset",
                IconLeft: X,
              },
              {
                onClick: saveUpdate,
                disabled: !edit,
                text: "save",
                IconLeft: Save,
              },
              {
                onClick: toggleEdit,
                text: edit ? "view" : "edit",
                IconLeft: edit ? Eye : Edit,
              },
            ]}
          />
        </div>
      </div>

      <RecipeDetailTable
        tweaks={tweaks}
        updateIngredient={updateIngredient}
        recipe={recipe}
        setRecipe={setRecipe}
      />
      <InstructionsListParser
        setSectionIngredients={(ings) => {
          setRecipe(replaceIngredients(recipe, ings));
          setEdit(true);
        }}
      />
      <h1>totals</h1>
      <div>
        calories: {totalCal}
        {quantity > 0 &&
          ` (${Math.round(totalCal / quantity)} per ${singular(unit)})`}
      </div>
      <div>
        grams: {totalGrams}
        {totalGrams > 0 &&
          ` (${Math.round(totalGrams / quantity)} per ${singular(unit)})`}
      </div>
      <h2>raw</h2>
      <pre>{w && encodeRecipe(recipe, w)}</pre>
      <h2>meals</h2>

      <Nutrition
        items={ingredientsWithNutrients}
        h={[...totalNutrients.keys()]}
      />
      {/* <Debug
        data={{
          loading,
          error,
          multiplier,
          override,
          recipe,
          foo: ingredientsSum,
        }}
      /> */}
    </div>
  );
};

export default RecipeDetail;
