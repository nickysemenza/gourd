import React from "react";
import { ActionMeta, ValueType } from "react-select";
import AsyncCreatableSelect from "react-select/async-creatable";
import { IngredientsApi, RecipesApi } from "../api/openapi-fetch";
import { getOpenapiFetchConfig } from "../config";
import { blankRecipeWrapper, blankIngredient } from "../util";
import { IngredientKind } from "./RecipeEditorUtils";

type Option = {
  label: string;
  value: string;
  kind: IngredientKind;
};

export const EntitySelector: React.FC<{
  value?: Option;
  onChange: (value: Option) => void;
  createKind?: IngredientKind;
  showKind?: IngredientKind[];
  placeholder?: string;
}> = ({
  value,
  onChange,
  createKind = "ingredient",
  placeholder = "Pick a Recipe/Ingredient...",
  showKind = ["ingredient", "recipe"],
}) => {
  const iApi = new IngredientsApi(getOpenapiFetchConfig());
  const rApi = new RecipesApi(getOpenapiFetchConfig());
  const loadOptions = async (
    inputValue: string,
    callback: (options: Option[]) => void
  ) => {
    const res = await iApi.search({ name: inputValue });
    const recipeOptions: Option[] = (res.recipes || [])
      .filter((r) => r.detail.is_latest_version)
      .map((r) => {
        return {
          label: `[r] ${r.detail.name} (v${r.detail.version})`,
          value: r.id,
          kind: "recipe",
        };
      });
    const ingredientOptions: Option[] = (res.ingredients || []).map((i) => {
      return { label: "[i]" + i.name, value: i.id, kind: "ingredient" };
    });
    callback([
      ...(showKind.includes("ingredient") ? ingredientOptions : []),
      ...(showKind.includes("recipe") ? recipeOptions : []),
    ]);
  };
  const onSelectChange = async (
    val: ValueType<Option, false>,
    action: ActionMeta<Option>
  ) => {
    console.log({ val, action });
    if (!val) return;

    switch (action.action) {
      case "select-option":
        console.log(`selected ${val?.label} (${val?.value})`);
        onChange(val);
        break;
      case "create-option":
        console.log(`creating ${createKind} ${val?.label}`);
        const name = val?.label || "";
        const newEntityId =
          createKind === "recipe"
            ? (
                await rApi.createRecipes({
                  recipeWrapper: blankRecipeWrapper(name),
                })
              ).detail.id
            : (
                await iApi.createIngredients({
                  ingredient: blankIngredient(name),
                })
              ).id;

        console.log(`created ${createKind} ${name} (${newEntityId})`);
        onChange({ label: val.label, value: newEntityId, kind: createKind });
        break;
    }
  };

  return (
    <div data-cy="name-input">
      <AsyncCreatableSelect
        placeholder={placeholder}
        classNamePrefix="react-select"
        loadOptions={loadOptions}
        value={value}
        formatCreateLabel={(val) => `create ${createKind}: ${val}`}
        // handleInputChange={(...a: any) => {
        //   console.log({ handle: a });
        // }}
        onChange={onSelectChange}
      />
    </div>
  );
};
