package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/icrowley/fake"
	"github.com/jinzhu/gorm"
	"github.com/nickysemenza/food/backend/app"
	"github.com/nickysemenza/food/backend/app/config"
	"github.com/nickysemenza/food/backend/app/model"
)

var db *gorm.DB

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r := app.SetupRouter(db)

	fmt.Println(db)

	r.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
func getAdminUser() *model.User {
	u := model.User{}
	u.Admin = true
	u.Email = fake.EmailAddress()
	db.Save(&u)
	return &u
}

func getGuestUser() *model.User {
	u := model.User{}
	u.Admin = false
	u.Email = fake.EmailAddress()
	db.Save(&u)
	return &u
}

func TestMain(m *testing.M) {
	c := &config.Config{
		DB: &config.DB{
			Dialect: "sqlite3",
			Name:    "/tmp/foo.sql",
		},
		Port: "7070",
	}
	a := app.NewApp(c)

	db = model.DBReset(a.Env.DB)
	db = model.DBMigrate(db)

	os.Exit(m.Run())
}
func makeRecipe(t *testing.T, slug string) *httptest.ResponseRecorder {
	newData := struct {
		Slug  string `json:"slug"`
		Title string `json:"title"`
	}{
		Slug: slug,
	}
	return doRequest(t, "POST", "/recipes", newData, getAdminUser())
}
func getRecipe(t *testing.T, slug string) model.Recipe {
	response := doRequest(t, "GET", "/recipes/"+slug, nil, nil)
	checkResponseCode(t, http.StatusOK, response.Code)

	decoder := json.NewDecoder(response.Body)
	var testRecipe model.Recipe
	decoder.Decode(&testRecipe)
	return testRecipe
}

func doRequest(t *testing.T, method string, url string, objToMarshall interface{}, user *model.User) *httptest.ResponseRecorder {
	var req *http.Request

	if objToMarshall == nil {
		req, _ = http.NewRequest(method, url, nil)
	} else {
		jsonValue, _ := json.Marshal(objToMarshall)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(jsonValue))
	}

	//add the JWT token as a header
	if user != nil {
		req.Header.Set("X-Jwt", user.GetJWTToken(db))
	}
	return executeRequest(req)

}
func TestAddRecipe(t *testing.T) {
	testRecipeSlug := "test-slug-TestAddRecipe"
	makeRecipe(t, testRecipeSlug)

	//make sure this recipe is in the list
	response := doRequest(t, "GET", "/recipes", nil, getAdminUser())
	checkResponseCode(t, http.StatusOK, response.Code)
	if !strings.Contains(response.Body.String(), testRecipeSlug) {
		t.Error("did not find " + testRecipeSlug + " in recipe list!")
	}
}

//TODO: fix auth and add this
func TestCannotModifyRecipeWithoutPermissions(t *testing.T) {
	testRecipeSlug := "test-slug-TestCannotAddRecipeWithoutPermissions"
	makeRecipe(t, testRecipeSlug)

	//grab the recipe from detail
	testRecipe := getRecipe(t, testRecipeSlug)

	//but we shouldn't be able to update it!
	response := doRequest(t, "PUT", "/recipes/"+testRecipeSlug, testRecipe, nil)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)

	response2 := doRequest(t, "PUT", "/recipes/"+testRecipeSlug, testRecipe, getGuestUser())
	checkResponseCode(t, http.StatusUnauthorized, response2.Code)
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
	response := doRequest(t, "PUT", "/recipes/"+testRecipeSlug, testRecipe, getAdminUser())
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
	doRequest(t, "PUT", "/recipes/"+testRecipeSlug, testRecipe, getAdminUser())

	//make sure it was really deleted
	testRecipe3 := getRecipe(t, testRecipeSlug)
	if numCategories := len(testRecipe3.Categories); numCategories != 1 {
		t.Errorf("found %d categories instead of 1", numCategories)
	}
}
