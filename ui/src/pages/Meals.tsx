import React, { useState } from "react";
import dayjs from "dayjs";
import Debug from "../components/ui/Debug";
import {
  fetchUpdateRecipesForMeal,
  useListMeals,
} from "../api/react-query/gourdApiComponents";
import PaginatedTable, {
  PaginationParameters,
} from "../components/ui/PaginatedTable";
import ProgressiveImage from "../components/ui/ProgressiveImage";
import { RecipeLink } from "../components/misc/Misc";
import { EntitySelector } from "../components/EntitySelector";
import { pushMealRecipe } from "../components/recipe/RecipeEditorUtils";
import update from "immutability-helper";
import queryString from "query-string";
import { Link } from "react-router-dom";
import { createColumnHelper } from "@tanstack/react-table";
import { Photo } from "../api/react-query/gourdApiSchemas";
const Meals: React.FC = () => {
  const initialParams: PaginationParameters = {
    offset: 0,
    limit: 2,
  };

  const [params, setParams] = useState(initialParams);

  const fetchData = React.useCallback((params: PaginationParameters) => {
    setParams(params);
  }, []);

  const {
    data,
    error,
    isLoading: loading,
  } = useListMeals({
    queryParams: params,
  });

  const meals = data?.meals || [];
  const [internalVal, setVal] = React.useState(meals);
  React.useEffect(() => {
    setVal(data?.meals || []);
  }, [data]);
  type i = (typeof meals)[0];
  const [checked, setChecked] = useState(new Set<string>());
  const columns = React.useMemo(() => {
    const columnHelper = createColumnHelper<i>();
    return [
      columnHelper.accessor((row) => row.ate_at, {
        id: "ate_at",
        header: () => <span>Ate on</span>,
        cell: (info) => {
          const ago = dayjs(info.getValue());

          // return <div>{ago.format("ddd, MMM D, YYYY h:mm A Z")}</div>;
          return <div>{ago.format("ddd, MMM D, YYYY")}</div>;
        },
      }),
      columnHelper.accessor((row) => row.id, {
        id: "select",
        cell: (info) => {
          // const { original } = info.row;
          const id = info.getValue();

          return (
            <div>
              <input
                type="checkbox"
                className="form-checkbox"
                checked={checked.has(id)}
                onClick={() =>
                  setChecked(
                    update(
                      checked,
                      checked.has(id) ? { $remove: [id] } : { $add: [id] }
                    )
                  )
                }
              />
            </div>
          );
        },
      }),
      columnHelper.accessor("name", {
        //accessorKey
        header: "Name",
      }),

      columnHelper.accessor((row) => row.recipes, {
        id: "recipes",
        cell: (info) => {
          return (
            <div className="w-64">
              <EntitySelector
                createKind="recipe"
                showKind={["recipe"]}
                placeholder="Pick a Recipe..."
                onChange={async (a) => {
                  console.log(a, info.row.index);
                  const res = await fetchUpdateRecipesForMeal({
                    pathParams: {
                      mealId: info.row.original.id,
                    },
                    body: {
                      multiplier: 1.0,
                      action: "add",
                      recipe_id: a.value,
                    },
                  });
                  console.log({ res });
                  setVal(
                    pushMealRecipe(internalVal, info.row.index, {
                      id: a.value,
                      name: a.label,
                      sections: [],
                      tags: [],
                      serving_info: {
                        quantity: 1,
                        unit: "",
                      },
                      meta: {
                        version: 0,
                        is_latest_version: false,
                      },
                      created_at: "",
                      sources: [],
                    })
                  );
                }}
              />
              {(info.getValue() || []).map((r) => (
                <div className="">
                  <RecipeLink recipe={r.recipe} multiplier={r.multiplier} />
                </div>
              ))}
              {/* {ago.format("dddd, MMMM D, YYYY h:mm A")} */}
            </div>
          );
        },
      }),
      columnHelper.accessor((row) => row.photos, {
        id: "photos",
        cell: (info) => {
          // https://developers.google.com/meals/library/guides/access-media-items#image-base-urls
          return (
            <div className="flex flex-row">
              {info.getValue().map((photo: Photo) => (
                <ProgressiveImage photo={photo} key={photo.id} />
              ))}
            </div>
          );
        },
      }),
    ];
  }, [internalVal, checked]);

  return (
    <div>
      <PaginatedTable
        columns={columns}
        data={internalVal}
        fetchData={fetchData}
        isLoading={loading}
        totalCount={data?.meta?.total_count || 0}
        pageCount={data?.meta?.page_count || 1}
      />
      <Debug data={{ error }} />
      {checked.size > 0 && (
        <>
          <Link
            to={`/diff?${queryString.stringify({
              recipes: internalVal
                .filter((m) => checked.has(m.id))
                .map((m) => (m.recipes || []).map((r) => r.recipe.id))
                .flat(),
            })}`}
          >
            Compare {checked.size} meals
          </Link>
        </>
      )}
    </div>
  );
};

export default Meals;
