// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
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
	// Get app config
	// (GET /config)
	GetConfig(ctx echo.Context) error
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
	// Converts an ingredient to a recipe, updating all recipes depending on it
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
	// load mappings
	// (POST /meta/load_ingredient_mappings)
	LoadIngredientMappings(ctx echo.Context) error
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
	// scrape a recipe by URL
	// (POST /recipes/scrape)
	ScrapeRecipe(ctx echo.Context) error
	// sum up recipes
	// (POST /recipes/sum)
	SumRecipes(ctx echo.Context) error
	// Info for a specific recipe
	// (GET /recipes/{recipe_id})
	GetRecipeById(ctx echo.Context, recipeId string) error
	// recipe as latex
	// (GET /recipes/{recipe_id}/latex)
	GetLatexByRecipeId(ctx echo.Context, recipeId string) error
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

	ctx.Set(BearerAuthScopes, []string{""})

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

// GetConfig converts echo context to params.
func (w *ServerInterfaceWrapper) GetConfig(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetConfig(ctx)
	return err
}

// RecipeDependencies converts echo context to params.
func (w *ServerInterfaceWrapper) RecipeDependencies(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.RecipeDependencies(ctx)
	return err
}

// GetFoodsByIds converts echo context to params.
func (w *ServerInterfaceWrapper) GetFoodsByIds(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

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

	ctx.Set(BearerAuthScopes, []string{""})

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

	err = runtime.BindStyledParameterWithLocation("simple", false, "fdc_id", runtime.ParamLocationPath, ctx.Param("fdc_id"), &fdcId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter fdc_id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetFoodById(ctx, fdcId)
	return err
}

// ListIngredients converts echo context to params.
func (w *ServerInterfaceWrapper) ListIngredients(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

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

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateIngredients(ctx)
	return err
}

// GetIngredientById converts echo context to params.
func (w *ServerInterfaceWrapper) GetIngredientById(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "ingredient_id" -------------
	var ingredientId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "ingredient_id", runtime.ParamLocationPath, ctx.Param("ingredient_id"), &ingredientId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter ingredient_id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetIngredientById(ctx, ingredientId)
	return err
}

// AssociateFoodWithIngredient converts echo context to params.
func (w *ServerInterfaceWrapper) AssociateFoodWithIngredient(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "ingredient_id" -------------
	var ingredientId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "ingredient_id", runtime.ParamLocationPath, ctx.Param("ingredient_id"), &ingredientId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter ingredient_id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

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

	err = runtime.BindStyledParameterWithLocation("simple", false, "ingredient_id", runtime.ParamLocationPath, ctx.Param("ingredient_id"), &ingredientId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter ingredient_id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ConvertIngredientToRecipe(ctx, ingredientId)
	return err
}

// MergeIngredients converts echo context to params.
func (w *ServerInterfaceWrapper) MergeIngredients(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "ingredient_id" -------------
	var ingredientId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "ingredient_id", runtime.ParamLocationPath, ctx.Param("ingredient_id"), &ingredientId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter ingredient_id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.MergeIngredients(ctx, ingredientId)
	return err
}

// ListMeals converts echo context to params.
func (w *ServerInterfaceWrapper) ListMeals(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

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

	err = runtime.BindStyledParameterWithLocation("simple", false, "meal_id", runtime.ParamLocationPath, ctx.Param("meal_id"), &mealId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter meal_id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetMealById(ctx, mealId)
	return err
}

// UpdateRecipesForMeal converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateRecipesForMeal(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "meal_id" -------------
	var mealId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "meal_id", runtime.ParamLocationPath, ctx.Param("meal_id"), &mealId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter meal_id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UpdateRecipesForMeal(ctx, mealId)
	return err
}

// LoadIngredientMappings converts echo context to params.
func (w *ServerInterfaceWrapper) LoadIngredientMappings(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.LoadIngredientMappings(ctx)
	return err
}

// ListPhotos converts echo context to params.
func (w *ServerInterfaceWrapper) ListPhotos(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

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

	ctx.Set(BearerAuthScopes, []string{""})

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

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateRecipes(ctx)
	return err
}

// GetRecipesByIds converts echo context to params.
func (w *ServerInterfaceWrapper) GetRecipesByIds(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

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

// ScrapeRecipe converts echo context to params.
func (w *ServerInterfaceWrapper) ScrapeRecipe(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ScrapeRecipe(ctx)
	return err
}

// SumRecipes converts echo context to params.
func (w *ServerInterfaceWrapper) SumRecipes(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SumRecipes(ctx)
	return err
}

// GetRecipeById converts echo context to params.
func (w *ServerInterfaceWrapper) GetRecipeById(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "recipe_id" -------------
	var recipeId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "recipe_id", runtime.ParamLocationPath, ctx.Param("recipe_id"), &recipeId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter recipe_id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetRecipeById(ctx, recipeId)
	return err
}

// GetLatexByRecipeId converts echo context to params.
func (w *ServerInterfaceWrapper) GetLatexByRecipeId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "recipe_id" -------------
	var recipeId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "recipe_id", runtime.ParamLocationPath, ctx.Param("recipe_id"), &recipeId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter recipe_id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetLatexByRecipeId(ctx, recipeId)
	return err
}

// Search converts echo context to params.
func (w *ServerInterfaceWrapper) Search(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

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
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/albums", wrapper.ListAllAlbums)
	router.POST(baseURL+"/auth", wrapper.AuthLogin)
	router.GET(baseURL+"/config", wrapper.GetConfig)
	router.GET(baseURL+"/data/recipe_dependencies", wrapper.RecipeDependencies)
	router.GET(baseURL+"/foods/bulk", wrapper.GetFoodsByIds)
	router.GET(baseURL+"/foods/search", wrapper.SearchFoods)
	router.GET(baseURL+"/foods/:fdc_id", wrapper.GetFoodById)
	router.GET(baseURL+"/ingredients", wrapper.ListIngredients)
	router.POST(baseURL+"/ingredients", wrapper.CreateIngredients)
	router.GET(baseURL+"/ingredients/:ingredient_id", wrapper.GetIngredientById)
	router.POST(baseURL+"/ingredients/:ingredient_id/associate_food", wrapper.AssociateFoodWithIngredient)
	router.POST(baseURL+"/ingredients/:ingredient_id/convert_to_recipe", wrapper.ConvertIngredientToRecipe)
	router.POST(baseURL+"/ingredients/:ingredient_id/merge", wrapper.MergeIngredients)
	router.GET(baseURL+"/meals", wrapper.ListMeals)
	router.GET(baseURL+"/meals/:meal_id", wrapper.GetMealById)
	router.PATCH(baseURL+"/meals/:meal_id/recipes", wrapper.UpdateRecipesForMeal)
	router.POST(baseURL+"/meta/load_ingredient_mappings", wrapper.LoadIngredientMappings)
	router.GET(baseURL+"/photos", wrapper.ListPhotos)
	router.GET(baseURL+"/recipes", wrapper.ListRecipes)
	router.POST(baseURL+"/recipes", wrapper.CreateRecipes)
	router.GET(baseURL+"/recipes/bulk", wrapper.GetRecipesByIds)
	router.POST(baseURL+"/recipes/scrape", wrapper.ScrapeRecipe)
	router.POST(baseURL+"/recipes/sum", wrapper.SumRecipes)
	router.GET(baseURL+"/recipes/:recipe_id", wrapper.GetRecipeById)
	router.GET(baseURL+"/recipes/:recipe_id/latex", wrapper.GetLatexByRecipeId)
	router.GET(baseURL+"/search", wrapper.Search)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"aBMpBQ2zsVH3N/f/HwAA//9T6avpE4MAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
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

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
