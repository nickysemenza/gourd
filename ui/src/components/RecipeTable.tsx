import React from "react";
import IngredientSearch from "./IngredientSearch";
import { Link } from "react-router-dom";
import {
  RecipeDetail,
  Ingredient,
  RecipeSection,
} from "../api/openapi-hooks/api";
import { getIngredient } from "../util";
export interface UpdateIngredientProps {
  sectionID: number;
  ingredientID: number;
  value: string;
  attr: "grams" | "name" | "amount" | "unit" | "adjective" | "optional";
}

export interface TableProps {
  recipe: RecipeDetail;
  updateIngredient: (i: UpdateIngredientProps) => void;
  updateIngredientInfo: (
    sectionID: number,
    ingredientID: number,
    ingredient: Pick<Ingredient, "id" | "name">,
    kind: "recipe" | "ingredient"
  ) => void;
  updateInstruction: (
    sectionID: number,
    instructionID: number,
    value: string
  ) => void;

  getIngredientValue: (
    sectionID: number,
    ingredientID: number,
    value: number
  ) => number;
  edit: boolean;
  addInstruction: (sectionID: number) => void;
  addIngredient: (sectionID: number) => void;
  addSection: () => void;
}
const RecipeTable: React.FC<TableProps> = ({
  recipe,
  updateIngredient,
  updateIngredientInfo,
  updateInstruction,
  getIngredientValue,
  edit,
  addInstruction,
  addIngredient,
  addSection,
}) => {
  // for baker's percentage cauclation we need the total mass of all flours (which together are '100%')
  const flourMass = (recipe.sections || []).reduce(
    (acc, section) =>
      acc +
      section.ingredients
        .filter((item) => item.ingredient?.name.includes("flour"))
        .reduce((acc, ingredient) => acc + ingredient?.grams, 0),
    0
  );

  const renderRow = (section: RecipeSection, x: number) => (
    <TableRow key={x}>
      <TableCell>
        <div className="inline-block bg-blue-200 text-blue-800 text-xs px-2 rounded-full uppercase font-semibold tracking-wide">
          {String.fromCharCode(65 + x)}
        </div>
      </TableCell>
      <TableCell>{section.minutes}</TableCell>
      <TableCell>
        {section.ingredients.map((ingredient, y) => (
          <div className="ing-table-row" key={y}>
            <TableInput
              data-cy="grams-input"
              edit={edit}
              softEdit
              value={getIngredientValue(x, y, ingredient.grams || 0)}
              onChange={(e) =>
                updateIngredient({
                  sectionID: x,
                  ingredientID: y,
                  value: e.target.value,
                  attr: "grams",
                })
              }
            />
            <div className="text-gray-600">
              g
              {flourMass > 0 && (
                <i> ({Math.round((ingredient.grams / flourMass) * 100)}%)</i>
              )}
            </div>
            {edit ? (
              <IngredientSearch
                initial={getIngredient(ingredient).name}
                callback={(item, kind) =>
                  updateIngredientInfo(x, y, item, kind)
                }
              />
            ) : (
              <div className="text-gray-600">
                {ingredient.kind === "recipe" ? (
                  <Link
                    to={`/recipe/${ingredient.recipe?.id}`}
                    className="link"
                  >
                    {ingredient.ingredient?.name}
                  </Link>
                ) : (
                  ingredient.ingredient?.name
                )}
              </div>
            )}
            <TableInput
              data-cy="amount-input"
              // width={16}
              edit={edit}
              softEdit
              value={getIngredientValue(x, y, ingredient.amount || 0)}
              onChange={(e) =>
                updateIngredient({
                  sectionID: x,
                  ingredientID: y,
                  value: e.target.value,
                  attr: "amount",
                })
              }
            />
            <TableInput
              data-cy="unit-input"
              width={16}
              edit={edit}
              value={ingredient.unit}
              onChange={(e) =>
                updateIngredient({
                  sectionID: x,
                  ingredientID: y,
                  value: e.target.value,
                  attr: "unit",
                })
              }
            />
            <TableInput
              data-cy="adjective-input"
              width={16}
              edit={edit}
              value={ingredient.adjective}
              onChange={(e) =>
                updateIngredient({
                  sectionID: x,
                  ingredientID: y,
                  value: e.target.value,
                  attr: "adjective",
                })
              }
            />
            {/* TODO: optional toggle */}
          </div>
        ))}
        {edit && (
          <div className="add-item" onClick={() => addIngredient(x)}>
            add ingredient
          </div>
        )}
      </TableCell>
      <TableCell>
        <ol className="list-decimal list-inside">
          {section.instructions.map((instruction, y) => (
            <li key={y}>
              <TableInput
                data-cy="instruction-input"
                width={16}
                edit={edit}
                value={instruction.instruction}
                onChange={(e) => updateInstruction(x, y, e.target.value)}
              />
            </li>
          ))}
        </ol>
        {edit && (
          <div className="add-item" onClick={() => addInstruction(x)}>
            add instruction
          </div>
        )}
      </TableCell>
    </TableRow>
  );

  return (
    <div className="border-gray-900 shadow-xl bg-gray-100">
      <TableRow header>
        <TableCell>Section</TableCell>
        <TableCell>Minutes</TableCell>
        <TableCell>
          <div className="ing-table-row">
            <div>x</div>
            <div>grams (BP)</div>
            <div>of y</div>
            <div>z</div>
            <div>units</div>
            <div>modifier</div>
          </div>
        </TableCell>
        <TableCell>Instructions</TableCell>
      </TableRow>
      {recipe.sections?.map((section, x) => renderRow(section, x))}
      {edit && (
        <div className="add-item" onClick={() => addSection()}>
          add section
        </div>
      )}
    </div>
  );
};
export default RecipeTable;

const TableCell: React.FC = ({ children }) => (
  <div className="border-solid border border-gray-600 p-1">{children}</div>
);
const TableRow: React.FC<{ header?: boolean }> = ({
  children,
  header = false,
}) => (
  <div className={`rec-table-row ${header && "font-semibold"}`}>{children}</div>
);

const TableInput: React.FC<{
  edit: boolean;
  softEdit?: boolean;
  value: string | number;
  width?: number;
  onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
}> = ({ edit, softEdit = false, width = 10, ...props }) =>
  edit || softEdit ? (
    <input
      {...props}
      className={`border-2 border-dashed p-0 h-6 w-${width} border-gray-200 hover:border-black ${
        softEdit && !edit && "bg-transparent"
      } focus:bg-gray-200`}
    />
  ) : (
    <div>{props.value}</div>
  );
