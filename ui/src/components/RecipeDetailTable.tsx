import React, { useCallback } from "react";
import IngredientSearch from "./IngredientSearch";
import { Link } from "react-router-dom";
import { RecipeWrapper, RecipeSection } from "../api/openapi-hooks/api";
import {
  formatText,
  formatTimeRange,
  getIngredient,
  parseTimeRange,
} from "../util";
import { DndProvider } from "react-dnd";
import { HTML5Backend } from "react-dnd-html5-backend";
import update from "immutability-helper";
import { ButtonGroup, Pill } from "./Button";
import { DragWrapper } from "./DragDrop";
import {
  addIngredient,
  addInstruction,
  addSection,
  canMoveI,
  delI,
  getIngredientValue,
  getGlobalInstructionNumber,
  I,
  isOverride,
  moveI,
  RecipeTweaks,
  updateIngredientInfo,
  updateInstruction,
  updateTimeRange,
} from "./RecipeEditorUtils";
import { ArrowDown, ArrowUp, XSquare } from "react-feather";

export interface UpdateIngredientProps {
  sectionID: number;
  ingredientID: number;
  value: string;
  attr: "grams" | "name" | "amount" | "unit" | "adjective" | "optional";
}

export interface TableProps {
  recipe: RecipeWrapper;
  updateIngredient: (i: UpdateIngredientProps) => void;
  setRecipe: React.Dispatch<React.SetStateAction<RecipeWrapper | null>>;
  tweaks: RecipeTweaks;
}
const RecipeDetailTable: React.FC<TableProps> = ({
  recipe,
  updateIngredient,
  setRecipe,
  tweaks,
}) => {
  const { edit } = tweaks;
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

  const iActions = (x: number, y: number, i: I) =>
    edit && (
      <ButtonGroup
        compact
        buttons={[
          {
            onClick: () => {
              setRecipe(moveI(recipe, x, y, true, i));
            },
            disabled: !canMoveI(recipe, x, y, true, i),
            IconLeft: ArrowUp,
          },
          {
            onClick: () => {
              setRecipe(moveI(recipe, x, y, false, i));
            },
            disabled: !canMoveI(recipe, x, y, false, i),
            IconLeft: ArrowDown,
          },
          {
            onClick: () => {
              setRecipe(delI(recipe, x, y, i));
            },
            IconLeft: XSquare,
          },
        ]}
      />
    );

  const renderRow = (section: RecipeSection, x: number) => (
    <TableRow key={x}>
      <TableCell>
        <Pill>{String.fromCharCode(65 + x)}</Pill>
      </TableCell>
      <TableCell>
        <TableInput
          width={40}
          data-cy="grams-input"
          edit={edit}
          value={formatTimeRange(section.duration)}
          blur
          onChange={(e) =>
            setRecipe(
              updateTimeRange(recipe, x, parseTimeRange(e) || section.duration)
            )
          }
        />
      </TableCell>
      <TableCell>
        {section.ingredients.map((ingredient, y) => {
          const bp = Math.round((ingredient.grams / flourMass) * 100);
          const { edit } = tweaks;
          return (
            <div className="flex flex-col">
              <div className="ing-table-row" key={y}>
                <TableInput
                  width={14}
                  data-cy="grams-input"
                  edit={edit}
                  softEdit
                  value={getIngredientValue(
                    tweaks,
                    x,
                    y,
                    ingredient.grams || 0,
                    "grams"
                  )}
                  blur
                  highlight={isOverride(tweaks, x, y, "grams")}
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
                      setRecipe(updateIngredientInfo(recipe, x, y, item, kind))
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
                  highlight={isOverride(tweaks, x, y, "amount")}
                  value={getIngredientValue(
                    tweaks,
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
                  value={ingredient.unit || ""}
                  onChange={(e) =>
                    updateIngredient({
                      sectionID: x,
                      ingredientID: y,
                      value: e,
                      attr: "unit",
                    })
                  }
                />
                <div className="flex">
                  <TableInput
                    data-cy="adjective-input"
                    width={16}
                    edit={edit}
                    value={ingredient.adjective || ""}
                    onChange={(e) =>
                      updateIngredient({
                        sectionID: x,
                        ingredientID: y,
                        value: e,
                        attr: "adjective",
                      })
                    }
                  />
                </div>
                <div>{iActions(x, y, "ingredients")}</div>
                {/* TODO: optional toggle */}
              </div>
              {!!ingredient.original && (
                <div className="italic text-xs pb-2">
                  original: {ingredient.original}
                </div>
              )}
            </div>
          );
        })}
        {edit && (
          <div
            className="add-item"
            onClick={() => setRecipe(addIngredient(recipe, x))}
          >
            add ingredient
          </div>
        )}
      </TableCell>
      <TableCell>
        {/* <ol className="list-decimal list-inside"> */}
        {section.instructions.map((instruction, y) => (
          <div key={y} className="flex font-serif">
            <div className="mr-4 w-4">
              <Pill>{getGlobalInstructionNumber(recipe, x, y)}</Pill>
            </div>
            <TableInput
              data-cy="instruction-input"
              width={72}
              tall
              edit={edit}
              value={instruction.instruction}
              onChange={(e) =>
                setRecipe(updateInstruction(recipe, tweaks, x, y, e))
              }
            />
            <div>{iActions(x, y, "instructions")}</div>
          </div>
        ))}
        {/* </ol> */}
        {edit && (
          <div
            className="add-item"
            onClick={() => setRecipe(addInstruction(recipe, x))}
          >
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
          <TableCell></TableCell>
          <TableCell>Time</TableCell>
          <TableCell>
            <div className="ing-table-row font-mono">
              <div>x</div>
              <div>g {showBP && "(BP)"}</div>
              <div>of y</div>
              <div>[z</div>
              <div>units],</div>
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
          <div
            className="add-item"
            onClick={() => setRecipe(addSection(recipe))}
          >
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
