import React, { useState } from "react";
import { useCreateRecipeMutation } from "../generated/graphql";
import Debug from "../components/Debug";

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
    <form>
      <Debug data={{ name, data, loading, error }} />
      <input
        className="bg-white focus:outline-none focus:shadow-outline border border-gray-300 rounded-lg py-2 px-4 block w-full appearance-none leading-normal"
        placeholder="recipe name"
        data-cy="name-input"
        value={name}
        onChange={(e) => {
          setName(e.target.value);
        }}
      />
      <button
        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
        type="button"
        onClick={create}
      >
        Create Recipe
      </button>
    </form>
  );
};

export default CreateRecipe;
