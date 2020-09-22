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

// Error defines model for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// Ingredient defines model for Ingredient.
type Ingredient struct {

	// Ingredients that are equivalent
	Children *[]Ingredient `json:"children,omitempty"`

	// UUID
	Id string `json:"id"`

	// Ingredient name
	Name string `json:"name"`

	// Recipes referencing this ingredient
	Recipes *[]Recipe `json:"recipes,omitempty"`
}

// List defines model for List.
type List struct {

	// How many items were requested for this page
	Limit int `json:"limit"`

	// todo
	Offset int `json:"offset"`

	// Total number of pages available
	PageCount int `json:"page_count"`

	// What number page this is
	PageNumber int `json:"page_number"`

	// Total number of items across all pages
	TotalCount int `json:"total_count"`
}

// PaginatedIngredients defines model for PaginatedIngredients.
type PaginatedIngredients struct {
	Ingredients *[]Ingredient `json:"ingredients,omitempty"`
	Meta        *List         `json:"meta,omitempty"`
}

// PaginatedRecipes defines model for PaginatedRecipes.
type PaginatedRecipes struct {
	Meta    *List     `json:"meta,omitempty"`
	Recipes *[]Recipe `json:"recipes,omitempty"`
}

// Recipe defines model for Recipe.
type Recipe struct {

	// UUID
	Id string `json:"id"`

	// recipe name
	Name string `json:"name"`

	// serving quantity
	Quantity int `json:"quantity"`

	// num servings
	Servings *int `json:"servings,omitempty"`

	// book or website? deprecated?
	Source *string `json:"source,omitempty"`

	// todo
	TotalMinutes *int `json:"total_minutes,omitempty"`

	// serving unit
	Unit string `json:"unit"`
}

// RecipeDetail defines model for RecipeDetail.
type RecipeDetail struct {

	// A recipe
	Recipe Recipe `json:"recipe"`

	// sections of the recipe
	Sections []RecipeSection `json:"sections"`
}

// RecipeSection defines model for RecipeSection.
type RecipeSection struct {

	// UUID
	Id string `json:"id"`

	// x
	Ingredients []SectionIngredient `json:"ingredients"`

	// x
	Instructions []SectionInstruction `json:"instructions"`

	// How many minutes the step takes, approximately (todo - make this a range)
	Minutes int `json:"minutes"`
}

// SectionIngredient defines model for SectionIngredient.
type SectionIngredient struct {

	// weight in grams
	Grams *float32 `json:"grams,omitempty"`

	// UUID
	Id string `json:"id"`

	// An Ingredient
	Ingredient *Ingredient `json:"ingredient,omitempty"`

	// what kind of ingredient
	Kind string `json:"kind"`

	// A recipe
	Recipe *Recipe `json:"recipe,omitempty"`
}

// SectionInstruction defines model for SectionInstruction.
type SectionInstruction struct {

	// UUID
	Id string `json:"id"`

	// instruction
	Instruction string `json:"instruction"`
}

// LimitParam defines model for limitParam.
type LimitParam int

// OffsetParam defines model for offsetParam.
type OffsetParam int

// ListIngredientsParams defines parameters for ListIngredients.
type ListIngredientsParams struct {

	// The number of items to skip before starting to collect the result set.
	Offset *OffsetParam `json:"offset,omitempty"`

	// The numbers of items to return.
	Limit *LimitParam `json:"limit,omitempty"`
}

// ListRecipesParams defines parameters for ListRecipes.
type ListRecipesParams struct {

	// The number of items to skip before starting to collect the result set.
	Offset *OffsetParam `json:"offset,omitempty"`

	// The numbers of items to return.
	Limit *LimitParam `json:"limit,omitempty"`
}

// CreateRecipesJSONBody defines parameters for CreateRecipes.
type CreateRecipesJSONBody RecipeDetail

// CreateRecipesRequestBody defines body for CreateRecipes for application/json ContentType.
type CreateRecipesJSONRequestBody CreateRecipesJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List all ingredients
	// (GET /ingredients)
	ListIngredients(ctx echo.Context, params ListIngredientsParams) error
	// List all recipes
	// (GET /recipes)
	ListRecipes(ctx echo.Context, params ListRecipesParams) error
	// Create a recipe
	// (POST /recipes)
	CreateRecipes(ctx echo.Context) error
	// Info for a specific recipe
	// (GET /recipes/{recipe_id})
	GetRecipeById(ctx echo.Context, recipeId string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ListIngredients converts echo context to params.
func (w *ServerInterfaceWrapper) ListIngredients(ctx echo.Context) error {
	var err error

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

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ListIngredients(ctx, params)
	return err
}

// ListRecipes converts echo context to params.
func (w *ServerInterfaceWrapper) ListRecipes(ctx echo.Context) error {
	var err error

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

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateRecipes(ctx)
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

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetRecipeById(ctx, recipeId)
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

	router.GET("/ingredients", wrapper.ListIngredients)
	router.GET("/recipes", wrapper.ListRecipes)
	router.POST("/recipes", wrapper.CreateRecipes)
	router.GET("/recipes/:recipe_id", wrapper.GetRecipeById)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"X/UVu5Nsc7a/2/8VAAD//0eboQrZGAAA",
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
