import React, { useEffect } from "react";
import {
  EntitySummary,
  IngredientKind,
  RecipesApi,
  SumsResponse,
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

interface SIWithMultiplier {
  si: SectionIngredient | undefined;
  multiplier: number;
}
const RecipeDiffView: React.FC<{ entitiesToDiff: EntitySummary[] }> = ({
  entitiesToDiff,
}) => {
  const { data } = useGetRecipesByIds({
    queryParamStringifyOptions: { arrayFormat: "repeat" }, // https://github.com/contiamo/restful-react/issues/313
    queryParams: {
      recipe_id: entitiesToDiff.map((x) => x.id),
    },
    // lazy: true,
  });

  const [showBP, setShow] = React.useState(false);
  const [sumResp, setSumResp] = React.useState<SumsResponse>();

  useEffect(() => {
    async function fetchMyAPI() {
      const rAPI = new RecipesApi(getOpenapiFetchConfig());
      let recipeSumResp = await rAPI.sumRecipes({
        inlineObject: {
          inputs: entitiesToDiff.map((id, x) => {
            let foo: EntitySummary = {
              id: id.id,
              kind: IngredientKind.RECIPE,
              multiplier: id.multiplier,
              name: "",
            };
            return foo;
          }),
        },
      });
      setSumResp(recipeSumResp);
    }
    fetchMyAPI();
  }, [entitiesToDiff]);

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

  let sectionIngredientByID: Record<string, SIWithMultiplier[]> = {};
  allSIIDs.forEach((eachId) => {
    let res: SIWithMultiplier[] = [];
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
            {entitiesToDiff.map((id, i) => (
              <th className={thClass} key={`h-${i}`}>
                <EntitySelector
                  showKind={["recipe"]}
                  placeholder={entitiesToDiff[i].id || `"Pick a Recipe..."`}
                  onChange={async (a) => {
                    console.log(a);
                  }}
                />
              </th>
            ))}
          </tr>
          <tr>
            {recipes.map((r, i) => {
              const MULTIPLIER_TODO = entitiesToDiff[i].multiplier;
              return (
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
              );
            })}
          </tr>
        </thead>
        <tbody>
          {!ingredientDetails && (
            <tr>
              <td colSpan={entitiesToDiff.length + 2}>loading...</td>
            </tr>
          )}
          {Object.keys(sectionIngredientByID).map((sectionIngID) => {
            let s = sumResp?.sums || [];
            // sums that are related to the current ingredient, or a child of it.
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
                    {s.length === 0 && "loading..."}
                    {filterIngSums(s, sectionIngID, ing_hints).map((s, x) => (
                      <UsageValueShow
                        key={`${sectionIngID}-sums-${x}`}
                        uv={s}
                      />
                    ))}
                  </div>
                </td>
                {sectionIngredientByID[sectionIngID].map((si, columnIndex) => {
                  if (!si.si) {
                    return (
                      <td
                        className={`${tdClass} text-gray-500 bg-gray-100`}
                        key={`${columnIndex}-${sectionIngID}-nobp`}
                      >
                        &mdash;
                      </td>
                    );
                  }
                  const grams = getGramsFromSI(si.si) || 0;
                  const bpRaw =
                    (grams /
                      totalFlourMass(recipes[columnIndex].detail.sections)) *
                      100 || 0;
                  const bp = scaledRound(bpRaw);
                  return (
                    <td
                      className={`${tdClass} text-gray-500 
                    ${bpRaw === 100 && showBP ? "bg-green-100" : ""} 
                     `}
                      key={`${columnIndex}-${sectionIngID}-bp`}
                    >
                      <div className="flex content-start flex-col">
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
                              .map(
                                (a) => `${a.value * si.multiplier} ${a.unit}`
                              )
                              .join(" | ")}
                          </div>
                        </div>
                        <hr />
                        <div>
                          {sumResp &&
                            filterIngSums(
                              sumResp.by_recipe[recipes[columnIndex].detail.id],
                              sectionIngID,
                              ing_hints
                            ).map((x, i) => <UsageValueShow uv={x} key={i} />)}
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

const UsageValueShow: React.FC<{ uv: UsageValue }> = ({ uv }) => (
  <div>
    {uv.sum.map((a) => `${a.value} ${a.unit}`).join(" + ")}
    <ul className="list-disc list-outside pl-4">
      {uv.ings.map((iu) => (
        <li>
          <div className="flex">
            <div className="flex m-1">
              {iu.multiplier}×
              <div className="flex">
                {iu.amounts.map((a) => `${a.value} ${a.unit}`).join(" or ")}
              </div>
            </div>
            <div className="flex">
              [
              {iu.required_by.map((r, x) => (
                <div className="flex">
                  <RecipeLink
                    multiplier={r.multiplier}
                    recipe={{
                      name: r.name,
                      id: r.id,
                      version: -1,
                      is_latest_version: true,
                    }}
                  />
                  {iu.required_by.length - 1 > x && (
                    <div className="text-sm italic px-1">includes</div>
                  )}
                </div>
              ))}
              ]
            </div>
          </div>
        </li>
      ))}
    </ul>
    {/* <Debug data={uv} compact /> */}
  </div>
);

const filterIngSums = (
  s: UsageValue[],
  sectionIngID: string,
  ing_hints: IngDetailsById
) =>
  s.filter(
    (s) =>
      s.meta.id === sectionIngID ||
      ing_hints[sectionIngID]?.children
        ?.map((c) => c.ingredient.id)
        .includes(s.meta.id)
  );
