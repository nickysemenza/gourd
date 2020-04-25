import React from "react";
import { Recipe } from "../generated/graphql";
import { Box, Flex, Text } from "rebass";
import { Input } from "@rebass/forms";
import { InputProps } from "theme-ui";

export interface TableProps {
  recipe: Partial<Recipe>;
  updateIngredient: (
    sectionID: number,
    ingredientID: number,
    value: string,
    attr: "grams" | "name"
  ) => void;
  updateInstruction: (
    sectionID: number,
    instructionID: number,
    value: string
  ) => void;

  getIngredientValue: (
    sectionID: number,
    ingredientID: number,
    value: number
  ) => number;
  edit: boolean;
  addInstruction: (sectionID: number) => void;
  addIngredient: (sectionID: number) => void;
}
const RecipeTable: React.FC<TableProps> = ({
  recipe,
  updateIngredient,
  updateInstruction,
  getIngredientValue,
  edit,
  addInstruction,
  addIngredient,
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
      <TableRow key={x}>
        <TableCell>{section.minutes}</TableCell>
        <TableCell>
          {section.ingredients.map((ingredient, y) => (
            <Flex>
              <TableInput
                data-cy="grams-input"
                edit={edit}
                softEdit
                value={getIngredientValue(x, y, ingredient.grams || 0)}
                onChange={(e) =>
                  updateIngredient(x, y, e.target.value, "grams")
                }
              />{" "}
              <Flex pl={1} width={1 / 2}>
                <Text pr={1} color="gray">
                  g
                </Text>
                <TableInput
                  width={"128px"}
                  edit={edit}
                  value={ingredient.info.name}
                  onChange={(e) =>
                    updateIngredient(x, y, e.target.value, "name")
                  }
                />
              </Flex>
            </Flex>
          ))}
          {edit && <Text onClick={() => addIngredient(x)}>add</Text>}
        </TableCell>
        <TableCell>
          <ol style={{ margin: 0 }}>
            {section.instructions.map((instruction, y) => (
              <li>
                <TableInput
                  width={"128px"}
                  edit={edit}
                  value={instruction.instruction}
                  onChange={(e) => updateInstruction(x, y, e.target.value)}
                />{" "}
              </li>
            ))}
          </ol>
          {edit && <Text onClick={() => addInstruction(x)}>add</Text>}
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

const TableInput: React.FC<{
  edit: boolean;
  softEdit?: boolean;
  value: string | number;
  width?: InputProps["width"];
  onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
}> = ({ edit, softEdit = false, width = "64px", ...props }) => (
  <Input
    {...props}
    padding={0}
    // value={value}
    // onChange={onChange}
    width={width}
    sx={{
      textAlign: softEdit ? "end" : "begin",
      ":not(:focus)": {
        borderColor: edit ? "text" : "transparent",
      },
      ":hover": {
        borderColor: softEdit ? "text" : "transparent",
        borderStyle: "dashed",
      },
      borderRadius: 0,
    }}
  />
);
