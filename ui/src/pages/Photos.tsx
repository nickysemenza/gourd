import React, { useState } from "react";
import { CellProps } from "react-table";
import dayjs from "dayjs";
import Debug from "../components/Debug";
import { useListPhotos } from "../api/openapi-hooks/api";
import PaginatedTable, {
  PaginationParameters,
} from "../components/PaginatedTable";

const Photos: React.FC = () => {
  let initialParams: PaginationParameters = {
    offset: 0,
    limit: 2,
  };

  const [params, setParams] = useState(initialParams);

  const fetchData = React.useCallback((params: PaginationParameters) => {
    setParams(params);
  }, []);

  const { data, error } = useListPhotos({
    queryParams: params,
  });

  const columns = React.useMemo(
    () => [
      //   {
      //     Header: "id",
      //     accessor: "id",
      //   },
      {
        Header: "Created",
        accessor: "created",
        Cell: (cell: CellProps<any>) => {
          const { created } = cell.row.original;
          const ago = dayjs(created);

          return <div>{ago.format("dddd, MMMM D, YYYY h:mm A")}</div>;
        },
      },
      {
        Header: "test",
        accessor: "test",
        Cell: (cell: CellProps<any>) => (
          // https://developers.google.com/photos/library/guides/access-media-items#image-base-urls
          <img src={`${cell.row.original.base_url}=w200`} alt="todo" />
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
        isLoading={false}
        totalCount={data?.meta?.total_count || 0}
        pageCount={data?.meta?.page_count || 1}
      />
      <Debug data={{ error }} />
    </div>
  );
};

export default Photos;
