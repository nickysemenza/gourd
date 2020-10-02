/* tslint:disable */
/* eslint-disable */
/**
 * Gourd Recipe Database
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
    RecipeDetail,
    RecipeDetailFromJSON,
    RecipeDetailToJSON,
} from '../models';

export interface CreateRecipesRequest {
    recipeDetail: RecipeDetail;
}

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
    async createRecipesRaw(requestParameters: CreateRecipesRequest): Promise<runtime.ApiResponse<RecipeDetail>> {
        if (requestParameters.recipeDetail === null || requestParameters.recipeDetail === undefined) {
            throw new runtime.RequiredError('recipeDetail','Required parameter requestParameters.recipeDetail was null or undefined when calling createRecipes.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/recipes`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: RecipeDetailToJSON(requestParameters.recipeDetail),
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => RecipeDetailFromJSON(jsonValue));
    }

    /**
     * Create a recipe
     */
    async createRecipes(requestParameters: CreateRecipesRequest): Promise<RecipeDetail> {
        const response = await this.createRecipesRaw(requestParameters);
        return await response.value();
    }

    /**
     * Info for a specific recipe
     */
    async getRecipeByIdRaw(requestParameters: GetRecipeByIdRequest): Promise<runtime.ApiResponse<RecipeDetail>> {
        if (requestParameters.recipeId === null || requestParameters.recipeId === undefined) {
            throw new runtime.RequiredError('recipeId','Required parameter requestParameters.recipeId was null or undefined when calling getRecipeById.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/recipes/{recipe_id}`.replace(`{${"recipe_id"}}`, encodeURIComponent(String(requestParameters.recipeId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => RecipeDetailFromJSON(jsonValue));
    }

    /**
     * Info for a specific recipe
     */
    async getRecipeById(requestParameters: GetRecipeByIdRequest): Promise<RecipeDetail> {
        const response = await this.getRecipeByIdRaw(requestParameters);
        return await response.value();
    }

    /**
     * List all recipes
     */
    async listRecipesRaw(requestParameters: ListRecipesRequest): Promise<runtime.ApiResponse<PaginatedRecipes>> {
        const queryParameters: any = {};

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
