import React, { useEffect, useState } from "react";
import { PlusCircle } from "react-feather";
import { FoodInfo, useSearchFoods } from "../api/openapi-hooks/api";
import { Code } from "../util";
import { ButtonGroup } from "./Button";
import { UnitMappingList } from "./Misc";

const FoodSearch: React.FC<{
  name: string;
  highlightId?: number;
  limit?: number;
  enableSearch?: boolean;
  addon?: FoodInfo;
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

  const placeholder: FoodInfo = {
    wrapper: {
      fdc_id: 0,
      description: "loading...",
      data_type: "branded_food",
      nutrients: [],
      branded_info: undefined,
    },
    unit_mappings: [],
  };
  const results = foods?.foods || [];
  const reallyLoading = loading && results.length === 0;
  const items: FoodInfo[] = [
    ...(reallyLoading ? Array(5).fill(placeholder) : []),
    ...(addon &&
    results.filter((f) => f.wrapper.fdc_id === addon.wrapper.fdc_id).length ===
      0
      ? [addon]
      : []),
    ...results,
  ];

  // force search on if no results
  const showSearch = enableSearch || results.length === 0;
  return (
    <div className="w-full">
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
        const isHighlighted = highlightId === r.wrapper.fdc_id;
        return (
          <FoodRow
            info={r}
            isHighlighted={isHighlighted}
            onLink={onLink}
            x={x}
            key={x}
            loading={reallyLoading}
          />
        );
      })}
    </div>
  );
};
export const FoodRow: React.FC<{
  info: FoodInfo;
  isHighlighted?: boolean;
  x?: number;
  onLink?: (fdc_id: number) => void;
  loading: boolean;
  wide?: boolean;
  descriptionComponent?: JSX.Element;
  brandOwnerComponent?: JSX.Element;
}> = ({
  info,
  isHighlighted = false,
  x = 0,
  onLink,
  loading,
  wide = false,
  descriptionComponent,
  brandOwnerComponent,
}) => {
  const loadingClass =
    (loading && "h-2 bg-gray-400 rounded animate-pulse") || "";
  const food = info.wrapper;
  return (
    <div
      style={{ gridTemplateColumns: "1fr 3fr " }}
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
      <div className={`flex ${wide ? "flex-row" : "flex-col"} p-1`}>
        <div>
          <div className="flex whitespace-normal">
            {descriptionComponent || food.description}
          </div>
          <div className="flex flex-row">
            <p className="font-mono text-xs">{food.data_type}</p>
            <p className="pl-1 text-xs">{food.nutrients?.length} nutrients</p>
          </div>
          <UnitMappingList unit_mappings={info.unit_mappings} />
        </div>

        {(food.branded_info || loading) && (
          <div
            className={`flex ${
              wide ? "flex-row" : "flex-col w-80"
            }  ${loadingClass}`}
          >
            <div>
              {brandOwnerComponent || food.branded_info?.brand_owner} <br />
              <p className={`text-sm italic ${loadingClass}`}>
                {food.branded_info?.branded_food_category}
              </p>
            </div>
            <div
              className={`text-xs text-gray-500 whitespace-normal ${loadingClass}`}
            >
              {food.branded_info?.ingredients}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default FoodSearch;
