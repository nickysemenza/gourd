import React, { useEffect, useState } from "react";
import { PlusCircle } from "react-feather";
import { Food, useSearchFoods } from "../api/openapi-hooks/api";
import { Code } from "../util";
import { ButtonGroup } from "./Button";
import { UnitMappingList } from "./Misc";

const FoodSearch: React.FC<{
  name: string;
  highlightId?: number;
  limit?: number;
  enableSearch?: boolean;
  addon?: Food;
  onLink?: (fdc_id: number) => void;
}> = ({
  name,
  highlightId,
  onLink,
  limit = 5,
  enableSearch = false,
  addon,
}) => {
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

  const placeholder: Food = {
    fdc_id: 0,
    description: "loading...",
    data_type: "branded_food",
    unit_mappings: [],
    nutrients: [],
    branded_info: undefined,
  };
  const results = foods?.foods || [];
  const items = [
    ...(loading && results.length === 0 ? Array(5).fill(placeholder) : []),
    ...(addon && results.filter((f) => f.fdc_id === addon.fdc_id).length === 0
      ? [addon]
      : []),
    ...results,
  ];

  // force search on if no results
  const showSearch = enableSearch || results.length === 0;
  return (
    <div className="w-3/12">
      {showSearch && (
        <input
          type="text"
          className="border-2 border-gray-300"
          value={ingredientName}
          onChange={(e) => {
            setIngredientName(e.target.value);
          }}
        />
      )}
      {items.map((r, x) => {
        const isHighlighted = highlightId === r.fdc_id;
        return (
          <FoodRow
            food={r}
            isHighlighted={isHighlighted}
            onLink={onLink}
            x={x}
            key={x}
          />
        );
      })}
      {/* <Debug data={foods.foods} /> */}
    </div>
  );
};
export const FoodRow: React.FC<{
  food: Food;
  isHighlighted?: boolean;
  x?: number;
  onLink?: (fdc_id: number) => void;
}> = ({ food, isHighlighted = false, x = 0, onLink }) => (
  <div
    style={{ gridTemplateColumns: "1fr 3fr 4fr" }}
    className={`border ${
      isHighlighted ? "border-red-600 " : "border-indigo-600"
    } ${isHighlighted && "bg-indigo-200"} grid p-1 text-sm`}
    key={`${food.fdc_id}@${x}`}
  >
    <div className="flex flex-col p-1">
      <Code>{food.fdc_id}</Code>
      <a
        href={`https://fdc.nal.usda.gov/fdc-app.html#/food-details/${food.fdc_id}/nutrients`}
        target="_blank"
        rel="noopener noreferrer"
        className="text-sm pr-1 underline text-blue-800"
      >
        view
      </a>
      {onLink !== undefined && (
        <ButtonGroup
          compact
          buttons={[
            {
              onClick: () => {
                onLink(food.fdc_id);
              },
              text: "link",
              disabled: isHighlighted,
              IconLeft: PlusCircle,
            },
          ]}
        />
      )}
    </div>
    <div className="flex flex-col p-1">
      <div className="flex whitespace-normal">{food.description}</div>
      <p className="font-mono text-xs text-gray-500">{food.data_type}</p>
      <p className="text-sm">{food.nutrients?.length} nutrients</p>
      <UnitMappingList unit_mappings={food.unit_mappings} />
      {!!food.branded_info && (
        <div className="flex flex-col w-80">
          {food.branded_info.brand_owner} <br />
          <p className="text-sm italic">
            {food.branded_info.branded_food_category}
          </p>
          <p className="text-xs text-gray-500 whitespace-normal">
            {food.branded_info.ingredients}
          </p>
        </div>
      )}
    </div>
  </div>
);

export default FoodSearch;
