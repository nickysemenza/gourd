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
    base_url: string;
    /**
     * blur hash
     * @type {string}
     * @memberof GooglePhoto
     */
    blur_hash?: string;
    /**
     * when it was taken
     * @type {Date}
     * @memberof GooglePhoto
     */
    created: Date;
    /**
     * width px
     * @type {number}
     * @memberof GooglePhoto
     */
    width: number;
    /**
     * height px
     * @type {number}
     * @memberof GooglePhoto
     */
    height: number;
    /**
     * where the photo came from
     * @type {string}
     * @memberof GooglePhoto
     */
    source: GooglePhotoSourceEnum;
}

/**
* @export
* @enum {string}
*/
export enum GooglePhotoSourceEnum {
    GOOGLE = 'google',
    NOTION = 'notion'
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
        'base_url': json['base_url'],
        'blur_hash': !exists(json, 'blur_hash') ? undefined : json['blur_hash'],
        'created': (new Date(json['created'])),
        'width': json['width'],
        'height': json['height'],
        'source': json['source'],
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
        'base_url': value.base_url,
        'blur_hash': value.blur_hash,
        'created': (value.created.toISOString()),
        'width': value.width,
        'height': value.height,
        'source': value.source,
    };
}


