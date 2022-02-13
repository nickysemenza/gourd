import { useContext } from "react";
import { Link } from "react-router-dom";
import { Photo, RecipeWrapper } from "../api/openapi-hooks/api";
import { formatTimeRange, getTotalDuration } from "../util";
import { WasmContext } from "../wasm";
import ProgressiveImage from "./ProgressiveImage";
import { sumIngredients } from "./RecipeEditorUtils";

export const RecipeGrid: React.FC<{
  recipes: RecipeWrapper[];
}> = ({ recipes }) => {
  const w = useContext(WasmContext);
  if (!w) {
    return null;
  }
  return (
    <div className="grid gap-5 row-gap-5 mb-8 lg:grid-cols-6 md:grid-cols-4 sm:grid-cols-2">
      {recipes.map((recipe) => {
        let lp = recipe.linked_photos || [];
        let rs = (recipe.detail.sources || []).filter(
          (s) => s.image_url !== undefined
        );
        let photo: Photo | undefined =
          lp.length > 0
            ? lp[0]
            : rs.length > 0
            ? ({ base_url: rs[0].image_url || "" } as Photo)
            : undefined;

        const ing = Object.keys(
          sumIngredients(recipe.detail.sections).ingredients
        );
        const rec = Object.keys(sumIngredients(recipe.detail.sections).recipes);

        const totalDuration = getTotalDuration(w, recipe.detail.sections);

        return (
          <Link
            key={`recipegrid-${recipe.detail.id}`}
            to={`/recipe/${recipe.detail.id}`}
            aria-label="View Item"
            className="inline-block overflow-hidden duration-300 transform bg-white dark:bg-stone-400 shadow hover:-translate-y-2"
          >
            <div className="flex flex-col h-full">
              {photo !== undefined ? (
                <ProgressiveImage
                  photo={photo}
                  className="object-cover w-full h-40"
                />
              ) : (
                <div className="object-cover w-full h-40 bg-gradient-to-r from-green-300 to-blue-300 dark:from-violet-800 dark:to-fuchsia-400" />
              )}
              <div className="flex-grow border dark:border-stone-400 border-t-0">
                <div className="p-5 flex flex-col justify-between h-full">
                  <div>
                    <h1 className="mb-2 font-semibold leading-5 text-m">
                      {recipe.detail.name}
                    </h1>
                    <p className="text-sm text-gray-900">
                      {ing.length} ingredients | {rec.length} recipes
                    </p>
                  </div>
                  {totalDuration.value > 0 && (
                    <div className="text-xs text-gray-700 ">
                      takes {formatTimeRange(w, totalDuration)}
                    </div>
                  )}
                  <ul className="list-disc">
                    {(recipe.linked_meals || []).map((m) => (
                      <li className="text-xs"> {m.name}</li>
                    ))}
                  </ul>
                </div>
              </div>
            </div>
          </Link>
        );
      })}
    </div>
  );
};
