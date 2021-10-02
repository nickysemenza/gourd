ALTER TABLE "ingredients"
ADD COLUMN "fdc_id" INT4;
ALTER TABLE "ingredients"
ADD COLUMN "parent" TEXT references ingredients(id) DEFAULT NULL;
-- https://www.alibabacloud.com/blog/postgresql-fuzzy-search-best-practices-single-word-double-word-and-multi-word-fuzzy-search-methods_595635
-- create index idx on usda_food(description collate "C");
-- https://scoutapm.com/blog/how-to-make-text-searches-in-postgresql-faster-with-trigram-similarity
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx2 ON usda_food USING GIN(description gin_trgm_ops);