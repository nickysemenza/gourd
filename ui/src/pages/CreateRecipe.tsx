import React from "react";
import { EntitySelector } from "../components/EntitySelector";
import { useHistory } from "react-router-dom";

const CreateRecipe: React.FC = () => {
  let history = useHistory();
  return (
    <div>
      <EntitySelector
        createKind="recipe"
        showKind={[]}
        onChange={(a) => {
          history.push(`/recipe/${a.value}`);
          console.log(a);
        }}
      />
    </div>
  );
};

export default CreateRecipe;
