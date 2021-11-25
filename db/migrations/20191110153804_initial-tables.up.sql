CREATE TABLE IF NOT EXISTS "ingredients" (
  "id" TEXT NOT NULL UNIQUE,
  "name" TEXT NOT NULL UNIQUE,
  PRIMARY KEY ("id")
);
CREATE TABLE IF NOT EXISTS "recipes" (
  "id" TEXT NOT NULL UNIQUE,
  "created_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
CREATE TABLE IF NOT EXISTS "recipe_details" (
  "id" TEXT NOT NULL UNIQUE,
  "recipe_id" TEXT references recipes(id) NOT NULL,
  "name" TEXT NOT NULL,
  "equipment" TEXT,
  "source" JSONB,
  "servings" INTEGER,
  "quantity" INTEGER,
  "unit" TEXT,
  "version" INTEGER NOT NULL,
  "is_latest_version" BOOLEAN DEFAULT FALSE,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "updated_at" timestamp NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  unique("recipe_id", "version")
);
-- https://stackoverflow.com/a/11014977
-- create unique index one_latest_revision_of_recipe on recipe_details (recipe_id)
-- where is_latest_version;
-- breaks https://github.com/volatiletech/sqlboiler/issues/698
CREATE TABLE IF NOT EXISTS "recipe_sections" (
  "id" TEXT NOT NULL UNIQUE,
  "recipe_detail_id" TEXT references recipe_details(id) NOT NULL,
  "sort" INTEGER,
  "duration_timerange" JSONB,
  PRIMARY KEY ("id")
);
CREATE TABLE IF NOT EXISTS "recipe_section_instructions" (
  "id" TEXT NOT NULL UNIQUE,
  "section_id" TEXT references recipe_sections(id) NOT NULL,
  "sort" INTEGER,
  "instruction" TEXT,
  PRIMARY KEY ("id")
);
CREATE TABLE IF NOT EXISTS "recipe_section_ingredients" (
  "id" TEXT NOT NULL UNIQUE,
  "section_id" TEXT references recipe_sections(id) NOT NULL,
  "sort" INTEGER,
  --   ingredient can be an `ingredient` or a `recipe`
  "ingredient_id" TEXT references ingredients(id),
  "recipe_id" TEXT references recipes(id),
  "amounts" JSON NOT NULL,
  "adjective" TEXT,
  "original" TEXT,
  "optional" boolean default false,
  "sub_for_ingredient_id" TEXT references recipe_section_ingredients(id),
  PRIMARY KEY ("id"),
  constraint check_ingredient check (
    ingredient_id is not null
    or recipe_id is not null
  )
);