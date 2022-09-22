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
import type { Items } from './Items';
import {
    ItemsFromJSON,
    ItemsFromJSONTyped,
    ItemsToJSON,
} from './Items';
import type { Meal } from './Meal';
import {
    MealFromJSON,
    MealFromJSONTyped,
    MealToJSON,
} from './Meal';

/**
 * pages of Meal
 * @export
 * @interface PaginatedMeals
 */
export interface PaginatedMeals {
    /**
     * 
     * @type {Array<Meal>}
     * @memberof PaginatedMeals
     */
    meals?: Array<Meal>;
    /**
     * 
     * @type {Items}
     * @memberof PaginatedMeals
     */
    meta: Items;
}

/**
 * Check if a given object implements the PaginatedMeals interface.
 */
export function instanceOfPaginatedMeals(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "meta" in value;

    return isInstance;
}

export function PaginatedMealsFromJSON(json: any): PaginatedMeals {
    return PaginatedMealsFromJSONTyped(json, false);
}

export function PaginatedMealsFromJSONTyped(json: any, ignoreDiscriminator: boolean): PaginatedMeals {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'meals': !exists(json, 'meals') ? undefined : ((json['meals'] as Array<any>).map(MealFromJSON)),
        'meta': ItemsFromJSON(json['meta']),
    };
}

export function PaginatedMealsToJSON(value?: PaginatedMeals | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'meals': value.meals === undefined ? undefined : ((value.meals as Array<any>).map(MealToJSON)),
        'meta': ItemsToJSON(value.meta),
    };
}

