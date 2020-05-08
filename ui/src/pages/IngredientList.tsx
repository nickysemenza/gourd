import React from "react";
import { useGetIngredientsQuery } from "../generated/graphql";
import styled from "styled-components";
import { useTable, Column, CellProps } from "react-table";
import { Link } from "react-router-dom";
import { Box } from "rebass";

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
    <table {...getTableProps()}>
      <thead>
        {headerGroups.map((headerGroup) => (
          <tr {...headerGroup.getHeaderGroupProps()}>
            {headerGroup.headers.map((column) => (
              <th {...column.getHeaderProps()}>{column.render("Header")}</th>
            ))}
          </tr>
        ))}
      </thead>
      <tbody {...getTableBodyProps()}>
        {rows.map((row) => {
          prepareRow(row);
          return (
            <tr {...row.getRowProps()}>
              {row.cells.map((cell) => {
                return <td {...cell.getCellProps()}>{cell.render("Cell")}</td>;
              })}
            </tr>
          );
        })}
      </tbody>
    </table>
  );
};

const Styles = styled.div`
  padding: 1rem;

  table {
    border-spacing: 0;
    border: 1px solid black;

    tr {
      :last-child {
        td {
          border-bottom: 0;
        }
      }
    }

    th,
    td {
      margin: 0;
      padding: 0.5rem;
      border-bottom: 1px solid black;
      border-right: 1px solid black;

      :last-child {
        border-right: 0;
      }
    }
  }
`;

const IngredientList: React.FC = () => {
  const { data } = useGetIngredientsQuery({});

  const ingredients = data?.ingredients || [];
  type i = Partial<typeof ingredients[0]>;

  const columns: Array<Column<i>> = React.useMemo(
    () => [
      {
        Header: "UUID",
        // accessor: ({ uuid }: CellProps<i>) => uuid,
        accessor: "uuid",
        // Cell: ({
        //   row: {
        //     original: { uuid },
        //   },
        // }: CellProps<i>) => uuid,
      },
      {
        Header: "Name",
        accessor: "name",
      },
      {
        Header: "test",
        id: "test",
        Cell: ({ row: { original } }: CellProps<i>) => (
          <Box>
            <ul>
              {(original.recipes || []).map((r) => (
                <li>
                  <Link to={`recipe/${r.uuid}`}>{r.name}</Link>
                </li>
              ))}
            </ul>
            {/* <Debug data={original} /> */}
          </Box>
        ),
      },
    ],
    []
  );

  return (
    <Styles>
      <Table<i> columns={columns} data={ingredients} />
    </Styles>
  );
};

export default IngredientList;
