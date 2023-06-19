import React from "react";
import { assert, expect, test } from "vitest";
import renderer from "react-test-renderer";
import { wasm } from "../util/wasmContext";
import { formatRichText } from "../util/rich";
const RichFormatter: React.FC<{
  text: string;
  ingredients: string[];
  w: wasm;
}> = ({ text, ingredients, w }) => formatRichText(w, w.rich(text, ingredients));

function toJson(component: renderer.ReactTestRenderer) {
  const result = component.toJSON();
  expect(result).toBeDefined();
  expect(result).toBeInstanceOf(Array);
  return result as renderer.ReactTestRendererJSON[];
}

const w = await import("../wasm/gourd_wasm");

test("rich text is displayed properly", () => {
  const component = renderer.create(
    <RichFormatter
      text="heat oven to 150°F for 5 minutes and put 1 cup flour in a bowl"
      w={w}
      ingredients={["flour"]}
    />
  );

  const tree = toJson(component);
  expect(tree).toMatchSnapshot();

  const children = tree.map((t) => t.children).filter((t) => t !== undefined);
  console.log(tree);
  assert.deepEqual(children, [["150 °f"], ["5 minutes"], ["1 cup"], ["flour"]]);
});
