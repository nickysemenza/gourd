import React from "react";
import { PlusCircle } from "react-feather";
import { Code } from "../components/Code";
import { ButtonGroup } from "./ui/ButtonGroup";
import { UnitMappingList } from "./misc/Misc";
import { useForm } from "react-hook-form";
import Loading from "./ui/Loading";
import { useSearchFoods } from "../api/react-query/gourdApiComponents";
import { TempFood } from "../api/react-query/gourdApiSchemas";
const FoodSearch: React.FC<{
  name: string;
  highlightId?: number;
  limit?: number;
  enableSearch?: boolean;
  addon?: TempFood;
  onLink?: (fdcId: number) => void;
}> = ({
  name,
  highlightId,
  onLink,
  limit = 5,
  enableSearch = false,
  addon,
}) => {
  const { register, watch } = useForm({ defaultValues: { name } });

  const { isLoading, data: foods } = useSearchFoods({
    queryParams: {
      name: watch("name"),
      limit,
    },
  });

  const placeholder: TempFood = {
    wrapper: {
      fdcId: 0,
      description: "loading...",
      dataType: "branded_food",
      //   nutrients: [],
      //   branded_info: undefined,
    },
    unit_mappings: [],
  };
  const results = foods?.foods || [];
  const reallyLoading = isLoading && results.length === 0;
  const items: TempFood[] = [
    ...(reallyLoading ? Array(5).fill(placeholder) : []),
    ...(addon &&
    results.filter((f) => f.wrapper.fdcId === addon.wrapper.fdcId).length === 0
      ? [addon]
      : []),
    ...results,
  ];

  // force search on if no results
  const showSearch = enableSearch || results.length === 0;

  return (
    <div className="w-full">
      <form>
        <input
          type="text"
          className="border-2 border-gray-300"
          defaultValue={name}
          {...register("name")}
          disabled={!showSearch}
        />
      </form>
      {items.map((r, x) => {
        const isHighlighted = highlightId === r.wrapper.fdcId;
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
  info: TempFood;
  isHighlighted?: boolean;
  x?: number;
  onLink?: (fdcId: number) => void;
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
  const food = info.wrapper;
  const brand = info.branded_food;
  return (
    <div
      style={{ gridTemplateColumns: "1fr 4fr " }}
      className={`border ${
        isHighlighted ? "border-red-600 " : "border-indigo-600"
      } ${isHighlighted && "bg-indigo-200"} grid p-1 text-sm`}
      key={`${food.fdcId}@${x}`}
    >
      <div className="flex flex-col p-1">
        <Code>{food.fdcId}</Code>
        <a
          href={`https://fdc.nal.usda.gov/fdc-app.html#/food-details/${food.fdcId}/nutrients`}
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
                  onLink(food.fdcId);
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
          <div className="flex flex-row justify-between">
            <div className="flex whitespace-normal text-slate-700 font-bold w-1/2">
              {descriptionComponent || food.description}
            </div>
            <Loading loading={loading}>
              <UnitMappingList unit_mappings={info.unit_mappings} />
            </Loading>
          </div>
          <Loading loading={loading}>
            <div className="flex flex-row justify-between">
              <p className="font-mono text-xs">{food?.dataType}</p>
              <p className="pl-1 text-xs">
                {info.foodNutrients?.length} nutrients
              </p>
            </div>
          </Loading>
        </div>

        <div className={`flex ${wide ? "flex-row" : "flex-col w-full"}`}>
          <div>
            <Loading loading={loading}>
              <div className="flex flex-row justify-between">
                {brandOwnerComponent || brand?.brandOwner}
                <p className={`text-sm italic`}>{brand?.brandedFoodCategory}</p>
              </div>
            </Loading>
          </div>
          <Loading loading={loading}>
            <div
              className={`text-xs text-emerald-600 whitespace-normal lowercase`}
            >
              {brand?.ingredients}
            </div>
          </Loading>
        </div>
      </div>
    </div>
  );
};

export default FoodSearch;
