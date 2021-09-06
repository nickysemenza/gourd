// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List all albums
	// (GET /albums)
	ListAllAlbums(ctx echo.Context) error
	// Google Login callback
	// (POST /auth)
	AuthLogin(ctx echo.Context, params AuthLoginParams) error
	// Get foods
	// (GET /data/recipe_dependencies)
	RecipeDependencies(ctx echo.Context) error
	// Get foods
	// (GET /foods/bulk)
	GetFoodsByIds(ctx echo.Context, params GetFoodsByIdsParams) error
	// Search foods
	// (GET /foods/search)
	SearchFoods(ctx echo.Context, params SearchFoodsParams) error
	// get a FDC entry by id
	// (GET /foods/{fdc_id})
	GetFoodById(ctx echo.Context, fdcId int) error
	// List all ingredients
	// (GET /ingredients)
	ListIngredients(ctx echo.Context, params ListIngredientsParams) error
	// Create a ingredient
	// (POST /ingredients)
	CreateIngredients(ctx echo.Context) error
	// Get a specific ingredient
	// (GET /ingredients/{ingredient_id})
	GetIngredientById(ctx echo.Context, ingredientId string) error
	// Assosiates a food with a given ingredient
	// (POST /ingredients/{ingredient_id}/associate_food)
	AssociateFoodWithIngredient(ctx echo.Context, ingredientId string, params AssociateFoodWithIngredientParams) error
	// Converts an ingredient to a recipe, updating all recipes depending on it.
	// (POST /ingredients/{ingredient_id}/convert_to_recipe)
	ConvertIngredientToRecipe(ctx echo.Context, ingredientId string) error
	// Merges the provide ingredients in the body into the param
	// (POST /ingredients/{ingredient_id}/merge)
	MergeIngredients(ctx echo.Context, ingredientId string) error
	// List all meals
	// (GET /meals)
	ListMeals(ctx echo.Context, params ListMealsParams) error
	// Info for a specific meal
	// (GET /meals/{meal_id})
	GetMealById(ctx echo.Context, mealId string) error
	// Update the recipes associated with a given meal
	// (PATCH /meals/{meal_id}/recipes)
	UpdateRecipesForMeal(ctx echo.Context, mealId string) error
	// List all photos
	// (GET /photos)
	ListPhotos(ctx echo.Context, params ListPhotosParams) error
	// List all recipes
	// (GET /recipes)
	ListRecipes(ctx echo.Context, params ListRecipesParams) error
	// Create a recipe
	// (POST /recipes)
	CreateRecipes(ctx echo.Context) error
	// Get recipes
	// (GET /recipes/bulk)
	GetRecipesByIds(ctx echo.Context, params GetRecipesByIdsParams) error
	// Info for a specific recipe
	// (GET /recipes/{recipe_id})
	GetRecipeById(ctx echo.Context, recipeId string) error
	// Search recipes and ingredients
	// (GET /search)
	Search(ctx echo.Context, params SearchParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ListAllAlbums converts echo context to params.
func (w *ServerInterfaceWrapper) ListAllAlbums(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ListAllAlbums(ctx)
	return err
}

// AuthLogin converts echo context to params.
func (w *ServerInterfaceWrapper) AuthLogin(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params AuthLoginParams
	// ------------- Required query parameter "code" -------------

	err = runtime.BindQueryParameter("form", true, true, "code", ctx.QueryParams(), &params.Code)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter code: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AuthLogin(ctx, params)
	return err
}

// RecipeDependencies converts echo context to params.
func (w *ServerInterfaceWrapper) RecipeDependencies(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.RecipeDependencies(ctx)
	return err
}

// GetFoodsByIds converts echo context to params.
func (w *ServerInterfaceWrapper) GetFoodsByIds(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetFoodsByIdsParams
	// ------------- Required query parameter "fdc_id" -------------

	err = runtime.BindQueryParameter("form", true, true, "fdc_id", ctx.QueryParams(), &params.FdcId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter fdc_id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetFoodsByIds(ctx, params)
	return err
}

// SearchFoods converts echo context to params.
func (w *ServerInterfaceWrapper) SearchFoods(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params SearchFoodsParams
	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Required query parameter "name" -------------

	err = runtime.BindQueryParameter("form", true, true, "name", ctx.QueryParams(), &params.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter name: %s", err))
	}

	// ------------- Optional query parameter "data_types" -------------

	err = runtime.BindQueryParameter("form", true, false, "data_types", ctx.QueryParams(), &params.DataTypes)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter data_types: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SearchFoods(ctx, params)
	return err
}

// GetFoodById converts echo context to params.
func (w *ServerInterfaceWrapper) GetFoodById(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "fdc_id" -------------
	var fdcId int

	err = runtime.BindStyledParameter("simple", false, "fdc_id", ctx.Param("fdc_id"), &fdcId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter fdc_id: %s", err))
	}

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetFoodById(ctx, fdcId)
	return err
}

// ListIngredients converts echo context to params.
func (w *ServerInterfaceWrapper) ListIngredients(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListIngredientsParams
	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "ingredient_id" -------------

	err = runtime.BindQueryParameter("form", true, false, "ingredient_id", ctx.QueryParams(), &params.IngredientId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter ingredient_id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ListIngredients(ctx, params)
	return err
}

// CreateIngredients converts echo context to params.
func (w *ServerInterfaceWrapper) CreateIngredients(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateIngredients(ctx)
	return err
}

// GetIngredientById converts echo context to params.
func (w *ServerInterfaceWrapper) GetIngredientById(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "ingredient_id" -------------
	var ingredientId string

	err = runtime.BindStyledParameter("simple", false, "ingredient_id", ctx.Param("ingredient_id"), &ingredientId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter ingredient_id: %s", err))
	}

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetIngredientById(ctx, ingredientId)
	return err
}

// AssociateFoodWithIngredient converts echo context to params.
func (w *ServerInterfaceWrapper) AssociateFoodWithIngredient(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "ingredient_id" -------------
	var ingredientId string

	err = runtime.BindStyledParameter("simple", false, "ingredient_id", ctx.Param("ingredient_id"), &ingredientId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter ingredient_id: %s", err))
	}

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params AssociateFoodWithIngredientParams
	// ------------- Required query parameter "fdc_id" -------------

	err = runtime.BindQueryParameter("form", true, true, "fdc_id", ctx.QueryParams(), &params.FdcId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter fdc_id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AssociateFoodWithIngredient(ctx, ingredientId, params)
	return err
}

// ConvertIngredientToRecipe converts echo context to params.
func (w *ServerInterfaceWrapper) ConvertIngredientToRecipe(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "ingredient_id" -------------
	var ingredientId string

	err = runtime.BindStyledParameter("simple", false, "ingredient_id", ctx.Param("ingredient_id"), &ingredientId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter ingredient_id: %s", err))
	}

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ConvertIngredientToRecipe(ctx, ingredientId)
	return err
}

// MergeIngredients converts echo context to params.
func (w *ServerInterfaceWrapper) MergeIngredients(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "ingredient_id" -------------
	var ingredientId string

	err = runtime.BindStyledParameter("simple", false, "ingredient_id", ctx.Param("ingredient_id"), &ingredientId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter ingredient_id: %s", err))
	}

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.MergeIngredients(ctx, ingredientId)
	return err
}

// ListMeals converts echo context to params.
func (w *ServerInterfaceWrapper) ListMeals(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListMealsParams
	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ListMeals(ctx, params)
	return err
}

// GetMealById converts echo context to params.
func (w *ServerInterfaceWrapper) GetMealById(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "meal_id" -------------
	var mealId string

	err = runtime.BindStyledParameter("simple", false, "meal_id", ctx.Param("meal_id"), &mealId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter meal_id: %s", err))
	}

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetMealById(ctx, mealId)
	return err
}

// UpdateRecipesForMeal converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateRecipesForMeal(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "meal_id" -------------
	var mealId string

	err = runtime.BindStyledParameter("simple", false, "meal_id", ctx.Param("meal_id"), &mealId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter meal_id: %s", err))
	}

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UpdateRecipesForMeal(ctx, mealId)
	return err
}

// ListPhotos converts echo context to params.
func (w *ServerInterfaceWrapper) ListPhotos(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListPhotosParams
	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ListPhotos(ctx, params)
	return err
}

// ListRecipes converts echo context to params.
func (w *ServerInterfaceWrapper) ListRecipes(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListRecipesParams
	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ListRecipes(ctx, params)
	return err
}

// CreateRecipes converts echo context to params.
func (w *ServerInterfaceWrapper) CreateRecipes(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateRecipes(ctx)
	return err
}

// GetRecipesByIds converts echo context to params.
func (w *ServerInterfaceWrapper) GetRecipesByIds(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetRecipesByIdsParams
	// ------------- Required query parameter "recipe_id" -------------

	err = runtime.BindQueryParameter("form", true, true, "recipe_id", ctx.QueryParams(), &params.RecipeId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter recipe_id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetRecipesByIds(ctx, params)
	return err
}

// GetRecipeById converts echo context to params.
func (w *ServerInterfaceWrapper) GetRecipeById(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "recipe_id" -------------
	var recipeId string

	err = runtime.BindStyledParameter("simple", false, "recipe_id", ctx.Param("recipe_id"), &recipeId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter recipe_id: %s", err))
	}

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetRecipeById(ctx, recipeId)
	return err
}

// Search converts echo context to params.
func (w *ServerInterfaceWrapper) Search(ctx echo.Context) error {
	var err error

	ctx.Set("bearerAuth.Scopes", []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params SearchParams
	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Required query parameter "name" -------------

	err = runtime.BindQueryParameter("form", true, true, "name", ctx.QueryParams(), &params.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter name: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.Search(ctx, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET("/albums", wrapper.ListAllAlbums)
	router.POST("/auth", wrapper.AuthLogin)
	router.GET("/data/recipe_dependencies", wrapper.RecipeDependencies)
	router.GET("/foods/bulk", wrapper.GetFoodsByIds)
	router.GET("/foods/search", wrapper.SearchFoods)
	router.GET("/foods/:fdc_id", wrapper.GetFoodById)
	router.GET("/ingredients", wrapper.ListIngredients)
	router.POST("/ingredients", wrapper.CreateIngredients)
	router.GET("/ingredients/:ingredient_id", wrapper.GetIngredientById)
	router.POST("/ingredients/:ingredient_id/associate_food", wrapper.AssociateFoodWithIngredient)
	router.POST("/ingredients/:ingredient_id/convert_to_recipe", wrapper.ConvertIngredientToRecipe)
	router.POST("/ingredients/:ingredient_id/merge", wrapper.MergeIngredients)
	router.GET("/meals", wrapper.ListMeals)
	router.GET("/meals/:meal_id", wrapper.GetMealById)
	router.PATCH("/meals/:meal_id/recipes", wrapper.UpdateRecipesForMeal)
	router.GET("/photos", wrapper.ListPhotos)
	router.GET("/recipes", wrapper.ListRecipes)
	router.POST("/recipes", wrapper.CreateRecipes)
	router.GET("/recipes/bulk", wrapper.GetRecipesByIds)
	router.GET("/recipes/:recipe_id", wrapper.GetRecipeById)
	router.GET("/search", wrapper.Search)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
