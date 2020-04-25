import React, { useState } from "react";
import { useGetRecipeByUuidQuery } from "../generated/graphql";

import { Box, Button } from "rebass";
import { useParams } from "react-router-dom";
import RecipeTable from "../components/RecipeTable";
import Debug from "../components/Debug";
import RecipeCard from "../components/RecipeCard";

type override = {
  sectionID: number;
  ingredientID: number;
  value: number;
};
const RecipeDetail: React.FC = () => {
  let { uuid } = useParams();
  const { loading, error, data } = useGetRecipeByUuidQuery({
    variables: { uuid: uuid || "" },
  });
  const [multiplier, setMultiplier] = useState(1.0);
  const [override, setOverride] = useState<override>();
  const [edit, setEdit] = useState(false);
  const recipe = data?.recipe;
  if (error) {
    console.error({ error });
    return (
      <Box color="primary" fontSize={4}>
        {error.message}
      </Box>
    );
  }
  if (!recipe) return null;

  const updateIngredient = (
    sectionID: number,
    ingredientID: number,
    value: string
  ) => {
    const newValue = parseFloat(value.endsWith(".") ? value + "0" : value);
    console.log(newValue);
    setOverride({
      sectionID,
      ingredientID,
      value: newValue,
    });
    const o = recipe!.sections[sectionID]!.ingredients[ingredientID]!.grams;
    if (o && value) {
      setMultiplier(Math.round((newValue / o + Number.EPSILON) * 100) / 100);
    }
  };

  const getIngredientValue = (
    sectionID: number,
    ingredientID: number,
    value: number
  ) => {
    if (
      override?.ingredientID === ingredientID &&
      override.sectionID === sectionID
    )
      return override.value;
    return value * multiplier;
  };

  return (
    <div>
      <Button onClick={() => setMultiplier(1)}>Reset</Button>
      <Button onClick={() => setEdit(!edit)}>{edit ? "edit" : "view"}</Button>
      <RecipeCard recipe={recipe} />
      <RecipeTable
        updateIngredient={updateIngredient}
        recipe={recipe}
        getIngredientValue={getIngredientValue}
        edit={edit}
      />
      <Debug data={{ loading, error, data, multiplier, override }} />
    </div>
  );
};

export default RecipeDetail;
