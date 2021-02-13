import React from "react";
import { CellProps, Column } from "react-table";
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
        return val ? Math.round(val) : "";
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
    //   {
    //     Header: "Name",
    //     // accessor: "name",
    //     Cell: ({ row: { original } }: CellProps<i>) => {
    //       const { versions } = original;
    //       return (
    //         <div>
    //           <ul>
    //             {(versions || []).map((i) => (
    //               <li>
    //                 <div className="flex">
    //                   <RecipeLink recipe={i} />
    //                 </div>
    //               </li>
    //             ))}
    //           </ul>
    //         </div>
    //       );
    //     },
    //   },
  ];

  return (
    <div>
      <PaginatedTable
        columns={columns}
        data={items}
        fetchData={() => {}}
        isLoading={false}
        totalCount={0}
        pageCount={1}
      />
    </div>
  );
};
export default Nutrient;
