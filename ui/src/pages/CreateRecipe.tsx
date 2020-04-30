import React, { useState } from "react";
import { useCreateRecipeMutation } from "../generated/graphql";
import Debug from "../components/Debug";
import { Box, Button } from "rebass";
import { Input } from "@rebass/forms";

const CreateRecipe: React.FC = () => {
  const [name, setName] = useState("");
  const [
    createRecipeMutation,
    { data, loading, error },
  ] = useCreateRecipeMutation({
    variables: {
      recipe: { name },
    },
  });

  const create = () => {
    createRecipeMutation();
  };
  return (
    <Box>
      <Debug data={{ data, loading, error }} />
      <Input
        data-cy="name-input"
        value={name}
        onChange={(e) => {
          setName(e.target.value);
        }}
      />
      <Button onClick={create}>Create Recipe</Button>
    </Box>
  );
};

export default CreateRecipe;
