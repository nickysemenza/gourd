ALTER TABLE "meals" DROP COLUMN "ate_at";
alter table meals drop constraint meals_name_key;
DROP TABLE meal_photo;
DROP TABLE gphotos_photos;
DROP TABLE gphotos_albums;