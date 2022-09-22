import React, { useState } from "react";
import { EntitySelector } from "../components/EntitySelector";
import { useNavigate } from "react-router-dom";
import { ButtonGroup } from "../components/Button";
import { RecipesApi } from "../api/openapi-fetch";
import { getOpenapiFetchConfig } from "../config";

const CreateRecipe: React.FC = () => {
  let history = useNavigate();
  const [ingredientName, setIngredientName] = useState("");
  return (
    <div className="w-1/2">
      <EntitySelector
        tall
        createKind="recipe"
        showKind={["recipe"]}
        onChange={(a) => {
          history(`/recipe/${a.value}`);
          console.log(a);
        }}
      />
      <form onSubmit={(e) => e.preventDefault()}>
        <input
          type="url"
          className="border-2 border-gray-300 w-64"
          value={ingredientName}
          onChange={(e) => {
            setIngredientName(e.target.value);
          }}
        />
        <ButtonGroup
          compact
          buttons={[
            {
              onClick: () => {
                const foo = async () => {
                  const bar = new RecipesApi(getOpenapiFetchConfig());
                  const recipe = await bar.scrapeRecipe({
                    scrapeRecipeRequest: { url: ingredientName },
                  });

                  console.log({ recipe });
                  history(`/recipe/${recipe.detail.id}`);
                };
                foo();
              },
              text: "scrape",
            },
          ]}
        />
      </form>
    </div>
  );
};

export default CreateRecipe;
