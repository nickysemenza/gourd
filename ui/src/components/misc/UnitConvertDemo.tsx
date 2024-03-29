import React, { useContext } from "react";

import {
  IngredientWrapper,
  UnitConversionRequest,
  UnitMapping,
} from "../../api/react-query/gourdApiSchemas";
import { wasm, WasmContext, Amount } from "../../util/wasmContext";
import Debug from "../ui/Debug";
import { TableInput } from "../ui/Input";

type UnitConvertDemoProps = { detail: IngredientWrapper };
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
    result = try_convert(w, detail.unit_mappings, ing, "weight");
    console.log("success");
  } catch (e) {
    console.error({ e });
  }

  return (
    <div>
      <Debug data={{ result }} />
      <TableInput
        data-cy="grams-input"
        placeholder="grams"
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
  target: UnitConversionRequest["target"]
): Amount | undefined => {
  const foo: UnitConversionRequest = {
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
