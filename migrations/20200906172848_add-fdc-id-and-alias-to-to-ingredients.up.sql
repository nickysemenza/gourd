ALTER TABLE "ingredients" ADD COLUMN "fdc_id" INT4;
ALTER TABLE "ingredients" ADD COLUMN "same_as" TEXT references ingredients(uuid) DEFAULT NULL;