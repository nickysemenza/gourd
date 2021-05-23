import React, { useContext, useState } from "react";
import { PlusCircle } from "react-feather";
import { RecipeDetail } from "../api/openapi-fetch";
import { WasmContext } from "../wasm";
import { ButtonGroup } from "./Button";

const InstructionsListParser: React.FC<{
  setDetail: (d: RecipeDetail) => void;
}> = ({ setDetail }) => {
  const w = useContext(WasmContext);
  const [area, setArea] = useState("");

  if (!w) return null;

  const si = w.decode_recipe_text(area) as RecipeDetail;
  const si2 = w.encode_recipe_to_compact_json(si);
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
        <div>
          <div className="flex flex-col">
            {si2.map((sec) => (
              <>
                {sec.map((line) => (
                  <>
                    {line.Ins && (
                      <div className="text-gray-600 italic">{line.Ins}</div>
                    )}
                    {line.Ing && (
                      <div className="flex">
                        {line.Ing.amounts.map((a, x, z) => (
                          <>
                            <div className="text-purple-800 px-1">
                              {a.value}
                            </div>
                            <div className="text-pink-800 pr-1">{a.unit}</div>
                            {x < z.length - 1 && <>/</>}
                          </>
                        ))}

                        <div className="text-green-800">{line.Ing.name}</div>
                        {line.Ing.modifier && (
                          <>
                            ,
                            <div className="text-yellow-400 pl-1">
                              {line.Ing.modifier}
                            </div>
                          </>
                        )}
                      </div>
                    )}
                  </>
                ))}
                <div>&nbsp;</div>
              </>
            ))}
          </div>
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
    </div>
  );
};
export default InstructionsListParser;
