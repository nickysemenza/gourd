# gourd

**Go**(lang) **U**niversal **R**ecipe **D**atabase

Before modern technology, people would use hollowed out gourds as food storage containers.

---

## this project is a WIP

[![codecov](https://codecov.io/gh/nickysemenza/gourd/branch/master/graph/badge.svg)](https://codecov.io/gh/nickysemenza/gourd) [![Go Report Card](https://goreportcard.com/badge/github.com/nickysemenza/gourd)](https://goreportcard.com/report/github.com/nickysemenza/gourd) [![CircleCI](https://circleci.com/gh/nickysemenza/gourd.svg?style=svg)](https://circleci.com/gh/nickysemenza/gourd) [![GoDoc](https://godoc.org/github.com/nickysemenza/gourd?status.svg)](https://pkg.go.dev/github.com/nickysemenza/gourd)

## features

This is comprised mulitple components:

1. **parser** for extracting structured information from freetext ingredient line items.
   - e.g. `1 1/2 cups flour (180g), sifted` will be parsed into `ingredient: flour`, `cup:1.5`, `grams:180`, `modifier:sifted`.
   - This allows for recipes to be scaled up and down, and for ingredients to be tied together.
2. **scraper** for saving recipes from websites (nytimes, seriouseats) for later analysis.
   - This currently works on websites using `json+ld` with [the appropriate schema](https://schema.org/Recipe)
3. **usda** is used for mapping ingredients to their USDA database equivalent, which has very detailed nutrition information
   - e.g. [plain strawberries](https://fdc.nal.usda.gov/fdc-app.html#/food-details/747448/nutrients) or [C&H brown sugar](https://fdc.nal.usda.gov/fdc-app.html#/food-details/392083/nutrients)
4. **ui** for viewing/editing/creating recipes.
5. **cli** for interacting with the server instead of the API
6. **server** for exposing graphQL + http APIs.
7. **notion** for extracting images from my Notion food log.

### todo

1. add authentication
2. deploy UI + backend properly
3. allow merging ingredients together
4. weights per-ingredient
5. meals w/ multiple recipes at set quantities
6. slugs for recipe urls
7. associate ingredients with fdc ids

range: min max approx unit kind (time, amount)
