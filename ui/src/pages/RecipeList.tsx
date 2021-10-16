import React, { useState } from "react";
import { CellProps, Column } from "react-table";
import Debug from "../components/Debug";
import { useListRecipes } from "../api/openapi-hooks/api";
import PaginatedTable, {
  PaginationParameters,
} from "../components/PaginatedTable";
import { RecipeLink } from "../components/Misc";
import { Code } from "../util";
import { Helmet } from "react-helmet";
import update from "immutability-helper";
import { Link } from "react-router-dom";
import queryString from "query-string";
import { ButtonGroup } from "../components/Button";
import { PlusCircle } from "react-feather";

const RecipeList: React.FC = () => {
  // const { data, error } = useGetRecipesQuery({});
  // const queryParams = React.useMemo(
  //   () => queryString.stringify(params as any),
  //   [params]
  // );
  const [checked, setChecked] = useState(new Set<string>());
  let initialParams: PaginationParameters = {
    offset: 0,
    limit: 2,
  };

  const [params, setParams] = useState(initialParams);
  const [showOlder, setShowOlder] = useState(false);

  const fetchData = React.useCallback((params: PaginationParameters) => {
    setParams(params);
  }, []);

  const { data, error } = useListRecipes({
    queryParams: params,
  });

  const recipes = data?.recipes || [];
  type i = typeof recipes[0];

  const columns: Array<Column<i>> = React.useMemo(
    () => [
      {
        Header: "id",
        Cell: ({ row: { original } }: CellProps<i>) => {
          return <Code>{original.id} </Code>;
        },
      },
      {
        Header: "Name",
        // accessor: "name",
        Cell: ({ row: { original } }: CellProps<i>) => {
          const { versions } = original;
          return (
            <div>
              <ul>
                {(versions || [])
                  .filter((v) => showOlder || v.is_latest_version !== showOlder)
                  .map((i) => (
                    <li>
                      <div className="flex">
                        <RecipeLink recipe={i} />
                        <input
                          type="checkbox"
                          className="form-checkbox"
                          checked={checked.has(i.id)}
                          onClick={() =>
                            setChecked(
                              update(
                                checked,
                                checked.has(i.id)
                                  ? { $remove: [i.id] }
                                  : { $add: [i.id] }
                              )
                            )
                          }
                        />
                      </div>
                    </li>
                  ))}
              </ul>
            </div>
          );
        },
      },
      // {
      //   Header: "test",
      //   accessor: "test",
      //   Cell: (cell: CellProps<any>) => (
      //     <Link to={`recipe/${cell.row.original.id}`} className="link">
      //       details
      //     </Link>
      //   ),
      // },
    ],
    [checked, showOlder]
  );

  return (
    <div>
      <Helmet>
        <title>recipes | gourd</title>
      </Helmet>
      <PaginatedTable
        columns={columns}
        data={data?.recipes || []}
        fetchData={fetchData}
        isLoading={false}
        totalCount={data?.meta?.total_count || 0}
        pageCount={data?.meta?.page_count || 1}
      />
      <div>
        <ButtonGroup
          compact
          buttons={[
            {
              onClick: () => {
                setShowOlder(!showOlder);
              },
              text: "toggle older",
              IconLeft: PlusCircle,
            },
          ]}
        />
        {checked.size > 0 && (
          <>
            <Link
              to={`/playground?${queryString.stringify({
                recipes: [...checked.keys()],
              })}`}
            >
              Compare {checked.size} recipes
            </Link>
          </>
        )}
      </div>
      <Debug data={{ error }} />
    </div>
  );
};

export default RecipeList;
