import React, { useCallback, useContext } from "react";
import {
  RecipeWrapper,
  RecipeSection,
  SectionIngredient,
} from "../api/openapi-hooks/api";
import { formatTimeRange, getIngredient, parseTimeRange } from "../util";
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
  adjustIngredientValue,
  getGlobalInstructionNumber,
  I,
  isOverride,
  moveI,
  RecipeTweaks,
  updateIngredientInfo,
  updateInstruction,
  updateTimeRange,
  FoodsById,
  IngDetailsById,
  inferGrams,
  getCal2,
} from "./RecipeEditorUtils";
import { ArrowDown, ArrowUp, PlusCircle, XSquare } from "react-feather";
import { RecipeLink } from "./Misc";
import { EntitySelector } from "./EntitySelector";
import { WasmContext } from "../wasm";
import { TableInput } from "./Input";
import { UnitConversionRequestTargetEnum } from "../api/openapi-fetch";

export interface UpdateIngredientProps {
  sectionID: number;
  ingredientID: number;
  subIndex?: number;
  value: string;
  attr: "grams" | "name" | "amount" | "unit" | "adjective" | "optional";
}

export interface TableProps {
  recipe: RecipeWrapper;
  updateIngredient: (i: UpdateIngredientProps) => void;
  setRecipe: React.Dispatch<React.SetStateAction<RecipeWrapper | null>>;
  tweaks: RecipeTweaks;
  hints: FoodsById;
  ing_hints: IngDetailsById;
}
const RecipeDetailTable: React.FC<TableProps> = ({
  recipe,
  updateIngredient,
  setRecipe,
  tweaks,
  hints,
  ing_hints,
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

  const instance = useContext(WasmContext);

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

  const renderIngredientItem = (
    ingredient: SectionIngredient,
    x: number,
    y: number,
    subIndex?: number
  ) => {
    const bp = Math.round((ingredient.grams / flourMass) * 100);
    const { edit } = tweaks;
    const isSub = subIndex !== undefined;
    return (
      <div className="flex flex-col">
        <div className={`ing-table-row`} key={y}>
          <div className="flex space-x-0.5">
            {isSub && (
              <div className="text-sm uppercase self-center text-gray-900">
                or
              </div>
            )}
            <TableInput
              width={"full"}
              data-cy="grams-input"
              edit={edit}
              softEdit
              value={adjustIngredientValue(
                tweaks,
                x,
                y,
                subIndex,
                ingredient.grams || 0,
                "grams"
              )}
              pValue={
                instance &&
                inferGrams(instance, ingredient, ing_hints) &&
                tweaks.multiplier *
                  (inferGrams(instance, ingredient, ing_hints) || 0)
              }
              blur
              highlight={isOverride(tweaks, x, y, subIndex, "grams")}
              onChange={(e) =>
                updateIngredient({
                  sectionID: x,
                  ingredientID: y,
                  subIndex,
                  value: e,
                  attr: "grams",
                })
              }
            />
            {/* {instance &&
              inferGrams(instance, ingredient, ing_hints) &&
              scaledRound(inferGrams(instance, ingredient, ing_hints) || 0)} */}
          </div>
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
            <EntitySelector
              createKind="ingredient"
              value={{
                value: "",
                label: getIngredient(ingredient).name,
                kind: "recipe",
                fdc_id: undefined,
              }}
              onChange={(a) => {
                setRecipe(
                  updateIngredientInfo(
                    recipe,
                    x,
                    y,
                    { id: a.value, name: a.label, fdc_id: a.fdc_id },
                    a.kind
                  )
                );
              }}
            />
          ) : (
            <div className="text-gray-600">
              {ingredient.kind === "recipe" && ingredient.recipe ? (
                <RecipeLink recipe={ingredient.recipe} />
              ) : (
                ingredient.ingredient?.name
              )}
            </div>
          )}
          <TableInput
            data-cy="amount-input"
            width={"full"}
            edit={edit}
            softEdit
            highlight={isOverride(tweaks, x, y, subIndex, "amount")}
            value={adjustIngredientValue(
              tweaks,
              x,
              y,
              subIndex,
              ingredient.amount || 0,
              "amount"
            )}
            onChange={(e) =>
              updateIngredient({
                sectionID: x,
                ingredientID: y,
                subIndex,
                value: e,
                attr: "amount",
              })
            }
          />
          <TableInput
            data-cy="unit-input"
            width={"full"}
            edit={edit}
            value={ingredient.unit || ""}
            onChange={(e) =>
              updateIngredient({
                sectionID: x,
                ingredientID: y,
                subIndex,
                value: e,
                attr: "unit",
              })
            }
          />
          <div className="flex space-x-1">
            <TableInput
              data-cy="adjective-input"
              width={"full"}
              edit={edit}
              value={ingredient.adjective || ""}
              onChange={(e) =>
                updateIngredient({
                  sectionID: x,
                  ingredientID: y,
                  subIndex,
                  value: e,
                  attr: "adjective",
                })
              }
            />
            {!edit && ingredient.optional && (
              <span className="italic">(optional)</span>
            )}
          </div>
          <div>
            {/* <div>{getCal(ingredient, hints, tweaks.multiplier)} kcal</div> */}
            <div>
              {instance &&
                getCal2(
                  instance,
                  ingredient,
                  ing_hints,
                  tweaks.multiplier,
                  UnitConversionRequestTargetEnum.CALORIES
                )}
              kcal
            </div>
            {/* <div>
              $
              {instance &&
                getCal2(
                  instance,
                  ingredient,
                  ing_hints,
                  UnitConversionRequestTargetEnum.MONEY
                )}
            </div> */}
            {!isSub && iActions(x, y, "ingredients")}
            {!isSub && edit && (
              <label className="flex items-center ml-1">
                <input
                  type="checkbox"
                  className="form-checkbox"
                  checked={ingredient.optional}
                  onClick={() =>
                    updateIngredient({
                      sectionID: x,
                      ingredientID: y,
                      subIndex,
                      value: ingredient.optional ? "false" : "true",
                      attr: "optional",
                    })
                  }
                />
                <span className="ml-1">Optional</span>
              </label>
            )}
          </div>
        </div>
        {!!ingredient.original && (
          <div className="italic text-xs pb-2">
            original: {ingredient.original} (
            {instance && instance.parse(ingredient.original)})
          </div>
        )}
      </div>
    );
  };
  const renderRow = (section: RecipeSection, x: number) => (
    <TableRow key={x}>
      <TableCell>
        <Pill>{String.fromCharCode(65 + x)}</Pill>
      </TableCell>
      <TableCell>
        <TableInput
          width={40}
          data-cy="time-input"
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
        {section.ingredients.map((ingredient, y) => (
          <div key={y}>
            {renderIngredientItem(ingredient, x, y)}{" "}
            {(ingredient.substitutes || []).map((sub, z) =>
              renderIngredientItem(sub, x, y, z)
            )}
          </div>
        ))}
        {edit && (
          <ButtonGroup
            compact
            buttons={[
              {
                onClick: () => {
                  setRecipe(addIngredient(recipe, x));
                },
                text: "add ingredient",
                IconLeft: PlusCircle,
              },
            ]}
          />
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
          <ButtonGroup
            compact
            buttons={[
              {
                onClick: () => {
                  setRecipe(addInstruction(recipe, x));
                },
                text: "add instruction",
                IconLeft: PlusCircle,
              },
            ]}
          />
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
          <ButtonGroup
            compact
            buttons={[
              {
                onClick: () => {
                  setRecipe(addSection(recipe));
                },
                text: "add section",
                IconLeft: PlusCircle,
              },
            ]}
          />
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
