import React from "react";
import {
  Recipe,
  Section,
  Ingredient,
  SectionIngredientKind,
} from "../generated/graphql";
import { Box, Text } from "rebass";
import { Input } from "@rebass/forms";
import { InputProps } from "theme-ui";
import IngredientSearch from "./IngredientSearch";
import { Link } from "react-router-dom";
export interface UpdateIngredientProps {
  sectionID: number;
  ingredientID: number;
  value: string;
  attr: "grams" | "name" | "amount" | "unit" | "adjective" | "optional";
}

export interface TableProps {
  recipe: Partial<Recipe>;
  updateIngredient: (i: UpdateIngredientProps) => void;
  updateIngredientInfo: (
    sectionID: number,
    ingredientID: number,
    ingredient: Pick<Ingredient, "uuid" | "name">,
    kind: SectionIngredientKind
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
  addSection: () => void;
}
const RecipeTable: React.FC<TableProps> = ({
  recipe,
  updateIngredient,
  updateIngredientInfo,
  updateInstruction,
  getIngredientValue,
  edit,
  addInstruction,
  addIngredient,
  addSection,
}) => {
  // for baker's percentage cauclation we need the total mass of all flours (which together are '100%')
  const flourMass = (recipe.sections || []).reduce(
    (acc, section) =>
      acc +
      section.ingredients
        .filter((item) => item.info.name.includes("flour"))
        .reduce((acc, ingredient) => acc + ingredient.grams, 0),
    0
  );

  const renderRow = (section: Section, x: number) => (
    <TableRow key={x}>
      <TableCell>
        <Box
          sx={{
            display: "inline-block",
            color: "highlight",
            bg: "primary",
            px: 1,
            py: 0,
            fontSize: 12,
            borderRadius: "50%",
          }}
        >
          {String.fromCharCode(65 + x)}
        </Box>
      </TableCell>
      <TableCell>{section.minutes}</TableCell>
      <TableCell>
        {section.ingredients.map((ingredient, y) => (
          <Box
            key={y}
            sx={{
              display: "grid",
              gridTemplateColumns: "1fr 70px 2fr 1fr 4fr 4fr",
              borderBottomWidth: "1px",
              borderBottomStyle: "solid",
              borderBottomColor: "green",
            }}
          >
            <TableInput
              data-cy="grams-input"
              edit={edit}
              softEdit
              value={getIngredientValue(x, y, ingredient.grams || 0)}
              onChange={(e) =>
                updateIngredient({
                  sectionID: x,
                  ingredientID: y,
                  value: e.target.value,
                  attr: "grams",
                })
              }
            />
            <Text pr={1} color="gray">
              g{" "}
              {flourMass > 0 && (
                <i>({Math.round((ingredient.grams / flourMass) * 100)}%)</i>
              )}
            </Text>
            {/* <TableInput
              data-cy="name-input"
              width={"128px"}
              edit={edit}
              value={ingredient.info.name}
              onChange={(e) =>
                updateIngredient({
                  sectionID: x,
                  ingredientID: y,
                  value: e.target.value,
                  attr: "name",
                })
              }
            /> */}
            {edit ? (
              <IngredientSearch
                initial={ingredient.info.name}
                callback={(item, kind) =>
                  // console.log({ item, kind })
                  updateIngredientInfo(x, y, item, kind)
                }
              />
            ) : (
              <Text pr={1} color="black">
                {ingredient.kind === SectionIngredientKind.Recipe ? (
                  <Link to={`/recipe/${ingredient.info.uuid}`}>
                    {ingredient.info.name}
                  </Link>
                ) : (
                  ingredient.info.name
                )}
              </Text>
            )}
            <TableInput
              data-cy="amount-input"
              width={"128px"}
              edit={edit}
              softEdit
              value={getIngredientValue(x, y, ingredient.amount || 0)}
              onChange={(e) =>
                updateIngredient({
                  sectionID: x,
                  ingredientID: y,
                  value: e.target.value,
                  attr: "amount",
                })
              }
            />
            <TableInput
              data-cy="unit-input"
              width={"64px"}
              edit={edit}
              value={ingredient.unit}
              onChange={(e) =>
                updateIngredient({
                  sectionID: x,
                  ingredientID: y,
                  value: e.target.value,
                  attr: "unit",
                })
              }
            />
            <TableInput
              data-cy="adjective-input"
              width={"128px"}
              edit={edit}
              value={ingredient.adjective}
              onChange={(e) =>
                updateIngredient({
                  sectionID: x,
                  ingredientID: y,
                  value: e.target.value,
                  attr: "adjective",
                })
              }
            />
            {/* TODO: optional toggle */}
          </Box>
        ))}
        {edit && <Text onClick={() => addIngredient(x)}>add ingredient</Text>}
      </TableCell>
      <TableCell>
        <ol style={{ margin: 0 }}>
          {section.instructions.map((instruction, y) => (
            <li key={y}>
              <TableInput
                data-cy="instruction-input"
                width={"128px"}
                edit={edit}
                value={instruction.instruction}
                onChange={(e) => updateInstruction(x, y, e.target.value)}
              />{" "}
            </li>
          ))}
        </ol>
        {edit && <Text onClick={() => addInstruction(x)}>add instruction</Text>}
      </TableCell>
    </TableRow>
  );

  return (
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
        <TableCell>Section</TableCell>
        <TableCell>Minutes</TableCell>
        <TableCell>
          Ingredients: x grams (BP) of y (z units, modifier)
        </TableCell>
        <TableCell>Instructions</TableCell>
      </TableRow>
      {recipe.sections?.map((section, x) => renderRow(section, x))}
      {edit && <Text onClick={() => addSection()}>add section</Text>}
    </Box>
  );
};
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
      gridTemplateColumns: "1fr 1fr 2fr 2fr",
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
}> = ({ edit, softEdit = false, width = "64px", ...props }) =>
  edit || softEdit ? (
    <Input
      {...props}
      padding={0}
      // width={width}
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
  ) : (
    <Text>{props.value}</Text>
  );
