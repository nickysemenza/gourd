import React from "react";
import { useTable, Column, usePagination } from "react-table";

export interface PaginationParameters {
  offset?: number;
  limit?: number;
}

interface TableProps<T extends object> {
  columns: Column<T>[];
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
  // Use the state and functions returned from useTable to build your UI
  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    rows,
    prepareRow,
    canPreviousPage,
    canNextPage,
    pageOptions,
    nextPage,
    previousPage,
    state: { pageIndex, pageSize },
  } = useTable(
    {
      columns,
      data,
      initialState: { pageIndex: 0, pageSize: 50 }, // Pass our hoisted table state
      manualPagination: true, // Tell the usePagination
      // hook that we'll handle our own data fetching
      // This means we'll also have to provide our own
      // pageCount.
      pageCount: controlledPageCount,
    },
    usePagination
  );

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
      <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
        <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
          <div className="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
            <table
              className="min-w-full divide-y divide-gray-200 w-full"
              {...getTableProps()}
              data-cy="recipe-table"
            >
              <thead className="bg-gray-50 dark:bg-gray-900">
                {headerGroups.map((headerGroup) => (
                  <tr {...headerGroup.getHeaderGroupProps()}>
                    {headerGroup.headers.map((column) => (
                      <th
                        className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                        {...column.getHeaderProps()}
                      >
                        {column.render("Header")}
                      </th>
                    ))}
                  </tr>
                ))}
              </thead>
              <tbody
                className="bg-white divide-y divide-gray-200 dark:divide-gray-700"
                {...getTableBodyProps()}
              >
                {isLoading && (
                  <tr>
                    <td colSpan={10} className="w-100 text-xl text-center h-16">
                      loading...
                    </td>
                  </tr>
                )}
                {rows.map((row) => {
                  prepareRow(row);
                  return (
                    <tr
                      {...row.getRowProps()}
                      className="bg-gray-100 odd:bg-gray-200 dark:bg-slate-500 dark:odd:bg-slate-400"
                    >
                      {row.cells.map((cell) => {
                        return (
                          <td
                            className="px-6 py-4 whitespace-nowrap"
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
          </div>
        </div>
      </div>
      <nav className="relative z-0 inline-flex shadow-sm">
        <button
          // href="#prev"
          className="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm leading-5 font-medium text-gray-500 hover:text-gray-400 focus:z-10 focus:outline-none focus:border-blue-300 focus:shadow-outline-blue active:bg-gray-100 active:text-gray-500 transition ease-in-out duration-150"
          aria-label="Previous"
          onClick={() => previousPage()}
          disabled={!canPreviousPage}
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
          onClick={() => nextPage()}
          disabled={!canNextPage}
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
          {pageIndex + 1} of {pageOptions.length}
        </strong>
      </span>
    </div>
  );
};

export default PaginatedTable;
