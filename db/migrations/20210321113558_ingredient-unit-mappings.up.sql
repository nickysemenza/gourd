CREATE TABLE "ingredient_units" (
    "id" SERIAL PRIMARY KEY,
    "ingredient" TEXT references ingredients(id) NOT NULL,
    "unit_a" text,
    "amount_a" numeric(10, 2),
    "unit_b" text,
    "amount_b" numeric(10, 2),
    "source" text
);