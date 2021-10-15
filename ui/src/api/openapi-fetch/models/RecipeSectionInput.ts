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
    SectionIngredientInput,
    SectionIngredientInputFromJSON,
    SectionIngredientInputFromJSONTyped,
    SectionIngredientInputToJSON,
    SectionInstructionInput,
    SectionInstructionInputFromJSON,
    SectionInstructionInputFromJSONTyped,
    SectionInstructionInputToJSON,
    TimeRange,
    TimeRangeFromJSON,
    TimeRangeFromJSONTyped,
    TimeRangeToJSON,
} from './';

/**
 * A step in the recipe
 * @export
 * @interface RecipeSectionInput
 */
export interface RecipeSectionInput {
    /**
     * 
     * @type {TimeRange}
     * @memberof RecipeSectionInput
     */
    duration?: TimeRange;
    /**
     * x
     * @type {Array<SectionInstructionInput>}
     * @memberof RecipeSectionInput
     */
    instructions: Array<SectionInstructionInput>;
    /**
     * x
     * @type {Array<SectionIngredientInput>}
     * @memberof RecipeSectionInput
     */
    ingredients: Array<SectionIngredientInput>;
}

export function RecipeSectionInputFromJSON(json: any): RecipeSectionInput {
    return RecipeSectionInputFromJSONTyped(json, false);
}

export function RecipeSectionInputFromJSONTyped(json: any, ignoreDiscriminator: boolean): RecipeSectionInput {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'duration': !exists(json, 'duration') ? undefined : TimeRangeFromJSON(json['duration']),
        'instructions': ((json['instructions'] as Array<any>).map(SectionInstructionInputFromJSON)),
        'ingredients': ((json['ingredients'] as Array<any>).map(SectionIngredientInputFromJSON)),
    };
}

export function RecipeSectionInputToJSON(value?: RecipeSectionInput | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'duration': TimeRangeToJSON(value.duration),
        'instructions': ((value.instructions as Array<any>).map(SectionInstructionInputToJSON)),
        'ingredients': ((value.ingredients as Array<any>).map(SectionIngredientInputToJSON)),
    };
}

