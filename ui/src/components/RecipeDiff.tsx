import React, { useEffect, useMemo } from "react";
import {
  EntitySummary,
  IngredientKind,
  RecipeDetail,
  RecipesApi,
  UsageValue,
} from "../api/openapi-fetch";
import {
  SectionIngredient,
  useListIngredients,
  useGetRecipesByIds,
} from "../api/openapi-hooks/api";
import {
  flatIngredients,
  getGramsFromSI,
  IngDetailsById,
  totalFlourMass,
} from "./RecipeEditorUtils";
import { EntitySelector } from "./EntitySelector";
import { RecipeLink } from "./Misc";
import { scaledRound } from "../util";
import { useLocation } from "react-router-dom";
import queryString from "query-string";
import { getOpenapiFetchConfig } from "../config";
import Debug from "./Debug";

const RecipeDiff: React.FC<{ details: RecipeDetail[] }> = ({ details }) => {
  const loc = useLocation();
  const ids = useMemo(() => {
    const url = queryString.parse(loc.search).recipes;
    return url ? (Array.isArray(url) ? url : [url]) : [];
  }, [loc]);

  const { data } = useGetRecipesByIds({
    queryParamStringifyOptions: { arrayFormat: "repeat" }, // https://github.com/contiamo/restful-react/issues/313
    queryParams: {
      recipe_id: ids,
    },
    // lazy: true,
  });
  const [sums, setSums] = React.useState<UsageValue[]>([]);

  useEffect(() => {
    async function fetchMyAPI() {
      const rAPI = new RecipesApi(getOpenapiFetchConfig());
      let foo = await rAPI.sumRecipes({
        inlineObject: {
          inputs: ids.map((id) => {
            let foo: EntitySummary = {
              id: id,
              kind: IngredientKind.RECIPE,
              multiplier: 1.0,
              name: "",
            };
            return foo;
          }),
        },
      });
      setSums(foo.sums);
    }
    fetchMyAPI();
  }, [ids]);

  const recipes = data?.recipes || [];
  // let d1 = details[0];

  const recipesIngredients = recipes.map((r) =>
    flatIngredients(r.detail.sections)
  );
  const allIds = recipesIngredients
    .map((r) =>
      r.map(
        (si) => si.ingredient?.ingredient.parent || si.ingredient?.ingredient.id
      )
    )
    .flat()
    .filter(function (item, pos, a) {
      return a.indexOf(item) === pos && item !== undefined;
    });

  let byId: Record<string, (SectionIngredient | undefined)[]> = {};
  allIds.forEach((eachId) => {
    let res: (SectionIngredient | undefined)[] = [];
    recipesIngredients.forEach((r) => {
      let result: SectionIngredient | undefined = undefined;
      r.forEach((si) => {
        const id =
          si.ingredient?.ingredient.parent || si.ingredient?.ingredient.id;
        if (id === eachId) {
          result = si;
        }
      });
      res.push(result);
    });
    if (eachId) {
      byId[eachId] = res;
    }
  });

  const { data: ingredientDetails } = useListIngredients({
    queryParamStringifyOptions: { arrayFormat: "repeat" }, // https://github.com/contiamo/restful-react/issues/313
    queryParams: {
      ingredient_id: [...Object.keys(byId)],
    },
    // lazy: true,
  });

  const ing_hints: IngDetailsById = Object.assign(
    {},
    ...(ingredientDetails?.ingredients || []).map((s) => ({
      [s.ingredient.id]: s,
    }))
  );

  return (
    <div className="flex flex-col">
      <table className="table-auto border-collapse border-1 border-gray-500 w-full">
        <thead>
          <tr>
            <th className="border border-gray-400">ingredient</th>
            {ids.map((id, i) => (
              <th className="border border-gray-400" key={i}>
                <EntitySelector
                  showKind={["recipe"]}
                  placeholder={ids[i] || `"Pick a Recipe..."`}
                  onChange={async (a) => {
                    console.log(a);
                    // setIds(update(ids, { [i]: { $set: a.rd || "" } }));
                  }}
                />
              </th>
            ))}
          </tr>
          <th>recipe</th>
          {recipes.map((r, i) => (
            <th className="border border-gray-400" key={i}>
              <RecipeLink recipe={r.detail} />
            </th>
          ))}
        </thead>
        <tbody>
          {Object.keys(byId).map((eachId) => (
            <tr>
              <td className="border border-gray-400">
                {/* {eachId} */}
                {ing_hints[eachId]?.ingredient.name}
              </td>
              {byId[eachId].map((si, x) => {
                if (!si) {
                  return <td className="border border-gray-400"></td>;
                }
                const grams = getGramsFromSI(si) || 0;
                const bp = scaledRound(
                  (grams / totalFlourMass(recipes[x].detail.sections)) * 100
                );
                return <td className="border border-gray-400">{bp}%</td>;
              })}
            </tr>
          ))}
        </tbody>
      </table>
      <div className="">
        <ul className="list-disc list-outside pl-4">
          {sums.map((s, x) => (
            <li key={x}>
              {s.ing.name} (
              {s.sum.map((a) => `${a.value} ${a.unit}`).join(" + ")})<ul></ul>
              <ul className="list-disc list-outside pl-4">
                {s.ings.map((si, y) => (
                  <li>{si.required_by.map((b) => b.name).join(" <- ")}</li>
                ))}
              </ul>
            </li>
          ))}
        </ul>
      </div>
      <Debug data={sums} />
    </div>
  );
};

export default RecipeDiff;
