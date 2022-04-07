/* tslint:disable */
/* eslint-disable */
/**
 * Gourd Recipe Database
 * API for https://github.com/nickysemenza/gourd
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: n@nickysemenza.com
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import * as runtime from '../runtime';
import {
    InlineObject,
    InlineObjectFromJSON,
    InlineObjectToJSON,
    InlineObject1,
    InlineObject1FromJSON,
    InlineObject1ToJSON,
    InlineResponse2001,
    InlineResponse2001FromJSON,
    InlineResponse2001ToJSON,
    InlineResponse2002,
    InlineResponse2002FromJSON,
    InlineResponse2002ToJSON,
    PaginatedRecipeWrappers,
    PaginatedRecipeWrappersFromJSON,
    PaginatedRecipeWrappersToJSON,
    RecipeDetail,
    RecipeDetailFromJSON,
    RecipeDetailToJSON,
    RecipeWrapper,
    RecipeWrapperFromJSON,
    RecipeWrapperToJSON,
    RecipeWrapperInput,
    RecipeWrapperInputFromJSON,
    RecipeWrapperInputToJSON,
    SearchResult,
    SearchResultFromJSON,
    SearchResultToJSON,
} from '../models';

export interface RecipesApiConvertIngredientToRecipeRequest {
    ingredientId: string;
}

export interface RecipesApiCreateRecipesRequest {
    recipeWrapperInput: RecipeWrapperInput;
}

export interface RecipesApiGetLatexByRecipeIdRequest {
    recipeId: string;
}

export interface RecipesApiGetRecipeByIdRequest {
    recipeId: string;
}

export interface RecipesApiGetRecipesByIdsRequest {
    recipeId: Array<string>;
}

export interface RecipesApiListRecipesRequest {
    offset?: number;
    limit?: number;
}

export interface RecipesApiScrapeRecipeRequest {
    inlineObject1: InlineObject1;
}

export interface RecipesApiSearchRequest {
    name: string;
    offset?: number;
    limit?: number;
}

export interface RecipesApiSumRecipesRequest {
    inlineObject: InlineObject;
}

/**
 * 
 */
export class RecipesApi extends runtime.BaseAPI {

