// enables intelligent code completion for Cypress commands
// https://on.cypress.io/intelligent-code-completion
/// <reference types="cypress" />

import { navItems } from "../../src/components/navItems";
import { COOKIE_NAME } from "../../src/config";

context("smoke tests", () => {
  beforeEach(() => {
    cy.visit("/");
  });
  it("navbar smoke test", () => {
    cy.setCookie(
      COOKIE_NAME,
      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2luZm8iOnsiZW1haWwiOiIxNG5pY2hvbGFzc2VAZ21haWwuY29tIiwiZmFtaWx5X25hbWUiOiJTZW1lbnphIiwiZ2l2ZW5fbmFtZSI6Ik5pY2t5IiwiaWQiOiIxMDIxOTE5NDU1ODU0MjYwNjQ0OTciLCJsb2NhbGUiOiJlbiIsIm5hbWUiOiJOaWNreSBTZW1lbnphIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hLS9BT2gxNEdpZ0JQOHotREViMWFfMF9XNUM5eTZKTFd5aG10OXNCTXVxRUc1cTBjVT1zOTYtYyIsInZlcmlmaWVkX2VtYWlsIjp0cnVlfSwiZXhwIjo5NjEwMzI4OTM3LCJpYXQiOjE2MDk0NjQ5Mzd9.c68hryTGinXyBh0HHeZIOas79F7hDAvXZJ37rAcf1og"
    );

    navItems
      .filter((i) => i.title !== "graph")
      .forEach((item) => {
        cy.contains(item.title).click();
        cy.url().should("include", item.path);
      });
  });
  // more examples
  //
  // https://github.com/cypress-io/cypress-example-todomvc
  // https://github.com/cypress-io/cypress-example-kitchensink
  // https://on.cypress.io/writing-your-first-test
});
