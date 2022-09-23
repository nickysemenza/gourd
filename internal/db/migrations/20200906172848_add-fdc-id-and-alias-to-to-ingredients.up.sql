ALTER TABLE "ingredients"
ADD COLUMN "fdc_id" INT4;
ALTER TABLE "ingredients"
ADD COLUMN "parent_ingredient_id" TEXT references ingredients(id) DEFAULT NULL;