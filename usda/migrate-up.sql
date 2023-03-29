-- schema adapted from https://github.com/hogand/USDA-FoodData-SQLite3/blob/master/import_food_data.py#L89-L404
-- todo: 
-- ERROR:  insert or update on table "usda_food" violates foreign key constraint "usda_food_food_category_id_fkey"
-- DETAIL:  Key (food_category_id)=(9602) is not present in table "usda_food_category".
CREATE TABLE IF NOT EXISTS usda_food_category (
  "id" INT NOT NULL PRIMARY KEY,
  "code" TEXT,
  "description" TEXT
);
CREATE TABLE IF NOT EXISTS usda_food (
  fdc_id INT NOT NULL PRIMARY KEY,
  data_type TEXT,
  description TEXT,
  food_category_id INT,
  -- REFERENCES usda_food_category(id),
  publication_date TEXT
);
CREATE INDEX IF NOT EXISTS idx_food_data_type ON usda_food (data_type);
CREATE INDEX IF NOT EXISTS idx_food_food_category_id ON usda_food (food_category_id);
CREATE TABLE IF NOT EXISTS usda_food_attribute_type (
  "id" INT NOT NULL PRIMARY KEY,
  "name" TEXT,
  "description" TEXT
);
CREATE TABLE IF NOT EXISTS usda_acquisition_sample (
  "fdc_id_of_sample_food" INT REFERENCES usda_food(fdc_id),
  "fdc_id_of_acquisition_food" INT REFERENCES usda_food(fdc_id),
  PRIMARY KEY(
    fdc_id_of_sample_food,
    fdc_id_of_acquisition_food
  )
);
CREATE TABLE IF NOT EXISTS usda_agricultural_acquisition (
  "fdc_id" INT NOT NULL PRIMARY KEY,
  "acquisition_date" TEXT,
  "market_class" TEXT,
  "treatment" TEXT,
  "state" TEXT
);
CREATE TABLE IF NOT EXISTS usda_branded_food (
  "fdc_id" INT NOT NULL PRIMARY KEY REFERENCES usda_food(fdc_id),
  "brand_owner" TEXT,
  -- XXX Inconsistent names
  "brand_name" TEXT,
  "subbrand_name" TEXT,
  "not_a_significant_source_of" TEXT,
  "gtin_upc" TEXT,
  "ingredients" TEXT,
  "serving_size" double precision,
  "serving_size_unit" TEXT CHECK(serving_size_unit IN ('g', 'ml')),
  "household_serving_fulltext" TEXT,
  "branded_food_category" TEXT,
  "data_source" TEXT,
  "package_weight" TEXT,
  "modified_date" TEXT,
  "available_date" TEXT,
  "discontinued_date" TEXT,
  "market_country" TEXT
);
CREATE INDEX IF NOT EXISTS idx_branded_food_gtin_upc ON usda_branded_food (gtin_upc);
CREATE INDEX IF NOT EXISTS idx_branded_food_branded_food_category ON usda_branded_food (branded_food_category);
CREATE TABLE IF NOT EXISTS usda_food_attribute (
  "id" INT NOT NULL PRIMARY KEY,
  "fdc_id" INT REFERENCES usda_food(fdc_id),
  "seq_num" INT,
  "food_attribute_type_id" INT REFERENCES usda_food_attribute_type(id),
  "name" TEXT,
  "value" TEXT
);
CREATE INDEX IF NOT EXISTS idx_food_attribute_fdc_id ON usda_food_attribute (fdc_id);
CREATE INDEX IF NOT EXISTS idx_food_attribute_food_attribute_type_id ON usda_food_attribute (food_attribute_type_id);
CREATE TABLE IF NOT EXISTS usda_food_nutrient_conversion_factor (
  "id" INT NOT NULL PRIMARY KEY,
  "fdc_id" INT REFERENCES usda_food(fdc_id)
);
CREATE INDEX IF NOT EXISTS idx_food_nutrient_conversion_factor_fdc_id ON usda_food_nutrient_conversion_factor (fdc_id);
CREATE TABLE IF NOT EXISTS usda_food_calorie_conversion_factor (
  "food_nutrient_conversion_factor_id" INT NOT NULL PRIMARY KEY,
  -- REFERENCES -0 usda_food_nutrient_conversion_factor(id),
  "protein_value" double precision,
  "fat_value" double precision,
  "carbohydrate_value" double precision
);
CREATE TABLE IF NOT EXISTS usda_food_component (
  "id" INT NOT NULL PRIMARY KEY,
  "fdc_id" INT REFERENCES usda_food(fdc_id),
  "name" TEXT,
  "pct_weight" double precision,
  "is_refuse" TEXT CHECK(is_refuse IN ('Y', 'N')),
  "gram_weight" double precision,
  "data_points" BIGINT,
  "min_year_acquired" TEXT
);
-- XXX Field Descriptions describes "food_fat_conversion_factor" but there is no table for it.
-- XXX File is missing?
CREATE TABLE IF NOT EXISTS usda_nutrient (
  "id" INT NOT NULL PRIMARY KEY,
  "name" TEXT,
  "unit_name" TEXT,
  "nutrient_nbr" TEXT,
  "rank" TEXT -- XXX Not documented
);
CREATE TABLE IF NOT EXISTS usda_food_nutrient_source (
  "id" INT NOT NULL PRIMARY KEY,
  "code" INT UNIQUE,
  -- Code for source (4=calculated).  XXX FK to ?
  "description" TEXT
);
CREATE TABLE IF NOT EXISTS usda_food_nutrient_derivation (
  "id" INT NOT NULL PRIMARY KEY,
  "code" TEXT,
  "description" TEXT,
  "source_id" INT REFERENCES usda_food_nutrient_source(id)
);
CREATE INDEX IF NOT EXISTS idx_food_nutrient_derivation_source_id ON usda_food_nutrient_derivation (source_id);
CREATE TABLE IF NOT EXISTS usda_food_nutrient_raw (
  "id" INT NOT NULL PRIMARY KEY,
  "fdc_id" INT,
  -- todo fk
  "nutrient_id" INT,
  -- todo fk
  "amount" double precision,
  "data_points" INT,
  "derivation_id" INT,
  -- XXX Missing standard_error from Field Descriptions
  "min" double precision,
  "max" double precision,
  "median" double precision,
  "loq" double precision,
  "footnote" TEXT,
  "min_year_acquired" TEXT
);
CREATE TABLE IF NOT EXISTS usda_food_nutrient (
  "id" INT NOT NULL PRIMARY KEY,
  "fdc_id" INT REFERENCES usda_food(fdc_id),
  "nutrient_id" INT REFERENCES usda_nutrient(id),
  "amount" double precision,
  "data_points" INT,
  "derivation_id" INT REFERENCES usda_food_nutrient_derivation(id),
  -- XXX Missing standard_error from Field Descriptions
  "min" double precision,
  "max" double precision,
  "median" double precision,
  "loq" double precision,
  "footnote" TEXT,
  "min_year_acquired" TEXT
);
CREATE INDEX IF NOT EXISTS idx_food_nutrient_nutrient_id ON usda_food_nutrient (nutrient_id);
CREATE INDEX IF NOT EXISTS idx_food_nutrient_derivation_id ON usda_food_nutrient (derivation_id);
CREATE INDEX IF NOT EXISTS idx_food_nutrient_fcd_id ON usda_food_nutrient (fdc_id);
-- new
CREATE TABLE IF NOT EXISTS usda_measure_unit (
  "id" INT NOT NULL PRIMARY KEY,
  "name" TEXT UNIQUE
);
CREATE TABLE IF NOT EXISTS usda_food_portion (
  "id" INT NOT NULL PRIMARY KEY,
  "fdc_id" INT REFERENCES usda_food(fdc_id),
  "seq_num" INT,
  "amount" double precision,
  "measure_unit_id" INT REFERENCES usda_measure_unit(id),
  "portion_description" TEXT,
  "modifier" TEXT,
  "gram_weight" double precision,
  "data_points" INT,
  "footnote" TEXT,
  "min_year_acquired" TEXT
);
CREATE INDEX IF NOT EXISTS idx_food_portion_measure_unit_id ON usda_food_portion (measure_unit_id);
CREATE INDEX IF NOT EXISTS idx_food_portion_fcd_id ON usda_food_portion (fdc_id);
-- new
CREATE TABLE IF NOT EXISTS usda_food_protein_conversion_factor (
  "food_nutrient_conversion_factor_id" INT NOT NULL PRIMARY KEY REFERENCES usda_food_nutrient_conversion_factor(id),
  "value" double precision
);
CREATE TABLE IF NOT EXISTS usda_foundation_food (
  "fdc_id" INT NOT NULL PRIMARY KEY REFERENCES usda_food(fdc_id),
  "ndb_number" INT,
  -- not rlly unique
  "footnote" TEXT
);
CREATE TABLE IF NOT EXISTS usda_input_food (
  "id" INT NOT NULL PRIMARY KEY,
  "fdc_id" INT REFERENCES usda_food(fdc_id),
  "fdc_id_of_input_food" INT REFERENCES usda_food(fdc_id),
  "seq_num" INT,
  "amount" double precision,
  "sr_code" INT,
  -- NDB code of SR food XXX but not a FK
  "sr_description" TEXT,
  "unit" TEXT,
  -- Unit of measure (but inconsistent)
  "portion_code" INT,
  -- Code for portion description XXX FK?
  "portion_description" TEXT,
  "gram_weight" double precision,
  "retention_code" INT,
  "survey_flag" INT
);
CREATE INDEX IF NOT EXISTS idx_input_food_fdc_id ON usda_input_food (fdc_id);
CREATE INDEX IF NOT EXISTS idx_input_food_fdc_id_of_input_food ON usda_input_food (fdc_id_of_input_food);
CREATE TABLE IF NOT EXISTS usda_lab_method (
  "id" INT NOT NULL PRIMARY KEY,
  "description" TEXT,
  "technique" TEXT
);
CREATE TABLE IF NOT EXISTS usda_lab_method_code (
  "id" INT NOT NULL PRIMARY KEY,
  "lab_method_id" INT REFERENCES usda_lab_method(id),
  "code" TEXT
);
CREATE INDEX IF NOT EXISTS idx_lab_method_code_lab_method_id ON usda_lab_method_code (lab_method_id);
CREATE TABLE IF NOT EXISTS usda_lab_method_nutrient (
  "id" INT NOT NULL PRIMARY KEY,
  "lab_method_id" INT REFERENCES usda_lab_method(id),
  "nutrient_id" INT -- XXX this constraint fails: REFERENCES usda_nutrient(id)
);
CREATE INDEX IF NOT EXISTS idx_lab_method_nutrient_lab_method_id ON usda_lab_method_nutrient (lab_method_id);
CREATE TABLE IF NOT EXISTS usda_market_acquisition (
  "fdc_id" INT NOT NULL PRIMARY KEY REFERENCES usda_food(fdc_id),
  "brand_description" TEXT,
  "expiration_date" TEXT,
  "label_weight" double precision,
  "location" TEXT,
  "acquisition_date" TEXT,
  "sales_type" TEXT,
  "sample_lot_nbr" TEXT,
  "sell_by_date" TEXT,
  "store_city" TEXT,
  "store_name" TEXT,
  "store_state" TEXT,
  "upc_code" TEXT
);
-- XXX Missing table nutrient_analysis_details per Field Descriptions
CREATE TABLE IF NOT EXISTS usda_nutrient_incoming_name (
  "id" INT NOT NULL PRIMARY KEY,
  "name" TEXT,
  "nutrient_id" INT REFERENCES usda_nutrient(id)
);
CREATE INDEX IF NOT EXISTS idx_nutrient_incoming_name_nutrient_id ON usda_nutrient_incoming_name (nutrient_id);
CREATE TABLE IF NOT EXISTS usda_retention_factor (
  "id" INT NOT NULL PRIMARY KEY,
  "code" TEXT,
  "food_group_id" INT REFERENCES usda_food_category(id),
  "description" TEXT
);
CREATE INDEX IF NOT EXISTS idx_retention_factor_food_group_id ON usda_retention_factor (food_group_id);
CREATE TABLE IF NOT EXISTS usda_sample_food (
  "fdc_id" INT NOT NULL PRIMARY KEY REFERENCES usda_food(fdc_id)
);
CREATE TABLE IF NOT EXISTS usda_sr_legacy_food (
  "fdc_id" INT NOT NULL PRIMARY KEY REFERENCES usda_food(fdc_id),
  "ndb_number" INT UNIQUE -- XXX doc says starts at 100k but not in practice
);
CREATE TABLE IF NOT EXISTS usda_sub_sample_food (
  "fdc_id" INT NOT NULL PRIMARY KEY REFERENCES usda_food(fdc_id),
  "fdc_id_of_sample_food" INT REFERENCES usda_food(fdc_id)
);
CREATE INDEX IF NOT EXISTS idx_sub_sample_food_fdc_id_of_sample_food ON usda_sub_sample_food (fdc_id_of_sample_food);
CREATE TABLE IF NOT EXISTS usda_sub_sample_result (
  "food_nutrient_id" INT NOT NULL PRIMARY KEY REFERENCES usda_food_nutrient(id),
  "adjusted_amount" double precision,
  "lab_method_id" INT REFERENCES usda_lab_method(id),
  -- XXX cannot use this because of broken refs: REFERENCES usda_lab_method(id),
  "nutrient_name" TEXT
);
CREATE INDEX IF NOT EXISTS idx_sub_sample_result_lab_method_id ON usda_sub_sample_result (lab_method_id);
CREATE TABLE IF NOT EXISTS usda_wweia_food_category (
  "wweia_food_category_code" INT NOT NULL PRIMARY KEY,
  "wweia_food_category_description" TEXT
);
CREATE TABLE IF NOT EXISTS usda_survey_fndds_food (
  "fdc_id" INT NOT NULL PRIMARY KEY REFERENCES usda_food(fdc_id),
  "food_code" INT UNIQUE,
  "wweia_category_code" INT REFERENCES usda_wweia_food_category(wweia_food_category_code),
  "start_date" TEXT,
  "end_date" TEXT
);
CREATE INDEX IF NOT EXISTS idx_survey_fndds_food_wweia_category_code ON usda_survey_fndds_food (wweia_category_code);
-- https://www.alibabacloud.com/blog/postgresql-fuzzy-search-best-practices-single-word-double-word-and-multi-word-fuzzy-search-methods_595635
-- create index idx on usda_food(description collate "C");
-- https://scoutapm.com/blog/how-to-make-text-searches-in-postgresql-faster-with-trigram-similarity
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx2 ON usda_food USING GIN(description gin_trgm_ops);