CREATE TABLE IF NOT EXISTS "meals" (
  "uuid" TEXT NOT NULL UNIQUE,
  "name" TEXT NOT NULL UNIQUE,
  "notion_link" TEXT,
  PRIMARY KEY ("uuid")
);
CREATE TABLE IF NOT EXISTS "meal_recipe" (
  "meal_uuid" TEXT references meals(uuid) NOT NULL,
  "recipe_uuid" TEXT references recipes(uuid) NOT NULL,
  "multiplier" numeric(10, 2) DEFAULT 1.0,
  unique (meal_uuid, recipe_uuid)
);