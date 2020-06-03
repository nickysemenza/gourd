import React from "react";
import { Recipe } from "../generated/graphql";
import { Box, Heading } from "rebass";

export interface Props {
  recipe: Partial<Recipe>;
}
const RecipeCard: React.FC<Props> = ({ recipe }) => (
  <Box>
    <Heading fontSize={[4, 5, 6]} color="secondary" data-cy="recipe-name">
      {recipe.name}
    </Heading>
    Makes x {recipe.unit}. {recipe.totalMinutes} minutes.
  </Box>
);
export default RecipeCard;
