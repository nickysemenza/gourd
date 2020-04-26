// enables intelligent code completion for Cypress commands
// https://on.cypress.io/intelligent-code-completion
/// <reference types="cypress" />

context("Basic Create List test", () => {
  beforeEach(() => {
    // usually we recommend setting baseUrl in cypress.json
    // but for simplicity of this example we just use it here
    // https://on.cypress.io/visit
    cy.visit("/");
  });

  it("creates a new one", function () {
    const newName = "cy-" + new Date().getTime();
    cy.contains("Create").click();
    cy.get("input[data-cy=name-input]").first().type(newName);
    cy.contains("Create Recipe").click();
    // cy.contains("Recipes").click();
  });
  it("updates the recipe", function () {
    const newName = "ingredient-" + new Date().getTime();
    cy.contains("Recipes").click();
    cy.get("[data-cy=recipe-table]").should("be.visible");
    cy.contains("details").click();
    cy.url().should("include", "/recipe/");

    // cy.get("[data-cy=recipe-name]").invoke("text").as("text");

    cy.contains("edit").click();
    cy.contains("add section").click();
    cy.contains("add ingredient").click();
    cy.get("input[data-cy=grams-input]")
      .first()
      .type("{selectall}4");
    cy.get("input[data-cy=name-input]").first().type(newName);
    cy.contains("add instruction").click();
    cy.get("input[data-cy=instruction-input]")
      .first()
      .type("mix");
    cy.contains("save").click();
    cy.contains("Ingredients").click();
    cy.contains(newName);
  });

  // more examples
  //
  // https://github.com/cypress-io/cypress-example-todomvc
  // https://github.com/cypress-io/cypress-example-kitchensink
  // https://on.cypress.io/writing-your-first-test
});
