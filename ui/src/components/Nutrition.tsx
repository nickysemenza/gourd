import React from "react";
import { CellProps, Column } from "react-table";
import { scaledRound } from "../util";
import PaginatedTable from "./PaginatedTable";

const Nutrient: React.FC<{
  h: string[];
  items: Array<{
    ingredient: string;
    nutrients: Map<string, number>;
  }>;
}> = ({ items, h }) => {
  type i = typeof items[0];

  const foo = h.map((n) => {
    const res: Column<i> = {
      Header: n,
      Cell: ({ row: { original } }: CellProps<i>) => {
        const val = original.nutrients.get(n);
        return val ? scaledRound(val) : "";
      },
    };
    return res;
  });
  const columns: Array<Column<i>> = [
    {
      Header: "ingredient",
      Cell: ({ row: { original } }: CellProps<i>) => {
        return original.ingredient;
      },
    },
    ...foo,
  ];

  return (
    <PaginatedTable
      columns={columns}
      data={items}
      fetchData={() => {}}
      isLoading={false}
      totalCount={0}
      pageCount={1}
    />
  );
};
export default Nutrient;
