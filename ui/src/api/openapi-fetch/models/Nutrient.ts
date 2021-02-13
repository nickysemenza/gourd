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
import {
    FoodNutrientUnit,
    FoodNutrientUnitFromJSON,
    FoodNutrientUnitFromJSONTyped,
    FoodNutrientUnitToJSON,
} from './';

/**
 * todo
 * @export
 * @interface Nutrient
 */
export interface Nutrient {
    /**
     * todo
     * @type {number}
     * @memberof Nutrient
     */
    id: number;
    /**
     * todo
     * @type {string}
     * @memberof Nutrient
     */
    name: string;
    /**
     * 
     * @type {FoodNutrientUnit}
     * @memberof Nutrient
     */
    unit_name: FoodNutrientUnit;
}

export function NutrientFromJSON(json: any): Nutrient {
    return NutrientFromJSONTyped(json, false);
}

export function NutrientFromJSONTyped(json: any, ignoreDiscriminator: boolean): Nutrient {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'name': json['name'],
        'unit_name': FoodNutrientUnitFromJSON(json['unit_name']),
    };
}

export function NutrientToJSON(value?: Nutrient | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'name': value.name,
        'unit_name': FoodNutrientUnitToJSON(value.unit_name),
    };
}


