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
    RecipeDependency,
    RecipeDependencyFromJSON,
    RecipeDependencyFromJSONTyped,
    RecipeDependencyToJSON,
} from './';

/**
 * 
 * @export
 * @interface InlineResponse2002
 */
export interface InlineResponse2002 {
    /**
     * all
     * @type {Array<RecipeDependency>}
     * @memberof InlineResponse2002
     */
    items?: Array<RecipeDependency>;
}

export function InlineResponse2002FromJSON(json: any): InlineResponse2002 {
    return InlineResponse2002FromJSONTyped(json, false);
}

export function InlineResponse2002FromJSONTyped(json: any, ignoreDiscriminator: boolean): InlineResponse2002 {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'items': !exists(json, 'items') ? undefined : ((json['items'] as Array<any>).map(RecipeDependencyFromJSON)),
    };
}

export function InlineResponse2002ToJSON(value?: InlineResponse2002 | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'items': value.items === undefined ? undefined : ((value.items as Array<any>).map(RecipeDependencyToJSON)),
    };
}

