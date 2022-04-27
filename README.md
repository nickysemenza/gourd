# gourd

**Go**(lang) **U**niversal **R**ecipe **D**atabase

> etymology: Before modern technology, people would use hollowed out gourds as food storage containers.

---

Gourd is a recipe database that can be used for meal planning and generating ingredient lists.

## this project is a WIP

[![codecov](https://codecov.io/gh/nickysemenza/gourd/branch/master/graph/badge.svg)](https://codecov.io/gh/nickysemenza/gourd)
[![Go Report Card](https://goreportcard.com/badge/github.com/nickysemenza/gourd)](https://goreportcard.com/report/github.com/nickysemenza/gourd)
[![GoDoc](https://godoc.org/github.com/nickysemenza/gourd?status.svg)](https://pkg.go.dev/github.com/nickysemenza/gourd)

## features

This is comprised mulitple components:

1. **api** OpenAPI 3.0 REST API, defined in `api/openapi.yaml`
1. **ui** in React + Typescript, using generated openapi clients
1. **scraper** for saving recipes from websites (nytimes, seriouseats) for later analysis.
   - This currently works on websites using `json+ld` with [the appropriate schema](https://schema.org/Recipe)
1. **parser** for extracting structured information from freetext ingredient line items.
   - e.g. `1 1/2 cups flour (180g), sifted` will be parsed into `{ingredient: flour, amount: 1.5, unit: cup, modifier: sifted`.
   - This leverages [nickysemenza/ingredient-parser](https://github.com/nickysemenza/ingredient-parser) and is exposed to the UI via WebAssembly.
1. **usda** is used for mapping ingredients to their USDA database equivalent, which has very detailed nutrition information
   - e.g. [plain strawberries](https://fdc.nal.usda.gov/fdc-app.html#/food-details/747448/nutrients) or [C&H brown sugar](https://fdc.nal.usda.gov/fdc-app.html#/food-details/392083/nutrients)
   - This dataset also contains the imperial to metric mappings (e.g. the data from the back of the flour bag that says `1/4 cup = 30 grams`)
1. **cli** for interacting with the api instead using the UI
   - this is used for importing/exporting recipes, as well as loading metadata
