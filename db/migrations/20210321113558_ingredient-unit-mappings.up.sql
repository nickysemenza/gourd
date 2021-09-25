CREATE TABLE "ingredient_units" (
    "id" SERIAL PRIMARY KEY,
    "ingredient" TEXT references ingredients(id) NOT NULL,
    "unit_a" text,
    "amount_a" numeric(10, 2),
    "unit_b" text,
    "amount_b" numeric(10, 2),
    "source" text
);
CREATE UNIQUE INDEX ingredient_unit_attrs_unique ON ingredient_units (ingredient, unit_a, amount_a, unit_b, amount_b);
ALTER TABLE ingredient_units
ADD CONSTRAINT ingredient_unit_attrs_unique UNIQUE USING INDEX ingredient_unit_attrs_unique;