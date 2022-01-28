#!/bin/bash
date
./testdata/seed.sh
./usda/import.sh ~/Downloads/FoodData_Central_csv_2021-10-28/
LOG_LEVEL=info ./bin/gourd sync
date