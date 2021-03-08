import React, { useState } from "react";
import { Column, CellProps } from "react-table";
import PaginatedTable, {
  PaginationParameters,
} from "../components/PaginatedTable";
import { useListIngredients } from "../api/openapi-hooks/api";
import { IngredientsApi } from "../api/openapi-fetch";
import { getOpenapiFetchConfig } from "../config";
import { toast } from "react-toastify";
import { ButtonGroup } from "../components/Button";
import { RecipeLink } from "../components/Misc";
import { AlertTriangle, PlusCircle } from "react-feather";
import { Code } from "../util";
import FoodSearch from "../components/FoodSearch";
import { Link } from "react-router-dom";

const IngredientList: React.FC = () => {
  let initialParams: PaginationParameters = {
    offset: 0,
    limit: 2,
  };

  const [params, setParams] = useState(initialParams);

  const fetchData = React.useCallback((params: PaginationParameters) => {
    setParams(params);
  }, []);

  const { data } = useListIngredients({
    queryParams: params,
  });

  const [onlyMissingFDC, setOnlyMissingFDC] = useState(false);

  const ingredients = (data?.ingredients || []).filter((i) =>
    onlyMissingFDC ? !i.food : true
  );

  type i = typeof ingredients[0];

  const columns: Array<Column<i>> = React.useMemo(() => {
    const iApi = new IngredientsApi(getOpenapiFetchConfig());

    const convertToRecipe = async (id: string) => {
      let res = await iApi.convertIngredientToRecipe({ ingredientId: id });
      toast.success(`created recipe ${res.id} for ${res.name}`);
    };
    const linkFoodToIngredient = async (
      ingredientId: string,
      fdcId: number
    ) => {
      await iApi.associateFoodWithIngredient({ ingredientId, fdcId });
      toast.success(`linked ${ingredientId} to food ${fdcId}`);
    };

    return [
      {
        Header: "id",
        Cell: ({ row: { original } }: CellProps<i>) => (
          <Code>{original.ingredient.id}</Code>
        ),
      },
      {
        Header: "Name",
        // accessor: "name",
        Cell: ({ row: { original } }: CellProps<i>) => {
          const { ingredient, children, food } = original;
          return (
            <div className="flex flex-col w-64">
              {ingredient.name}
              <ButtonGroup
                compact
                buttons={[
                  {
                    onClick: () => {
                      convertToRecipe(ingredient.id);
                    },
                    text: "make Recipe",
                    IconLeft: PlusCircle,
                  },
                ]}
              />
              <ul>
                {(children || []).map((i) => (
                  <li className="pl-6 flex">
                    aka. <div className="italic pl-1">{i.ingredient.name}</div>
                  </li>
                ))}
              </ul>
              <div>fdc: {!!food ? food.description : "n/a"}</div>
              <Link to={`/ingredients/${ingredient.id}`} className="link">
                <div className="text-blue-800">view</div>
              </Link>
            </div>
          );
        },
      },

      {
        Header: "Recipes",
        id: "recipes",
        Cell: ({ row: { original } }: CellProps<i>) => {
          const recipes = original.recipes || [];
          const children = original.children || [];
          return (
            <div>
              <ul className="list-disc list-outside pl-4">
                {recipes.map((r) => (
                  <li key={`${original.ingredient.id}@${r.name}@${r.version}`}>
                    <RecipeLink recipe={r} />
                  </li>
                ))}
                {children.map((r) => (
                  <div>
                    <li className="italic">{r.ingredient.name}</li>
                    <ul className="list-disc list-outside pl-4">
                      {(r.recipes || []).map((r) => (
                        <li
                          // className="pl-6"
                          key={`${original.ingredient.id}@${r.name}@${r.version}`}
                        >
                          <RecipeLink recipe={r} />
                        </li>
                      ))}
                    </ul>
                  </div>
                ))}
              </ul>
              {/* <Debug data={original} /> */}
            </div>
          );
        },
      },
      {
        Header: "USDA Food",
        id: "food",
        Cell: ({ row: { original } }: CellProps<i>) => {
          return (
            <FoodSearch
              enableSearch={onlyMissingFDC}
              name={original.ingredient.name}
              highlightId={original.food?.fdc_id}
              onLink={(fdcId: number) => {
                linkFoodToIngredient(original.ingredient.id, fdcId);
              }}
            />
          );
        },
      },
    ];
  }, [onlyMissingFDC]);

  return (
    <>
      <ButtonGroup
        compact
        buttons={[
          {
            text: onlyMissingFDC ? "show all" : "only show missing FDC",
            onClick: () => {
              setOnlyMissingFDC(!onlyMissingFDC);
            },
            IconLeft: AlertTriangle,
          },
        ]}
      />
      <PaginatedTable
        columns={columns}
        data={ingredients}
        fetchData={fetchData}
        isLoading={false}
        totalCount={data?.meta?.total_count || 0}
        pageCount={data?.meta?.page_count || 1}
      />
    </>
  );
};

export default IngredientList;
