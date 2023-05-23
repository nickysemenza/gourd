import React from "react";
import {
  useReactTable,
  getCoreRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  PaginationState,
  flexRender,
  ColumnDef,
} from "@tanstack/react-table";

export interface PaginationParameters {
  offset?: number;
  limit?: number;
}

interface TableProps<T extends object> {
  columns: any | ColumnDef<T>[];
  data: T[];
  fetchData: (params: PaginationParameters) => void;
  isLoading: boolean;
  totalCount: number;
  pageCount: number;
}

const PaginatedTable = <T extends object>({
  columns,
  data,
  pageCount: controlledPageCount,
  fetchData,
  isLoading,
}: TableProps<T>) => {
  const [{ pageIndex, pageSize }, setPagination] =
    React.useState<PaginationState>({
      pageIndex: 0,
      pageSize: 10,
    });

  const pagination = React.useMemo(
    () => ({
      pageIndex,
      pageSize,
    }),
    [pageIndex, pageSize]
  );

  // Use the state and functions returned from useTable to build your UI
  const table = useReactTable({
    columns,
    data,
    pageCount: controlledPageCount ?? -1,
    state: {
      pagination,
    },
    manualPagination: true, // Tell the usePagination
    // hook that we'll handle our own data fetching
    // This means we'll also have to provide our own
    // pageCount.
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    onPaginationChange: setPagination,
  });

  React.useEffect(() => {
    let params: PaginationParameters = {
      limit: pageSize,
      offset: pageSize * pageIndex,
    };
    fetchData(params);
  }, [fetchData, pageIndex, pageSize]);

  // Render the UI for your table
  return (
    <div className="flex flex-col">
      {/* <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8"> */}
      {/* <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8"> */}
      {/* <div className="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg"> */}
      <table>
        <thead>
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                return (
                  <th key={header.id} colSpan={header.colSpan}>
                    {header.isPlaceholder ? null : (
                      <div>
                        {flexRender(
                          header.column.columnDef.header,
                          header.getContext()
                        )}
                      </div>
                    )}
                  </th>
                );
              })}
            </tr>
          ))}
        </thead>
        <tbody className="bg-white divide-y divide-gray-200 dark:divide-gray-700">
          {isLoading && (
            <tr>
              <td colSpan={10} className="w-100 text-xl text-center h-16">
                loading...
              </td>
            </tr>
          )}
          {table.getRowModel().rows.map((row) => {
            return (
              <tr
                key={row.id}
                className="bg-gray-100 odd:bg-gray-200 dark:bg-slate-500 dark:odd:bg-slate-400"
              >
                {row.getVisibleCells().map((cell) => {
                  return (
                    <td key={cell.id} className="px-6 py-4 whitespace-nowrap">
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </td>
                  );
                })}
              </tr>
            );
          })}
        </tbody>
      </table>
      {/* </div> */}
      {/* </div> */}
      {/* </div> */}
      <nav className="relative z-0 inline-flex shadow-sm">
        <button
          // href="#prev"
          className="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm leading-5 font-medium text-gray-500 hover:text-gray-400 focus:z-10 focus:outline-none focus:border-blue-300 focus:shadow-outline-blue active:bg-gray-100 active:text-gray-500 transition ease-in-out duration-150"
          aria-label="Previous"
          onClick={() => table.previousPage()}
          disabled={!table.getCanPreviousPage()}
        >
          {/* <!-- Heroicon name: chevron-left --> */}
          <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path
              fillRule="evenodd"
              d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z"
              clipRule="evenodd"
            />
          </svg>
        </button>
        <button
          disabled
          className="-ml-px relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm leading-5 font-medium text-gray-700 hover:text-gray-500 focus:z-10 focus:outline-none focus:border-blue-300 focus:shadow-outline-blue active:bg-gray-100 active:text-gray-700 transition ease-in-out duration-150"
        >
          {pageIndex + 1}
        </button>
        <button
          className="-ml-px relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm leading-5 font-medium text-gray-500 hover:text-gray-400 focus:z-10 focus:outline-none focus:border-blue-300 focus:shadow-outline-blue active:bg-gray-100 active:text-gray-500 transition ease-in-out duration-150"
          aria-label="Next"
          onClick={() => table.nextPage()}
          disabled={!table.getCanNextPage()}
        >
          {/* <!-- Heroicon name: chevron-right --> */}
          <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path
              fillRule="evenodd"
              d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
              clipRule="evenodd"
            />
          </svg>
        </button>
      </nav>
      <span>
        Page{" "}
        <strong>
          {table.getState().pagination.pageIndex + 1} of {table.getPageCount()}
        </strong>
      </span>
    </div>
  );
};

export default PaginatedTable;
