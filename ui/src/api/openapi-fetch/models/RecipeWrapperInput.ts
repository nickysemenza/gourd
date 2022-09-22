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
import type { RecipeDetailInput } from './RecipeDetailInput';
import {
    RecipeDetailInputFromJSON,
    RecipeDetailInputFromJSONTyped,
    RecipeDetailInputToJSON,
} from './RecipeDetailInput';

/**
 * A recipe with subcomponents
 * @export
 * @interface RecipeWrapperInput
 */
export interface RecipeWrapperInput {
    /**
     * id
     * @type {string}
     * @memberof RecipeWrapperInput
     */
    id?: string;
    /**
     * 
     * @type {RecipeDetailInput}
     * @memberof RecipeWrapperInput
     */
    detail: RecipeDetailInput;
}

/**
 * Check if a given object implements the RecipeWrapperInput interface.
 */
export function instanceOfRecipeWrapperInput(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "detail" in value;

    return isInstance;
}

export function RecipeWrapperInputFromJSON(json: any): RecipeWrapperInput {
    return RecipeWrapperInputFromJSONTyped(json, false);
}

export function RecipeWrapperInputFromJSONTyped(json: any, ignoreDiscriminator: boolean): RecipeWrapperInput {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': !exists(json, 'id') ? undefined : json['id'],
        'detail': RecipeDetailInputFromJSON(json['detail']),
    };
}

export function RecipeWrapperInputToJSON(value?: RecipeWrapperInput | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'detail': RecipeDetailInputToJSON(value.detail),
    };
}

