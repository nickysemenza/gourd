import React, { useContext, useState } from "react";
import { PlusCircle } from "react-feather";
import { IngredientKind, SectionIngredient } from "../api/openapi-fetch";
import { blankIngredient } from "../util";
import { wasm, WasmContext } from "../wasm";
import { ButtonGroup } from "./Button";

const InstructionsListParser: React.FC<{
  setSectionIngredients: (si: SectionIngredient[][]) => void;
}> = ({ setSectionIngredients }) => {
  const instance = useContext(WasmContext);
  const [area, setArea] = useState("");

  const sections = area.split("\n\n");

  if (!instance) return null;

  // const a = foo(sections, instance, true);
  const b = foo(sections, instance, false);
  // const si =
  //   a.length > 0 &&
  //   a[0].length > 0 &&
  //   b.length > 0 &&
  //   b[0].length > 0 &&
  //   a[0][0].grams === 0 &&
  //   b[0][0].grams !== 0
  //     ? b
  //     : a;
  const si = b;
  return (
    <div>
      <div className="grid grid-cols-2 gap-4">
        <textarea
          className="border-2 border-blue-400"
          value={area}
          onChange={(e) => {
            setArea(e.target.value);
          }}
        />
        {/* <div> */}
        {/* {sectionLines.flat().map((line) => (
            <p>{instance.parse(line)}</p>
          ))} */}
        <ButtonGroup
          compact
          buttons={[
            {
              onClick: () => {
                setSectionIngredients(si);
              },
              text: "inject ingredients",
              IconLeft: PlusCircle,
            },
          ]}
        />
        {/* </div> */}
      </div>
    </div>
  );
};
export default InstructionsListParser;

const foo = (sections: string[], instance: wasm, doubles: boolean) => {
  const sectionLines = sections.map((section) =>
    section
      .split("\n")
      .map((l, x, t) => {
        if (doubles) {
          //todo: better way to determine if lines need to be joined - if [0] and  [1] fail to parse but together succeed
          return x % 2 === 0 ? `${t[x]} ${t[x + 1]}` : "";
        }
        return l;
      })
      .filter((x) => !!x && x !== "")
  );
  return sectionLines.map((section) => {
    return section.map((line) => {
      console.log(line);
      const res = instance.parse2(line);
      //todo: use wasm for parsing out grams
      const other = res.amounts
        .filter((amount) => amount.unit !== "g" && amount.unit !== "grams")
        .shift();
      const si: SectionIngredient = {
        id: "",
        grams:
          res.amounts
            .filter((amount) => amount.unit === "g" || amount.unit === "grams")
            .shift()?.value || 0,
        amount: other?.value || undefined,
        unit: other?.unit || undefined,
        kind: IngredientKind.INGREDIENT,
        ingredient: blankIngredient(res.name),
      };
      return si;
    });
  });
};
