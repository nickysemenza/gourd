import React, { useContext } from "react";
import { RecipeDetail } from "../api/openapi-fetch";
import { useGetRecipeById } from "../api/openapi-hooks/api";
import { WasmContext } from "../wasm";
import Debug from "./Debug";
import { sumIngredients } from "./RecipeEditorUtils";
import ReactDiffViewer from "react-diff-viewer";

const RecipeDiff: React.FC<{ details: RecipeDetail[] }> = ({ details }) => {
  const w = useContext(WasmContext);
  const { data: r1 } = useGetRecipeById({
    recipe_id: "rd_ea3ca186",
  });
  const { data: r2 } = useGetRecipeById({
    recipe_id: "rd_16eff17f",
  });
  if (!r1 || !r2 || !w) return <div />;
  // let d1 = details[0];
  return (
    <div>
      foo
      <Debug data={{ details }} />
      <pre>{w.encode_recipe_text(r1.detail)}</pre>
      <pre>{w.encode_recipe_text(r2.detail)}</pre>
      <ReactDiffViewer
        leftTitle={`${r1.detail.id} - v${r1.detail.version}`}
        rightTitle={`${r2.detail.id} - v${r2.detail.version}`}
        oldValue={w.encode_recipe_text(r1.detail)}
        newValue={w.encode_recipe_text(r2.detail)}
        splitView={true}
      />
    </div>
  );
};

export default RecipeDiff;
