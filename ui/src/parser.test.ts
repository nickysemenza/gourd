import { parseRecipe, parseLine, encodeIngredient } from "./parser";
import { SectionIngredientKind } from "./generated/graphql";

test("parse with adjective", () => {
  const [grams, ingredient, adjective] = parseLine("30g butter, melted");
  expect(grams).toEqual(30);
  expect(ingredient).toEqual("butter");
  expect(adjective).toEqual("melted");
});

test("parse without adjective", () => {
  const [grams, ingredient, adjective] = parseLine("30g butter");
  expect(grams).toEqual(30);
  expect(ingredient).toEqual("butter");
  expect(adjective).toEqual("");
});

test("encode igredient", () => {
  expect(
    encodeIngredient({
      info: { name: "egg", uuid: "" },
      uuid: "",
      kind: SectionIngredientKind.Ingredient,
      grams: 27,
      amount: 0,
      unit: "",
      optional: false,
      adjective: "",
    })
  ).toEqual("27g egg");
  expect(
    encodeIngredient({
      info: { name: "apple", uuid: "" },
      uuid: "",
      kind: SectionIngredientKind.Ingredient,
      grams: 2,
      amount: 0,
      unit: "",
      optional: false,
      adjective: "chopped",
    })
  ).toEqual("2g apple, chopped");
});
