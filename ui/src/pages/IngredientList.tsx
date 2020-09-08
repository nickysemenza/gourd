import React from "react";
import { useGetIngredientsQuery } from "../generated/graphql";
import { useTable, Column, CellProps } from "react-table";
import { Link } from "react-router-dom";
import IngredientSearch from "../components/IngredientSearch";

interface TableProps<T extends object> {
  columns: Array<Column<T>>;
  data: T[];
}

const Table = <T extends object>({ columns, data }: TableProps<T>) => {
  // Use the state and functions returned from useTable to build your UI
  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    rows,
    prepareRow,
  } = useTable<T>({
    columns,
    data,
  });

  // Render the UI for your table
  return (
    <table
      className="table-auto border-collapse border-1 border-gray-500"
      {...getTableProps()}
    >
      <thead>
        {headerGroups.map((headerGroup) => (
          <tr {...headerGroup.getHeaderGroupProps()}>
            {headerGroup.headers.map((column) => (
              <th
                className="border border-gray-400"
                {...column.getHeaderProps()}
              >
                {column.render("Header")}
              </th>
            ))}
          </tr>
        ))}
      </thead>
      <tbody {...getTableBodyProps()}>
        {rows.map((row) => {
          prepareRow(row);
          return (
            <tr {...row.getRowProps()} className="bg-white odd:bg-gray-200">
              {row.cells.map((cell) => {
                return (
                  <td
                    className="border border-gray-400 p-2"
                    {...cell.getCellProps()}
                  >
                    {cell.render("Cell")}
                  </td>
                );
              })}
            </tr>
          );
        })}
      </tbody>
    </table>
  );
};

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
      <Table<i> columns={columns} data={ingredients} />
    </>
  );
};

export default IngredientList;
