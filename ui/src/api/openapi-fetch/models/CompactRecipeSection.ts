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
 * @interface CompactRecipeSection
 */
export interface CompactRecipeSection {
    /**
     * 
     * @type {Array<string>}
     * @memberof CompactRecipeSection
     */
    ingredients: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof CompactRecipeSection
     */
    instructions: Array<string>;
}

/**
 * Check if a given object implements the CompactRecipeSection interface.
 */
export function instanceOfCompactRecipeSection(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "ingredients" in value;
    isInstance = isInstance && "instructions" in value;

    return isInstance;
}

export function CompactRecipeSectionFromJSON(json: any): CompactRecipeSection {
    return CompactRecipeSectionFromJSONTyped(json, false);
}

export function CompactRecipeSectionFromJSONTyped(json: any, ignoreDiscriminator: boolean): CompactRecipeSection {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'ingredients': json['ingredients'],
        'instructions': json['instructions'],
    };
}

export function CompactRecipeSectionToJSON(value?: CompactRecipeSection | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'ingredients': value.ingredients,
        'instructions': value.instructions,
    };
}

