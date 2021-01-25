import React from "react";
import { useParams } from "react-router-dom";
import { toast } from "react-toastify";
import { IngredientsApi } from "../api/openapi-fetch";
import { useGetIngredientById } from "../api/openapi-hooks/api";
import Debug from "../components/Debug";
import FoodSearch from "../components/FoodSearch";
import { getOpenapiFetchConfig } from "../config";

const IngredientDetail: React.FC = () => {
  let { id } = useParams() as { id?: string };

  const { loading, data } = useGetIngredientById({
    ingredient_id: id || "",
  });

  const iApi = new IngredientsApi(getOpenapiFetchConfig());

  const linkFoodToIngredient = async (ingredientId: string, fdcId: number) => {
    await iApi.associateFoodWithIngredient({ ingredientId, fdcId });
    toast.success(`linked ${ingredientId} to food ${fdcId}`);
  };

  if (loading) return <div>loading</div>;
  if (!data) return <div>not found</div>;
  const { ingredient, food } = data;
  return (
    <div>
      <h1>Ingredient</h1>
      <h2>{id}</h2>
      <div className="flex flex-row">
        <div>
          <Debug data={data} />
        </div>
        <div>
          <FoodSearch
            name={ingredient.name}
            limit={20}
            enableSearch
            highlightId={food?.fdc_id}
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
