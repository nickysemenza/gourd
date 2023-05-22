import React, { useState } from "react";
import dayjs from "dayjs";
import Debug from "../components/Debug";
import { Photo, useListPhotos } from "../api/openapi-hooks/api";
import PaginatedTable, {
  PaginationParameters,
} from "../components/PaginatedTable";
import ProgressiveImage from "../components/ProgressiveImage";
import { createColumnHelper } from "@tanstack/react-table";

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

  const columnHelper = createColumnHelper<Photo>();
  const columns = [
    columnHelper.accessor((row) => row.taken_at, {
      id: "taken",
      cell: (info) => {
        const ago = dayjs(info.getValue());

        return (
          <div>
            {ago.format("dddd, MMMM D, YYYY h:mm A")}
            <br />
            {ago.fromNow()}
          </div>
        );
      },
    }),

    columnHelper.accessor((row) => row, {
      id: "test",
      cell: (info) => <ProgressiveImage photo={info.row.original} />,
    }),
  ];

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
