CREATE TABLE IF NOT EXISTS "images" (
	"id" text NOT NULL,
	"blur_hash" text NOT NULL,
	"source" text NOT NULL,
	"taken_at" timestamp,
	PRIMARY KEY ("id")
);
ALTER TABLE "gphotos_photos"
ADD COLUMN "image_id" text references images(id) NOT NULL;
CREATE TABLE IF NOT EXISTS "notion_recipe" (
	"page_id" text NOT NULL,
	"page_title" text NOT NULL,
	"meta" json,
	"last_seen" timestamp NOT NULL DEFAULT now(),
	"recipe_id" TEXT references recipes(id),
	"ate_at" timestamp,
	"scale" numeric(10, 2),
	"deleted_at" timestamp,
	PRIMARY KEY ("page_id")
);
CREATE TABLE IF NOT EXISTS "notion_image" (
	"block_id" text NOT NULL,
	"page_id" text references notion_recipe(page_id) NOT NULL,
	"last_seen" timestamp NOT NULL DEFAULT now(),
	"image_id" text references images(id) NOT NULL,
	primary key (block_id, page_id),
	unique (block_id, page_id)
);
CREATE TABLE IF NOT EXISTS "notion_meal" (
	"meal_id" TEXT references meals(id) NOT NULL,
	"notion_recipe" TEXT references notion_recipe(page_id) NOT NULL,
	primary key (meal_id, notion_recipe)
);