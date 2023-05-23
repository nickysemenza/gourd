import React from "react";
import { scaledRound } from "../util/util";
import PaginatedTable from "./ui/PaginatedTable";
import { createColumnHelper } from "@tanstack/react-table";

const Nutrient: React.FC<{
  h: string[];
  items: Array<{
    ingredient: string;
    nutrients: Map<string, number>;
  }>;
}> = ({ items, h }) => {
  type i = (typeof items)[0];

  const columnHelper = createColumnHelper<i>();
  const foo = h.map((n) => {
    const res = columnHelper.accessor((row) => row, {
      id: n,
      cell: (info) => {
        const val = info.row.original.nutrients.get(n);
        return val ? scaledRound(val) : "";
      },
    });

    return res;
  });
  const columns = [
    columnHelper.accessor((row) => row.ingredient, {
      id: "ingredient",
    }),
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
