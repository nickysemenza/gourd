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
  const w = useContext(WasmContext);

  if (!w) return <div />;

  let result: Amount | undefined = undefined;
  try {
    const ing = w.parse_amount(input);
    // let foo: UnitConversionRequest = {
    //   target: UnitConversionRequestTargetEnum.WEIGHT,
    //   unit_mappings: detail.unit_mappings,
    //   input: ing,
    // };
    // result = w.dolla(foo);
    result = try_convert(
      w,
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
  w: wasm,
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
  // let err;
  try {
    result = w.dolla(foo);
  } catch (e) {
    // console.error(e, { input, target, unit_mappings });
    // err = e;
  }
  // console.log("try_convert:" + msg, {
  //   input,
  //   target,
  //   unit_mappings,
  //   result,
  //   err,
  // });
  return result;
};