    /**
     * todo
     * Converts an ingredient to a recipe, updating all recipes depending on it
     */
    async convertIngredientToRecipeRaw(requestParameters: RecipesApiConvertIngredientToRecipeRequest): Promise<runtime.ApiResponse<RecipeDetail>> {
        if (requestParameters.ingredientId === null || requestParameters.ingredientId === undefined) {
            throw new runtime.RequiredError('ingredientId','Required parameter requestParameters.ingredientId was null or undefined when calling convertIngredientToRecipe.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/ingredients/{ingredient_id}/convert_to_recipe`.replace(`{${"ingredient_id"}}`, encodeURIComponent(String(requestParameters.ingredientId))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => RecipeDetailFromJSON(jsonValue));
    }

    /**
     * todo
     * Converts an ingredient to a recipe, updating all recipes depending on it
     */
    async convertIngredientToRecipe(requestParameters: RecipesApiConvertIngredientToRecipeRequest): Promise<RecipeDetail> {
        const response = await this.convertIngredientToRecipeRaw(requestParameters);
        return await response.value();
    }

    /**
     * todo
     * Create a recipe
     */
    async createRecipesRaw(requestParameters: RecipesApiCreateRecipesRequest): Promise<runtime.ApiResponse<RecipeWrapper>> {
        if (requestParameters.recipeWrapperInput === null || requestParameters.recipeWrapperInput === undefined) {
            throw new runtime.RequiredError('recipeWrapperInput','Required parameter requestParameters.recipeWrapperInput was null or undefined when calling createRecipes.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/recipes`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: RecipeWrapperInputToJSON(requestParameters.recipeWrapperInput),
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => RecipeWrapperFromJSON(jsonValue));
    }

    /**
     * todo
     * Create a recipe
     */
    async createRecipes(requestParameters: RecipesApiCreateRecipesRequest): Promise<RecipeWrapper> {
        const response = await this.createRecipesRaw(requestParameters);
        return await response.value();
    }

    /**
     * todo
     * recipe as latex
     */
    async getLatexByRecipeIdRaw(requestParameters: RecipesApiGetLatexByRecipeIdRequest): Promise<runtime.ApiResponse<Blob>> {
        if (requestParameters.recipeId === null || requestParameters.recipeId === undefined) {
            throw new runtime.RequiredError('recipeId','Required parameter requestParameters.recipeId was null or undefined when calling getLatexByRecipeId.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/recipes/{recipe_id}/latex`.replace(`{${"recipe_id"}}`, encodeURIComponent(String(requestParameters.recipeId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.BlobApiResponse(response);
    }

    /**
     * todo
     * recipe as latex
     */
    async getLatexByRecipeId(requestParameters: RecipesApiGetLatexByRecipeIdRequest): Promise<Blob> {
        const response = await this.getLatexByRecipeIdRaw(requestParameters);
        return await response.value();
    }

    /**
     * todo
     * Info for a specific recipe
     */
    async getRecipeByIdRaw(requestParameters: RecipesApiGetRecipeByIdRequest): Promise<runtime.ApiResponse<RecipeWrapper>> {
        if (requestParameters.recipeId === null || requestParameters.recipeId === undefined) {
            throw new runtime.RequiredError('recipeId','Required parameter requestParameters.recipeId was null or undefined when calling getRecipeById.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/recipes/{recipe_id}`.replace(`{${"recipe_id"}}`, encodeURIComponent(String(requestParameters.recipeId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => RecipeWrapperFromJSON(jsonValue));
    }

    /**
     * todo
     * Info for a specific recipe
     */
    async getRecipeById(requestParameters: RecipesApiGetRecipeByIdRequest): Promise<RecipeWrapper> {
        const response = await this.getRecipeByIdRaw(requestParameters);
        return await response.value();
    }

    /**
     * get recipes by ids
     * Get recipes
     */
    async getRecipesByIdsRaw(requestParameters: RecipesApiGetRecipesByIdsRequest): Promise<runtime.ApiResponse<PaginatedRecipeWrappers>> {
        if (requestParameters.recipeId === null || requestParameters.recipeId === undefined) {
            throw new runtime.RequiredError('recipeId','Required parameter requestParameters.recipeId was null or undefined when calling getRecipesByIds.');
        }

        const queryParameters: any = {};

        if (requestParameters.recipeId) {
            queryParameters['recipe_id'] = requestParameters.recipeId;
        }

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/recipes/bulk`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => PaginatedRecipeWrappersFromJSON(jsonValue));
    }

    /**
     * get recipes by ids
     * Get recipes
     */
    async getRecipesByIds(requestParameters: RecipesApiGetRecipesByIdsRequest): Promise<PaginatedRecipeWrappers> {
        const response = await this.getRecipesByIdsRaw(requestParameters);
        return await response.value();
    }

    /**
     * todo
     * List all recipes
     */
    async listRecipesRaw(requestParameters: RecipesApiListRecipesRequest): Promise<runtime.ApiResponse<PaginatedRecipeWrappers>> {
        const queryParameters: any = {};

        if (requestParameters.offset !== undefined) {
            queryParameters['offset'] = requestParameters.offset;
        }

        if (requestParameters.limit !== undefined) {
            queryParameters['limit'] = requestParameters.limit;
        }

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/recipes`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => PaginatedRecipeWrappersFromJSON(jsonValue));
    }

    /**
     * todo
     * List all recipes
     */
    async listRecipes(requestParameters: RecipesApiListRecipesRequest): Promise<PaginatedRecipeWrappers> {
        const response = await this.listRecipesRaw(requestParameters);
        return await response.value();
    }

    /**
     * recipe dependencies
     * Get foods
     */
    async recipeDependenciesRaw(): Promise<runtime.ApiResponse<InlineResponse2002>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/data/recipe_dependencies`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => InlineResponse2002FromJSON(jsonValue));
    }

    /**
     * recipe dependencies
     * Get foods
     */
    async recipeDependencies(): Promise<InlineResponse2002> {
        const response = await this.recipeDependenciesRaw();
        return await response.value();
    }

    /**
     * todo
     * scrape a recipe by URL
     */
    async scrapeRecipeRaw(requestParameters: RecipesApiScrapeRecipeRequest): Promise<runtime.ApiResponse<RecipeWrapper>> {
        if (requestParameters.inlineObject1 === null || requestParameters.inlineObject1 === undefined) {
            throw new runtime.RequiredError('inlineObject1','Required parameter requestParameters.inlineObject1 was null or undefined when calling scrapeRecipe.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/recipes/scrape`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: InlineObject1ToJSON(requestParameters.inlineObject1),
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => RecipeWrapperFromJSON(jsonValue));
    }

    /**
     * todo
     * scrape a recipe by URL
     */
    async scrapeRecipe(requestParameters: RecipesApiScrapeRecipeRequest): Promise<RecipeWrapper> {
        const response = await this.scrapeRecipeRaw(requestParameters);
        return await response.value();
    }

    /**
     * todo
     * Search recipes and ingredients
     */
    async searchRaw(requestParameters: RecipesApiSearchRequest): Promise<runtime.ApiResponse<SearchResult>> {
        if (requestParameters.name === null || requestParameters.name === undefined) {
            throw new runtime.RequiredError('name','Required parameter requestParameters.name was null or undefined when calling search.');
        }

        const queryParameters: any = {};

        if (requestParameters.offset !== undefined) {
            queryParameters['offset'] = requestParameters.offset;
        }

        if (requestParameters.limit !== undefined) {
            queryParameters['limit'] = requestParameters.limit;
        }

        if (requestParameters.name !== undefined) {
            queryParameters['name'] = requestParameters.name;
        }

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/search`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => SearchResultFromJSON(jsonValue));
    }

    /**
     * todo
     * Search recipes and ingredients
     */
    async search(requestParameters: RecipesApiSearchRequest): Promise<SearchResult> {
        const response = await this.searchRaw(requestParameters);
        return await response.value();
    }

    /**
     * sums up the given recipes
     * sum up recipes
     */
    async sumRecipesRaw(requestParameters: RecipesApiSumRecipesRequest): Promise<runtime.ApiResponse<InlineResponse2001>> {
        if (requestParameters.inlineObject === null || requestParameters.inlineObject === undefined) {
            throw new runtime.RequiredError('inlineObject','Required parameter requestParameters.inlineObject was null or undefined when calling sumRecipes.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/recipes/sum`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: InlineObjectToJSON(requestParameters.inlineObject),
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => InlineResponse2001FromJSON(jsonValue));
    }

    /**
     * sums up the given recipes
     * sum up recipes
     */
    async sumRecipes(requestParameters: RecipesApiSumRecipesRequest): Promise<InlineResponse2001> {
        const response = await this.sumRecipesRaw(requestParameters);
        return await response.value();
    }

}
