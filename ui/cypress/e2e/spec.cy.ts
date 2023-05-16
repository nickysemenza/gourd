// enables intelligent code completion for Cypress commands
// https://on.cypress.io/intelligent-code-completion
/// <reference types="cypress" />

context("Basic Create, List, Edit test", () => {
  beforeEach(() => {
    cy.visit("/");
  });

  const newName = "cy-" + new Date().getTime();

  it("creates a new recipe", function () {
    // https://stackoverflow.com/a/63844955
    cy.intercept("POST", "/api/recipes").as("makeRecipeFromDropdown");

    cy.contains("create").click();
    cy.get("div[data-cy=name-input]")
      .find(".react-select__control")
      .first()
      .type(newName);

    cy.get("div[data-cy=name-input]")
      .find(".react-select__option") // find all options
      .first()
      .click(); // click on first options
    cy.wait("@makeRecipeFromDropdown")
      .its("response.statusCode")
      .should("eq", 201);
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

    cy.intercept("/api/search?name=" + newIngredient).as("ingredientDropdown");
    let foo = cy.get("div[data-cy=name-input]");
    foo
      .find(".react-select__input-container")
      .click()
      // foo.find("input").first().
      .type(`${newIngredient}`);
    cy.wait("@ingredientDropdown").its("response.statusCode").should("eq", 200);
    cy.contains(`create ingredient: ${newIngredient}`);

    cy.intercept("POST", "/api/ingredients").as("makeIngredient");
    cy.get("div[data-cy=name-input]")
      .find(".react-select__option") // find all options
      // .find(".react-select__single-value") // find all options
      .first()
      .click(); // click on first options

    // cy.get("div[data-cy=name-input]").first().type(`{enter}`);
    cy.wait("@makeIngredient").its("response.statusCode").should("eq", 201);
    cy.contains("add instruction").click();
    cy.get("textarea[data-cy=instruction-input]").first().type("mix");
    cy.contains("save").click();
    cy.wait(500);
    cy.contains("ingredients").click();
    cy.reload();
    // cy.contains(newIngredient);
  });
});
