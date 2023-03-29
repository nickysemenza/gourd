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
/**
 * 
 * @export
 * @interface WweiaFoodCategory
 */
export interface WweiaFoodCategory {
    /**
     * 
     * @type {number}
     * @memberof WweiaFoodCategory
     */
    wweia_food_category_code?: number;
    /**
     * 
     * @type {string}
     * @memberof WweiaFoodCategory
     */
    wweia_food_category_description?: string;
}

/**
 * Check if a given object implements the WweiaFoodCategory interface.
 */
export function instanceOfWweiaFoodCategory(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function WweiaFoodCategoryFromJSON(json: any): WweiaFoodCategory {
    return WweiaFoodCategoryFromJSONTyped(json, false);
}

export function WweiaFoodCategoryFromJSONTyped(json: any, ignoreDiscriminator: boolean): WweiaFoodCategory {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'wweia_food_category_code': !exists(json, 'wweiaFoodCategoryCode') ? undefined : json['wweiaFoodCategoryCode'],
        'wweia_food_category_description': !exists(json, 'wweiaFoodCategoryDescription') ? undefined : json['wweiaFoodCategoryDescription'],
    };
}

export function WweiaFoodCategoryToJSON(value?: WweiaFoodCategory | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'wweiaFoodCategoryCode': value.wweia_food_category_code,
        'wweiaFoodCategoryDescription': value.wweia_food_category_description,
    };
}

