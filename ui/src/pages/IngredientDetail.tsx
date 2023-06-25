import React from "react";
import { useParams } from "react-router-dom";
import { toast } from "react-toastify";
import {
  fetchAssociateFoodWithIngredient,
  useGetIngredientById,
} from "../api/react-query/gourdApiComponents";
import Debug from "../components/ui/Debug";
import FoodSearch from "../components/FoodSearch";
import { UnitConvertDemo } from "../components/misc/UnitConvertDemo";

const IngredientDetail: React.FC = () => {
  const { id } = useParams() as { id?: string };

  const { isLoading: loading, data } = useGetIngredientById({
    pathParams: { ingredientId: id || "" },
  });

  const linkFoodToIngredient = async (ingredientId: string, fdcId: number) => {
    await fetchAssociateFoodWithIngredient({
      pathParams: { ingredientId },
      queryParams: { fdc_id: fdcId },
    });
    toast.success(`linked ${ingredientId} to food ${fdcId}`);
  };

  if (loading) return <div>loading</div>;
  if (!data) return <div>not found</div>;
  const { ingredient, food } = data;
  return (
    <div>
      <div className=" flex flex-col my-1">
        <h1 className="text-2xl font-extrabold text-gray-900 tracking-tight">
          {ingredient.name}
        </h1>
        <span className="text-gray-700 text-sm">{id}</span>
      </div>
      <div className="flex flex-row">
        <div className="w-6/12">
          <UnitConvertDemo detail={data} />
          <Debug data={data} />
        </div>
        <div className="w-6/12">
          <FoodSearch
            name={ingredient.name}
            limit={20}
            enableSearch
            highlightId={food?.wrapper.fdcId}
            onLink={(fdcId: number) => {
              linkFoodToIngredient(ingredient.id, fdcId);
            }}
          />
        </div>
      </div>
    </div>
  );
};
export default IngredientDetail;
