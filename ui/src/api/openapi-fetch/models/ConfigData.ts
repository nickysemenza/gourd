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
 * config data
 * @export
 * @interface ConfigData
 */
export interface ConfigData {
    /**
     * 
     * @type {string}
     * @memberof ConfigData
     */
    google_scopes: string;
    /**
     * 
     * @type {string}
     * @memberof ConfigData
     */
    google_client_id: string;
}

export function ConfigDataFromJSON(json: any): ConfigData {
    return ConfigDataFromJSONTyped(json, false);
}

export function ConfigDataFromJSONTyped(json: any, ignoreDiscriminator: boolean): ConfigData {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'google_scopes': json['google_scopes'],
        'google_client_id': json['google_client_id'],
    };
}

export function ConfigDataToJSON(value?: ConfigData | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'google_scopes': value.google_scopes,
        'google_client_id': value.google_client_id,
    };
}

