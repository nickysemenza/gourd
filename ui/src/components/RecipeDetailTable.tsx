import React, { useCallback, useContext } from "react";
import {
  RecipeWrapper,
  RecipeSection,
  SectionIngredient,
} from "../api/openapi-hooks/api";
import {
  formatRichText,
  formatTimeRange,
  getIngredient,
  scaledRound,
} from "../util";
import { DndProvider } from "react-dnd";
import { HTML5Backend } from "react-dnd-html5-backend";
import update from "immutability-helper";
import { ButtonGroup, PillLabel } from "./Button";
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
  getHint,
  getStats,
  totalFlourMass,
  getGramsFromSI,
  isGram,
  isVolume,
  flatIngredients,
  getMultiplierFromRecipe,
} from "./RecipeEditorUtils";
import {
  ArrowDown,
  ArrowRight,
  ArrowUp,
  PlusCircle,
  XSquare,
} from "react-feather";
import { RecipeLink } from "./Misc";
import { EntitySelector } from "./EntitySelector";
import { WasmContext } from "../wasmContext";
import { TableInput } from "./Input";
import IngredientPopover from "./IngredientPopover";

export interface UpdateIngredientProps {
  sectionID: number;
  ingredientID: number;
  subIndex?: number;
  value: string;
  attr: "grams" | "name" | "amount" | "unit" | "adjective" | "optional";
  amountIndex?: number;
}

