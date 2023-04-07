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
 * @interface FoodComponent
 */
export interface FoodComponent {
    /**
     * 
     * @type {number}
     * @memberof FoodComponent
     */
    id?: number;
    /**
     * 
     * @type {string}
     * @memberof FoodComponent
     */
    name?: string;
    /**
     * 
     * @type {number}
     * @memberof FoodComponent
     */
    data_points?: number;
    /**
     * 
     * @type {number}
     * @memberof FoodComponent
     */
    gram_weight?: number;
    /**
     * 
     * @type {boolean}
     * @memberof FoodComponent
     */
    is_refuse?: boolean;
    /**
     * 
     * @type {number}
     * @memberof FoodComponent
     */
    min_year_acquired?: number;
    /**
     * 
     * @type {number}
     * @memberof FoodComponent
     */
    percent_weight?: number;
}

/**
 * Check if a given object implements the FoodComponent interface.
 */
export function instanceOfFoodComponent(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function FoodComponentFromJSON(json: any): FoodComponent {
    return FoodComponentFromJSONTyped(json, false);
}

export function FoodComponentFromJSONTyped(json: any, ignoreDiscriminator: boolean): FoodComponent {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': !exists(json, 'id') ? undefined : json['id'],
        'name': !exists(json, 'name') ? undefined : json['name'],
        'data_points': !exists(json, 'dataPoints') ? undefined : json['dataPoints'],
        'gram_weight': !exists(json, 'gramWeight') ? undefined : json['gramWeight'],
        'is_refuse': !exists(json, 'isRefuse') ? undefined : json['isRefuse'],
        'min_year_acquired': !exists(json, 'minYearAcquired') ? undefined : json['minYearAcquired'],
        'percent_weight': !exists(json, 'percentWeight') ? undefined : json['percentWeight'],
    };
}

export function FoodComponentToJSON(value?: FoodComponent | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'name': value.name,
        'dataPoints': value.data_points,
        'gramWeight': value.gram_weight,
        'isRefuse': value.is_refuse,
        'minYearAcquired': value.min_year_acquired,
        'percentWeight': value.percent_weight,
    };
}
