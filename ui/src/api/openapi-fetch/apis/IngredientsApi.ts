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
    Ingredient,
    IngredientFromJSON,
    IngredientToJSON,
    IngredientDetail,
    IngredientDetailFromJSON,
    IngredientDetailToJSON,
    IngredientMappingsPayload,
    IngredientMappingsPayloadFromJSON,
    IngredientMappingsPayloadToJSON,
    InlineObject2,
    InlineObject2FromJSON,
    InlineObject2ToJSON,
    InlineResponse2002,
    InlineResponse2002FromJSON,
    InlineResponse2002ToJSON,
    PaginatedIngredients,
    PaginatedIngredientsFromJSON,
    PaginatedIngredientsToJSON,
    RecipeDetail,
    RecipeDetailFromJSON,
    RecipeDetailToJSON,
    SearchResult,
    SearchResultFromJSON,
    SearchResultToJSON,
} from '../models';

export interface IngredientsApiAssociateFoodWithIngredientRequest {
    ingredientId: string;
    fdcId: number;
}

export interface IngredientsApiConvertIngredientToRecipeRequest {
    ingredientId: string;
}

export interface IngredientsApiCreateIngredientsRequest {
    ingredient: Ingredient;
}

export interface IngredientsApiGetIngredientByIdRequest {
    ingredientId: string;
}

export interface IngredientsApiListIngredientsRequest {
    offset?: number;
    limit?: number;
    ingredientId?: Array<string>;
}

export interface IngredientsApiLoadIngredientMappingsRequest {
    ingredientMappingsPayload: IngredientMappingsPayload;
}

export interface IngredientsApiMergeIngredientsRequest {
    ingredientId: string;
    inlineObject2: InlineObject2;
}

export interface IngredientsApiSearchRequest {
    name: string;
    offset?: number;
    limit?: number;
}

/**
 * 
 */
export class IngredientsApi extends runtime.BaseAPI {

