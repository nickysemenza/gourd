#!/bin/sh
# data from https://fdc.nal.usda.gov/download-datasets.html
# e.g. https://fdc.nal.usda.gov/fdc-datasets/FoodData_Central_csv_2020-04-29.zip
# usage: ./import.sh ./FoodData_Central_csv_2020-04-29/
# takes: 6-8 min
set -euf -o pipefail
tables=(food_category food food_attribute_type acquisition_sample agricultural_acquisition branded_food food_attribute food_nutrient_conversion_factor food_calorie_conversion_factor food_component nutrient food_nutrient_source food_nutrient_derivation food_nutrient measure_unit food_portion food_protein_conversion_factor foundation_food input_food lab_method lab_method_code lab_method_nutrient market_acquisition nutrient_incoming_name retention_factor sample_food sr_legacy_food sub_sample_food sub_sample_result wweia_food_category survey_fndds_food)

alias p='psql "postgresql://gourd:gourd@localhost:5555/food"'

p -c "select count(*) from usda_food";
n=${#tables[*]}
for (( i = n-1; i >= 0; i-- ))
do
    p -c "truncate table usda_${tables[i]} cascade;"
done
for f in ${tables[@]}; do
    echo $f
    headers=$(head -n1 $1$f.csv | tr -d '"')
    tmp="$f:tmp"
    sed 's/""/NULL/g' $1$f.csv > $tmp.csv
    p -c "\copy usda_$f($headers) from '$tmp.csv' (format csv, null \"NULL\", DELIMITER ',', HEADER);"
done

rm "*.csv"