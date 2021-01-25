import React from "react";
import { useParams } from "react-router-dom";

const IngredientDetail: React.FC = () => {
  let { id } = useParams() as { id?: string };
  return (
    <div>
      <h1>Ingredient</h1>
      <h2>{id}</h2>
    </div>
  );
};
export default IngredientDetail;
