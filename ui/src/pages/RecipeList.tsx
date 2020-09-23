import React, { useState } from "react";
import { CellProps } from "react-table";
import { Link } from "react-router-dom";
import Debug from "../components/Debug";
import { useListRecipes } from "../api/openapi-hooks/api";
import PaginatedTable, {
  PaginationParameters,
} from "../components/PaginatedTable";

const RecipeList: React.FC = () => {
  // const { data, error } = useGetRecipesQuery({});
  // const queryParams = React.useMemo(
  //   () => queryString.stringify(params as any),
  //   [params]
  // );

  let initialParams: PaginationParameters = {
    offset: 0,
    limit: 2,
  };

  const [params, setParams] = useState(initialParams);

  const fetchData = React.useCallback((params: PaginationParameters) => {
    setParams(params);
  }, []);

  const { data, error } = useListRecipes({
    base: "http://localhost:4242/api",
    queryParams: params,
  });

  const columns = React.useMemo(
    () => [
      {
        Header: "id",
        accessor: "id",
      },
      {
        Header: "Name",
        accessor: "name",
      },
      {
        Header: "test",
        accessor: "test",
        Cell: (cell: CellProps<any>) => (
          <Link to={`recipe/${cell.row.original.id}`} className="link">
            details
          </Link>
        ),
      },
    ],
    []
  );

  return (
    <div>
      <PaginatedTable
        columns={columns}
        data={data?.recipes || []}
        fetchData={fetchData}
        isLoading={false}
        totalCount={data?.meta?.total_count || 0}
        pageCount={data?.meta?.page_count || 1}
      />
      <Debug data={{ error }} />
    </div>
  );
};

export default RecipeList;
