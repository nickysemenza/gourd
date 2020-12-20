import React, { useState, useEffect } from "react";

import { useParams } from "react-router-dom";
import RecipeTable, { UpdateIngredientProps } from "../components/RecipeTable";
import Debug from "../components/Debug";
import update from "immutability-helper";
import { useHotkeys } from "react-hotkeys-hook";
import { encodeRecipe } from "../parser";
import {
  useGetRecipeById,
  Ingredient,
  useCreateRecipes,
} from "../api/openapi-hooks/api";

type override = {
  sectionID: number;
  ingredientID: number;
  value: number;
};
const RecipeDetail: React.FC = () => {
  let { uuid } = useParams() as { uuid?: string };

  const { loading, error, data } = useGetRecipeById({
    recipe_id: uuid || "",
  });

  const [multiplier, setMultiplier] = useState(1.0);
  const [override, setOverride] = useState<override>();
  const [edit, setEdit] = useState(false);
  const [recipe, setRecipe] = useState(data);

  const { mutate: post } = useCreateRecipes({
    onMutate: (_, data) => {
      // setRecipe(data);
    },
  });

  const resetMultiplier = () => setMultiplier(1);
  const toggleEdit = () => {
    resetMultiplier();
    setEdit(!edit);
  };
  const saveUpdate = async () => {
    recipe && post(recipe);
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
    if (data?.recipe) {
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

  const updateIngredientInfo = (
    sectionID: number,
    ingredientID: number,
    ingredient: Ingredient,
    kind: "recipe" | "ingredient"
  ) => {
    const { id, name } = ingredient;
    setRecipe(
      update(recipe, {
        sections: {
          [sectionID]: {
            ingredients: {
              [ingredientID]: {
                recipe: {
                  $set:
                    kind === "recipe"
                      ? { id, name, quantity: 0, unit: "" }
                      : undefined,
                },
                ingredient: {
                  $set: kind === "ingredient" ? { id, name } : undefined,
                },
                kind: { $set: kind },
              },
            },
          },
        },
      })
    );
  };
  const updateIngredient = ({
    sectionID,
    ingredientID,
    value,
    attr,
  }: UpdateIngredientProps) => {
    const newValue = parseFloat(value.endsWith(".") ? value + "0" : value);
    attr === "grams" &&
      !edit &&
      setOverride({
        sectionID,
        ingredientID,
        value: newValue,
      });
    const { grams } = recipe.sections[sectionID]!.ingredients[ingredientID];
    edit
      ? setRecipe(
          update(recipe, {
            sections: {
              [sectionID]: {
                ingredients: {
                  [ingredientID]: {
                    // info: {
                    //   id: "",
                    //   name: {
                    //     $apply: (v) => (attr === "name" ? value : v),
                    //   },
                    // },
                    // yikes
                    grams: {
                      $apply: (v) => (attr === "grams" ? parseFloat(value) : v),
                    },
                    amount: {
                      $apply: (v) =>
                        attr === "amount" ? parseFloat(value) : v,
                    },
                    unit: {
                      $apply: (v) => (attr === "unit" ? value : v),
                    },
                    adjective: {
                      $apply: (v) => (attr === "adjective" ? value : v),
                    },
                    optional: {
                      $apply: (v) =>
                        attr === "optional" ? value === "true" : v,
                    },
                  },
                },
              },
            },
          })
        )
      : setMultiplier(
          Math.round((newValue / grams + Number.EPSILON) * 100) / 100
        );
  };

  const getIngredientValue = (
    sectionID: number,
    ingredientID: number,
    value: number
  ) => {
    if (
      override?.ingredientID === ingredientID &&
      override.sectionID === sectionID
    )
      return override.value;
    return value * multiplier;
  };

  const updateInstruction = (
    sectionID: number,
    instructionID: number,
    value: string
  ) => {
    edit &&
      setRecipe(
        update(recipe, {
          sections: {
            [sectionID]: {
              instructions: {
                [instructionID]: { instruction: { $set: value } },
              },
            },
          },
        })
      );
  };

  const updateRecipeName = (value: string) => {
    setRecipe(
      update(recipe, {
        recipe: { name: { $set: value } },
      })
    );
  };

  const addInstruction = (sectionID: number) => {
    setRecipe(
      update(recipe, {
        sections: {
          [sectionID]: {
            instructions: {
              $push: [{ id: "", instruction: "" }],
            },
          },
        },
      })
    );
  };

  const addIngredient = (sectionID: number) => {
    setRecipe(
      update(recipe, {
        sections: {
          [sectionID]: {
            ingredients: {
              $push: [
                {
                  id: "",
                  grams: 1,
                  kind: "ingredient",
                  // info: { name: "", id: "", __typename: "Ingredient" },
                  amount: 0,
                  unit: "",
                  adjective: "",
                  optional: false,
                },
              ],
            },
          },
        },
      })
    );
  };
  const addSection = () => {
    setRecipe(
      update(recipe, {
        sections: {
          $push: [{ id: "", minutes: 0, ingredients: [], instructions: [] }],
        },
      })
    );
  };
  const info = recipe.recipe;
  return (
    <div>
      <div className="lg:flex lg:items-center lg:justify-between mb-2 ">
        <div>
          {edit ? (
            <input
              className="border-2 w-96"
              value={info.name}
              onChange={(e) => updateRecipeName(e.target.value)}
            ></input>
          ) : (
            <h2 className="text-2xl font-bold leading-7 text-gray-900">
              {info.name}
            </h2>
          )}

          <div className="flex">
            {/* {info.source && (
              <div className="text-sm text-gray-600">
                <RecipeSource source={info.source} />
              </div>
            )} */}
            {info.unit !== "" && (
              <div className="text-sm text-gray-600">
                Makes x {info.unit}. {info.total_minutes} minutes.
              </div>
            )}
          </div>
        </div>
        <div className="inline-flex">
          <button
            onClick={resetMultiplier}
            className="bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-2 px-4 rounded-l"
          >
            Reset
          </button>
          <button
            onClick={saveUpdate}
            className="bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-2 px-4"
          >
            save
          </button>
          <button
            onClick={toggleEdit}
            className="bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-2 px-4 rounded-r"
          >
            {edit ? "view" : "edit"}
          </button>
        </div>
      </div>

      <RecipeTable
        updateIngredient={updateIngredient}
        updateIngredientInfo={updateIngredientInfo}
        recipe={recipe}
        getIngredientValue={getIngredientValue}
        edit={edit}
        addInstruction={addInstruction}
        addIngredient={addIngredient}
        updateInstruction={updateInstruction}
        addSection={addSection}
      />
      <h2>raw</h2>
      <pre>{encodeRecipe(recipe)}</pre>
      <h2>meals</h2>
      <Debug data={{ recipe, loading, error, data, multiplier, override }} />
    </div>
  );
};

export default RecipeDetail;
