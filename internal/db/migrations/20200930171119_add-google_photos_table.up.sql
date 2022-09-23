CREATE TABLE "gphotos_albums" (
    "id" text NOT NULL,
    "usecase" text NOT NULL,
    PRIMARY KEY ("id")
);
CREATE TABLE "gphotos_photos" (
    "id" text NOT NULL,
    "album_id" text references gphotos_albums(id) NOT NULL,
    "creation_time" timestamp NOT NULL,
    "last_seen" timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);
CREATE TABLE IF NOT EXISTS "meal_gphoto" (
    "meal_id" TEXT references meals(id) NOT NULL,
    "gphotos_id" TEXT references gphotos_photos(id) NOT NULL,
    "highlight_recipe_id" TEXT references recipes(id),
    primary key (meal_id, gphotos_id)
);
ALTER TABLE "meals"
ADD COLUMN "ate_at" timestamp NOT NULL;