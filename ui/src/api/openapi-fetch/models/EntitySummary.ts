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
    IngredientKind,
    IngredientKindFromJSON,
    IngredientKindFromJSONTyped,
    IngredientKindToJSON,
} from './';

/**
 * name and id of something
 * @export
 * @interface EntitySummary
 */
export interface EntitySummary {
    /**
     * recipe_detail or ingredient id
     * @type {string}
     * @memberof EntitySummary
     */
    id: string;
    /**
     * recipe or ingredient name
     * @type {string}
     * @memberof EntitySummary
     */
    name: string;
    /**
     * multiplier
     * @type {number}
     * @memberof EntitySummary
     */
    multiplier: number;
    /**
     * 
     * @type {IngredientKind}
     * @memberof EntitySummary
     */
    kind: IngredientKind;
}

export function EntitySummaryFromJSON(json: any): EntitySummary {
    return EntitySummaryFromJSONTyped(json, false);
}

export function EntitySummaryFromJSONTyped(json: any, ignoreDiscriminator: boolean): EntitySummary {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'name': json['name'],
        'multiplier': json['multiplier'],
        'kind': IngredientKindFromJSON(json['kind']),
    };
}

export function EntitySummaryToJSON(value?: EntitySummary | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'name': value.name,
        'multiplier': value.multiplier,
        'kind': IngredientKindToJSON(value.kind),
    };
}

