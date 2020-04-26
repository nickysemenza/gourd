import React from "react";
import { render } from "@testing-library/react";
import App from "./App";
import { recipeToRecipeInput } from "./util";

test("renders learn react link", () => {
  // render(<App />);
  // const linkElement = getByText(/learn react/i);
  // expect(linkElement).toBeInTheDocument();
});

test("utils", () => {
  const { name } = recipeToRecipeInput({ name: "foo" });
  expect(name).toEqual("fooa");
});
