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
import type { Meal } from './Meal';
import {
    MealFromJSON,
    MealFromJSONTyped,
    MealToJSON,
} from './Meal';
import type { Photo } from './Photo';
import {
    PhotoFromJSON,
    PhotoFromJSONTyped,
    PhotoToJSON,
} from './Photo';
import type { RecipeDetail } from './RecipeDetail';
import {
    RecipeDetailFromJSON,
    RecipeDetailFromJSONTyped,
    RecipeDetailToJSON,
} from './RecipeDetail';

/**
 * A recipe with subcomponents, including some "generated" fields to enhance data
 * @export
 * @interface RecipeWrapper
 */
export interface RecipeWrapper {
    /**
     * id
     * @type {string}
     * @memberof RecipeWrapper
     */
    id: string;
    /**
     * 
     * @type {RecipeDetail}
     * @memberof RecipeWrapper
     */
    detail: RecipeDetail;
    /**
     * 
     * @type {Array<Meal>}
     * @memberof RecipeWrapper
     */
    linked_meals?: Array<Meal>;
    /**
     * 
     * @type {Array<Photo>}
     * @memberof RecipeWrapper
     */
    linked_photos?: Array<Photo>;
    /**
     * Other versions
     * @type {Array<RecipeDetail>}
     * @memberof RecipeWrapper
     */
    other_versions?: Array<RecipeDetail>;
}

/**
 * Check if a given object implements the RecipeWrapper interface.
 */
export function instanceOfRecipeWrapper(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "detail" in value;

    return isInstance;
}

export function RecipeWrapperFromJSON(json: any): RecipeWrapper {
    return RecipeWrapperFromJSONTyped(json, false);
}

export function RecipeWrapperFromJSONTyped(json: any, ignoreDiscriminator: boolean): RecipeWrapper {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'detail': RecipeDetailFromJSON(json['detail']),
        'linked_meals': !exists(json, 'linked_meals') ? undefined : ((json['linked_meals'] as Array<any>).map(MealFromJSON)),
        'linked_photos': !exists(json, 'linked_photos') ? undefined : ((json['linked_photos'] as Array<any>).map(PhotoFromJSON)),
        'other_versions': !exists(json, 'other_versions') ? undefined : ((json['other_versions'] as Array<any>).map(RecipeDetailFromJSON)),
    };
}

export function RecipeWrapperToJSON(value?: RecipeWrapper | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'detail': RecipeDetailToJSON(value.detail),
        'linked_meals': value.linked_meals === undefined ? undefined : ((value.linked_meals as Array<any>).map(MealToJSON)),
        'linked_photos': value.linked_photos === undefined ? undefined : ((value.linked_photos as Array<any>).map(PhotoToJSON)),
        'other_versions': value.other_versions === undefined ? undefined : ((value.other_versions as Array<any>).map(RecipeDetailToJSON)),
    };
}

