/* tslint:disable */
/* eslint-disable */
/**
 * Swagger Recipestore
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import * as runtime from '../runtime';
import {
    PaginatedRecipes,
    PaginatedRecipesFromJSON,
    PaginatedRecipesToJSON,
    Recipe,
    RecipeFromJSON,
    RecipeToJSON,
} from '../models';

export interface GetRecipeByIdRequest {
    recipeId: string;
}

export interface ListRecipesRequest {
    offset?: number;
    limit?: number;
}

/**
 * 
 */
export class RecipesApi extends runtime.BaseAPI {

    /**
     * Create a recipe
     */
    async createRecipesRaw(): Promise<runtime.ApiResponse<Recipe>> {
        const queryParameters: runtime.HTTPQuery = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/recipes`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => RecipeFromJSON(jsonValue));
    }

    /**
     * Create a recipe
     */
    async createRecipes(): Promise<Recipe> {
        const response = await this.createRecipesRaw();
        return await response.value();
    }

    /**
     * Info for a specific recipe
     */
    async getRecipeByIdRaw(requestParameters: GetRecipeByIdRequest): Promise<runtime.ApiResponse<Recipe>> {
        if (requestParameters.recipeId === null || requestParameters.recipeId === undefined) {
            throw new runtime.RequiredError('recipeId','Required parameter requestParameters.recipeId was null or undefined when calling getRecipeById.');
        }

        const queryParameters: runtime.HTTPQuery = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/recipes/{recipe_id}`.replace(`{${"recipe_id"}}`, encodeURIComponent(String(requestParameters.recipeId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => RecipeFromJSON(jsonValue));
    }

    /**
     * Info for a specific recipe
     */
    async getRecipeById(requestParameters: GetRecipeByIdRequest): Promise<Recipe> {
        const response = await this.getRecipeByIdRaw(requestParameters);
        return await response.value();
    }

    /**
     * List all recipes
     */
    async listRecipesRaw(requestParameters: ListRecipesRequest): Promise<runtime.ApiResponse<PaginatedRecipes>> {
        const queryParameters: runtime.HTTPQuery = {};

        if (requestParameters.offset !== undefined) {
            queryParameters['offset'] = requestParameters.offset;
        }

        if (requestParameters.limit !== undefined) {
            queryParameters['limit'] = requestParameters.limit;
        }

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/recipes`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => PaginatedRecipesFromJSON(jsonValue));
    }

    /**
     * List all recipes
     */
    async listRecipes(requestParameters: ListRecipesRequest): Promise<PaginatedRecipes> {
        const response = await this.listRecipesRaw(requestParameters);
        return await response.value();
    }

}
