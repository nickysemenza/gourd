import React, { useState } from "react";
import Debug from "../components/Debug";
import { RecipesApi } from "../api/openapi-fetch";
import { getOpenapiFetchConfig } from "../config";

const CreateRecipe: React.FC = () => {
  const [name, setName] = useState("");
  const [resp, setResp] = useState<any>();

  const api = new RecipesApi(getOpenapiFetchConfig());

  const create = async () => {
    const r = await api.createRecipes({
      recipeDetail: {
        recipe: { name, id: "", quantity: 0, unit: "" },
        sections: [],
        id: "",
      },
    });
    setResp(r);
    // createRecipeMutation();
  };
  return (
    <form>
      <Debug data={{ name, resp }} />
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
