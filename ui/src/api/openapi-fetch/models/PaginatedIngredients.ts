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

import { exists, mapValues } from '../runtime';
import type { IngredientDetail } from './IngredientDetail';
import {
    IngredientDetailFromJSON,
    IngredientDetailFromJSONTyped,
    IngredientDetailToJSON,
} from './IngredientDetail';
import type { Items } from './Items';
import {
    ItemsFromJSON,
    ItemsFromJSONTyped,
    ItemsToJSON,
} from './Items';

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
     * @type {Items}
     * @memberof PaginatedIngredients
     */
    meta: Items;
}

/**
 * Check if a given object implements the PaginatedIngredients interface.
 */
export function instanceOfPaginatedIngredients(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "meta" in value;

    return isInstance;
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
        'meta': ItemsFromJSON(json['meta']),
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
        'meta': ItemsToJSON(value.meta),
    };
}

