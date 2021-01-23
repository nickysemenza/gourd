./bin/gourd scrape http://cooking.nytimes.com/recipes/1018180-perfect-buttermilk-pancakes
./bin/gourd scrape http://www.seriouseats.com/recipes/2013/12/the-food-lab-best-chocolate-chip-cookie-recipe.html
./bin/gourd scrape http://cooking.nytimes.com/recipes/1017060-doughnuts
./bin/gourd scrape http://cooking.nytimes.com/recipes/1017456-three-cup-chicken
./bin/gourd scrape http://www.seriouseats.com/recipes/2011/08/grilled-naan-recipe.html
./bin/gourd import testdata/cookies_2.yaml
./bin/gourd import testdata/cookies_1.yaml

./usda/import.sh ~/Downloads/FoodData_Central_csv_2020-04-29/
./bin/gourd load-ingredients testdata/ingredient_fdc_mapping.yaml