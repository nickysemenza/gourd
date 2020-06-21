import React, { useState } from "react";
import { Search, SearchProps, SearchResultData } from "semantic-ui-react";
import {
  useSearchIngredientsAndRecipesQuery,
  useCreateIngredientMutation,
  Ingredient,
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
    ingredient: Pick<Ingredient, "uuid" | "name">,
    kind: SectionIngredientKind
  ) => void;
  initial?: string;
}> = ({ callback, initial }) => {
  const [value, setValue] = useState(initial || "");
  const [createIngredientMutation] = useCreateIngredientMutation({
    variables: {
      name: value,
    },
  });

  const { data, loading } = useSearchIngredientsAndRecipesQuery({
    variables: {
      searchQuery: value, // value for 'searchQuery'
    },
  });
  const results: Results = {
    ingredients: {
      name: "ingredients",
      results: [
        ...(data?.ingredients.map((i) => ({
          title: i.name,
          uuid: i.uuid,
          kind: SectionIngredientKind.Ingredient,
        })) || []),
        {
          title: `${value} (create)`,
          uuid: "",
          kind: SectionIngredientKind.Ingredient,
        },
      ],
    },
    recipes: {
      name: "recipes",
      results:
        data?.recipes.map((i) => ({
          title: i.name,
          uuid: i.uuid,
          kind: SectionIngredientKind.Recipe,
        })) || [],
    },
  };

  const handleSearchChange = (
    event: React.MouseEvent<HTMLElement, MouseEvent>,
    data: SearchProps
  ) => {
    setValue(data.value || "");
  };
  const handleResultSelect = async (
    event: React.MouseEvent<HTMLDivElement, MouseEvent>,
    data: SearchResultData
  ) => {
    const selection = data.result as ResultItem;
    if (selection.uuid === "") {
      let res = (await createIngredientMutation()).data;
      if (res) {
        setValue(res.createIngredient.name);
        callback(res.createIngredient, SectionIngredientKind.Ingredient);
      }
    } else {
      setValue(selection.title);
      callback({ name: selection.title, uuid: selection.uuid }, selection.kind);
    }
  };
  return (
    <Search
      category
      loading={loading}
      onResultSelect={handleResultSelect}
      onSearchChange={handleSearchChange}
      results={results}
      value={value}
      size="mini"
    />
  );
};
export default IngredientSearch;
