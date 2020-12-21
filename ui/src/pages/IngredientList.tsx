import React, { useState } from "react";
import { Column, CellProps } from "react-table";
import { Link } from "react-router-dom";
import IngredientSearch from "../components/IngredientSearch";
import PaginatedTable, {
  PaginationParameters,
} from "../components/PaginatedTable";
import { useListIngredients } from "../api/openapi-hooks/api";

const IngredientList: React.FC = () => {
  let initialParams: PaginationParameters = {
    offset: 0,
    limit: 2,
  };

  const [params, setParams] = useState(initialParams);

  const fetchData = React.useCallback((params: PaginationParameters) => {
    setParams(params);
  }, []);

  const { data } = useListIngredients({
    queryParams: params,
  });

  const ingredients = data?.ingredients || [];
  type i = typeof ingredients[0];

  const columns: Array<Column<i>> = React.useMemo(
    () => [
      {
        Header: "UUID",
        Cell: ({
          row: {
            original: {
              ingredient: { id },
            },
          },
        }: CellProps<i>) => <code>{id}</code>,
      },
      {
        Header: "Name",
        // accessor: "name",
        Cell: ({ row: { original } }: CellProps<i>) => {
          const { ingredient, children } = original;
          return (
            <div>
              {ingredient.name}
              <ul>
                {(children || []).map((i) => (
                  <li>{i.name}</li>
                ))}
              </ul>
            </div>
          );
        },
      },
      {
        Header: "Recipes",
        id: "recipes",
        Cell: ({ row: { original } }: CellProps<i>) => {
          const recipes = original.recipes || [];
          return (
            <div>
              <ul>
                {recipes.map((r) => (
                  <li>
                    <Link to={`recipe/${r.id}`} className="link">
                      <div className="flex">
                        {r.name} (<div className="font-mono">v{r.version}</div>)
                      </div>
                    </Link>
                  </li>
                ))}
              </ul>
              {/* <Debug data={original} /> */}
            </div>
          );
        },
      },
      // {
      //   Header: "USDA food",
      //   // accessor: "name",
      //   Cell: ({ row: { original } }: CellProps<i>) => {
      //     const { usdaFood } = original;
      //     return <div>{usdaFood?.description}</div>;
      //   },
      // },
    ],
    []
  );

  return (
    <>
      <IngredientSearch
        initial="eg"
        callback={(item, kind) => console.log({ item, kind })}
      />
      <PaginatedTable
        columns={columns}
        data={data?.ingredients || []}
        fetchData={fetchData}
        isLoading={false}
        totalCount={data?.meta?.total_count || 0}
        pageCount={data?.meta?.page_count || 1}
      />
    </>
  );
};

export default IngredientList;
