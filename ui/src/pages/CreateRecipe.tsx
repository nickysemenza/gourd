import React, { useState } from "react";
import { EntitySelector } from "../components/EntitySelector";
import { useHistory } from "react-router-dom";
import { ButtonGroup } from "../components/Button";
import { RecipesApi } from "../api/openapi-fetch";
import { getOpenapiFetchConfig } from "../config";

const CreateRecipe: React.FC = () => {
  let history = useHistory();
  const [ingredientName, setIngredientName] = useState("");
  return (
    <div>
      <EntitySelector
        createKind="recipe"
        showKind={["recipe"]}
        onChange={(a) => {
          history.push(`/recipe/${a.value}`);
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
                    inlineObject1: { url: ingredientName },
                  });

                  console.log({ recipe });
                  history.push(`/recipe/${recipe.detail.id}`);
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
