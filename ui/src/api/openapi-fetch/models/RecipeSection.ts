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
    Amount,
    AmountFromJSON,
    AmountFromJSONTyped,
    AmountToJSON,
    SectionIngredient,
    SectionIngredientFromJSON,
    SectionIngredientFromJSONTyped,
    SectionIngredientToJSON,
    SectionInstruction,
    SectionInstructionFromJSON,
    SectionInstructionFromJSONTyped,
    SectionInstructionToJSON,
} from './';

/**
 * A step in the recipe
 * @export
 * @interface RecipeSection
 */
export interface RecipeSection {
    /**
     * id
     * @type {string}
     * @memberof RecipeSection
     */
    id: string;
    /**
     * 
     * @type {Amount}
     * @memberof RecipeSection
     */
    duration?: Amount;
    /**
     * x
     * @type {Array<SectionInstruction>}
     * @memberof RecipeSection
     */
    instructions: Array<SectionInstruction>;
    /**
     * x
     * @type {Array<SectionIngredient>}
     * @memberof RecipeSection
     */
    ingredients: Array<SectionIngredient>;
}

export function RecipeSectionFromJSON(json: any): RecipeSection {
    return RecipeSectionFromJSONTyped(json, false);
}

export function RecipeSectionFromJSONTyped(json: any, ignoreDiscriminator: boolean): RecipeSection {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'duration': !exists(json, 'duration') ? undefined : AmountFromJSON(json['duration']),
        'instructions': ((json['instructions'] as Array<any>).map(SectionInstructionFromJSON)),
        'ingredients': ((json['ingredients'] as Array<any>).map(SectionIngredientFromJSON)),
    };
}

export function RecipeSectionToJSON(value?: RecipeSection | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'duration': AmountToJSON(value.duration),
        'instructions': ((value.instructions as Array<any>).map(SectionInstructionToJSON)),
        'ingredients': ((value.ingredients as Array<any>).map(SectionIngredientToJSON)),
    };
}


