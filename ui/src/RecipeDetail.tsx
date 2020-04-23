import React from "react";
import { gql } from "apollo-boost";
import { useGetRecipeByUuidQuery } from "./generated/graphql";

import { Box } from "rebass";
import { useParams } from "react-router-dom";

const _ = gql`
  query getRecipeByUUID($uuid: String!) {
    recipe(uuid: $uuid) {
      uuid
      name
      total_minutes
      unit
      sections {
        minutes
        ingredients {
          uuid
          info {
            name
          }
          grams
        }
        instructions {
          instruction
          uuid
        }
      }
    }
  }
`;

const RecipeDetail: React.FC = () => {
  let { uuid } = useParams();
  const { loading, error, data } = useGetRecipeByUuidQuery({
    variables: { uuid: uuid || "" },
  });

  const foo = data?.recipe?.name;
  return (
    <div>
      {/* <Input value={uuid} onChange={(e) => setUUID(e.target.value)} /> */}
      <Box
        sx={{
          borderWidth: "1px",
          borderStyle: "solid",
          borderColor: "highlight",
        }}
      >
        <pre>{JSON.stringify({ loading, error, data, foo }, null, 2)}</pre>
      </Box>

      <Box
        sx={{
          borderWidth: "1px",
          borderStyle: "solid",
          borderColor: "muted",
          boxShadow: "0 0 16px rgba(0, 0, 0, .25)",
        }}
        bg="muted"
      >
        {data?.recipe?.sections.map((section) => (
          <Box
            sx={{
              display: "grid",
              gridTemplateColumns: "1fr 2fr 2fr",
              borderBottomWidth: "1px",
              borderBottomStyle: "solid",
              borderBottomColor: "muted",
            }}
          >
            <Box> {section?.minutes}</Box>
            <Box>
              <ul>
                {section?.ingredients.map((ingredient) => (
                  <li>
                    {ingredient?.grams} grams {ingredient?.info.name}{" "}
                  </li>
                ))}
              </ul>
            </Box>
            <Box>
              <ol>
                {section?.instructions.map((instruction) => (
                  <li>{instruction?.instruction} </li>
                ))}
              </ol>
            </Box>
          </Box>
        ))}
      </Box>
    </div>
  );
};

export default RecipeDetail;
