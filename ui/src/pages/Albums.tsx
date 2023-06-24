import React from "react";
import Debug from "../components/ui/Debug";
import PaginatedTable from "../components/ui/PaginatedTable";
import { createColumnHelper } from "@tanstack/react-table";
import { Code } from "../components/Code";
import { useListAllAlbums } from "../api/react-query/gourdApiComponents";
import { GooglePhotosAlbum } from "../api/react-query/gourdApiSchemas";

const Albums: React.FC = () => {
  const { data, error, isLoading } = useListAllAlbums({});

  const columns = React.useMemo(() => {
    const columnHelper = createColumnHelper<GooglePhotosAlbum>();
    return [
      columnHelper.accessor((row) => row.title, {
        id: "Name",
        cell: (info) => {
          return <div>{info.getValue()}</div>;
        },
      }),
      columnHelper.accessor((row) => row.usecase, {
        id: "Use case",
        cell: (info) => {
          return <div>{info.getValue()}</div>;
        },
      }),
      columnHelper.accessor((row) => row.id, {
        id: "id",
        cell: (info) => {
          return <Code>{info.getValue()}</Code>;
        },
      }),
    ];
  }, []);

  return (
    <div>
      <PaginatedTable
        columns={columns}
        data={data?.albums || []}
        fetchData={() => null}
        isLoading={isLoading}
        totalCount={0}
        pageCount={1}
      />
      <Debug data={{ error }} />
    </div>
  );
};

export default Albums;
