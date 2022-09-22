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
import type { GooglePhotosAlbum } from './GooglePhotosAlbum';
import {
    GooglePhotosAlbumFromJSON,
    GooglePhotosAlbumFromJSONTyped,
    GooglePhotosAlbumToJSON,
} from './GooglePhotosAlbum';

/**
 * 
 * @export
 * @interface ListAllAlbums200Response
 */
export interface ListAllAlbums200Response {
    /**
     * The list of albums
     * @type {Array<GooglePhotosAlbum>}
     * @memberof ListAllAlbums200Response
     */
    albums?: Array<GooglePhotosAlbum>;
}

/**
 * Check if a given object implements the ListAllAlbums200Response interface.
 */
export function instanceOfListAllAlbums200Response(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function ListAllAlbums200ResponseFromJSON(json: any): ListAllAlbums200Response {
    return ListAllAlbums200ResponseFromJSONTyped(json, false);
}

export function ListAllAlbums200ResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): ListAllAlbums200Response {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'albums': !exists(json, 'albums') ? undefined : ((json['albums'] as Array<any>).map(GooglePhotosAlbumFromJSON)),
    };
}

export function ListAllAlbums200ResponseToJSON(value?: ListAllAlbums200Response | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'albums': value.albums === undefined ? undefined : ((value.albums as Array<any>).map(GooglePhotosAlbumToJSON)),
    };
}
