import React, { useState } from "react";
import AsyncCreatableSelect from "react-select/async-creatable";
import { Ingredient } from "../api/openapi-fetch";

import {
  useSearchIngredientsAndRecipesQuery,
  useCreateIngredientMutation,
  SectionIngredientKind,
} from "../generated/graphql";

export interface Results {
  ingredients: ResultType;
  recipes: ResultType;
}

export interface ResultItem {
  title: string;
  uuid: string;
  kind: SectionIngredientKind;
}
export interface ResultType {
  name: "ingredients" | "recipes";
  results: ResultItem[];
}

const IngredientSearch: React.FC<{
  callback: (
    ingredient: Pick<Ingredient, "id" | "name">,
    kind: SectionIngredientKind
  ) => void;
  initial?: string;
}> = ({ callback, initial }) => {
  const i = initial || "";
  const [value, setValue] = useState(i);

  const [v, setV] = useState<any>({ label: i });
  const [createIngredientMutation] = useCreateIngredientMutation({
    variables: {
      name: value,
      kind: SectionIngredientKind.Ingredient,
    },
  });

  const { data } = useSearchIngredientsAndRecipesQuery({
    variables: {
      searchQuery: value, // value for 'searchQuery'
    },
  });

  const handleCreate = async (inputValue: any) => {
    console.log("foo", inputValue);
    let res = (await createIngredientMutation()).data;
    if (res) {
      callback(
        { id: res.upsertIngredient, name: inputValue },
        SectionIngredientKind.Ingredient
      );
    }
  };

  const handleChange = async (newValue: any, actionMeta: any) => {
    console.group("Value Changed");
    console.log(newValue);
    console.log(`action: ${actionMeta.action}`);
    console.groupEnd();
    if (newValue.__isNew__) {
      let res = (await createIngredientMutation()).data;
      if (res) {
        callback(
          { id: res.upsertIngredient, name: newValue.label },
          SectionIngredientKind.Ingredient
        );
      }
    } else {
      callback({ name: newValue.label, id: newValue.uuid }, newValue.kind);
    }
    setV(newValue);
  };

  const loadOptions = (inputValue: string, callback: any) => {
    setValue(inputValue || "");

    callback([
      ...(data?.ingredients || []).map((i) => ({
        label: i.name,
        kind: SectionIngredientKind.Ingredient,
        uuid: i.uuid,
      })),
      ...(data?.recipes || []).map((i) => ({
        label: i.name + " (Recipe)",
        kind: SectionIngredientKind.Recipe,
        uuid: i.uuid,
      })),
    ]);
  };

  return (
    <div data-cy="name-input">
      <AsyncCreatableSelect
        onChange={handleChange}
        loadOptions={loadOptions}
        onCreateOption={handleCreate}
        value={v}
      />
    </div>
  );
};
export default IngredientSearch;
