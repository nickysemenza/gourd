import React, { useState } from "react";
import { CellProps } from "react-table";
import dayjs from "dayjs";
import Debug from "../components/Debug";
import { useListPhotos } from "../api/openapi-hooks/api";
import PaginatedTable, {
  PaginationParameters,
} from "../components/PaginatedTable";
import ProgressiveImage from "../components/ProgressiveImage";

const Photos: React.FC = () => {
  let initialParams: PaginationParameters = {
    offset: 0,
    limit: 2,
  };

  const [params, setParams] = useState(initialParams);

  const fetchData = React.useCallback((params: PaginationParameters) => {
    setParams(params);
  }, []);

  const { data, error, loading } = useListPhotos({
    queryParams: params,
  });

  const columns = React.useMemo(
    () => [
      {
        Header: "Taken",
        accessor: "Taken",
        Cell: (cell: CellProps<any>) => {
          const { taken_at } = cell.row.original;
          const ago = dayjs(taken_at);

          return (
            <div>
              {ago.format("dddd, MMMM D, YYYY h:mm A")}
              <br />
              {ago.fromNow()}
            </div>
          );
        },
      },
      {
        Header: "test",
        accessor: "test",
        Cell: (cell: CellProps<any>) => (
          <ProgressiveImage photo={cell.row.original} />
        ),
      },
    ],
    []
  );

  return (
    <div>
      <PaginatedTable
        columns={columns}
        data={data?.photos || []}
        fetchData={fetchData}
        isLoading={loading}
        totalCount={data?.meta?.total_count || 0}
        pageCount={data?.meta?.page_count || 1}
      />
      <Debug data={{ error }} />
    </div>
  );
};

export default Photos;
