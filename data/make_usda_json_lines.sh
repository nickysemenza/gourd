#!/bin/bash

cat FoodData_Central_survey_food_json_2022-10-28.json | jq -c '.SurveyFoods[]' >surveyfoods.ndjson
cat FoodData_Central_sr_legacy_food_json_2021-10-28.json | jq -c '.SRLegacyFoods[]' >legacyfoods.ndjson
cat FoodData_Central_foundation_food_json_2022-04-28.json | jq -c '.FoundationFoods[]' >foundationfoods.ndjson
cat FoodData_Central_branded_food_json_2022-04-28.json | jq -c '.BrandedFoods[]' >brandedfoods.ndjson
