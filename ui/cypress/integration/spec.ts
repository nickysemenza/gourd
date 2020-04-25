// enables intelligent code completion for Cypress commands
// https://on.cypress.io/intelligent-code-completion
/// <reference types="cypress" />

context("Example Cypress TodoMVC test", () => {
  beforeEach(() => {
    // usually we recommend setting baseUrl in cypress.json
    // but for simplicity of this example we just use it here
    // https://on.cypress.io/visit
    cy.visit("http://localhost:3000/recipes");
  });

  it("adds 2 todos", function () {
    cy.get("[data-cy=recipe-table]").should("be.visible");
    cy.contains("details").click();
    cy.url().should("include", "/recipe/");
    cy.get("input[data-cy=grams-input]")
      .first()
      .type("{selectall}1");
  });

  // more examples
  //
  // https://github.com/cypress-io/cypress-example-todomvc
  // https://github.com/cypress-io/cypress-example-kitchensink
  // https://on.cypress.io/writing-your-first-test
});
