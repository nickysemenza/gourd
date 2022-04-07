import React, { useEffect } from "react";
import {
  EntitySummary,
  IngredientKind,
  RecipesApi,
  UsageValue,
} from "../api/openapi-fetch";
import {
  SectionIngredient,
  useListIngredients,
  useGetRecipesByIds,
  RecipeDetail,
} from "../api/openapi-hooks/api";
import {
  flatIngredients,
  getMeasureUnitsFromSI,
  getGramsFromSI,
  IngDetailsById,
  totalFlourMass,
  extractIngredientID,
  getMultiplierFromRecipe,
} from "./RecipeEditorUtils";
import { EntitySelector } from "./EntitySelector";
import { RecipeLink } from "./Misc";
import { scaledRound } from "../util";
import { getOpenapiFetchConfig } from "../config";
import { HideShowButton } from "./Button";
import Debug from "./Debug";

interface Foo {
  si: SectionIngredient | undefined;
  multiplier: number;
}
const RecipeDiffView: React.FC<{ ids: string[] }> = ({ ids }) => {
  const { data } = useGetRecipesByIds({
    queryParamStringifyOptions: { arrayFormat: "repeat" }, // https://github.com/contiamo/restful-react/issues/313
    queryParams: {
      recipe_id: ids,
    },
    // lazy: true,
  });

  // const MULTIPLIER_TODO = 1.0;
  const MULTIPLIER_TODO = 0.5;

  const [showBP, setShow] = React.useState(false);
  const [sums, setSums] = React.useState<UsageValue[]>([]);

  useEffect(() => {
    async function fetchMyAPI() {
      const rAPI = new RecipesApi(getOpenapiFetchConfig());
      let recipeSumResp = await rAPI.sumRecipes({
        inlineObject: {
          inputs: ids.map((id) => {
            let foo: EntitySummary = {
              id: id,
              kind: IngredientKind.RECIPE,
              multiplier: MULTIPLIER_TODO,
              name: "",
            };
            return foo;
          }),
        },
      });
      setSums(recipeSumResp.sums);
    }
    fetchMyAPI();
  }, [ids]);

  const recipes = data?.recipes || [];
  // let d1 = details[0];

  const recipesIngredients = recipes.map((r) =>
    flatIngredients(r.detail.sections)
  );
  const allSIIDs = recipesIngredients
    .map((sectionIngredients) =>
      sectionIngredients.map((si) => extractIngredientID(si, true))
    )
    .flat()
    .filter(function (item, pos, a) {
      return a.indexOf(item) === pos && item !== undefined;
    });
  console.log({ recipesIngredients, allSIIDs });

  let sectionIngredientByID: Record<string, Foo[]> = {};
  allSIIDs.forEach((eachId) => {
    let res: Foo[] = [];
    recipesIngredients.forEach((r) => {
      let result: SectionIngredient | undefined = undefined;
      let multiplier = 1.0;
      r.forEach((si) => {
        if (si.kind === "recipe") {
          // if its a recipe, scale the child ingredients accordingly
          multiplier = getMultiplierFromRecipe(si, 1);
        }
        const ingredientID = extractIngredientID(si);
        if (ingredientID === eachId) {
          result = si;
        }
      });
      res.push({ si: result, multiplier });
    });
    if (eachId) {
      sectionIngredientByID[eachId] = res;
    }
  });

  const ingIds = [...Object.keys(sectionIngredientByID)];
  const { data: ingredientDetails } = useListIngredients({
    queryParamStringifyOptions: { arrayFormat: "repeat" }, // https://github.com/contiamo/restful-react/issues/313
    queryParams: {
      ingredient_id: ingIds,
      limit: ingIds.length || 0,
    },
    // lazy: true,
  });

  const ing_hints: IngDetailsById = Object.assign(
    {},
    ...(ingredientDetails?.ingredients || []).map((s) => ({
      [s.ingredient.id]: s,
    }))
  );

  const thClass =
    "border p-4 dark:border-dark-5 whitespace-nowrap font-normal text-gray-900";
  const tdClass = "border mx-4 px-2 py-1 dark:border-dark-5";
  return (
    <div className="flex flex-col mb-1 sm:mb-0 justify-between w-100">
      <h2 className="text-2xl leading-tight ">Recipe Diff View</h2>
      <h4 className="text-xs uppercase">comparing {recipes.length} recipes</h4>
      <HideShowButton show={showBP} setVal={setShow} />
      <table className="table-auto p-4 bg-white shadow rounded-lg w-full">
        <thead>
          <tr>
            <th rowSpan={2} className={thClass}>
              ingredient
            </th>
            <th rowSpan={2} className={thClass}>
              total
            </th>
            {ids.map((id, i) => (
              <th className={thClass} key={`h-${i}`}>
                <EntitySelector
                  showKind={["recipe"]}
                  placeholder={ids[i] || `"Pick a Recipe..."`}
                  onChange={async (a) => {
                    console.log(a);
                  }}
                />
              </th>
            ))}
          </tr>
          <tr>
            {recipes.map((r, i) => (
              <th className={thClass} key={i}>
                <RecipeLink recipe={r.detail} multiplier={MULTIPLIER_TODO} />
                <div className="">
                  {recipesIngredients[i]
                    .filter((i) => i.kind === "recipe")
                    .map((si) => (
                      <div className="text-xs" key={si.id}>
                        <div className="italic">includes</div>
                        <RecipeLink
                          recipe={si.recipe as unknown as RecipeDetail}
                          multiplier={getMultiplierFromRecipe(
                            si,
                            MULTIPLIER_TODO
                          )}
                        />
                      </div>
                    ))}
                </div>
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {!ingredientDetails && (
            <tr>
              <td colSpan={ids.length + 2}>loading...</td>
            </tr>
          )}
          {Object.keys(sectionIngredientByID).map((sectionIngID) => {
            // sums that are related to the current ingredient, or a child of it.
            const ingSums = sums.filter(
              (s) =>
                s.ing.id === sectionIngID ||
                ing_hints[sectionIngID]?.children
                  ?.map((c) => c.ingredient.id)
                  .includes(s.ing.id)
            );
            return (
              <tr key={sectionIngID} className="text-gray-700">
                <td className={tdClass} key={sectionIngID}>
                  {ing_hints[sectionIngID]?.ingredient.name}
                </td>
                <td
                  className={`${tdClass} text-gray-500 bg-gray-100`}
                  key={`${sectionIngID}-total`}
                >
                  <div>
                    {sums.length === 0 && "loading..."}
                    {ingSums.map((s, x) => (
                      <div key={`${sectionIngID}-sums-${x}`}>
                        {s.sum.map((a) => `${a.value} ${a.unit}`).join(" + ")}
                        <ul className="list-disc list-outside pl-4">
                          {s.ings.map((iu) => (
                            <li>
                              {/* {iu.required_by.map((r) => r.name).join(" via ")} */}
                              <div className="flex">
                                {iu.required_by.map((r) => (
                                  <div className="mx-1 flex">
                                    <RecipeLink
                                      multiplier={r.multiplier}
                                      recipe={{
                                        name: r.name,
                                        id: r.id,
                                        version: 0,
                                        is_latest_version: false,
                                      }}
                                    />
                                  </div>
                                ))}
                              </div>
                              <ul className="list-disc list-outside pl-4">
                                {iu.amounts.map((a) => (
                                  <li>
                                    {`${a.value} ${a.unit}`} @ {iu.multiplier}x
                                  </li>
                                ))}
                              </ul>
                            </li>
                          ))}
                        </ul>
                        {/* <Debug data={s} compact /> */}
                      </div>
                    ))}
                  </div>
                </td>
                {sectionIngredientByID[sectionIngID].map((si, x) => {
                  if (!si.si) {
                    return (
                      <td
                        className={`${tdClass} text-gray-500 bg-gray-100`}
                        key={`${x}-${sectionIngID}-nobp`}
                      >
                        &mdash;
                      </td>
                    );
                  }
                  const grams = getGramsFromSI(si.si) || 0;
                  const bpRaw =
                    (grams / totalFlourMass(recipes[x].detail.sections)) *
                      100 || 0;
                  const bp = scaledRound(bpRaw);
                  return (
                    <td
                      className={`${tdClass} text-gray-500 
                    ${bpRaw === 100 && showBP ? "bg-green-100" : ""}
                     `}
                      key={`${x}-${sectionIngID}-bp`}
                    >
                      <div className="flex justify-between">
                        {showBP && (
                          <div
                            className={`${
                              bpRaw === 0 ? "text-yellow-500" : ""
                            }`}
                          >
                            {bp}%
                          </div>
                        )}
                        <div>
                          {getMeasureUnitsFromSI(si.si)
                            .map((a) => `${a.value * si.multiplier} ${a.unit}`)
                            .join(" | ")}
                        </div>
                      </div>
                    </td>
                  );
                })}
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
};

export default RecipeDiffView;
