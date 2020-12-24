import React, { useCallback, useRef } from "react";
import IngredientSearch from "./IngredientSearch";
import { Link } from "react-router-dom";
import {
  RecipeWrapper,
  Ingredient,
  RecipeSection,
} from "../api/openapi-hooks/api";
import { getIngredient } from "../util";
import {
  useDrag,
  useDrop,
  DropTargetMonitor,
  XYCoord,
  DndProvider,
} from "react-dnd";
import { HTML5Backend } from "react-dnd-html5-backend";
import update from "immutability-helper";

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
  edit,
  addInstruction,
  addIngredient,
  addSection,
  setRecipe,
}) => {
  // for baker's percentage cauclation we need the total mass of all flours (which together are '100%')
  const flourMass = (recipe.detail.sections || []).reduce(
    (acc, section) =>
      acc +
      section.ingredients
        .filter((item) => item.ingredient?.name.includes("flour"))
        .reduce((acc, ingredient) => acc + ingredient?.grams, 0),
    0
  );
  const showBP = flourMass > 0;

  const renderRow = (section: RecipeSection, x: number) => (
    <TableRow key={x}>
      <TableCell>
        <div className="inline-block bg-blue-200 text-blue-800 text-xs px-2 rounded-full uppercase font-semibold tracking-wide">
          {String.fromCharCode(65 + x)}
        </div>
      </TableCell>
      <TableCell>{section.minutes}</TableCell>
      <TableCell>
        {section.ingredients.map((ingredient, y) => {
          const bp = Math.round((ingredient.grams / flourMass) * 100);
          return (
            <div className="ing-table-row" key={y}>
              <TableInput
                data-cy="grams-input"
                edit={edit}
                softEdit
                value={getIngredientValue(x, y, ingredient.grams || 0, "grams")}
                onChange={(e) =>
                  updateIngredient({
                    sectionID: x,
                    ingredientID: y,
                    value: e.target.value,
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
              onChange={(e) => updateInstruction(x, y, e.target.value)}
            />
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

  const { sections } = recipe.detail;

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
          <Card key={section.id} index={x} id={section.id} moveCard={moveCard}>
            {renderRow(section, x)}
          </Card>
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

const TableInput: React.FC<{
  edit: boolean;
  softEdit?: boolean;
  value: string | number;
  width?: number;
  tall?: boolean;
  onChange: (
    event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => void;
}> = ({ edit, softEdit = false, width = 10, tall, ...props }) => {
  const className = `border-2 border-dashed p-0 h-${
    tall ? 18 : 6
  } w-${width} border-gray-200 disabled:border-red-100 hover:border-black ${
    softEdit && !edit && "bg-transparent"
  } focus:bg-gray-200`;
  return edit || softEdit ? (
    tall ? (
      <textarea {...props} className={className} rows={3} />
    ) : (
      <input
        {...props}
        className={className}
        disabled={!edit && props.value === 0}
      />
    )
  ) : (
    <p className="flex flex-wrap">{formatText(props.value)}</p>
  );
};
const re = /[\d]* ?F/g;
const formatText = (text: React.ReactText) => {
  if (typeof text === "number") {
    return text;
  }

  let pairs = [];
  const matches = [...text.matchAll(re)];
  if (matches.length === 0) {
    return text;
  }

  console.log(matches);
  // matches.next
  let lastProcessed = 0;
  for (const match of matches) {
    const matchStart = match.index || 0;
    const matchEnd = matchStart + match[0].length;
    pairs.push(text.substring(lastProcessed, matchStart));
    pairs.push(
      <code className="text-red-800 mx-1">
        {text.substring(matchStart, matchEnd)}
      </code>
    );
    // pairs.push()
    lastProcessed = matchEnd;
    // pairs.push([, ]);
  }
  pairs.push(text.substring(lastProcessed));
  // let res = [];
  // for
  return pairs;

  // console.log(pairs);
};

export interface CardProps {
  id: any;
  index: number;
  moveCard: (dragIndex: number, hoverIndex: number) => void;
}

interface DragItem {
  index: number;
  id: string;
  type: string;
}

export const Card: React.FC<CardProps> = ({
  id,
  index,
  moveCard,
  children,
}) => {
  const ref = useRef<HTMLDivElement>(null);
  const [, drop] = useDrop({
    accept: "card1",
    hover(item: DragItem, monitor: DropTargetMonitor) {
      if (!ref.current) {
        return;
      }
      const dragIndex = item.index;
      const hoverIndex = index;

      // Don't replace items with themselves
      if (dragIndex === hoverIndex) {
        return;
      }

      // Determine rectangle on screen
      const hoverBoundingRect = ref.current?.getBoundingClientRect();

      // Get vertical middle
      const hoverMiddleY =
        (hoverBoundingRect.bottom - hoverBoundingRect.top) / 2;

      // Determine mouse position
      const clientOffset = monitor.getClientOffset();

      // Get pixels to the top
      const hoverClientY = (clientOffset as XYCoord).y - hoverBoundingRect.top;

      // Only perform the move when the mouse has crossed half of the items height
      // When dragging downwards, only move when the cursor is below 50%
      // When dragging upwards, only move when the cursor is above 50%

      // Dragging downwards
      if (dragIndex < hoverIndex && hoverClientY < hoverMiddleY) {
        return;
      }

      // Dragging upwards
      if (dragIndex > hoverIndex && hoverClientY > hoverMiddleY) {
        return;
      }

      // Time to actually perform the action
      moveCard(dragIndex, hoverIndex);

      // Note: we're mutating the monitor item here!
      // Generally it's better to avoid mutations,
      // but it's good here for the sake of performance
      // to avoid expensive index searches.
      item.index = hoverIndex;
    },
  });

  const [{ isDragging }, drag] = useDrag({
    item: { type: "card1", id, index },
    collect: (monitor: any) => ({
      isDragging: monitor.isDragging(),
    }),
  });

  const opacity = isDragging ? 0 : 1;
  drag(drop(ref));
  return (
    <div ref={ref} style={{ opacity }}>
      {children}
    </div>
  );
};
