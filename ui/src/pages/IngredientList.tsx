import React, { useState } from "react";
import PaginatedTable, {
  PaginationParameters,
} from "../components/ui/PaginatedTable";
import { toast } from "react-toastify";
import { ButtonGroup } from "../components/ui/ButtonGroup";
import { RecipeLink, UnitMappingList } from "../components/misc/Misc";
import { AlertTriangle, PlusCircle } from "react-feather";
import { Code } from "../components/Code";
import FoodSearch from "../components/FoodSearch";
import { Link } from "react-router-dom";
import { createColumnHelper } from "@tanstack/react-table";
import PageWrapper from "../components/ui/PageWrapper";
import {
  fetchAssociateFoodWithIngredient,
  fetchConvertIngredientToRecipe,
  useListIngredients,
} from "../api/react-query/gourdApiComponents";

const IngredientList: React.FC = () => {
  const showIDs = false;

  const initialParams: PaginationParameters = {
    offset: 0,
    limit: 2,
  };

  const [params, setParams] = useState(initialParams);

  const fetchData = React.useCallback((params: PaginationParameters) => {
    setParams(params);
  }, []);

  const { data, isLoading } = useListIngredients({
    queryParams: params,
  });

  const [onlyMissingFDC, setOnlyMissingFDC] = useState(false);
  const [justLinked, setJustLinked] = useState<string[]>([]);

  const ingredients = (data?.ingredients || []).filter((i) =>
    onlyMissingFDC ? !i.food && !justLinked.includes(i.ingredient.id) : true
  );

  type i = (typeof ingredients)[0];

  const columns = React.useMemo(() => {
    const columnHelper = createColumnHelper<i>();

    const convertToRecipe = async (id: string) => {
      const res = await fetchConvertIngredientToRecipe({
        pathParams: { ingredientId: id },
      });
      toast.success(`created recipe ${res.id} for ${res.name}`);
    };
    const linkFoodToIngredient = async (
      ingredientId: string,
      fdcId: number
    ) => {
      await fetchAssociateFoodWithIngredient({
        pathParams: { ingredientId },
        queryParams: { fdc_id: fdcId },
      });
      toast.success(`linked ${ingredientId} to food ${fdcId}`);
    };

    return [
      columnHelper.accessor((row) => row, {
        id: "name",
        cell: (info) => {
          const { original } = info.row;
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
      }),
      columnHelper.accessor((row) => row, {
        id: "recipes",
        cell: (info) => {
          const { original } = info.row;
          const recipes = original.recipes || [];
          const children = original.children || [];
          return (
            <div className="w-52">
              <ul className="list-disc list-outside pl-4">
                {recipes.map((r) => (
                  <li
                    key={`${original.ingredient.id}@${r.name}@${r.meta.version}`}
                  >
                    <RecipeLink recipe={r} />
                  </li>
                ))}
                {children.map((r) => (
                  <div key={`${r.ingredient.id}`}>
                    <li className="italic">{r.ingredient.name}</li>
                    <ul className="list-disc list-outside pl-4">
                      {(r.recipes || []).map((r) => (
                        <li
                          key={`${original.ingredient.id}@${r.name}@${r.meta.version}`}
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
      }),
      columnHelper.accessor((row) => row.unit_mappings, {
        id: "units",
        cell: (info) => {
          return (
            <div className="flex flex-col">
              <UnitMappingList unit_mappings={info.getValue()} includeDot />
            </div>
          );
        },
      }),
      columnHelper.accessor((row) => row, {
        id: "meals",
        cell: (info) => {
          const { original } = info.row;
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
      }),

      columnHelper.accessor((row) => row, {
        header: () => <span>USDA Food</span>,
        id: "food",
        cell: (info) => {
          const { original } = info.row;
          return (
            <div className="flex flex-col w-full">
              <FoodSearch
                enableSearch={onlyMissingFDC}
                name={original.ingredient.name}
                highlightId={original.food?.wrapper.fdcId}
                onLink={(fdcId: number) => {
                  linkFoodToIngredient(original.ingredient.id, fdcId);
                  setJustLinked([...justLinked, original.ingredient.id]);
                }}
                addon={original.food}
              />
            </div>
          );
        },
      }),
    ];
  }, [onlyMissingFDC, showIDs, justLinked]);

  return (
    <PageWrapper title="ingredients">
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
        isLoading={isLoading}
        totalCount={data?.meta?.total_count || 0}
        pageCount={data?.meta?.page_count || 1}
      />
    </PageWrapper>
  );
};

export default IngredientList;
