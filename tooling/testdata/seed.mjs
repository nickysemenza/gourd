#!/usr/bin/env zx
await Promise.all([
  $`./bin/gourd scrape http://cooking.nytimes.com/recipes/1018180-perfect-buttermilk-pancakes`,
  $`./bin/gourd scrape http://www.seriouseats.com/recipes/2013/12/the-food-lab-best-chocolate-chip-cookie-recipe.html`,
  $`./bin/gourd scrape http://cooking.nytimes.com/recipes/1017060-doughnuts`,
  $`./bin/gourd scrape http://cooking.nytimes.com/recipes/1017456-three-cup-chicken`,
  $`./bin/gourd scrape http://www.seriouseats.com/recipes/2011/08/grilled-naan-recipe.html`,
  $`./bin/gourd scrape https://cooking.nytimes.com/recipes/1018177-cornmeal-blueberry-pancakes`,
]);

await Promise.all([
  $`./bin/gourd import tooling/testdata/cookies_2.txt`,
  $`./bin/gourd import tooling/testdata/cookies_2b.txt`,
  $`./bin/gourd import tooling/testdata/cookies_1.yaml`,
  $`./bin/gourd import tooling/testdata/pasta.txt`,
]);
await $`./bin/gourd load-ingredients tooling/ingredient_fdc_mapping.yaml`;
