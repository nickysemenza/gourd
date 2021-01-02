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

import { exists, mapValues } from '../runtime';
import {
    IngredientDetail,
    IngredientDetailFromJSON,
    IngredientDetailFromJSONTyped,
    IngredientDetailToJSON,
    List,
    ListFromJSON,
    ListFromJSONTyped,
    ListToJSON,
} from './';

/**
 * pages of IngredientDetail
 * @export
 * @interface PaginatedIngredients
 */
export interface PaginatedIngredients {
    /**
     * 
     * @type {Array<IngredientDetail>}
     * @memberof PaginatedIngredients
     */
    ingredients?: Array<IngredientDetail>;
    /**
     * 
     * @type {List}
     * @memberof PaginatedIngredients
     */
    meta?: List;
}

export function PaginatedIngredientsFromJSON(json: any): PaginatedIngredients {
    return PaginatedIngredientsFromJSONTyped(json, false);
}

export function PaginatedIngredientsFromJSONTyped(json: any, ignoreDiscriminator: boolean): PaginatedIngredients {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'ingredients': !exists(json, 'ingredients') ? undefined : ((json['ingredients'] as Array<any>).map(IngredientDetailFromJSON)),
        'meta': !exists(json, 'meta') ? undefined : ListFromJSON(json['meta']),
    };
}

export function PaginatedIngredientsToJSON(value?: PaginatedIngredients | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'ingredients': value.ingredients === undefined ? undefined : ((value.ingredients as Array<any>).map(IngredientDetailToJSON)),
        'meta': ListToJSON(value.meta),
    };
}


