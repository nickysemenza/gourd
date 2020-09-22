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
import {
    Recipe,
    RecipeFromJSON,
    RecipeFromJSONTyped,
    RecipeToJSON,
    RecipeSection,
    RecipeSectionFromJSON,
    RecipeSectionFromJSONTyped,
    RecipeSectionToJSON,
} from './';

/**
 * A recipe with subcomponents
 * @export
 * @interface RecipeDetail
 */
export interface RecipeDetail {
    /**
     * sections of the recipe
     * @type {Array<RecipeSection>}
     * @memberof RecipeDetail
     */
    sections: Array<RecipeSection>;
    /**
     * 
     * @type {Recipe}
     * @memberof RecipeDetail
     */
    recipe: Recipe;
}

export function RecipeDetailFromJSON(json: any): RecipeDetail {
    return RecipeDetailFromJSONTyped(json, false);
}

export function RecipeDetailFromJSONTyped(json: any, ignoreDiscriminator: boolean): RecipeDetail {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'sections': ((json['sections'] as Array<any>).map(RecipeSectionFromJSON)),
        'recipe': RecipeFromJSON(json['recipe']),
    };
}

export function RecipeDetailToJSON(value?: RecipeDetail | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'sections': ((value.sections as Array<any>).map(RecipeSectionToJSON)),
        'recipe': RecipeToJSON(value.recipe),
    };
}


