import { assert, expect, test } from "vitest";
import { getTotalDuration } from "../util/time";
// import { getTotalDuration } from "util/time";
// Edit an assertion and save to see HMR in action
const w = await import("../wasm/gourd_wasm");

test("Math.sqrt()", () => {
  expect(Math.sqrt(4)).toBe(2);
  expect(Math.sqrt(144)).toBe(12);
  expect(Math.sqrt(2)).toBe(Math.SQRT2);
});

test("JSON", () => {
  const input = {
    foo: "hello",
    bar: "world",
  };

  const output = JSON.stringify(input);

  expect(output).eq('{"foo":"hello","bar":"world"}');
  assert.deepEqual(JSON.parse(output), input, "matches original");
});

test("time addition", () => {
  const res = getTotalDuration(w, [
    {
      id: "a",
      ingredients: [],
      instructions: [],
      duration: { value: 1, unit: "hour" },
    },
    {
      id: "a",
      ingredients: [],
      instructions: [],
      duration: { value: 15, unit: "minute" },
    },
  ]);
  assert.equal(res.value, 1.25);
  assert.equal(res.unit, "Hour");

  assert.equal(getTotalDuration(w, []).value, 0);
});
