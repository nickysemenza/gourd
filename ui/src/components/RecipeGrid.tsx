import { Link } from "react-router-dom";
import { RecipeWrapper } from "../api/openapi-hooks/api";
import { formatTimeRange, getTotalDuration } from "../util";
import ProgressiveImage from "./ProgressiveImage";
import { sumIngredients } from "./RecipeEditorUtils";

export const RecipeGrid: React.FC<{
  recipes: RecipeWrapper[];
}> = ({ recipes }) => (
  <div className="grid gap-5 row-gap-5 mb-8 lg:grid-cols-6 sm:grid-cols-2">
    {recipes.map((recipe) => {
      let lp = recipe.linked_photos || [];
      let photo = lp.length > 0 ? lp[0] : undefined;
      console.log({ photo });

      const ing = Object.keys(
        sumIngredients(recipe.detail.sections).ingredients
      );
      const rec = Object.keys(sumIngredients(recipe.detail.sections).recipes);

      const totalDuration = getTotalDuration(recipe.detail.sections);

      return (
        <Link
          key={`recipegrid-${recipe.detail.id}`}
          to={`/recipe/${recipe.detail.id}`}
          aria-label="View Item"
          className="inline-block overflow-hidden duration-300 transform bg-white dark:bg-stone-400 rounded shadow hover:-translate-y-2"
        >
          <div className="flex flex-col h-full">
            {photo !== undefined ? (
              <ProgressiveImage
                photo={photo}
                className="object-cover w-full h-48"
              />
            ) : (
              <div className="object-cover w-full h-48 bg-gradient-to-bl from-violet-400 to-fuchsia-200 dark:from-violet-800 dark:to-fuchsia-400" />
            )}
            <div className="flex-grow border dark:border-stone-400 border-t-0 rounded-b">
              <div className="p-5">
                <h6 className="mb-2 font-semibold leading-5 text-lg">
                  {recipe.detail.name}
                </h6>
                <p className="text-sm text-gray-900">
                  {ing.length} ingredients | {rec.length} recipes
                </p>
                <div className="text-sm">
                  Takes {formatTimeRange(totalDuration)}
                </div>
              </div>
            </div>
          </div>
        </Link>
      );
    })}
  </div>
);
