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
 * @interface BrandedFoodItemLabelNutrientsPotassium
 */
export interface BrandedFoodItemLabelNutrientsPotassium {
    /**
     * 
     * @type {number}
     * @memberof BrandedFoodItemLabelNutrientsPotassium
     */
    value?: number;
}

/**
 * Check if a given object implements the BrandedFoodItemLabelNutrientsPotassium interface.
 */
export function instanceOfBrandedFoodItemLabelNutrientsPotassium(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function BrandedFoodItemLabelNutrientsPotassiumFromJSON(json: any): BrandedFoodItemLabelNutrientsPotassium {
    return BrandedFoodItemLabelNutrientsPotassiumFromJSONTyped(json, false);
}

export function BrandedFoodItemLabelNutrientsPotassiumFromJSONTyped(json: any, ignoreDiscriminator: boolean): BrandedFoodItemLabelNutrientsPotassium {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'value': !exists(json, 'value') ? undefined : json['value'],
    };
}

export function BrandedFoodItemLabelNutrientsPotassiumToJSON(value?: BrandedFoodItemLabelNutrientsPotassium | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'value': value.value,
    };
}
