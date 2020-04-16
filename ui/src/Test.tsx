import React, { useState } from "react";
import { gql } from "apollo-boost";
import { useGetRecipeByUuidQuery } from "./generated/graphql";

const RecipeByUUID = gql`
  query getRecipeByUUID($uuid: String!) {
    recipe(uuid: $uuid) {
      uuid
      name
      total_minutes
      unit
      sections {
        minutes
        ingredients {
          name
          grams
        }
        instructions {
          instruction
        }
      }
    }
  }
`;

const Test: React.FC = () => {
  const [uuid, setUUID] = useState("9f089d13-ac16-41f2-a490-66352d181e7f");
  const { loading, error, data } = useGetRecipeByUuidQuery({
    variables: { uuid: uuid },
  });

  const foo = data?.recipe?.name;
  return (
    <div>
      <input value={uuid} onChange={(e) => setUUID(e.target.value)} />
      <pre>{JSON.stringify({ loading, error, data, foo }, null, 2)}</pre>
    </div>
  );
};

export default Test;
