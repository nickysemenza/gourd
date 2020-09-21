import React from "react";
import { useGetIngredientsQuery } from "../generated/graphql";
import { Column, CellProps } from "react-table";
import { Link } from "react-router-dom";
import IngredientSearch from "../components/IngredientSearch";
import PaginatedTable from "../components/PaginatedTable";

const IngredientList: React.FC = () => {
  const { data } = useGetIngredientsQuery({});

  const ingredients = data?.ingredients || [];
  type i = Partial<typeof ingredients[0]>;

  const columns: Array<Column<i>> = React.useMemo(
    () => [
      {
        Header: "UUID",
        Cell: ({
          row: {
            original: { uuid },
          },
        }: CellProps<i>) => <code>{uuid}</code>,
      },
      {
        Header: "Name",
        // accessor: "name",
        Cell: ({ row: { original } }: CellProps<i>) => {
          const { name, same } = original;
          return (
            <div>
              {name}
              <ul>
                {(same || []).map((i) => (
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
                    <Link to={`recipe/${r.uuid}`} className="link">
                      {r.name}
                    </Link>
                  </li>
                ))}
              </ul>
              {/* <Debug data={original} /> */}
            </div>
          );
        },
      },
      {
        Header: "USDA food",
        // accessor: "name",
        Cell: ({ row: { original } }: CellProps<i>) => {
          const { usdaFood } = original;
          return <div>{usdaFood?.description}</div>;
        },
      },
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
        data={ingredients}
        fetchData={() => {}}
        isLoading={false}
        pageCount={1}
        totalCount={ingredients.length}
      />
    </>
  );
};

export default IngredientList;
