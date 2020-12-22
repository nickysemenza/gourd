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
/**
 * An Ingredient
 * @export
 * @interface Ingredient
 */
export interface Ingredient {
    /**
     * id
     * @type {string}
     * @memberof Ingredient
     */
    id: string;
    /**
     * Ingredient name
     * @type {string}
     * @memberof Ingredient
     */
    name: string;
}

export function IngredientFromJSON(json: any): Ingredient {
    return IngredientFromJSONTyped(json, false);
}

export function IngredientFromJSONTyped(json: any, ignoreDiscriminator: boolean): Ingredient {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'name': json['name'],
    };
}

export function IngredientToJSON(value?: Ingredient | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'name': value.name,
    };
}