    /**
     * todo
     * Assosiates a food with a given ingredient
     */
    async associateFoodWithIngredientRaw(requestParameters: IngredientsApiAssociateFoodWithIngredientRequest): Promise<runtime.ApiResponse<RecipeDetail>> {
        if (requestParameters.ingredientId === null || requestParameters.ingredientId === undefined) {
            throw new runtime.RequiredError('ingredientId','Required parameter requestParameters.ingredientId was null or undefined when calling associateFoodWithIngredient.');
        }

        if (requestParameters.fdcId === null || requestParameters.fdcId === undefined) {
            throw new runtime.RequiredError('fdcId','Required parameter requestParameters.fdcId was null or undefined when calling associateFoodWithIngredient.');
        }

        const queryParameters: any = {};

        if (requestParameters.fdcId !== undefined) {
            queryParameters['fdc_id'] = requestParameters.fdcId;
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
            path: `/ingredients/{ingredient_id}/associate_food`.replace(`{${"ingredient_id"}}`, encodeURIComponent(String(requestParameters.ingredientId))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => RecipeDetailFromJSON(jsonValue));
    }

    /**
     * todo
     * Assosiates a food with a given ingredient
     */
    async associateFoodWithIngredient(requestParameters: IngredientsApiAssociateFoodWithIngredientRequest): Promise<RecipeDetail> {
        const response = await this.associateFoodWithIngredientRaw(requestParameters);
        return await response.value();
    }

    /**
     * todo
     * Converts an ingredient to a recipe, updating all recipes depending on it
     */
    async convertIngredientToRecipeRaw(requestParameters: IngredientsApiConvertIngredientToRecipeRequest): Promise<runtime.ApiResponse<RecipeDetail>> {
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
    async convertIngredientToRecipe(requestParameters: IngredientsApiConvertIngredientToRecipeRequest): Promise<RecipeDetail> {
        const response = await this.convertIngredientToRecipeRaw(requestParameters);
        return await response.value();
    }

    /**
     * todo
     * Create a ingredient
     */
    async createIngredientsRaw(requestParameters: IngredientsApiCreateIngredientsRequest): Promise<runtime.ApiResponse<Ingredient>> {
        if (requestParameters.ingredient === null || requestParameters.ingredient === undefined) {
            throw new runtime.RequiredError('ingredient','Required parameter requestParameters.ingredient was null or undefined when calling createIngredients.');
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
            path: `/ingredients`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: IngredientToJSON(requestParameters.ingredient),
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => IngredientFromJSON(jsonValue));
    }

    /**
     * todo
     * Create a ingredient
     */
    async createIngredients(requestParameters: IngredientsApiCreateIngredientsRequest): Promise<Ingredient> {
        const response = await this.createIngredientsRaw(requestParameters);
        return await response.value();
    }

    /**
     * todo
     * Get a specific ingredient
     */
    async getIngredientByIdRaw(requestParameters: IngredientsApiGetIngredientByIdRequest): Promise<runtime.ApiResponse<IngredientDetail>> {
        if (requestParameters.ingredientId === null || requestParameters.ingredientId === undefined) {
            throw new runtime.RequiredError('ingredientId','Required parameter requestParameters.ingredientId was null or undefined when calling getIngredientById.');
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
            path: `/ingredients/{ingredient_id}`.replace(`{${"ingredient_id"}}`, encodeURIComponent(String(requestParameters.ingredientId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => IngredientDetailFromJSON(jsonValue));
    }

    /**
     * todo
     * Get a specific ingredient
     */
    async getIngredientById(requestParameters: IngredientsApiGetIngredientByIdRequest): Promise<IngredientDetail> {
        const response = await this.getIngredientByIdRaw(requestParameters);
        return await response.value();
    }

    /**
     * todo
     * List all ingredients
     */
    async listIngredientsRaw(requestParameters: IngredientsApiListIngredientsRequest): Promise<runtime.ApiResponse<PaginatedIngredients>> {
        const queryParameters: any = {};

        if (requestParameters.offset !== undefined) {
            queryParameters['offset'] = requestParameters.offset;
        }

        if (requestParameters.limit !== undefined) {
            queryParameters['limit'] = requestParameters.limit;
        }

        if (requestParameters.ingredientId) {
            queryParameters['ingredient_id'] = requestParameters.ingredientId;
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
            path: `/ingredients`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => PaginatedIngredientsFromJSON(jsonValue));
    }

    /**
     * todo
     * List all ingredients
     */
    async listIngredients(requestParameters: IngredientsApiListIngredientsRequest): Promise<PaginatedIngredients> {
        const response = await this.listIngredientsRaw(requestParameters);
        return await response.value();
    }

    /**
     * todo
     * load mappings
     */
    async loadIngredientMappingsRaw(requestParameters: IngredientsApiLoadIngredientMappingsRequest): Promise<runtime.ApiResponse<object>> {
        if (requestParameters.ingredientMappingsPayload === null || requestParameters.ingredientMappingsPayload === undefined) {
            throw new runtime.RequiredError('ingredientMappingsPayload','Required parameter requestParameters.ingredientMappingsPayload was null or undefined when calling loadIngredientMappings.');
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
            path: `/meta/load_ingredient_mappings`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: IngredientMappingsPayloadToJSON(requestParameters.ingredientMappingsPayload),
        });

        return new runtime.JSONApiResponse<any>(response);
    }

    /**
     * todo
     * load mappings
     */
    async loadIngredientMappings(requestParameters: IngredientsApiLoadIngredientMappingsRequest): Promise<object> {
        const response = await this.loadIngredientMappingsRaw(requestParameters);
        return await response.value();
    }

    /**
     * todo
     * Merges the provide ingredients in the body into the param
     */
    async mergeIngredientsRaw(requestParameters: IngredientsApiMergeIngredientsRequest): Promise<runtime.ApiResponse<Ingredient>> {
        if (requestParameters.ingredientId === null || requestParameters.ingredientId === undefined) {
            throw new runtime.RequiredError('ingredientId','Required parameter requestParameters.ingredientId was null or undefined when calling mergeIngredients.');
        }

        if (requestParameters.inlineObject2 === null || requestParameters.inlineObject2 === undefined) {
            throw new runtime.RequiredError('inlineObject2','Required parameter requestParameters.inlineObject2 was null or undefined when calling mergeIngredients.');
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
            path: `/ingredients/{ingredient_id}/merge`.replace(`{${"ingredient_id"}}`, encodeURIComponent(String(requestParameters.ingredientId))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: InlineObject2ToJSON(requestParameters.inlineObject2),
        });

        return new runtime.JSONApiResponse(response, (jsonValue) => IngredientFromJSON(jsonValue));
    }

    /**
     * todo
     * Merges the provide ingredients in the body into the param
     */
    async mergeIngredients(requestParameters: IngredientsApiMergeIngredientsRequest): Promise<Ingredient> {
        const response = await this.mergeIngredientsRaw(requestParameters);
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
     * Search recipes and ingredients
     */
    async searchRaw(requestParameters: IngredientsApiSearchRequest): Promise<runtime.ApiResponse<SearchResult>> {
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
    async search(requestParameters: IngredientsApiSearchRequest): Promise<SearchResult> {
        const response = await this.searchRaw(requestParameters);
        return await response.value();
    }

}
