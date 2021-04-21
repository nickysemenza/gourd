import React, { useState, useEffect, useContext } from "react";
import queryString from "query-string";
import { useHistory, useLocation, useParams } from "react-router-dom";
import RecipeDetailTable, {
  UpdateIngredientProps,
} from "../components/RecipeDetailTable";
import update, { Spec } from "immutability-helper";
import { useHotkeys } from "react-hotkeys-hook";
import {
  useGetRecipeById,
  useCreateRecipes,
  RecipeSource,
  SectionIngredient,
  useGetFoodsByIds,
  useListIngredients,
} from "../api/openapi-hooks/api";
import { formatTimeRange, scaledRound, sumTimeRanges } from "../util";
import {
  calCalc,
  countTotalGrams,
  flatIngredients,
  FoodsById,
  getFDCIds,
  IngDetailsById,
  Override,
  RecipeTweaks,
  setDetail,
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
  const { search } = useLocation();
  const values = queryString.parse(search) as { multiplier?: string };
  const { error, data } = useGetRecipeById({
    recipe_id: id || "",
  });

  const w = useContext(WasmContext);

  const [multiplier, setMultiplier] = useState(1.0);
  const [multiplierTouched, setMultiplierTouched] = useState(false);
  const [override, setOverride] = useState<Override>();
  const [edit, setEdit] = useState(false);
  const [recipe, setRecipe] = useState(data);

  const tweaks: RecipeTweaks = { override, multiplier, edit };
  const { mutate: post } = useCreateRecipes({
    onMutate: (_) => {
      // setRecipe(data);
    },
  });

  const { data: foods } = useGetFoodsByIds({
    queryParamStringifyOptions: { arrayFormat: "repeat" }, // https://github.com/contiamo/restful-react/issues/313
    queryParams: {
      fdc_id: [...getFDCIds(recipe ? recipe.detail.sections : []), 0],
    },
    // lazy: true,
  });

  const { data: ingredientDetails } = useListIngredients({
    queryParamStringifyOptions: { arrayFormat: "repeat" }, // https://github.com/contiamo/restful-react/issues/313
    queryParams: {
      ingredient_id: [
        ...flatIngredients(recipe?.detail.sections || []).map(
          (i) => i.ingredient?.same_as || i.ingredient?.id || ""
        ),
        "",
      ],
    },
    // lazy: true,
  });

  // on the (usually just first) load, seed multiplier state with the value from the URL
  // once the multipler has been 'touched', stop seeding from the URL,
  // and instead seed the URL based on state
  const multiplierParam =
    (!!values.multiplier && parseFloat(values.multiplier)) || 0;
  useEffect(() => {
    if (
      multiplierParam !== multiplier &&
      multiplierParam !== 0 &&
      !multiplierTouched
    ) {
      setMultiplier(multiplierParam);
      history.push(
        `/recipe/${id}?${queryString.stringify({
          multiplier: multiplierParam,
        })}`
      );
    }
  }, [id, history, multiplier, multiplierParam, multiplierTouched]);

  // wrapper to set touched, and sync the state
  const setMultiplierW = (m: number) => {
    setMultiplierTouched(true);
    setMultiplier(m);
    history.push(`/recipe/${id}?${queryString.stringify({ multiplier: m })}`);
  };

  //https://stackoverflow.com/a/26265095
  const hints: FoodsById = Object.assign(
    {},
    ...(foods?.foods || []).map((s) => ({ [s.fdc_id]: s }))
  );
  const ing_hints: IngDetailsById = Object.assign(
    {},
    ...(ingredientDetails?.ingredients || []).map((s) => ({
      [s.ingredient.id]: s,
    }))
  );

  const resetMultiplier = () => setMultiplierW(1);
  const toggleEdit = () => {
    resetMultiplier();
    setEdit(!edit);
  };
  const saveUpdate = async () => {
    if (recipe) {
      const updated = await post(recipe);
      setEdit(false);
      history.push(
        `/recipe/${updated.detail.id}?${queryString.stringify(values)}`
      );
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

  if (!recipe || !w) return null;

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

        setMultiplierW(
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

  const { totalCal, ingredientsWithNutrients, totalNutrients } = calCalc(
    recipe.detail.sections,
    hints,
    multiplier
  );
  const totalGrams = countTotalGrams(recipe.detail.sections, w, ing_hints);
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
        hints={hints}
        ing_hints={ing_hints}
        tweaks={tweaks}
        updateIngredient={updateIngredient}
        recipe={recipe}
        setRecipe={setRecipe}
      />
      <InstructionsListParser
        setDetail={(s) => {
          setRecipe(setDetail(recipe, s));
          setEdit(true);
        }}
      />
      <p className="text-lg font-bold">totals</p>
      <div>
        calories: {totalCal}
        {totalCal > 0 &&
          quantity > 0 &&
          ` (${scaledRound(totalCal / quantity)} per ${singular(unit)})`}
      </div>
      <div>
        grams: {totalGrams}
        {totalGrams > 0 &&
          quantity > 0 &&
          ` (${scaledRound(totalGrams / quantity)} per ${singular(unit)})`}
      </div>
      <p className="text-lg font-bold">raw</p>
      <pre>{w.encode_recipe_text(recipe.detail)}</pre>
      <p className="text-lg font-bold">meals</p>

      <Nutrition
        items={ingredientsWithNutrients}
        h={[...totalNutrients.keys()]}
      />
    </div>
  );
};

export default RecipeDetail;
