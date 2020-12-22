CREATE TABLE IF NOT EXISTS "meals" (
  "id" TEXT NOT NULL UNIQUE,
  "name" TEXT NOT NULL UNIQUE,
  "notion_link" TEXT,
  PRIMARY KEY ("id")
);
CREATE TABLE IF NOT EXISTS "meal_recipe" (
  "meal_id" TEXT references meals(id) NOT NULL,
  "recipe_id" TEXT references recipes(id) NOT NULL,
  "multiplier" numeric(10, 2) DEFAULT 1.0,
  unique (meal_id, recipe_id)
);