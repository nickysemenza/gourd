CREATE TABLE "ingredient_units" (
    "id" SERIAL PRIMARY KEY,
    "ingredient_id" TEXT references ingredients(id) NOT NULL,
    "unit_a" text NOT NULL,
    "amount_a" numeric(10, 2) NOT NULL,
    "unit_b" text NOT NULL,
    "amount_b" numeric(10, 2) NOT NULL,
    "source" text NOT NULL
);
CREATE UNIQUE INDEX ingredient_unit_attrs_unique ON ingredient_units (
    ingredient_id,
    unit_a,
    amount_a,
    unit_b,
    amount_b
);
ALTER TABLE ingredient_units
ADD CONSTRAINT ingredient_unit_attrs_unique UNIQUE USING INDEX ingredient_unit_attrs_unique;