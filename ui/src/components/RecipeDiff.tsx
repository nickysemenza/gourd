import React, { useContext } from "react";
import { RecipeDetail } from "../api/openapi-fetch";
import {
  SectionIngredient,
  useGetRecipeById,
  useListIngredients,
} from "../api/openapi-hooks/api";
import { WasmContext } from "../wasm";
import Debug from "./Debug";
import ReactDiffViewer from "react-diff-viewer";
import {
  flatIngredients,
  IngDetailsById,
  totalFlourMass,
} from "./RecipeEditorUtils";

const RecipeDiff: React.FC<{ details: RecipeDetail[] }> = ({ details }) => {
  const w = useContext(WasmContext);
  const { data: r1 } = useGetRecipeById({
    recipe_id: "rd_d041e06b",
  });
  const { data: r2 } = useGetRecipeById({
    recipe_id: "rd_3f5f67df",
  });
  const recipes = !r1 || !r2 || !w ? [] : [r1, r2];
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
    <div>
      DIFF
      <table className="table-auto border-collapse border-1 border-gray-500 w-full">
        <thead>
          <td>aa</td>
          {recipes.map((r) => (
            <td>{r.detail.name}</td>
          ))}
        </thead>
        <tbody>
          {Object.keys(byId).map((eachId) => (
            <tr>
              <td>
                {eachId} {ing_hints[eachId]?.ingredient.name}{" "}
              </td>
              {byId[eachId].map((r, x) => {
                if (!r) {
                  return null;
                }
                const bp = Math.round(
                  (r.grams / totalFlourMass(recipes[x].detail.sections)) * 100
                );
                return (
                  <td>
                    {r?.grams || "0"} BP: {bp}
                  </td>
                );
              })}
            </tr>
          ))}
        </tbody>
      </table>
      <Debug data={{ r1, r2, byId, ing_hints }} />
      {/* <pre>{w.encode_recipe_text(r1.detail)}</pre> */}
      {/* <pre>{w.encode_recipe_text(r2.detail)}</pre> */}
      {/* <ReactDiffViewer
        leftTitle={`${r1.detail.id} - v${r1.detail.version}`}
        rightTitle={`${r2.detail.id} - v${r2.detail.version}`}
        oldValue={w.encode_recipe_text(r1.detail)}
        newValue={w.encode_recipe_text(r2.detail)}
        splitView={true}
      /> */}
    </div>
  );
};

export default RecipeDiff;
