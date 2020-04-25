import React, { useState } from "react";
import { useGetRecipeByUuidQuery } from "../generated/graphql";

import { Box, Flex, Button } from "rebass";
import { useParams } from "react-router-dom";
import { Input } from "theme-ui";

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
      <Box
        sx={{
          borderWidth: "1px",
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
        {recipe?.sections.map((section, x) => (
          <TableRow>
            <TableCell>{section?.minutes}</TableCell>
            <TableCell>
              <ul style={{ margin: 0 }}>
                {section?.ingredients.map((ingredient, y) => (
                  <li>
                    <Flex>
                      <Input
                        padding={0}
                        value={getIngredientValue(x, y, ingredient?.grams || 0)}
                        onChange={(e) => updateIngredient(x, y, e.target.value)}
                        sx={{
                          textAlign: "end",
                          width: "80px",
                          ":not(:focus)": {
                            // outline: "none",
                            // border: "none",
                            borderColor: edit ? "text" : "transparent",
                          },
                          ":hover": {
                            borderColor: "text",
                            borderStyle: "dashed",
                          },
                          borderRadius: 0,
                        }}
                      />{" "}
                      <Box pl={1} width={1 / 2}>
                        grams {ingredient?.info.name}
                      </Box>
                    </Flex>
                  </li>
                ))}
              </ul>
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
      <Box
        sx={{
          borderWidth: "1px",
          borderStyle: "solid",
          borderColor: "highlight",
        }}
      >
        <pre>
          {JSON.stringify(
            { loading, error, data, multiplier, override },
            null,
            2
          )}
        </pre>
      </Box>
    </div>
  );
};

export default RecipeDetail;

const TableCell: typeof Box = ({ children }) => <Box>{children}</Box>;
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
