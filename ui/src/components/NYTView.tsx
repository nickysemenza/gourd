import { useContext } from "react";
import { RecipeWrapper } from "../api/openapi-hooks/api";
import { WasmContext } from "../wasmContext";
import { getGlobalInstructionNumber } from "./RecipeEditorUtils";

export const NYTView: React.FC<{
  recipe: RecipeWrapper;
}> = ({ recipe }) => {
  const w = useContext(WasmContext);
  if (!w) return null;
  return (
    <div className="flex md:flex-row w-full flex-col pt-2">
      <div className="md:w-5/12">
        <div className="text-m font-serif font-bold text-black uppercase">
          Ingredients
        </div>
        {recipe.detail.sections.map((section) =>
          section.ingredients.map((i) => (
            <div className="py-1 flex flex-row justify-center">
              <div className="w-1/2 flex justify-end pr-1 font-light text-gray-600">
                {i.amounts
                  .filter((a) => a.unit !== "$" && a.unit !== "kcal")
                  .map((a) => w.format_amount(a))
                  .join(" / ")}
              </div>
              <div className="w-1/2">{i.ingredient?.ingredient.name}</div>
            </div>
          ))
        )}
      </div>
      <div className="md:w-9/12">
        <div className="text-m font-serif font-bold text-black uppercase">
          Instructions
        </div>
        {recipe.detail.sections.map((section, x) =>
          section.instructions.map((i, y) => (
            <div className="py-2">
              <div className="text-l font-bold">
                Step {getGlobalInstructionNumber(recipe, x, y)}
              </div>
              {i.instruction}
            </div>
          ))
        )}
      </div>
    </div>
  );
};
