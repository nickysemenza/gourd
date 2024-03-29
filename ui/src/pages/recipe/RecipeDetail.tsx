import React, { useState, useEffect, useContext } from "react";
import queryString from "query-string";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import RecipeDetailTable, {
  UpdateIngredientProps,
} from "../../components/recipe/RecipeDetailTable";
import update, { Spec } from "immutability-helper";
import { useHotkeys } from "react-hotkeys-hook";
import {
  useGetRecipeById,
  useGetFoodsByIds,
  useListIngredients,
  fetchCreateRecipes,
} from "../../api/react-query/gourdApiComponents";
import {
  RecipeSource,
  SectionIngredient,
  Amount,
  RecipeWrapper,
  RecipeWrapperInput,
  SectionIngredientInput,
  SectionInstructionInput,
} from "../../api/react-query/gourdApiSchemas";
import { scaledRound } from "../../util/util";
import {
  calCalc,
  countTotals,
  extractIngredientID,
  flatIngredients,
  FoodsById,
  getFDCIds,
  IngDetailsById,
  isGram,
  Override,
  RecipeTweaks,
  updateRecipeName,
  updateRecipeSource,
} from "../../components/recipe/RecipeEditorUtils";
import {
  ButtonGroup,
  HideShowHOC,
  makeHideShowButton,
} from "../../components/ui/ButtonGroup";
import { Edit, Eye, Save, X } from "react-feather";
import { singular } from "pluralize";
import Nutrition from "../../components/Nutrition";
import { WasmContext } from "../../util/wasmContext";

import { RecipeLink } from "../../components/misc/Misc";
import { Alert } from "../../components/ui/Alert";
import ProgressiveImage from "../../components/ui/ProgressiveImage";
import Debug from "../../components/ui/Debug";
import { NYTView } from "../../components/recipe/NYTView";
import ErrorBoundary from "../../components/ui/ErrorBoundary";
import PageWrapper from "../../components/ui/PageWrapper";
import { formatTimeRange, getTotalDuration } from "../../util/time";
import { getAPIURL } from "../../util/urls";

