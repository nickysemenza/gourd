import React, { useState, useEffect } from "react";

import { useParams } from "react-router-dom";
import RecipeDetailTable, {
  UpdateIngredientProps,
} from "../components/RecipeDetailTable";
import Debug from "../components/Debug";
import update from "immutability-helper";
import { useHotkeys } from "react-hotkeys-hook";
import { encodeRecipe } from "../parser";
import { useGetRecipeById, useCreateRecipes } from "../api/openapi-hooks/api";
import { formatTimeRange, sumTimeRanges } from "../util";
import {
  Override,
  RecipeTweaks,
  updateRecipeName,
} from "../components/RecipeEditorUtils";
import { ButtonGroup } from "../components/Button";
import { Edit, Eye, Save, X } from "react-feather";

const RecipeDetail: React.FC = () => {
  let { id } = useParams() as { id?: string };

  const { loading, error, data } = useGetRecipeById({
    recipe_id: id || "",
  });

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

  const updateIngredient = ({
    sectionID,
    ingredientID,
    value,
    attr,
  }: UpdateIngredientProps) => {
    if (edit) {
      setRecipe(
        update(recipe, {
          detail: {
            sections: {
              [sectionID]: {
                ingredients: {
                  [ingredientID]: {
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
          },
        })
      );
    } else {
      const newValue = parseFloat(value.endsWith(".") ? value + "0" : value);
      const { grams, amount } = detail.sections[sectionID]!.ingredients[
        ingredientID
      ];

      if (attr === "grams" || attr === "amount") {
        setOverride({
          sectionID,
          ingredientID,
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

  const totalDuration = sumTimeRanges(detail.sections.map((s) => s.duration));

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
            <h2 className="text-2xl font-bold leading-7 text-gray-900">
              {detail.name}
            </h2>
          )}

          <div className="flex">
            {detail.unit !== "" && (
              <div className="text-sm text-gray-600">Makes x {detail.unit}</div>
            )}
            {formatTimeRange(totalDuration)}
          </div>
        </div>
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

      <RecipeDetailTable
        tweaks={tweaks}
        updateIngredient={updateIngredient}
        recipe={recipe}
        setRecipe={setRecipe}
      />
      <h2>raw</h2>
      <pre>{encodeRecipe(recipe)}</pre>
      <h2>meals</h2>
      <Debug data={{ loading, error, multiplier, override, recipe }} />
    </div>
  );
};

export default RecipeDetail;
