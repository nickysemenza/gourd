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
 * A generic list (for pagination use)
 * @export
 * @interface List
 */
export interface List {
    /**
     * What number page this is
     * @type {number}
     * @memberof List
     */
    page_number: number;
    /**
     * How many items were requested for this page
     * @type {number}
     * @memberof List
     */
    limit: number;
    /**
     * todo
     * @type {number}
     * @memberof List
     */
    offset: number;
    /**
     * Total number of items across all pages
     * @type {number}
     * @memberof List
     */
    total_count: number;
    /**
     * Total number of pages available
     * @type {number}
     * @memberof List
     */
    page_count: number;
}

export function ListFromJSON(json: any): List {
    return ListFromJSONTyped(json, false);
}

export function ListFromJSONTyped(json: any, ignoreDiscriminator: boolean): List {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'page_number': json['page_number'],
        'limit': json['limit'],
        'offset': json['offset'],
        'total_count': json['total_count'],
        'page_count': json['page_count'],
    };
}

export function ListToJSON(value?: List | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'page_number': value.page_number,
        'limit': value.limit,
        'offset': value.offset,
        'total_count': value.total_count,
        'page_count': value.page_count,
    };
}


