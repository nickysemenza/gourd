import React from "react";
import { gql } from "apollo-boost";
import { useQuery } from "@apollo/react-hooks";

const Test: React.FC = () => {
  const QUERY = gql`
    {
      recipe(uuid: "9f089d13-ac16-41f2-a490-66352d181e7f") {
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

  const { loading, error, data } = useQuery(QUERY);

  return <pre>{JSON.stringify({ loading, error, data }, null, 2)}</pre>;
};

export default Test;
