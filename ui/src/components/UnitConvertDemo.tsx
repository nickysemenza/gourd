import { Amount } from "gourd_rs";
import React, { useContext } from "react";
import {
  UnitConversionRequestTargetEnum,
  UnitMapping,
} from "../api/openapi-fetch";
import {
  IngredientDetail,
  UnitConversionRequest,
} from "../api/openapi-hooks/api";
import { wasm, WasmContext } from "../wasm";
import Debug from "./Debug";
import { TableInput } from "./Input";

type UnitConvertDemoProps = { detail: IngredientDetail };
export const UnitConvertDemo: React.FC<UnitConvertDemoProps> = ({ detail }) => {
  const [input, setInput] = React.useState("1 cup");
  const instance = useContext(WasmContext);

  if (!instance) return <div />;

  let result: Amount | undefined = undefined;
  try {
    const ing = instance.parse_amount(input);
    // let foo: UnitConversionRequest = {
    //   target: UnitConversionRequestTargetEnum.WEIGHT,
    //   unit_mappings: detail.unit_mappings,
    //   input: ing,
    // };
    // result = instance.dolla(foo);
    result = try_convert(
      instance,
      detail.unit_mappings,
      ing,
      UnitConversionRequestTargetEnum.WEIGHT
    );
    console.log("success");
  } catch (e) {
    console.error({ e });
  }

  return (
    <div>
      <Debug data={{ result }} />
      <TableInput
        data-cy="grams-input"
        edit={true}
        value={input}
        blur
        onChange={(e) => setInput(e)}
      />{" "}
      {result && `= ${result.value} ${result.unit}`}
    </div>
  );
};

export const try_convert = (
  instance: wasm,
  unit_mappings: UnitMapping[],
  input: Amount[],
  target: UnitConversionRequestTargetEnum,
  msg?: string
): Amount | undefined => {
  let foo: UnitConversionRequest = {
    target,
    unit_mappings,
    input,
  };
  let result: Amount | undefined = undefined;
  let err;
  try {
    result = instance.dolla(foo);
  } catch (e) {
    console.error(e);
    err = e;
  }
  console.log(msg, { input, target, unit_mappings, result, err });
  return result;
};
