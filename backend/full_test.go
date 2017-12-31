package main

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/nickysemenza/food/backend/app"
	"github.com/nickysemenza/food/backend/app/handler"
	"github.com/nickysemenza/food/backend/app/model"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var env *handler.Env

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r := *env.Router
	r.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestMain(m *testing.M) {
	godotenv.Load()
	globalConfig := app.GetConfig()
	mainApp := &app.App{}
	env = mainApp.Initialize(globalConfig)

	env.DB = model.DBReset(env.DB)
	env.DB = model.DBMigrate(env.DB)
	os.Exit(m.Run())
}
func makeRecipe(t *testing.T, slug string) *httptest.ResponseRecorder {
	newData := struct {
		Slug  string `json:"slug"`
		Title string `json:"title"`
	}{
		Slug: slug,
	}
	return doRequest(t, "POST", "/recipes", newData)
}
func getRecipe(t *testing.T, slug string) model.Recipe {
	response := doRequest(t, "GET", "/recipes/"+slug, nil)
	checkResponseCode(t, http.StatusOK, response.Code)

	decoder := json.NewDecoder(response.Body)
	var testRecipe model.Recipe
	decoder.Decode(&testRecipe)
	return testRecipe
}

func doRequest(t *testing.T, method string, url string, objToMarshall interface{}) *httptest.ResponseRecorder {
	if objToMarshall == nil {
		req, _ := http.NewRequest(method, url, nil)
		return executeRequest(req)
	} else {
		jsonValue, _ := json.Marshal(objToMarshall)
		req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonValue))
		return executeRequest(req)
	}
}
func TestAddRecipe(t *testing.T) {
	testRecipeSlug := "test-slug-TestAddRecipe"
	makeRecipe(t, testRecipeSlug)

	//make sure this recipe is in the list
	response := doRequest(t, "GET", "/recipes", nil)
	checkResponseCode(t, http.StatusOK, response.Code)
	if !strings.Contains(response.Body.String(), testRecipeSlug) {
		t.Error("did not find " + testRecipeSlug + " in recipe list!")
	}
}

func TestAddDeleteCategories(t *testing.T) {
	testRecipeSlug := "test-slug-1"
	makeRecipe(t, testRecipeSlug)

	//grab the recipe from detail
	testRecipe := getRecipe(t, testRecipeSlug)

	//add some categories
	testRecipe.Categories = []model.Category{
		{Name: "cat1"},
		{Name: "cat2"},
	}
	response := doRequest(t, "PUT", "/recipes/"+testRecipeSlug, testRecipe)
	checkResponseCode(t, http.StatusOK, response.Code)

	//ensure the categories were added
	testRecipe2 := getRecipe(t, testRecipeSlug)
	if numCategories := len(testRecipe2.Categories); numCategories != 2 {
		t.Errorf("found %d categories instead of 2", numCategories)
	}

	//now delete a category
	testRecipe.Categories = []model.Category{
		{Name: "cat2"},
	}
	doRequest(t, "PUT", "/recipes/"+testRecipeSlug, testRecipe)

	//make sure it was really deleted
	testRecipe3 := getRecipe(t, testRecipeSlug)
	if numCategories := len(testRecipe3.Categories); numCategories != 1 {
		t.Errorf("found %d categories instead of 1", numCategories)
	}
}
