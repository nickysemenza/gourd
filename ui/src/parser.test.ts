import { encodeIngredient } from "./parser";

test("encode igredient", async () => {
  const wasm = await import("gourd_rs");
  expect(
    encodeIngredient(
      {
        ingredient: { name: "egg", id: "" },
        id: "",
        kind: "ingredient",
        grams: 27,
        amount: 0,
        unit: "",
        optional: false,
        adjective: "",
      },
      wasm.format_ingredient
    )
  ).toEqual("27g egg");
  expect(
    encodeIngredient(
      {
        ingredient: { name: "apple", id: "" },
        id: "",
        kind: "ingredient",
        grams: 2,
        amount: 0,
        unit: "",
        optional: false,
        adjective: "chopped",
      },
      wasm.format_ingredient
    )
  ).toEqual("2g apple, chopped");
});
