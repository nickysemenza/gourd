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
import { RecipeLink, UnitMappingList } from "../components/Misc";
import { AlertTriangle, PlusCircle } from "react-feather";
import { Code } from "../util";
import FoodSearch from "../components/FoodSearch";
import { Link } from "react-router-dom";
import { Helmet } from "react-helmet";

const IngredientList: React.FC = () => {
  const showIDs = false;

  let initialParams: PaginationParameters = {
    offset: 0,
    limit: 2,
  };

  const [params, setParams] = useState(initialParams);

  const fetchData = React.useCallback((params: PaginationParameters) => {
    setParams(params);
  }, []);

  const { data, loading } = useListIngredients({
    queryParams: params,
  });

  const [onlyMissingFDC, setOnlyMissingFDC] = useState(false);
  const [justLinked, setJustLinked] = useState<string[]>([]);

  const ingredients = (data?.ingredients || []).filter((i) =>
    onlyMissingFDC ? !i.food && !justLinked.includes(i.ingredient.id) : true
  );

  type i = (typeof ingredients)[0];

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
        Header: "Name",
        width: 30,
        Cell: ({ row: { original } }: CellProps<i>) => {
          const { ingredient, children } = original;
          return (
            <div className="flex flex-col w-20 whitespace-normal">
              {showIDs && <Code>{original.ingredient.id}</Code>}

              <div className="text-md font-medium text-gray-900">
                {ingredient.name}
              </div>

              {(children || []).map((i) => (
                <React.Fragment key={i.ingredient.id}>
                  <div className="text-sm text-gray-500 pl-1">
                    {i.ingredient.name}
                  </div>
                </React.Fragment>
              ))}
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
            <div className="w-52">
              <ul className="list-disc list-outside pl-4">
                {recipes.map((r) => (
                  <li key={`${original.ingredient.id}@${r.name}@${r.version}`}>
                    <RecipeLink recipe={r} />
                  </li>
                ))}
                {children.map((r) => (
                  <div key={`${r.ingredient.id}`}>
                    <li className="italic">{r.ingredient.name}</li>
                    <ul className="list-disc list-outside pl-4">
                      {(r.recipes || []).map((r) => (
                        <li
                          key={`${original.ingredient.id}@${r.name}@${r.version}`}
                        >
                          <RecipeLink recipe={r} />
                        </li>
                      ))}
                    </ul>
                  </div>
                ))}
              </ul>
            </div>
          );
        },
      },
      {
        Header: "Units",
        Cell: ({ row: { original } }: CellProps<i>) => {
          const { unit_mappings } = original;
          return (
            <div className="flex flex-col">
              <UnitMappingList unit_mappings={unit_mappings} includeDot />
            </div>
          );
        },
      },
      {
        Header: "Actions",
        id: "actions",
        Cell: ({ row: { original } }: CellProps<i>) => {
          const { ingredient } = original;
          return (
            <div className="flex flex-col">
              <ButtonGroup
                compact
                buttons={[
                  {
                    disabled: ingredient.fdc_id != null,
                    onClick: () => {
                      convertToRecipe(ingredient.id);
                    },
                    text: "convert to recipe",
                    IconLeft: PlusCircle,
                  },
                ]}
              />
              <Link to={`/ingredients/${ingredient.id}`} className="link">
                <div className="text-blue-800">view detail</div>
              </Link>
            </div>
          );
        },
      },

      {
        Header: "USDA Food",
        id: "food",
        Cell: ({ row: { original } }: CellProps<i>) => {
          return (
            <div className="flex flex-col w-full">
              <FoodSearch
                enableSearch={onlyMissingFDC}
                name={original.ingredient.name}
                highlightId={original.food?.wrapper.fdc_id}
                onLink={(fdcId: number) => {
                  linkFoodToIngredient(original.ingredient.id, fdcId);
                  setJustLinked([...justLinked, original.ingredient.id]);
                }}
                addon={original.food}
              />
            </div>
          );
        },
      },
    ];
  }, [onlyMissingFDC, showIDs, justLinked]);

  return (
    <>
      {/* @ts-ignore */}
      <Helmet>
        <title>ingredients | gourd</title>
      </Helmet>
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
        isLoading={loading}
        totalCount={data?.meta?.total_count || 0}
        pageCount={data?.meta?.page_count || 1}
      />
    </>
  );
};

export default IngredientList;
