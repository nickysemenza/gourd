CREATE TABLE IF NOT EXISTS "ingredients" (
  "uuid" TEXT NOT NULL,
  "name" TEXT NOT NULL UNIQUE,
  PRIMARY KEY ("uuid")
);
CREATE TABLE IF NOT EXISTS "recipes" (
  "uuid" TEXT NOT NULL,
  "name" TEXT NOT NULL UNIQUE,
  "total_minutes" serial,
  "equipment" text,
  "source" text,
  "servings" serial,
  "quantity" serial,
  "unit" TEXT,
  PRIMARY KEY ("uuid")
);
CREATE TABLE IF NOT EXISTS "recipe_sections" (
  "uuid" TEXT NOT NULL,
  "recipe" text references recipes(uuid) NOT NULL,
  "sort" serial,
  "minutes" serial,
  PRIMARY KEY ("uuid")
);
CREATE TABLE IF NOT EXISTS "recipe_section_instructions" (
  "uuid" TEXT NOT NULL,
  "section" text references recipe_sections(uuid) NOT NULL,
  "sort" serial,
  "instruction" text,
  PRIMARY KEY ("uuid")
);
CREATE TABLE IF NOT EXISTS "recipe_section_ingredients" (
  "uuid" TEXT NOT NULL,
  "section" text references recipe_sections(uuid) NOT NULL,
  "sort" serial,
  --   ingredient can be an `ingredient` or a `recipe`
  "ingredient" text references ingredients(uuid),
  "recipe" text references recipes(uuid),
  "grams" numeric(10, 2),
  "amount" numeric(10, 2),
  "unit" text,
  "adjective" text,
  "optional" boolean default false,
  PRIMARY KEY ("uuid"),
  constraint check_ingredient check (
    ingredient is not null
    or recipe is not null
  )
);