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
import { ButtonGroup, Pill2 } from "../components/Button";
import { PlusCircle } from "react-feather";
import ProgressiveImage from "../components/ProgressiveImage";
import { sumIngredients } from "../components/RecipeEditorUtils";
import dayjs from "dayjs";

const RecipeList: React.FC = () => {
  const showIds = false;
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
        Header: "Name",
        // accessor: "name",
        Cell: ({ row: { original } }: CellProps<i>) => {
          const olderVersions = original.detail.other_versions || [];
          const versions = [
            original.detail,
            ...(showOlder ? olderVersions : []),
          ];
          const ing = Object.keys(
            sumIngredients(original.detail.sections).ingredients
          );
          const rec = Object.keys(
            sumIngredients(original.detail.sections).recipes
          );
          return (
            <div>
              <ul>
                {versions.map((i) => (
                  <li key={i.id}>
                    <div className="flex">
                      <RecipeLink recipe={i} />
                      <input
                        type="checkbox"
                        className="form-checkbox"
                        checked={checked.has(i.id)}
                        onChange={(e) => {}}
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
              <Pill2 color={ing.length + rec.length > 0 ? "green" : "red"}>
                {ing.length} ing / {rec.length} rec
              </Pill2>
            </div>
          );
        },
      },
      {
        Header: "meals",
        Cell: ({ row: { original } }: CellProps<i>) => (
          <div className="w-9/12 flex ">
            {(original.linked_photos || []).map((p) => (
              <ProgressiveImage photo={p} />
            ))}
          </div>
        ),
        // return <Debug data={original.linked_meals} />;
      },
      {
        Header: "created at",
        Cell: ({ row: { original } }: CellProps<i>) => {
          const ago = dayjs(original.detail.created_at);

          // return <div>{ago.format("ddd, MMM D, YYYY h:mm A Z")}</div>;
          return <div>{ago.format("ddd, MMM D, YYYY")}</div>;
        },
      },
      ...(showIds
        ? [
            {
              Header: "edit",
              Cell: ({ row: { original } }: CellProps<i>) => {
                return <Code>{original.id} </Code>;
              },
            },
          ]
        : []),
    ],
    [checked, showOlder, showIds]
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
              to={`/diff?${queryString.stringify({
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
