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
import type { FoodInfo } from './FoodInfo';
import {
    FoodInfoFromJSON,
    FoodInfoFromJSONTyped,
    FoodInfoToJSON,
} from './FoodInfo';
import type { FoodResultByItem } from './FoodResultByItem';
import {
    FoodResultByItemFromJSON,
    FoodResultByItemFromJSONTyped,
    FoodResultByItemToJSON,
} from './FoodResultByItem';

/**
 * A meal, which bridges recipes to photos
 * @export
 * @interface FoodSearchResult
 */
export interface FoodSearchResult {
    /**
     * 
     * @type {Array<FoodInfo>}
     * @memberof FoodSearchResult
     */
    foods: Array<FoodInfo>;
    /**
     * 
     * @type {FoodResultByItem}
     * @memberof FoodSearchResult
     */
    results?: FoodResultByItem;
}

/**
 * Check if a given object implements the FoodSearchResult interface.
 */
export function instanceOfFoodSearchResult(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "foods" in value;

    return isInstance;
}

export function FoodSearchResultFromJSON(json: any): FoodSearchResult {
    return FoodSearchResultFromJSONTyped(json, false);
}

export function FoodSearchResultFromJSONTyped(json: any, ignoreDiscriminator: boolean): FoodSearchResult {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'foods': ((json['foods'] as Array<any>).map(FoodInfoFromJSON)),
        'results': !exists(json, 'results') ? undefined : FoodResultByItemFromJSON(json['results']),
    };
}

export function FoodSearchResultToJSON(value?: FoodSearchResult | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'foods': ((value.foods as Array<any>).map(FoodInfoToJSON)),
        'results': FoodResultByItemToJSON(value.results),
    };
}

