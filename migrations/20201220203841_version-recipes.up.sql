ALTER TABLE "recipes"
ADD COLUMN "version" INT DEFAULT 1;
ALTER TABLE "recipes"
ADD COLUMN "is_latest_version" boolean;
ALTER TABLE recipes DROP CONSTRAINT recipes_name_key;