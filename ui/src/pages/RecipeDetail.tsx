import React, { useState, useEffect } from "react";
import {
  useGetRecipeByUuidQuery,
  useUpdateRecipeMutation,
  RecipeInput,
  SectionIngredientKind,
  Ingredient,
} from "../generated/graphql";

import { useParams } from "react-router-dom";
import RecipeTable, { UpdateIngredientProps } from "../components/RecipeTable";
import Debug from "../components/Debug";
import { recipeToRecipeInput } from "../util";
import update from "immutability-helper";
import { useHotkeys } from "react-hotkeys-hook";
import { encodeRecipe } from "../parser";
import RecipeSource from "../components/RecipeSource";

type override = {
  sectionID: number;
  ingredientID: number;
  value: number;
};
const RecipeDetail: React.FC = () => {
  let { uuid } = useParams();
  const { loading, error, data, refetch } = useGetRecipeByUuidQuery({
    variables: { uuid: uuid || "" },
  });
  const [multiplier, setMultiplier] = useState(1.0);
  const [override, setOverride] = useState<override>();
  const [edit, setEdit] = useState(false);
  const [recipe, setRecipe] = useState(data?.recipe);

  const [recipeUpdate, setRecipeUpdate] = useState<RecipeInput>({
    name: "tmp",
    uuid: "tmp",
  });
  const [updateRecipeMutation, { error: saveError }] = useUpdateRecipeMutation({
    variables: {
      recipe: recipeUpdate,
    },
  });

  const resetMultiplier = () => setMultiplier(1);
  const toggleEdit = () => {
    resetMultiplier();
    setEdit(!edit);
  };
  const saveUpdate = async () => {
    await updateRecipeMutation();
    await refetch();
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
      setRecipe(data.recipe);
    }
  }, [data]);

  useEffect(() => {
    if (recipe) {
      setRecipeUpdate(recipeToRecipeInput(recipe));
    }
  }, [recipe]);

  const e = error || saveError;
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
    ingredient: Pick<Ingredient, "uuid" | "name">,
    kind: SectionIngredientKind
  ) => {
    const { uuid, name } = ingredient;
    setRecipe(
      update(recipe, {
        sections: {
          [sectionID]: {
            ingredients: {
              [ingredientID]: {
                info: {
                  $set:
                    kind === SectionIngredientKind.Recipe
                      ? {
                          uuid,
                          name,
                          __typename: "Recipe",
                        }
                      : {
                          uuid,
                          name,
                          __typename: "Ingredient",
                        },
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
                    //   uuid: "",
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

  const addInstruction = (sectionID: number) => {
    setRecipe(
      update(recipe, {
        sections: {
          [sectionID]: {
            instructions: {
              $push: [{ uuid: "x", instruction: "" }],
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
                  uuid: "x",
                  grams: 1,
                  kind: SectionIngredientKind.Ingredient,
                  info: { name: "", uuid: "", __typename: "Ingredient" },
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
          $push: [{ uuid: "", minutes: 0, ingredients: [], instructions: [] }],
        },
      })
    );
  };

  return (
    <div>
      <div className="lg:flex lg:items-center lg:justify-between mb-2 ">
        <div>
          <h2 className="text-2xl font-bold leading-7 text-gray-900">
            {recipe.name}
          </h2>

          <div className="flex">
            {recipe.source && (
              <div className="text-sm text-gray-600">
                <RecipeSource source={recipe.source} />
              </div>
            )}
            {recipe.unit !== "" && (
              <div className="text-sm text-gray-600">
                Makes x {recipe.unit}. {recipe.totalMinutes} minutes.
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
      {recipe.meals.map((snippet) => (
        <>
          {snippet.name}
          {/* {snippet.imageURLs.map((u) => (
            <Image src={u} size="small" />
          ))} */}
        </>
      ))}
      <Debug data={{ recipe, loading, error, data, multiplier, override }} />
    </div>
  );
};

export default RecipeDetail;