const toInput = (r: RecipeWrapper): RecipeWrapperInput => {
  return {
    id: r.id,
    detail: {
      name: r.detail.name,
      serving_info: r.detail.serving_info,
      sources: r.detail.sources,
      tags: r.detail.tags,
      sections: r.detail.sections.map((s) => {
        const instructions: Array<SectionInstructionInput> = s.instructions.map(
          (i) => ({
            instruction: i.instruction,
          })
        );
        const ingredients: Array<SectionIngredientInput> = s.ingredients.map(
          (i) => ({
            name: i.ingredient?.ingredient.name || i.recipe?.name || undefined,
            amounts: i.amounts,
            kind: i.kind,
            target_id: i.ingredient?.ingredient.id || i.recipe?.id || "",
            original: i.original,
            optional: i.optional,
            adjective: i.adjective,
          })
        );
        return {
          ingredients,
          instructions,
        };
      }),
    },
  };
};
const RecipeDetail: React.FC = () => {
  const { id } = useParams() as { id?: string };
  const history = useNavigate();
  const { search } = useLocation();
  const values = queryString.parse(search) as { multiplier?: string };
  const { error, data } = useGetRecipeById({
    pathParams: { recipeId: id || "" },
  });

  const w = useContext(WasmContext);

  const [multiplier, setMultiplier] = useState(1.0);
  const [multiplierTouched, setMultiplierTouched] = useState(false);
  const [override, setOverride] = useState<Override>();
  const [edit, setEdit] = useState(false);
  const [recipe, setRecipe] = useState(data);
  const [showOriginalLine, setshowOriginalLine] = useState(false);
  const [showKcalDollars, setshowKcalDollars] = useState(false);
  const [showAltView, setshowAltView] = useState(false);

  if (recipe)
    console.log({ recipe, a: toInput(recipe as unknown as RecipeWrapper) });

  const tweaks: RecipeTweaks = { override, multiplier, edit };

  const fdc_ids = [...getFDCIds(recipe ? recipe.detail.sections : [])];
  const { data: foods } = useGetFoodsByIds(
    {
      queryParams: {
        fdc_id: fdc_ids,
      },
    },
    { enabled: fdc_ids.length > 0 }
  );

  const foo = [
    ...flatIngredients(recipe?.detail.sections || []).map(
      (i) => extractIngredientID(i) || ""
    ),
  ];
  const { data: IngredientWrappers } = useListIngredients(
    {
      queryParams: {
        ingredient_id: foo,
      },
    },
    { enabled: foo.length > 0 }
  );

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
      history(
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
    history(`/recipe/${id}?${queryString.stringify({ multiplier: m })}`);
  };

  //https://stackoverflow.com/a/26265095
  const hints: FoodsById = Object.assign(
    {},
    ...(foods?.foods || []).map((s) => ({ [s.wrapper.fdcId]: s }))
  );
  const ing_hints: IngDetailsById = Object.assign(
    {},
    ...(IngredientWrappers?.ingredients || []).map((s) => ({
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
      console.log({ recipe });
      const updated = await fetchCreateRecipes({
        body: toInput(recipe as unknown as RecipeWrapper),
      });
      // const updated = await post(recipe);
      setEdit(false);
      history(`/recipe/${updated.detail.id}?${queryString.stringify(values)}`);
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
          <p>{e.status}</p>
          <Debug data={e.payload} />
        </div>
      </div>
    );
  }

  if (!recipe || !w) return null;

  const { detail } = recipe;
  const { quantity, unit } = detail.serving_info;

  const updateIngredient = ({
    sectionID,
    ingredientID,
    subIndex,
    value,
    attr,
    amountIndex,
  }: UpdateIngredientProps) => {
    const foo: Spec<SectionIngredient> = {
      amounts: {
        $apply: (v: Amount[]) => {
          console.log("update", { v, value, attr, amountIndex });
          if (amountIndex !== undefined) {
            if (amountIndex >= 0) {
              switch (attr) {
                case "grams":
                case "amount":
                  v[amountIndex].value = parseFloat(value);
                  break;
                case "unit":
                  v[amountIndex].unit = value;
              }
            } else {
              const toAdd: Amount = { unit: "", value: 0 };
              switch (attr) {
                case "grams":
                case "amount":
                  toAdd.value = parseFloat(value);
                  if (attr === "grams") {
                    toAdd.unit = "grams";
                  }
                  break;
                case "unit":
                  toAdd.unit = value;
              }
              v.push(toAdd);
            }
          }
          return v;
        },
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
      const { amounts } =
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
            (newValue /
              (amounts.filter((x) => isGram(x) === (attr === "grams")).pop()
                ?.value || 0) +
              Number.EPSILON) *
              100
          ) / 100
        );
      }
    }
  };

  const totalDuration = getTotalDuration(w, detail.sections);

  const sourceTypes: (keyof RecipeSource)[] = ["url", "title", "page"];

  const { totalCal, ingredientsWithNutrients, totalNutrients } = calCalc(
    recipe.detail.sections,
    hints,
    multiplier
  );
  const {
    grams: totalGrams,
    cents: totalCents,
    kcal: totalKCal,
  } = countTotals(recipe.detail.sections, w, ing_hints);

  const newerVersion = recipe.other_versions
    ?.filter((r) => r.meta.is_latest_version)
    .pop();

  const latexURL = `${getAPIURL()}/recipes/${recipe.detail.id}/latex`;
  return (
    <PageWrapper title={recipe.detail.name}>
      <div className="divide-y">
        <div className="lg:flex lg:items-center lg:justify-between mb-2">
          <div>
            {edit ? (
              <input
                type="text"
                className="border-2 w-96"
                value={detail.name}
                onChange={(e) =>
                  setRecipe(updateRecipeName(recipe, e.target.value))
                }
              ></input>
            ) : (
              <div className="text-gray-900 flex">
                <h2 className="text-2xl font-bold leading-7 ">{detail.name}</h2>
                {!!detail.meta.version && (
                  <h4 className="text-small self-end pl-1">
                    version {detail.meta.version}
                  </h4>
                )}
              </div>
            )}

            <div className="flex flex-col">
              {detail.serving_info.unit !== "" && (
                <div className="text-sm text-gray-600">
                  Makes {detail.serving_info.quantity}{" "}
                  {detail.serving_info.unit}
                </div>
              )}
              <div className="text-sm">
                Takes {formatTimeRange(w, totalDuration)}
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
                          type="text"
                          className="border-2 w-96"
                          value={source[key]}
                          placeholder={key}
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
          <div className="flex flex-col self-end justify-end items-end gap-1">
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
            <ButtonGroup
              buttons={[
                makeHideShowButton(
                  showOriginalLine,
                  setshowOriginalLine,
                  "original"
                ),
                makeHideShowButton(
                  showKcalDollars,
                  setshowKcalDollars,
                  "kcal/dollars"
                ),
                makeHideShowButton(showAltView, setshowAltView, "view"),
              ]}
            />
          </div>
        </div>

        {newerVersion && (
          <Alert
            title="Newer version available"
            line={
              <p>
                newest version is
                <RecipeLink recipe={newerVersion} multiplier={multiplier} />,
                this is version{recipe.detail.meta.version}
              </p>
            }
          />
        )}
        <ErrorBoundary>
          {showAltView ? (
            <NYTView recipe={recipe} />
          ) : (
            <RecipeDetailTable
              hints={hints}
              ing_hints={ing_hints}
              tweaks={tweaks}
              updateIngredient={updateIngredient}
              recipe={recipe}
              setRecipe={setRecipe}
              showOriginalLine={showOriginalLine}
              showKcalDollars={showKcalDollars}
            />
          )}
        </ErrorBoundary>
        {/* <InstructionsListParser
        setDetail={(s) => {
          setRecipe(setDetail(recipe, s));
          setEdit(true);
        }}
      /> */}
        <div>
          <p className="text-lg font-bold">totals</p>
          <div>
            calories: {totalCal}
            {totalCal > 0 &&
              quantity > 0 &&
              ` (${scaledRound(totalCal / quantity)} per ${singular(unit)})`}
          </div>
          <div>
            kcal: {totalKCal}
            {totalKCal &&
              quantity > 0 &&
              ` (${scaledRound(totalKCal / quantity)} per ${singular(unit)})`}
          </div>
          <div>
            cents: {totalCents}
            {totalCents &&
              quantity > 0 &&
              ` (${scaledRound(totalCents / quantity)} per ${singular(unit)})`}
          </div>
          <div>
            grams: {totalGrams}
            {totalGrams &&
              quantity > 0 &&
              ` (${scaledRound(totalGrams / quantity)} per ${singular(unit)})`}
          </div>
        </div>
        <div>
          <p className="text-lg font-bold">meals</p>

          <HideShowHOC>
            <Nutrition
              items={[
                ...ingredientsWithNutrients,
                ...[{ ingredient: "total", nutrients: totalNutrients }],
              ]}
              h={[...totalNutrients.keys()]}
            />
          </HideShowHOC>
        </div>
        <div>
          <div className="w-9/12 flex ">
            {(recipe.linked_photos || []).map((p) => (
              <ProgressiveImage photo={p} className="w-52" key={p.id} />
            ))}
          </div>
        </div>
      </div>
      <div className="flex flex-row w-full">
        <div className="w-1/2">
          <h3>other versions</h3>
          <ul>
            {recipe.other_versions?.map((v) => (
              <li key={`${v.id}@`}>
                <RecipeLink recipe={v} />
              </li>
            ))}
          </ul>
          <p className="text-lg font-bold">raw</p>
          <HideShowHOC>
            <pre className="whitespace-pre-wrap">
              {w.encode_recipe_text(
                toInput(recipe as unknown as RecipeWrapper).detail
              )}
            </pre>
          </HideShowHOC>
        </div>
        <div className="w-1/2">
          <a href={latexURL} target="_blank" rel="noreferrer" className="link">
            open latex
          </a>
          <iframe
            src={latexURL + "#navpanes=0&toolbar=0"}
            width="90%"
            height="900px"
            title="pdf view"
          ></iframe>
        </div>
      </div>
    </PageWrapper>
  );
};

export default RecipeDetail;
