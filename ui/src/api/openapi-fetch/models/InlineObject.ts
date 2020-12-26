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

import { exists, mapValues } from "../runtime";
/**
 *
 * @export
 * @interface InlineObject
 */
export interface InlineObject {
  /**
   *
   * @type {Array<string>}
   * @memberof InlineObject
   */
  ingredientIds: Array<string>;
}

export function InlineObjectFromJSON(json: any): InlineObject {
  return InlineObjectFromJSONTyped(json, false);
}

export function InlineObjectFromJSONTyped(
  json: any,
  ignoreDiscriminator: boolean
): InlineObject {
  if (json === undefined || json === null) {
    return json;
  }
  return {
    ingredientIds: json["ingredient_ids"],
  };
}

export function InlineObjectToJSON(value?: InlineObject | null): any {
  if (value === undefined) {
    return undefined;
  }
  if (value === null) {
    return null;
  }
  return {
    ingredient_ids: value.ingredientIds,
  };
}
