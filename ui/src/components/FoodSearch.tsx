import React, { useEffect, useState } from "react";
import { PlusCircle } from "react-feather";
import { useSearchFoods } from "../api/openapi-hooks/api";
import { Code } from "../util";
import { ButtonGroup } from "./Button";
import { getCalories } from "./RecipeEditorUtils";

const FoodSearch: React.FC<{
  name: string;
  highlightId?: number;
  limit?: number;
  enableSearch?: boolean;
  onLink?: (fdc_id: number) => void;
}> = ({ name, highlightId, onLink, limit = 5, enableSearch = false }) => {
  const [ingredientName, setIngredientName] = useState(name);

  useEffect(() => {
    setIngredientName(name);
  }, [name]);

  const { loading, data: foods } = useSearchFoods({
    queryParams: {
      name: ingredientName,
      limit,
      data_types: [
        "foundation_food",
        "sample_food",
        "market_acquisition",
        "survey_fndds_food",
        "sub_sample_food",
        "agricultural_acquisition",
        "sr_legacy_food",
        "branded_food",
      ],
    },
  });

  const items = foods?.foods || [];
  return (
    <div className="">
      <ul className="list-disc list-outside pl-4">
        {loading && <div className="w-full">loading...</div>}
        {enableSearch && (
          <div className="w-full">
            <input
              className="border-2 border-gray-300"
              value={ingredientName}
              onChange={(e) => {
                setIngredientName(e.target.value);
              }}
            />
          </div>
        )}
        {(items || []).map((r) => {
          const isHighlighted = highlightId === r.fdc_id;
          return (
            <div
              style={{ gridTemplateColumns: "5rem 15rem 5rem 5rem 5rem 5rem" }}
              className={`border ${
                isHighlighted ? "border-red-600 " : "border-indigo-600"
              } ${isHighlighted && "bg-indigo-200"} grid`}
              key={`${name}@${r.fdc_id}`}
            >
              <div className="flex flex-col">
                <Code>{r.fdc_id}</Code>
                <a
                  href={`https://fdc.nal.usda.gov/fdc-app.html#/food-details/${r.fdc_id}/nutrients`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm pr-1"
                >
                  (view)
                </a>
              </div>
              <div className="flex flex-col">
                <div className="">{r.description}</div>{" "}
                <Code>{r.data_type}</Code>
              </div>
              {!!r.branded_info && (
                <div className="italic">{r.branded_info.brand_owner}</div>
              )}
              {/* <div className="flex"> */}
              <div className="flex flex-col">
                <div className="font-bold flex">nutrients:</div>
                <Code>{r.nutrients?.length}</Code>
              </div>
              <div className="flex flex-col">
                <div className="font-bold flex ml-1">nutrition:</div>
                <div>{`${getCalories(r)} kcal/100g`}</div>
              </div>
              {/* </div> */}
              {onLink !== undefined && (
                <ButtonGroup
                  compact
                  buttons={[
                    {
                      onClick: () => {
                        onLink(r.fdc_id);
                      },
                      text: "link",
                      disabled: isHighlighted,
                      IconLeft: PlusCircle,
                    },
                  ]}
                />
              )}
            </div>
          );
        })}
      </ul>
      {/* <Debug data={foods.foods} /> */}
    </div>
  );
};
export default FoodSearch;
