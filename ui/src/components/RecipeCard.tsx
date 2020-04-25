import React from "react";
import { Recipe } from "../generated/graphql";
import { Box, Heading } from "rebass";

export interface Props {
  recipe: Partial<Recipe>;
}
const RecipeCard: React.FC<Props> = ({ recipe }) => (
  <Box>
    <Heading fontSize={[4, 5, 6]} color="secondary">
      {recipe.name}
    </Heading>
    Makes x {recipe.unit}. {recipe.total_minutes} minutes.
  </Box>
);
export default RecipeCard;
