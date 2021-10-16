import React, { useState } from "react";
import { CellProps, Column } from "react-table";
import dayjs from "dayjs";
import Debug from "../components/Debug";
import { GooglePhoto, useListMeals } from "../api/openapi-hooks/api";
import PaginatedTable, {
  PaginationParameters,
} from "../components/PaginatedTable";
import ProgressiveImage from "../components/ProgressiveImage";
import { RecipeLink } from "../components/Misc";
import { EntitySelector } from "../components/EntitySelector";
import { pushMealRecipe } from "../components/RecipeEditorUtils";
import { getOpenapiFetchConfig } from "../config";
import { MealRecipeUpdateActionEnum, MealsApi } from "../api/openapi-fetch";
const Meals: React.FC = () => {
  let initialParams: PaginationParameters = {
    offset: 0,
    limit: 2,
  };

  const [params, setParams] = useState(initialParams);

  const fetchData = React.useCallback((params: PaginationParameters) => {
    setParams(params);
  }, []);

  const { data, error } = useListMeals({
    queryParams: params,
  });

  const meals = data?.meals || [];
  const [internalVal, setVal] = React.useState(meals);
  React.useEffect(() => {
    setVal(data?.meals || []);
  }, [data]);
  type i = typeof meals[0];

  const columns: Array<Column<i>> = React.useMemo(() => {
    const mApi = new MealsApi(getOpenapiFetchConfig());
    return [
      {
        Header: "ate_at",
        accessor: "ate_at",
        Cell: (cell: CellProps<i>) => {
          const { ate_at } = cell.row.original;
          const ago = dayjs(ate_at);

          return <div>{ago.format("dddd, MMMM D, YYYY h:mm A")}</div>;
        },
      },
      {
        Header: "recipes",
        accessor: "recipes",
        Cell: (cell: CellProps<i>) => {
          const { recipes } = cell.row.original;

          return (
            <div className="w-64">
              <EntitySelector
                createKind="recipe"
                showKind={["recipe"]}
                placeholder="Pick a Recipe..."
                onChange={async (a) => {
                  console.log(a, cell.row.index);
                  let res = await mApi.updateRecipesForMeal({
                    mealId: cell.row.original.id,
                    mealRecipeUpdate: {
                      multiplier: 1.0,
                      action: MealRecipeUpdateActionEnum.ADD,
                      recipe_id: a.value,
                    },
                  });
                  console.log({ res });
                  setVal(
                    pushMealRecipe(internalVal, cell.row.index, {
                      id: a.value,
                      name: a.label,
                      sections: [],
                      quantity: 1,
                      unit: "",
                      version: 0,
                      is_latest_version: false,
                      created_at: "",
                    })
                  );
                }}
              />
              {(recipes || []).map((r) => (
                <div className="">
                  <RecipeLink recipe={r.recipe} multiplier={r.multiplier} />
                </div>
              ))}
              {/* {ago.format("dddd, MMMM D, YYYY h:mm A")} */}
            </div>
          );
        },
      },
      {
        Header: "Photos",
        accessor: "photos",
        Cell: (cell: CellProps<i>) => {
          const { photos } = cell.row.original;
          // https://developers.google.com/meals/library/guides/access-media-items#image-base-urls
          return (
            <div className="flex flex-wrap">
              {photos.map((photo: GooglePhoto) => (
                <ProgressiveImage photo={photo} />
                // <img
                //   onLoad={(x) => {
                //     console.log(x);
                //   }}
                //   key={photo.id}
                //   src={`${photo.base_url}=w120`}
                //   alt="todo"
                // />
              ))}
            </div>
          );
        },
      },
    ];
  }, [internalVal]);

  return (
    <div>
      <PaginatedTable
        columns={columns}
        data={internalVal}
        fetchData={fetchData}
        isLoading={false}
        totalCount={data?.meta?.total_count || 0}
        pageCount={data?.meta?.page_count || 1}
      />
      <Debug data={{ error }} />
    </div>
  );
};

export default Meals;
