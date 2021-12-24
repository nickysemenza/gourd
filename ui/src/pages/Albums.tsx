import React from "react";
import { CellProps } from "react-table";
import Debug from "../components/Debug";
import { useListAllAlbums } from "../api/openapi-hooks/api";
import PaginatedTable from "../components/PaginatedTable";
import { Code } from "../util";

const Albums: React.FC = () => {
  const { data, error, loading } = useListAllAlbums({});

  const columns = React.useMemo(
    () => [
      {
        Header: "Name",
        accessor: "name",
        Cell: (cell: CellProps<any>) => {
          const { title } = cell.row.original;
          return <div>{title}</div>;
        },
      },
      {
        Header: "Use case",
        accessor: "usecase",
        Cell: (cell: CellProps<any>) => {
          const { usecase } = cell.row.original;
          return <div>{usecase}</div>;
        },
      },
      {
        Header: "id",
        accessor: "id",
        Cell: (cell: CellProps<any>) => {
          const { id } = cell.row.original;
          return <Code>{id}</Code>;
        },
      },
    ],
    []
  );

  return (
    <div>
      <PaginatedTable
        columns={columns}
        data={data?.albums || []}
        fetchData={() => {}}
        isLoading={loading}
        totalCount={0}
        pageCount={1}
      />
      <Debug data={{ error }} />
    </div>
  );
};

export default Albums;
