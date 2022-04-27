#!/bin/bash

# data from https://fdc.nal.usda.gov/download-datasets.html
# e.g. https://fdc.nal.usda.gov/fdc-datasets/FoodData_Central_csv_2020-04-29.zip
# usage: ./import.sh ./FoodData_Central_csv_2020-04-29/
# takes: 6-8 min
set -euf -o pipefail
shopt -s expand_aliases

tables=(
  food_category
	food
	food_attribute_type
	# acquisition_sample
	# agricultural_acquisition
	branded_food
	food_attribute
	food_nutrient_conversion_factor
	# food_calorie_conversion_factor
	food_component
	nutrient
	food_nutrient_source
	food_nutrient_derivation
	# food_nutrient
	food_nutrient_raw
	measure_unit
	food_portion
	# food_protein_conversion_factor
	# foundation_food
	# input_food
	# lab_method
	# lab_method_code
	# lab_method_nutrient
	# market_acquisition
	# nutrient_incoming_name
	# retention_factor
	# sample_food
	# sr_legacy_food
	# sub_sample_food
	# sub_sample_result
	# wweia_food_category
	# survey_fndds_food
  )

start_time="$(date -u +%s)"

# https://stackoverflow.com/questions/39296472/how-to-check-if-an-environment-variable-exists-and-get-its-value
# https://stackoverflow.com/questions/7832080/test-if-a-variable-is-set-in-bash-when-using-set-o-nounset
if [[ -z "${DB_DSN:-}" ]]; then
  MY_SCRIPT_VARIABLE="postgresql://gourd:gourd@localhost:5556/usda"
else
  MY_SCRIPT_VARIABLE="${DB_DSN}"
fi

if [[ -n "${CI:-}" ]]; then
  echo "sleeping for 30"
  sleep 30
fi

cp "$1"food_nutrient.csv "$1"food_nutrient_raw.csv

alias p='psql "${MY_SCRIPT_VARIABLE}"'

p -c "select count(*) from usda_food";
n=${#tables[*]}
for (( i = n-1; i >= 0; i-- ))
do
    p -c "truncate table usda_${tables[i]} cascade;"
done

for f in "${tables[@]}"; do
	# f=${f%"$_raw"}
    echo "$f"
    headers=$(head -n1 "$1""$f".csv | tr -d '"')
    tmp="$f:tmp"
    sed 's/""/NULL/g' "$1""$f".csv > "$tmp".csv
    p -c "\copy usda_$f($headers) from '$tmp.csv' (format csv, null \"NULL\", DELIMITER ',', HEADER);"
    rm "$tmp".csv
done

# copy from temp table to final
p -c "insert into usda_food_nutrient (select * from usda_food_nutrient_raw where fdc_id in (select fdc_id from usda_food) and nutrient_id in (select id from usda_nutrient));"

end_time="$(date -u +%s)"
elapsed="$(($end_time-$start_time))"
echo "Total of $elapsed seconds elapsed for USDA import"
