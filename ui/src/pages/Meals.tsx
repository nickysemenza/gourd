import React, { useState } from "react";
import { CellProps, Column } from "react-table";
import dayjs from "dayjs";
import Debug from "../components/Debug";
import { GooglePhoto, useListMeals } from "../api/openapi-hooks/api";
import PaginatedTable, {
  PaginationParameters,
} from "../components/PaginatedTable";
import ProgressiveImage from "../components/ProgressiveImage";

const Meals: React.FC = () => {
  let initialParams: PaginationParameters = {
    offset: 0,
    limit: 2,
  };

  const [params, setParams] = useState(initialParams);

  const fetchData = React.useCallback((params: PaginationParameters) => {
    setParams(params);
  }, []);

  const { data, error } = useListMeals({
    queryParams: params,
  });

  const meals = data?.meals || [];
  type i = typeof meals[0];

  const columns: Array<Column<i>> = React.useMemo(
    () => [
      {
        Header: "ate_at",
        accessor: "ate_at",
        Cell: (cell: CellProps<i>) => {
          const { ate_at } = cell.row.original;
          const ago = dayjs(ate_at);

          return <div>{ago.format("dddd, MMMM D, YYYY h:mm A")}</div>;
        },
      },
      {
        Header: "recipes",
        accessor: "recipes",
        Cell: (cell: CellProps<i>) => {
          const { recipes } = cell.row.original;

          return (
            <div className="w-64">
              {(recipes || []).map((r) => (
                <div className="flex">
                  <div>{r.recipe.name}</div>
                  <div className="font-mono">@{r.multiplier}x</div>
                </div>
              ))}
              {/* {ago.format("dddd, MMMM D, YYYY h:mm A")} */}
            </div>
          );
        },
      },
      {
        Header: "Photos",
        accessor: "photos",
        Cell: (cell: CellProps<i>) => {
          const { photos } = cell.row.original;
          // https://developers.google.com/meals/library/guides/access-media-items#image-base-urls
          return (
            <div className="flex flex-wrap">
              {photos.map((photo: GooglePhoto) => (
                <ProgressiveImage photo={photo} />
                // <img
                //   onLoad={(x) => {
                //     console.log(x);
                //   }}
                //   key={photo.id}
                //   src={`${photo.base_url}=w120`}
                //   alt="todo"
                // />
              ))}
            </div>
          );
        },
      },
    ],
    []
  );

  return (
    <div>
      <PaginatedTable
        columns={columns}
        data={meals}
        fetchData={fetchData}
        isLoading={false}
        totalCount={data?.meta?.total_count || 0}
        pageCount={data?.meta?.page_count || 1}
      />
      <Debug data={{ error }} />
    </div>
  );
};

export default Meals;
