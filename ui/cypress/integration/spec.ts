// enables intelligent code completion for Cypress commands
// https://on.cypress.io/intelligent-code-completion
/// <reference types="cypress" />

context("Basic Create, List, Edit test", () => {
  beforeEach(() => {
    cy.visit("/");
  });

  it("creates a new recipe", function () {
    const newName = "cy-" + new Date().getTime();
    cy.contains("Create").click();
    cy.get("input[data-cy=name-input]").first().type(newName);
    cy.contains("Create Recipe").click();
    // cy.contains("Recipes").click();
  });
  it("updates the recipe", function () {
    const newIngredient = "ingredient-" + new Date().getTime();
    cy.contains("Recipes").click();
    cy.get("[data-cy=recipe-table]").should("be.visible");
    cy.contains("details").click();
    cy.url().should("include", "/recipe/");

    cy.contains("edit").click();
    cy.contains("add section").click();
    cy.contains("add ingredient").click();
    cy.get("input[data-cy=grams-input]")
      .first()
      .type("{selectall}4");
    cy.get("input[data-cy=name-input]").first().type(newIngredient);
    cy.contains("add instruction").click();
    cy.get("input[data-cy=instruction-input]")
      .first()
      .type("mix");
    cy.contains("save").click();
    cy.contains("Ingredients").click();
    cy.contains(newIngredient);
  });
  // more examples
  //
  // https://github.com/cypress-io/cypress-example-todomvc
  // https://github.com/cypress-io/cypress-example-kitchensink
  // https://on.cypress.io/writing-your-first-test
});
