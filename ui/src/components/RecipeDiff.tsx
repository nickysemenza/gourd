import React, { useContext } from "react";
import { RecipeDetail } from "../api/openapi-fetch";
import {
  SectionIngredient,
  useListIngredients,
  useGetRecipesByIds,
} from "../api/openapi-hooks/api";
import { WasmContext } from "../wasm";
import {
  flatIngredients,
  getGramsFromSI,
  inferGrams,
  IngDetailsById,
  totalFlourMass,
} from "./RecipeEditorUtils";
import { EntitySelector } from "./EntitySelector";
import { RecipeLink } from "./Misc";
import { scaledRound } from "../util";
import { useLocation } from "react-router-dom";
import queryString from "query-string";

const RecipeDiff: React.FC<{ details: RecipeDetail[] }> = ({ details }) => {
  const w = useContext(WasmContext);

  // const [ids, setIds] = useState(["rd_d041e06b", "rd_3f5f67df", ""]);
  const url = queryString.parse(useLocation().search).recipes;
  const ids = url ? (Array.isArray(url) ? url : [url]) : [];

  const { data } = useGetRecipesByIds({
    queryParamStringifyOptions: { arrayFormat: "repeat" }, // https://github.com/contiamo/restful-react/issues/313
    queryParams: {
      recipe_id: ids,
    },
    // lazy: true,
  });

  const recipes = data?.recipes || [];
  // let d1 = details[0];

  const recipesIngredients = recipes.map((r) =>
    flatIngredients(r.detail.sections)
  );
  const allIds = recipesIngredients
    .map((r) => r.map((si) => si.ingredient?.same_as || si.ingredient?.id))
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
        const id = si.ingredient?.same_as || si.ingredient?.id;
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
    <div className="flex">
      <table className="table-auto border-collapse border-1 border-gray-500 w-full">
        <thead>
          <tr>
            <th className="border border-gray-400">aa</th>
            {ids.map((id, i) => (
              <th className="border border-gray-400">
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
          <td>aa</td>
          {recipes.map((r) => (
            <td className="border border-gray-400">
              <RecipeLink recipe={r.detail} />
            </td>
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
                const grams =
                  getGramsFromSI(si) ||
                  (w && inferGrams(w, si, ing_hints)) ||
                  0;
                const bp = scaledRound(
                  (grams / totalFlourMass(recipes[x].detail.sections)) * 100
                );
                return <td className="border border-gray-400">{bp}%</td>;
              })}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default RecipeDiff;
