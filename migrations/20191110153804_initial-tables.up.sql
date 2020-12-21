CREATE TABLE IF NOT EXISTS "ingredients" (
  "uuid" TEXT NOT NULL UNIQUE,
  "name" TEXT NOT NULL UNIQUE,
  PRIMARY KEY ("uuid")
);
CREATE TABLE IF NOT EXISTS "recipes" (
  "uuid" TEXT NOT NULL UNIQUE,
  PRIMARY KEY ("uuid")
);
CREATE TABLE IF NOT EXISTS "recipe_details" (
  "uuid" TEXT NOT NULL UNIQUE,
  "recipe" TEXT references recipes(uuid) NOT NULL,
  "name" TEXT NOT NULL,
  "total_minutes" INTEGER,
  "equipment" TEXT,
  "source" JSONB,
  "servings" INTEGER,
  "quantity" INTEGER,
  "unit" TEXT,
  "version" INTEGER NOT NULL,
  PRIMARY KEY ("uuid"),
  unique("recipe", "version")
);
CREATE TABLE IF NOT EXISTS "recipe_sections" (
  "uuid" TEXT NOT NULL UNIQUE,
  "recipe_detail" TEXT references recipe_details(uuid) NOT NULL,
  "sort" INTEGER,
  "minutes" INTEGER,
  PRIMARY KEY ("uuid")
);
CREATE TABLE IF NOT EXISTS "recipe_section_instructions" (
  "uuid" TEXT NOT NULL UNIQUE,
  "section" TEXT references recipe_sections(uuid) NOT NULL,
  "sort" INTEGER,
  "instruction" TEXT,
  PRIMARY KEY ("uuid")
);
CREATE TABLE IF NOT EXISTS "recipe_section_ingredients" (
  "uuid" TEXT NOT NULL UNIQUE,
  "section" TEXT references recipe_sections(uuid) NOT NULL,
  "sort" INTEGER,
  --   ingredient can be an `ingredient` or a `recipe`
  "ingredient" TEXT references ingredients(uuid),
  "recipe" TEXT references recipes(uuid),
  "grams" numeric(10, 2),
  "amount" numeric(10, 2),
  "unit" TEXT,
  "adjective" TEXT,
  "optional" boolean default false,
  PRIMARY KEY ("uuid"),
  constraint check_ingredient check (
    ingredient is not null
    or recipe is not null
  )
);