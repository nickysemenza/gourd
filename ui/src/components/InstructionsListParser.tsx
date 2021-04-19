import React, { useContext, useState } from "react";
import { PlusCircle } from "react-feather";
import { RecipeDetail } from "../api/openapi-fetch";
import { WasmContext } from "../wasm";
import { ButtonGroup } from "./Button";

const InstructionsListParser: React.FC<{
  setDetail: (d: RecipeDetail) => void;
}> = ({ setDetail }) => {
  const instance = useContext(WasmContext);
  const [area, setArea] = useState("");

  if (!instance) return null;

  const si = instance.decode_recipe_text(area) as RecipeDetail;
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

        <ButtonGroup
          compact
          buttons={[
            {
              onClick: () => {
                setDetail(si);
              },
              text: "inject ingredients",
              IconLeft: PlusCircle,
            },
          ]}
        />
      </div>
    </div>
  );
};
export default InstructionsListParser;
