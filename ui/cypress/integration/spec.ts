// enables intelligent code completion for Cypress commands
// https://on.cypress.io/intelligent-code-completion
/// <reference types="cypress" />

import { navItems } from "../../src/components/navItems";
import { COOKIE_NAME } from "../../src/config";

context("Basic Create, List, Edit test", () => {
  beforeEach(() => {
    cy.visit("/");
  });

  const newName = "cy-" + new Date().getTime();

  it("creates a new recipe", function () {
    // https://stackoverflow.com/a/63844955
    cy.contains("create").click();
    cy.get("div[data-cy=name-input]")
      .find(".react-select__control")
      .first()
      .type(newName);

    cy.get("div[data-cy=name-input]")
      .find(".react-select__option") // find all options
      .first()
      .click(); // click on first options
    cy.wait(1000);
    cy.url().should("include", "/recipe/");
    // cy.contains(`create recipe: ${newName}`).click({ force: true });
  });
  it("updates the recipe", function () {
    const newIngredient = "ingredient-" + new Date().getTime();
    cy.contains("recipes").click();
    cy.contains(newName).click();
    cy.url().should("include", "/recipe/");

    cy.contains("edit").click();
    cy.contains("add section").click();
    cy.contains("add ingredient").click();
    cy.get("input[data-cy=grams-input]").first().type("{selectall}4");
    cy.get("div[data-cy=name-input]")
      .find("input")
      .first()
      .type(`${newIngredient}`);
    cy.wait(500);
    cy.contains(`create ingredient: ${newIngredient}`);

    cy.get("div[data-cy=name-input]")
      .find(".react-select__option") // find all options
      .first()
      .click(); // click on first options

    // cy.get("div[data-cy=name-input]").first().type(`{enter}`);
    cy.wait(500);
    cy.contains("add instruction").click();
    cy.get("textarea[data-cy=instruction-input]").first().type("mix");
    cy.contains("save").click();
    cy.wait(500);
    cy.contains("ingredients").click();
    cy.reload();
    // cy.contains(newIngredient);
  });
});
context("smoke tests", () => {
  it("navbar smoke test", () => {
    cy.setCookie(
      COOKIE_NAME,
      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2luZm8iOnsiZW1haWwiOiIxNG5pY2hvbGFzc2VAZ21haWwuY29tIiwiZmFtaWx5X25hbWUiOiJTZW1lbnphIiwiZ2l2ZW5fbmFtZSI6Ik5pY2t5IiwiaWQiOiIxMDIxOTE5NDU1ODU0MjYwNjQ0OTciLCJsb2NhbGUiOiJlbiIsIm5hbWUiOiJOaWNreSBTZW1lbnphIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hLS9BT2gxNEdpZ0JQOHotREViMWFfMF9XNUM5eTZKTFd5aG10OXNCTXVxRUc1cTBjVT1zOTYtYyIsInZlcmlmaWVkX2VtYWlsIjp0cnVlfSwiZXhwIjo5NjEwMzI4OTM3LCJpYXQiOjE2MDk0NjQ5Mzd9.c68hryTGinXyBh0HHeZIOas79F7hDAvXZJ37rAcf1og"
    );

    navItems
      .filter((i) => i.title !== "graph")
      .forEach((item) => {
        cy.contains(item.title).click();
        cy.wait(500);
        cy.url().should("include", item.path);
      });
  });
  // more examples
  //
  // https://github.com/cypress-io/cypress-example-todomvc
  // https://github.com/cypress-io/cypress-example-kitchensink
  // https://on.cypress.io/writing-your-first-test
});
