CREATE TABLE "gphotos_albums" (
    "id" text NOT NULL,
    "usecase" text NOT NULL,
    PRIMARY KEY ("id")
);
CREATE TABLE "gphotos_photos" (
    "id" text NOT NULL,
    "album_id" text references gphotos_albums(id) NOT NULL,
    "creation_time" timestamp NOT NULL,
    PRIMARY KEY ("id")
);