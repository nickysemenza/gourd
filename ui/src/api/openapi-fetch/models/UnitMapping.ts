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

/**
 * mappings
 * @export
 * @interface UnitMapping
 */
export interface UnitMapping {
    /**
     * 
     * @type {Amount}
     * @memberof UnitMapping
     */
    a: Amount;
    /**
     * 
     * @type {Amount}
     * @memberof UnitMapping
     */
    b: Amount;
    /**
     * source of the mapping
     * @type {string}
     * @memberof UnitMapping
     */
    source?: string;
}

/**
 * Check if a given object implements the UnitMapping interface.
 */
export function instanceOfUnitMapping(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "a" in value;
    isInstance = isInstance && "b" in value;

    return isInstance;
}

export function UnitMappingFromJSON(json: any): UnitMapping {
    return UnitMappingFromJSONTyped(json, false);
}

export function UnitMappingFromJSONTyped(json: any, ignoreDiscriminator: boolean): UnitMapping {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'a': AmountFromJSON(json['a']),
        'b': AmountFromJSON(json['b']),
        'source': !exists(json, 'source') ? undefined : json['source'],
    };
}

export function UnitMappingToJSON(value?: UnitMapping | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'a': AmountToJSON(value.a),
        'b': AmountToJSON(value.b),
        'source': value.source,
    };
}

