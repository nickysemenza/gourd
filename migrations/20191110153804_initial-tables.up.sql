CREATE TABLE IF NOT EXISTS "ingredients" (
  "id" TEXT NOT NULL UNIQUE,
  "name" TEXT NOT NULL UNIQUE,
  PRIMARY KEY ("id")
);
CREATE TABLE IF NOT EXISTS "recipes" (
  "id" TEXT NOT NULL UNIQUE,
  PRIMARY KEY ("id")
);
CREATE TABLE IF NOT EXISTS "recipe_details" (
  "id" TEXT NOT NULL UNIQUE,
  "recipe" TEXT references recipes(id) NOT NULL,
  "name" TEXT NOT NULL,
  "total_minutes" INTEGER,
  "equipment" TEXT,
  "source" JSONB,
  "servings" INTEGER,
  "quantity" INTEGER,
  "unit" TEXT,
  "version" INTEGER NOT NULL,
  PRIMARY KEY ("id"),
  unique("recipe", "version")
);
CREATE TABLE IF NOT EXISTS "recipe_sections" (
  "id" TEXT NOT NULL UNIQUE,
  "recipe_detail" TEXT references recipe_details(id) NOT NULL,
  "sort" INTEGER,
  "minutes" INTEGER,
  PRIMARY KEY ("id")
);
CREATE TABLE IF NOT EXISTS "recipe_section_instructions" (
  "id" TEXT NOT NULL UNIQUE,
  "section" TEXT references recipe_sections(id) NOT NULL,
  "sort" INTEGER,
  "instruction" TEXT,
  PRIMARY KEY ("id")
);
CREATE TABLE IF NOT EXISTS "recipe_section_ingredients" (
  "id" TEXT NOT NULL UNIQUE,
  "section" TEXT references recipe_sections(id) NOT NULL,
  "sort" INTEGER,
  --   ingredient can be an `ingredient` or a `recipe`
  "ingredient" TEXT references ingredients(id),
  "recipe" TEXT references recipes(id),
  "grams" numeric(10, 2),
  "amount" numeric(10, 2),
  "unit" TEXT,
  "adjective" TEXT,
  "optional" boolean default false,
  PRIMARY KEY ("id"),
  constraint check_ingredient check (
    ingredient is not null
    or recipe is not null
  )
);