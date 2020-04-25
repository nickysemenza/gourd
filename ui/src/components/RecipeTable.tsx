import React from "react";
import { Recipe } from "../generated/graphql";
import { Box, Flex, Text } from "rebass";
import { Input } from "@rebass/forms";

export interface Props {
  recipe: Partial<Recipe>;
  updateIngredient: (
    sectionID: number,
    ingredientID: number,
    value: string
  ) => void;
  getIngredientValue: (
    sectionID: number,
    ingredientID: number,
    value: number
  ) => number;
  edit: boolean;
}
const RecipeTable: React.FC<Props> = ({
  recipe,
  updateIngredient,
  getIngredientValue,
  edit,
}) => (
  <Box
    sx={{
      borderWidth: 1,
      borderStyle: "solid",
      borderColor: "highlight",
      boxShadow: "0 0 16px rgba(0, 0, 0, .25)",
    }}
    bg="muted"
  >
    <TableRow>
      <TableCell>Minutes</TableCell>
      <TableCell>Ingredients</TableCell>
      <TableCell>Instructions</TableCell>
    </TableRow>
    {recipe.sections?.map((section, x) => (
      <TableRow>
        <TableCell>{section?.minutes}</TableCell>
        <TableCell>
          {/* <ul style={{ margin: 0 }}> */}
          {section?.ingredients.map((ingredient, y) => (
            <Flex>
              <Input
                padding={0}
                value={getIngredientValue(x, y, ingredient?.grams || 0)}
                onChange={(e) => updateIngredient(x, y, e.target.value)}
                width={"64px"}
                sx={{
                  textAlign: "end",
                  ":not(:focus)": {
                    borderColor: edit ? "text" : "transparent",
                  },
                  ":hover": {
                    borderColor: "text",
                    borderStyle: "dashed",
                  },
                  borderRadius: 0,
                }}
              />{" "}
              <Flex pl={1} width={1 / 2}>
                <Text pr={1} color="gray">
                  g
                </Text>
                {ingredient?.info.name}
              </Flex>
            </Flex>
          ))}
          {/* </ul> */}
        </TableCell>
        <TableCell>
          <ol style={{ margin: 0 }}>
            {section?.instructions.map((instruction) => (
              <li>{instruction?.instruction} </li>
            ))}
          </ol>
        </TableCell>
      </TableRow>
    ))}
  </Box>
);
export default RecipeTable;

const TableCell: typeof Box = ({ children }) => (
  <Box
    sx={{
      borderLeftWidth: "1px",
      borderLeftStyle: "solid",
      borderLeftColor: "highlight",
    }}
  >
    {children}
  </Box>
);
const TableRow: typeof Box = ({ children }) => (
  <Box
    sx={{
      display: "grid",
      gridTemplateColumns: "1fr 2fr 2fr",
      borderBottomWidth: "1px",
      borderBottomStyle: "solid",
      borderBottomColor: "highlight",
    }}
  >
    {children}
  </Box>
);
