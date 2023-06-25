/* eslint-disable @typescript-eslint/no-unused-vars */
import React from "react";
import { ActionMeta, GroupBase, SingleValue, StylesConfig } from "react-select";
import AsyncCreatableSelect from "react-select/async-creatable";
import { IngredientKind } from "./recipe/RecipeEditorUtils";
import { blankIngredient, blankRecipeWrapperInput } from "../util/util";
import { Pill } from "./ui/Pill";
import {
  fetchCreateIngredients,
  fetchCreateRecipes,
  fetchSearch,
} from "../api/react-query/gourdApiComponents";

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
  const loadOptions = async (inputValue: string) => {
    const res = await fetchSearch({ queryParams: { name: inputValue } });
    const recipeOptions: Option[] = (res.recipes || [])
      .filter((r) => r.detail.meta.is_latest_version)
      .map((r) => {
        return {
          label: `${r.detail.name} (v${r.detail.meta.version})`,
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
    return [
      ...(showKind.includes("ingredient") ? ingredientOptions : []),
      ...(showKind.includes("recipe") ? recipeOptions : []),
    ];
  };

  const onSelectChange = async (
    val: SingleValue<Option>,
    action: ActionMeta<Option>
  ) => {
    console.log({ val, action });
    if (!val) return;

    switch (action.action) {
      case "select-option":
        console.log(`selected ${val?.label} (${val?.value})`);
        onChange(val);
        break;
      case "create-option": {
        console.log(`creating ${createKind} ${val?.label}`);
        const name = val?.label || "";
        const newEntityId =
          createKind === "recipe"
            ? (
                await fetchCreateRecipes({
                  body: blankRecipeWrapperInput(name),
                })
              ).detail.id
            : (
                await fetchCreateIngredients({
                  body: blankIngredient(name),
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
    }
  };

  const height = tall ? 40 : 24;
  const customStyles: StylesConfig<Option, false, GroupBase<Option>> = {
    control: (base) => ({
      ...base,
      height: height,
      minHeight: height,
    }),
    valueContainer: (provided, _state) => {
      return { ...provided, height: height, padding: "2px" };
    },
    dropdownIndicator: (provided, _state) => {
      return { ...provided, padding: "2px" };
    },
    indicatorsContainer: (provided, _state) => ({
      ...provided,
      height: height + "px",
    }),
    indicatorSeparator: (_state) => ({
      display: "none",
    }),
  };

  return (
    <div data-cy="name-input">
      <AsyncCreatableSelect
        styles={customStyles}
        placeholder={placeholder || `pick a ${showKind.join(" or ")}`}
        classNamePrefix="react-select"
        loadOptions={loadOptions}
        value={value}
        formatOptionLabel={(option, meta) => (
          <div className="flex flex-row justify-between">
            <div className="text-orange-600 font-bold pr-1">{option.label}</div>
            {option.kind && <Pill color="green">{option.kind}</Pill>}
          </div>
        )}
        onChange={onSelectChange}
        formatCreateLabel={(val) => (
          <div className="flex flex-row">
            <div className="text-green-600 font-bold pr-1">
              create {createKind}:
            </div>
            <div> {val}</div>
          </div>
        )}
      />
    </div>
  );
};
