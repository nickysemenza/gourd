package main

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/nickysemenza/food/backend/app"
	"github.com/nickysemenza/food/backend/app/handler"
	"github.com/nickysemenza/food/backend/app/model"
	"log"
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
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	globalConfig := app.GetConfig()
	mainApp := &app.App{}
	env = mainApp.Initialize(globalConfig)

	env.DB = model.DBReset(env.DB)
	env.DB = model.DBMigrate(env.DB)
	os.Exit(m.Run())
}

func TestGetRecipes(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/recipes", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	log.Print(response.Body)

	//create a new entry
	testRecipeSlug := "test-slug-1"
	values := model.Recipe{
		Slug:  testRecipeSlug,
		Title: "title1",
	}
	jsonValue, _ := json.Marshal(values)
	req2, _ := http.NewRequest("PUT", "/api/recipes/"+testRecipeSlug, bytes.NewBuffer(jsonValue))
	response2 := executeRequest(req2)
	checkResponseCode(t, http.StatusOK, response2.Code)
	log.Print(response2.Body)

	//ensure that entry was persisted
	req3, _ := http.NewRequest("GET", "/api/recipes", nil)
	response3 := executeRequest(req3)
	checkResponseCode(t, http.StatusOK, response3.Code)
	if !strings.Contains(response3.Body.String(), testRecipeSlug) {
		t.Error("did not find " + testRecipeSlug + " in recipe list!")
	}

}
