import React, { useState } from "react";
import { CellProps } from "react-table";
import dayjs from "dayjs";
import Debug from "../components/Debug";
import { useListAllAlbums } from "../api/openapi-hooks/api";
import PaginatedTable, {
  PaginationParameters,
} from "../components/PaginatedTable";
import { Blurhash } from "react-blurhash";
import ProgressiveImage from "../components/ProgressiveImage";

const Albums: React.FC = () => {
  let initialParams: PaginationParameters = {
    offset: 0,
    limit: 2,
  };

  const [params, setParams] = useState(initialParams);

  const fetchData = React.useCallback((params: PaginationParameters) => {
    setParams(params);
  }, []);

  const { data, error } = useListAllAlbums({});

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
          return <div>{id}</div>;
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
        fetchData={fetchData}
        isLoading={false}
        totalCount={0}
        pageCount={1}
      />
      <Debug data={{ error }} />
    </div>
  );
};

export default Albums;
