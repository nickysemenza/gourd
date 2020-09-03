CREATE TABLE IF NOT EXISTS "recipe_notes" (
  "uuid" TEXT NOT NULL UNIQUE,
  "recipe" TEXT references recipes(uuid) NOT NULL,
  "note" TEXT,
  "date" TIMESTAMP,
  PRIMARY KEY ("uuid")
);

CREATE TABLE IF NOT EXISTS "recipe_sources" (
  "uuid" TEXT NOT NULL UNIQUE,
  "recipe" TEXT references recipes(uuid) NOT NULL,
  "name" TEXT,
  "meta" TEXT,
  PRIMARY KEY ("uuid")
);