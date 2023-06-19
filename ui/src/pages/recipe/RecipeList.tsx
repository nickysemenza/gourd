import React, { useState } from "react";
import Debug from "../../components/ui/Debug";
import { RecipeWrapper, useListRecipes } from "../../api/openapi-hooks/api";
import PaginatedTable, {
  PaginationParameters,
} from "../../components/ui/PaginatedTable";
import { RecipeLink } from "../../components/misc/Misc";
import { Code } from "../../components/Code";
import { Helmet } from "react-helmet";
import { useNavigate } from "react-router-dom";
import queryString from "query-string";
import { ButtonGroup } from "../../components/ui/ButtonGroup";
import { Grid, List, PlusCircle } from "react-feather";
import ProgressiveImage from "../../components/ui/ProgressiveImage";
import { sumIngredients } from "../../components/recipe/RecipeEditorUtils";
import dayjs from "dayjs";
import { RecipeGrid } from "../../components/recipe/RecipeGrid";
import { createColumnHelper } from "@tanstack/react-table";
import { Pill } from "../../components/ui/Pill";
import { Button } from "../../components/ui/Button";
import { Checkbox } from "../../components/ui/Checkbox";

const RecipeList: React.FC = () => {
  const showIds = false;
  const initialParams: PaginationParameters = {
    offset: 0,
    limit: 180,
  };

  const [params, setParams] = useState(initialParams);
  const [showOlder, setShowOlder] = useState(false);
  const [showEmpty, setShowEmpty] = useState(false);
  const [grid, setGrid] = useState(true);

  const fetchData = React.useCallback((params: PaginationParameters) => {
    setParams(params);
  }, []);
  const navigate = useNavigate();

  const { data, error, loading } = useListRecipes({
    queryParams: params,
  });

  const recipes = (data?.recipes || []).filter(
    (r) =>
      showEmpty ||
      r.detail.sections.length > 0 ||
      r.detail.name.startsWith("cy-")
  );

  const future = new Set<RecipeWrapper>();
  const past = new Set<RecipeWrapper>();
  const other = new Set<RecipeWrapper>();
  recipes.forEach((r) => {
    let added = false;
    (r.linked_meals || []).forEach((m) => {
      const ago = dayjs(m.ate_at);
      if (ago.isAfter(dayjs())) {
        future.add(r);
        added = true;
      }
      if (ago.isBefore(dayjs())) {
        past.add(r);
        added = true;
      }
    });
    if (!added) {
      other.add(r);
    }
  });
  console.log(
    `future: ${future.size}, past: ${past.size}, other: ${other.size} (total ${recipes.length})`
  );

  type i = (typeof recipes)[0];
  const columnHelper = createColumnHelper<i>();

  const columns = [
    columnHelper.accessor((row) => row, {
      id: "select",
      header: ({ table }) => (
        <Checkbox
          checked={table.getIsAllPageRowsSelected()}
          onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
          aria-label="Select all"
        />
      ),
      cell: ({ row }) => (
        <Checkbox
          checked={row.getIsSelected()}
          onCheckedChange={(value) => row.toggleSelected(!!value)}
          aria-label="Select row"
        />
      ),
      enableSorting: false,
      enableHiding: false,
    }),
    columnHelper.accessor((row) => row, {
      id: "name",
      cell: (info) => {
        const { original } = info.row;
        const olderVersions = original.other_versions || [];
        const versions = [original.detail, ...(showOlder ? olderVersions : [])];
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
                  <RecipeLink recipe={i} />
                </li>
              ))}
            </ul>
            <Pill color={ing.length + rec.length > 0 ? "green" : "red"}>
              {ing.length} ing / {rec.length} rec
            </Pill>
          </div>
        );
      },
    }),
    columnHelper.accessor((row) => row, {
      id: "meals",
      cell: (info) => {
        const { original } = info.row;
        return (
          <div className="flex flex-row">
            {(original.linked_photos || []).map((p) => (
              <ProgressiveImage photo={p} key={p.id} className="w-1/6" />
            ))}
          </div>
        );
      },
    }),
    columnHelper.accessor((row) => row.detail.created_at, {
      id: "created at",
      cell: (info) => {
        const ago = dayjs(info.getValue());
        return <div>{ago.format("ddd, MMM D, YYYY")}</div>;
      },
      header: () => <span>Created At</span>,
    }),
    ...(showIds
      ? [
          columnHelper.accessor((row) => row.id, {
            id: "id",
            cell: (info) => {
              return <Code>{info.getValue()} </Code>;
            },
            header: () => <span>id</span>,
          }),
        ]
      : []),
  ];

  return (
    <div className="">
      <Helmet>
        <title>recipes | gourd</title>
      </Helmet>

      <div className="py-2">
        <ButtonGroup
          buttons={[
            {
              onClick: () => setGrid(true),
              disabled: grid,
              text: "grid",
              IconLeft: Grid,
            },
            {
              onClick: () => setGrid(false),
              disabled: !grid,
              text: "table",
              IconLeft: List,
            },
          ]}
        />
      </div>
      {grid ? (
        <>
          <RecipeGrid recipes={Array.from(future)} />
          <RecipeGrid recipes={Array.from(past)} />
          <RecipeGrid recipes={Array.from(other)} />
        </>
      ) : (
        <PaginatedTable
          columns={columns}
          data={recipes}
          fetchData={fetchData}
          isLoading={loading}
          totalCount={data?.meta?.total_count || 0}
          pageCount={data?.meta?.page_count || 1}
          withSelected={(selected) => {
            const ids = selected.map((r) => r.original.detail.id);
            return (
              <div>
                <Button
                  disabled={ids.length < 2}
                  onClick={() =>
                    navigate(
                      `/diff?${queryString.stringify({
                        recipes: ids,
                      })}`
                    )
                  }
                >
                  Compare {ids.length} recipes
                </Button>
              </div>
            );
          }}
        />
      )}

      <div>
        <ButtonGroup
          compact
          buttons={[
            {
              onClick: () => {
                setShowOlder(!showOlder);
              },
              text: showOlder ? "hide older" : "show older",
              IconLeft: PlusCircle,
            },
            {
              onClick: () => {
                setShowEmpty(!showEmpty);
              },
              text: showEmpty ? "hide empty" : "show empty",
              IconLeft: PlusCircle,
            },
          ]}
        />
      </div>
      <Debug data={{ error }} />
    </div>
  );
};

export default RecipeList;
