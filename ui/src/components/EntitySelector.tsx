import React from "react";
import { ActionMeta, Styles, ValueType } from "react-select";
import AsyncCreatableSelect from "react-select/async-creatable";
import { IngredientsApi, RecipesApi } from "../api/openapi-fetch";
import { getOpenapiFetchConfig } from "../config";
import { blankRecipeWrapperInput, blankIngredient } from "../util";
import { Pill2 } from "./Button";
import { IngredientKind } from "./RecipeEditorUtils";

type Option = {
  label: string;
  value: string;
  kind?: IngredientKind;
  fdc_id: number | undefined;
  rd?: string;
};

export const EntitySelector: React.FC<{
  value?: Option;
  onChange: (value: Option) => void;
  createKind?: IngredientKind;
  showKind?: IngredientKind[];
  placeholder?: string;
  tall?: boolean;
}> = ({
  value,
  onChange,
  createKind = "ingredient",
  placeholder,
  showKind = ["ingredient", "recipe"],
  tall = false,
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
          label: `${r.detail.name} (v${r.detail.version})`,
          value: r.id,
          kind: "recipe",
          fdc_id: undefined,
          rd: r.detail.id,
        };
      });
    const ingredientOptions: Option[] = (res.ingredients || []).map((i) => {
      return {
        label: i.name,
        value: i.id,
        kind: "ingredient",
        fdc_id: i.fdc_id,
      };
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
                  recipeWrapperInput: blankRecipeWrapperInput(name),
                })
              ).detail.id
            : (
                await iApi.createIngredients({
                  ingredient: blankIngredient(name),
                })
              ).id;

        console.log(`created ${createKind} ${name} (${newEntityId})`);
        onChange({
          label: val.label,
          value: newEntityId,
          kind: createKind,
          fdc_id: undefined,
        });
        break;
    }
  };

  const height = tall ? 40 : 24;
  const customStyles: Partial<Styles<Option, false>> = {
    // option: (provided, state) => ({
    //   ...provided,
    //   // borderBottom: "1px dotted pink",
    //   // color: state.isSelected ? "red" : "blue",
    //   padding: 20,
    //   width: "100%",
    // }),
    control: (base) => ({
      ...base,
      height: height,
      minHeight: height,
    }),
    // singleValue: (provided, state) => {
    //   const opacity = state.isDisabled ? 0.5 : 1;
    //   const transition = "opacity 300ms";

    //   return { ...provided, opacity, transition };
    // },
    valueContainer: (provided, state) => {
      return { ...provided, height: height, padding: "2px" };
    },
    dropdownIndicator: (provided, state) => {
      return { ...provided, padding: "2px" };
    },
    indicatorsContainer: (provided, state) => ({
      ...provided,
      height: height + "px",
    }),
    indicatorSeparator: (state) => ({
      display: "none",
    }),
  };

  return (
    <div data-cy="name-input">
      {/* @ts-ignore */}
      <AsyncCreatableSelect
        styles={customStyles}
        placeholder={placeholder || `pick a ${showKind.join(" or ")}`}
        classNamePrefix="react-select"
        loadOptions={loadOptions}
        value={value}
        formatOptionLabel={(option, meta) => (
          <div className="flex flex-row justify-between">
            <div className="text-orange-600 font-bold pr-1">{option.label}</div>
            {option.kind && <Pill2 color="green">{option.kind}</Pill2>}
          </div>
        )}
        formatCreateLabel={(val) => (
          <div className="flex flex-row">
            <div className="text-green-600 font-bold pr-1">
              create {createKind}:
            </div>
            <div> {val}</div>
          </div>
        )}
        // handleInputChange={(...a: any) => {
        //   console.log({ handle: a });
        // }}
        onChange={onSelectChange}
      />
    </div>
  );
};
