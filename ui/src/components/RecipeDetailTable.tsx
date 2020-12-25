import React, { useCallback } from "react";
import IngredientSearch from "./IngredientSearch";
import { Link } from "react-router-dom";
import {
  RecipeWrapper,
  Ingredient,
  RecipeSection,
} from "../api/openapi-hooks/api";
import { formatText, formatTimeRange, getIngredient } from "../util";
import { DndProvider } from "react-dnd";
import { HTML5Backend } from "react-dnd-html5-backend";
import update from "immutability-helper";
import { Button } from "./Button";
import { DragWrapper } from "./DragDrop";

export interface UpdateIngredientProps {
  sectionID: number;
  ingredientID: number;
  value: string;
  attr: "grams" | "name" | "amount" | "unit" | "adjective" | "optional";
}

export interface TableProps {
  recipe: RecipeWrapper;
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
    value: number,
    attr: "grams" | "amount"
  ) => number;
  isOverride: (
    sectionID: number,
    ingredientID: number,
    attr: "grams" | "amount"
  ) => boolean;
  edit: boolean;
  addInstruction: (sectionID: number) => void;
  addIngredient: (sectionID: number) => void;
  addSection: () => void;
  setRecipe: React.Dispatch<React.SetStateAction<RecipeWrapper | null>>;
}
const RecipeDetailTable: React.FC<TableProps> = ({
  recipe,
  updateIngredient,
  updateIngredientInfo,
  updateInstruction,
  getIngredientValue,
  isOverride,
  edit,
  addInstruction,
  addIngredient,
  addSection,
  setRecipe,
}) => {
  const { sections } = recipe.detail;
  // for baker's percentage cauclation we need the total mass of all flours (which together are '100%')
  const flourMass = (sections || []).reduce(
    (acc, section) =>
      acc +
      section.ingredients
        .filter((item) => item.ingredient?.name.includes("flour"))
        .reduce((acc, ingredient) => acc + ingredient?.grams, 0),
    0
  );
  const showBP = flourMass > 0;

  type I = "ingredients" | "instructions";
  const calculateMoveI = (
    sectionIndex: number,
    index: number,
    movingUp: boolean,
    i: I
  ) => {
    const numI = sections[sectionIndex][i].length;
    const numSections = sections.length;
    const firstInSection = index === 0;
    const lastInSection = index === numI - 1;

    let newSectionIndex = sectionIndex;
    let newInIndex: number;
    if (firstInSection && movingUp) {
      // needs to go to prior section
      newSectionIndex--;
      if (newSectionIndex < 0) {
        // out of bounds
        return null;
      }
      newInIndex = sections[newSectionIndex][i].length;
    } else if (!firstInSection && movingUp) {
      // prior row in same section
      newInIndex = index - 1;
    } else if (lastInSection && !movingUp) {
      // needs to go to next section
      newSectionIndex++;
      if (newSectionIndex > numSections - 1) {
        // out of bounds
        return null;
      }
      newInIndex = 0;
    } else {
      // next row in same section
      newInIndex = index + 1;
    }

    return { newSectionIndex, newInIndex };
  };
  const canMoveI = (
    sectionIndex: number,
    index: number,
    movingUp: boolean,
    i: I
  ) => !!calculateMoveI(sectionIndex, index, movingUp, i);
  const moveI = (
    sectionIndex: number,
    index: number,
    movingUp: boolean,
    i: I
  ) => {
    const coords = calculateMoveI(sectionIndex, index, movingUp, i);
    if (!coords) return;
    const { newSectionIndex, newInIndex } = coords;
    console.log("moving!", {
      sectionIndex,
      newSectionIndex,
      index,
      newInIndex,
    });
    const target = sections[sectionIndex][i][index];
    setRecipe(
      update(recipe, {
        detail: {
          sections:
            newSectionIndex === sectionIndex
              ? {
                  [sectionIndex]: {
                    [i]: {
                      $splice: [
                        [index, 1],
                        [newInIndex, 0, target],
                      ],
                    },
                  },
                }
              : {
                  [sectionIndex]: {
                    [i]: {
                      $splice: [[index, 1]],
                    },
                  },
                  [newSectionIndex]: {
                    [i]: {
                      $splice: [[newInIndex, 0, target]],
                    },
                  },
                },
        },
      })
    );
  };
  const delI = (sectionIndex: number, index: number, i: I) =>
    setRecipe(
      update(recipe, {
        detail: {
          sections: {
            [sectionIndex]: {
              [i]: {
                $splice: [[index, 1]],
              },
            },
          },
        },
      })
    );
  const iActions = (x: number, y: number, i: I) =>
    edit && (
      <div className="flex space-x-1">
        <Button
          onClick={() => {
            moveI(x, y, true, i);
          }}
          disabled={!canMoveI(x, y, true, i)}
          label={`up`}
        />
        <Button
          onClick={() => {
            moveI(x, y, false, i);
          }}
          disabled={!canMoveI(x, y, false, i)}
          label={`down`}
        />
        <Button
          onClick={() => {
            delI(x, y, i);
          }}
          label={`del`}
        />
      </div>
    );

  const renderRow = (section: RecipeSection, x: number) => (
    <TableRow key={x}>
      <TableCell>
        <div className="inline-block bg-blue-200 text-blue-800 text-xs px-2 rounded-full uppercase font-semibold tracking-wide">
          {String.fromCharCode(65 + x)}
        </div>
      </TableCell>
      <TableCell>{formatTimeRange(section.duration)}</TableCell>
      <TableCell>
        {section.ingredients.map((ingredient, y) => {
          const bp = Math.round((ingredient.grams / flourMass) * 100);
          return (
            <div className="ing-table-row" key={y}>
              <TableInput
                width={14}
                data-cy="grams-input"
                edit={edit}
                softEdit
                value={getIngredientValue(x, y, ingredient.grams || 0, "grams")}
                blur
                highlight={isOverride(x, y, "grams")}
                onChange={(e) =>
                  updateIngredient({
                    sectionID: x,
                    ingredientID: y,
                    value: e,
                    attr: "grams",
                  })
                }
              />
              <div className="flex space-x-0.5">
                <div className="text-gray-600">g</div>
                {showBP && (
                  <div
                    className={`${
                      bp > 0 ? "text-gray-600" : "text-red-300"
                    } italic`}
                  >
                    ({bp}%)
                  </div>
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
                      {ingredient.recipe?.name}
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
                highlight={isOverride(x, y, "amount")}
                value={getIngredientValue(
                  x,
                  y,
                  ingredient.amount || 0,
                  "amount"
                )}
                onChange={(e) =>
                  updateIngredient({
                    sectionID: x,
                    ingredientID: y,
                    value: e,
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
                    value: e,
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
                    value: e,
                    attr: "adjective",
                  })
                }
              />
              <div>{iActions(x, y, "ingredients")}</div>
              {/* TODO: optional toggle */}
            </div>
          );
        })}
        {edit && (
          <div className="add-item" onClick={() => addIngredient(x)}>
            add ingredient
          </div>
        )}
      </TableCell>
      <TableCell>
        {/* <ol className="list-decimal list-inside"> */}
        {section.instructions.map((instruction, y) => (
          <div key={y} className="flex font-serif">
            <div className="mr-2 w-4">{y + 1}. </div>
            <TableInput
              data-cy="instruction-input"
              width={72}
              tall
              edit={edit}
              value={instruction.instruction}
              onChange={(e) => updateInstruction(x, y, e)}
            />
            <div>{iActions(x, y, "instructions")}</div>
          </div>
        ))}
        {/* </ol> */}
        {edit && (
          <div className="add-item" onClick={() => addInstruction(x)}>
            add instruction
          </div>
        )}
      </TableCell>
    </TableRow>
  );

  const moveCard = useCallback(
    (dragIndex: number, hoverIndex: number) => {
      const dragCard = sections[dragIndex];
      console.log("drag", { dragIndex, hoverIndex, dragCard });

      setRecipe(
        update(recipe, {
          detail: {
            sections: {
              $splice: [
                [dragIndex, 1],
                [hoverIndex, 0, dragCard],
              ],
            },
          },
        })
      );
    },
    [sections, recipe, setRecipe]
  );

  return (
    <DndProvider backend={HTML5Backend}>
      <div className="border-gray-900 shadow-xl bg-gray-100">
        <TableRow header>
          <TableCell>Section</TableCell>
          <TableCell>Minutes</TableCell>
          <TableCell>
            <div className="ing-table-row font-mono">
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
        {sections?.map((section, x) => (
          <DragWrapper
            key={section.id}
            index={x}
            id={section.id}
            moveCard={moveCard}
            enable={edit}
          >
            {renderRow(section, x)}
          </DragWrapper>
        ))}
        {edit && (
          <div className="add-item" onClick={() => addSection()}>
            add section
          </div>
        )}
      </div>
    </DndProvider>
  );
};
export default RecipeDetailTable;

const TableCell: React.FC = ({ children }) => (
  <div className="border-solid border border-gray-600 p-1">{children}</div>
);
const TableRow: React.FC<{ header?: boolean }> = ({
  children,
  header = false,
}) => (
  <div className={`rec-table-row ${header && "font-semibold"}`}>{children}</div>
);

// <input> can't do onBlur?
type TallOrBlur =
  | {
      tall: true;
      blur?: false;
    }
  | {
      blur?: boolean;
      tall?: false;
    };

type TableInputProps = {
  edit: boolean;
  softEdit?: boolean;
  value: string | number;
  width?: number;
  highlight?: boolean;
  onChange: (event: string) => void;
} & TallOrBlur;

const TableInput: React.FC<TableInputProps> = ({
  edit,
  softEdit = false,
  width = 10,
  tall = false,
  blur = false,
  highlight = false,
  value,
  onChange,
  ...props
}) => {
  const controlledVal = value.toString();
  const [internalVal, setVal] = React.useState(controlledVal);
  React.useEffect(() => {
    setVal(controlledVal);
  }, [controlledVal]);

  const oC = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setVal(e.target.value);
    if (!blur) {
      onChange(e.target.value);
    }
  };
  const oB = (e: React.FocusEvent<HTMLInputElement>) => {
    if (!blur) {
      return;
    }
    if (internalVal !== controlledVal) {
      onChange(internalVal);
    }
  };

  const className = `border-2 border-dashedp-0 h-${tall ? 18 : 6} w-${width} ${
    highlight ? "border-blue-400" : "border-gray-200"
  } disabled:border-red-100 hover:border-black ${
    softEdit && !edit && "bg-transparent"
  } focus:bg-gray-200`;

  return edit || softEdit ? (
    tall ? (
      <textarea
        {...props}
        value={internalVal}
        onChange={oC}
        className={className}
        rows={3}
      />
    ) : (
      <input
        {...props}
        value={internalVal}
        onChange={oC}
        onBlur={oB}
        className={className}
        disabled={!edit && controlledVal === "0"}
      />
    )
  ) : (
    <p className="flex flex-wrap">{formatText(internalVal)}</p>
  );
};
