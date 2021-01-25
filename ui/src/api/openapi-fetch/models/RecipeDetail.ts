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
    RecipeSection,
    RecipeSectionFromJSON,
    RecipeSectionFromJSONTyped,
    RecipeSectionToJSON,
    RecipeSource,
    RecipeSourceFromJSON,
    RecipeSourceFromJSONTyped,
    RecipeSourceToJSON,
} from './';

/**
 * A revision of a recipe
 * @export
 * @interface RecipeDetail
 */
export interface RecipeDetail {
    /**
     * id
     * @type {string}
     * @memberof RecipeDetail
     */
    id: string;
    /**
     * sections of the recipe
     * @type {Array<RecipeSection>}
     * @memberof RecipeDetail
     */
    sections: Array<RecipeSection>;
    /**
     * recipe name
     * @type {string}
     * @memberof RecipeDetail
     */
    name: string;
    /**
     * book or websites
     * @type {Array<RecipeSource>}
     * @memberof RecipeDetail
     */
    sources?: Array<RecipeSource>;
    /**
     * num servings
     * @type {number}
     * @memberof RecipeDetail
     */
    servings?: number;
    /**
     * serving quantity
     * @type {number}
     * @memberof RecipeDetail
     */
    quantity: number;
    /**
     * serving unit
     * @type {string}
     * @memberof RecipeDetail
     */
    unit: string;
    /**
     * version of the recipe
     * @type {number}
     * @memberof RecipeDetail
     */
    version?: number;
    /**
     * whether or not it is the most recent version
     * @type {boolean}
     * @memberof RecipeDetail
     */
    is_latest_version?: boolean;
}

export function RecipeDetailFromJSON(json: any): RecipeDetail {
    return RecipeDetailFromJSONTyped(json, false);
}

export function RecipeDetailFromJSONTyped(json: any, ignoreDiscriminator: boolean): RecipeDetail {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'sections': ((json['sections'] as Array<any>).map(RecipeSectionFromJSON)),
        'name': json['name'],
        'sources': !exists(json, 'sources') ? undefined : ((json['sources'] as Array<any>).map(RecipeSourceFromJSON)),
        'servings': !exists(json, 'servings') ? undefined : json['servings'],
        'quantity': json['quantity'],
        'unit': json['unit'],
        'version': !exists(json, 'version') ? undefined : json['version'],
        'is_latest_version': !exists(json, 'is_latest_version') ? undefined : json['is_latest_version'],
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
        
        'id': value.id,
        'sections': ((value.sections as Array<any>).map(RecipeSectionToJSON)),
        'name': value.name,
        'sources': value.sources === undefined ? undefined : ((value.sources as Array<any>).map(RecipeSourceToJSON)),
        'servings': value.servings,
        'quantity': value.quantity,
        'unit': value.unit,
        'version': value.version,
        'is_latest_version': value.is_latest_version,
    };
}


