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
/**
 * A google photo
 * @export
 * @interface GooglePhoto
 */
export interface GooglePhoto {
    /**
     * id
     * @type {string}
     * @memberof GooglePhoto
     */
    id: string;
    /**
     * public image
     * @type {string}
     * @memberof GooglePhoto
     */
    baseUrl: string;
    /**
     * when it was taken
     * @type {Date}
     * @memberof GooglePhoto
     */
    created: Date;
}

export function GooglePhotoFromJSON(json: any): GooglePhoto {
    return GooglePhotoFromJSONTyped(json, false);
}

export function GooglePhotoFromJSONTyped(json: any, ignoreDiscriminator: boolean): GooglePhoto {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'baseUrl': json['base_url'],
        'created': (new Date(json['created'])),
    };
}

export function GooglePhotoToJSON(value?: GooglePhoto | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'base_url': value.baseUrl,
        'created': (value.created.toISOString()),
    };
}


