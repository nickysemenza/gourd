import React from "react";
import { useGetRecipesQuery } from "../generated/graphql";
import { useTable, Column, CellProps } from "react-table";
import { Link } from "react-router-dom";
import Debug from "../components/Debug";

interface TableProps<T extends object> {
  columns: Column<T>[];
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
  } = useTable({
    columns,
    data,
  });

  // Render the UI for your table
  return (
    <table
      className="table-auto border-collapse border-1 border-gray-500"
      {...getTableProps()}
      data-cy="recipe-table"
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

const RecipeList: React.FC = () => {
  const { data, error } = useGetRecipesQuery({});

  const columns = React.useMemo(
    () => [
      {
        Header: "UUID",
        accessor: "uuid",
      },
      {
        Header: "Name",
        accessor: "name",
      },
      {
        Header: "test",
        accessor: "test",
        Cell: (cell: CellProps<any>) => (
          <Link to={`recipe/${cell.row.original.uuid}`} className="link">
            details
          </Link>
        ),
      },
    ],
    []
  );

  return (
    <div>
      <Table columns={columns} data={data?.recipes || []} />
      <Debug data={{ error }} />
    </div>
  );
};

export default RecipeList;