export interface TableProps {
  recipe: RecipeWrapper;
  updateIngredient: (i: UpdateIngredientProps) => void;
  setRecipe: React.Dispatch<React.SetStateAction<RecipeWrapper | null>>;
  tweaks: RecipeTweaks;
  hints: FoodsById;
  ing_hints: IngDetailsById;
  showOriginalLine: boolean;
  showKcalDollars: boolean;
}
const RecipeDetailTable: React.FC<TableProps> = ({
  recipe,
  updateIngredient,
  setRecipe,
  tweaks,
  hints,
  ing_hints,
  showOriginalLine,
  showKcalDollars,
}) => {
  const { edit } = tweaks;
  const { sections } = recipe.detail;
  const flourMass = totalFlourMass(sections || []);
  const showBP = flourMass > 0;

  const w = useContext(WasmContext);

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
    const bp = Math.round((getGramsFromSI(ingredient) / flourMass) * 100);
    const { edit } = tweaks;
    const isSub = subIndex !== undefined;
    const hint = getHint(ingredient, ing_hints);

    if (!w) return <div>loading...</div>;

    const { kcal, dollars } = getStats(
      w,
      ingredient,
      ing_hints,
      tweaks.multiplier
    );

    let ingIndex = -1;
    let gramIndex = -1;
    ingredient.amounts.forEach((a, i) => {
      if (gramIndex === -1 && isGram(a)) {
        gramIndex = i;
      }
      if (ingIndex === -1 && isVolume(a)) {
        ingIndex = i;
      }
    });
    const placeholderGrams =
      (gramIndex >= 0 &&
        ingredient.amounts[gramIndex].source !== "db" &&
        ingredient.amounts[gramIndex].value) ||
      undefined;
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
                gramIndex === -1 ? 0 : ingredient.amounts[gramIndex].value,
                "grams"
              )}
              pValue={placeholderGrams && tweaks.multiplier * placeholderGrams}
              blur
              highlight={isOverride(tweaks, x, y, subIndex, "grams")}
              onChange={(e) =>
                updateIngredient({
                  sectionID: x,
                  ingredientID: y,
                  subIndex,
                  value: e,
                  attr: "grams",
                  amountIndex: gramIndex,
                })
              }
            />
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
                kind: getIngredient(ingredient).kind,
                fdc_id: undefined,
              }}
              onChange={(a) => {
                setRecipe(
                  updateIngredientInfo(
                    recipe,
                    x,
                    y,
                    { id: a.value, name: a.label, fdc_id: a.fdc_id },
                    a.kind || "ingredient"
                  )
                );
              }}
            />
          ) : (
            <div className="text-gray-600 dark:text-zinc-200">
              {ingredient.kind === "recipe" && ingredient.recipe ? (
                <RecipeLink
                  recipe={ingredient.recipe}
                  multiplier={getMultiplierFromRecipe(
                    ingredient,
                    tweaks.multiplier
                  )}
                />
              ) : (
                <div className="flex justify-between pr-1">
                  <p>{ingredient.ingredient?.ingredient.name}</p>
                  {hint && <IngredientPopover detail={hint} />}
                </div>
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
              ingIndex === -1 ? 0 : ingredient.amounts[ingIndex].value,
              "amount"
            )}
            onChange={(e) =>
              updateIngredient({
                sectionID: x,
                ingredientID: y,
                subIndex,
                value: e,
                attr: "amount",
                amountIndex: ingIndex,
              })
            }
          />
          <TableInput
            data-cy="unit-input"
            width={"full"}
            edit={edit}
            value={ingIndex === -1 ? "" : ingredient.amounts[ingIndex].unit}
            onChange={(e) =>
              updateIngredient({
                sectionID: x,
                ingredientID: y,
                subIndex,
                value: e,
                attr: "unit",
                amountIndex: ingIndex,
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
            {showKcalDollars && (
              <div>
                {kcal ? scaledRound(kcal) : "n/a"}
                kcal ${dollars ? scaledRound(dollars) : "n/a"}
              </div>
            )}

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
        {showOriginalLine && ingredient.original && (
          <div className="italic text-xs inline-flex">
            <div className="text-slate-700 mr-1">{ingredient.original}</div>
            <ArrowRight className="pb-2" width={10} />
            <div className="text-green-800">
              {w && w.parse(ingredient.original)}
            </div>
          </div>
        )}
      </div>
    );
  };
  const allIngredients: string[] = flatIngredients(sections)
    .map((i) =>
      [
        i.ingredient?.ingredient.name || "flour",
        i.ingredient?.children?.map((i) => i.ingredient.name || "flour") || [],
      ].flat()
    )
    .flat()
    .map((i) => i.toLowerCase());
  const renderRow = (section: RecipeSection, x: number) => (
    <TableRow key={x}>
      <TableCell>
        <PillLabel kind="letter" x={x} />
        <TableInput
          width={40}
          data-cy="time-input"
          edit={edit}
          value={formatTimeRange(w, section.duration)}
          blur
          onChange={(e) =>
            w &&
            setRecipe(
              updateTimeRange(
                recipe,
                x,
                w.parse_amount(e).pop() || section.duration
              )
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
        {section.instructions.map((instruction, y) => (
          <div key={y} className="flex font-serif">
            <div className="mr-4 w-4">
              <PillLabel
                kind="number"
                x={getGlobalInstructionNumber(recipe, x, y)}
              />
            </div>
            {edit && (
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
            )}
            <div>{iActions(x, y, "instructions")}</div>
            <div className="py-1">
              <div>
                {w &&
                  formatRichText(
                    w,
                    w.rich(instruction.instruction, allIngredients)
                  )}
              </div>
              <div className={w && "text-gray-500"}>
                {!w || (showOriginalLine && instruction.instruction)}
              </div>
            </div>
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
      <div className="border-gray-400 shadow-xl bg-white dark:bg-gray-700">
        <TableRow header>
          <TableCell></TableCell>
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
            enable={edit && sections.length > 1}
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

const TableCell: React.FC<{
  children?: React.ReactNode;
}> = ({ children }) => (
  <div className="border-solid border border-gray-300 p-1">{children}</div>
);
const TableRow: React.FC<{ header?: boolean; children?: React.ReactNode }> = ({
  children,
  header = false,
}) => (
  <div className={`rec-table-row ${header && "font-semibold"}`}>{children}</div>
);
