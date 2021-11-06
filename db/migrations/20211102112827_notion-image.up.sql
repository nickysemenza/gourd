CREATE TABLE IF NOT EXISTS "notion_recipe" (
	"page_id" text NOT NULL,
	"page_title" text NOT NULL,
	"meta" json,
	"last_seen" timestamp NOT NULL DEFAULT now(),
	"recipe" TEXT references recipes(id),
	-- NOT NULL,
	"ate_at" timestamp,
	-- NOT NULL,
	PRIMARY KEY ("page_id")
);
CREATE TABLE IF NOT EXISTS "notion_image" (
	"block_id" text NOT NULL,
	"page_id" text references notion_recipe(page_id) NOT NULL,
	"blur_hash" text,
	"last_seen" timestamp NOT NULL DEFAULT now(),
	PRIMARY KEY ("block_id"),
	unique (block_id, page_id)
);
CREATE TABLE IF NOT EXISTS "notion_meal" (
	"meal" TEXT references meals(id) NOT NULL,
	"notion_recipe" TEXT references notion_recipe(page_id) NOT NULL,
	unique (meal, notion_recipe)
);
-- 
CREATE TABLE IF NOT EXISTS "images" (
	"id" text NOT NULL,
	"blur_hash" text NOT NULL,
	"source" text NOT NULL,
	unique (id)
);
ALTER TABLE "notion_image"
ADD COLUMN "image" text references images(id) NOT NULL;
ALTER TABLE "notion_image" DROP COLUMN "blur_hash";
ALTER TABLE "gphotos_photos"
ADD COLUMN "image" text references images(id) NOT NULL;