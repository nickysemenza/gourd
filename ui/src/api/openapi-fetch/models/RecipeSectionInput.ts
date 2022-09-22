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
import type { Amount } from './Amount';
import {
    AmountFromJSON,
    AmountFromJSONTyped,
    AmountToJSON,
} from './Amount';
import type { SectionIngredientInput } from './SectionIngredientInput';
import {
    SectionIngredientInputFromJSON,
    SectionIngredientInputFromJSONTyped,
    SectionIngredientInputToJSON,
} from './SectionIngredientInput';
import type { SectionInstructionInput } from './SectionInstructionInput';
import {
    SectionInstructionInputFromJSON,
    SectionInstructionInputFromJSONTyped,
    SectionInstructionInputToJSON,
} from './SectionInstructionInput';

/**
 * A step in the recipe
 * @export
 * @interface RecipeSectionInput
 */
export interface RecipeSectionInput {
    /**
     * 
     * @type {Amount}
     * @memberof RecipeSectionInput
     */
    duration?: Amount;
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

/**
 * Check if a given object implements the RecipeSectionInput interface.
 */
export function instanceOfRecipeSectionInput(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "instructions" in value;
    isInstance = isInstance && "ingredients" in value;

    return isInstance;
}

export function RecipeSectionInputFromJSON(json: any): RecipeSectionInput {
    return RecipeSectionInputFromJSONTyped(json, false);
}

export function RecipeSectionInputFromJSONTyped(json: any, ignoreDiscriminator: boolean): RecipeSectionInput {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'duration': !exists(json, 'duration') ? undefined : AmountFromJSON(json['duration']),
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
        
        'duration': AmountToJSON(value.duration),
        'instructions': ((value.instructions as Array<any>).map(SectionInstructionInputToJSON)),
        'ingredients': ((value.ingredients as Array<any>).map(SectionIngredientInputToJSON)),
    };
}

