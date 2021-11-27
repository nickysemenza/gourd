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
    Amount,
    AmountFromJSON,
    AmountFromJSONTyped,
    AmountToJSON,
    EntitySummary,
    EntitySummaryFromJSON,
    EntitySummaryFromJSONTyped,
    EntitySummaryToJSON,
    IngredientUsage,
    IngredientUsageFromJSON,
    IngredientUsageFromJSONTyped,
    IngredientUsageToJSON,
} from './';

/**
 * 
 * @export
 * @interface UsageValue
 */
export interface UsageValue {
    /**
     * multiplier
     * @type {Array<IngredientUsage>}
     * @memberof UsageValue
     */
    ings: Array<IngredientUsage>;
    /**
     * amounts
     * @type {Array<Amount>}
     * @memberof UsageValue
     */
    sum: Array<Amount>;
    /**
     * 
     * @type {EntitySummary}
     * @memberof UsageValue
     */
    ing: EntitySummary;
}

export function UsageValueFromJSON(json: any): UsageValue {
    return UsageValueFromJSONTyped(json, false);
}

export function UsageValueFromJSONTyped(json: any, ignoreDiscriminator: boolean): UsageValue {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'ings': ((json['ings'] as Array<any>).map(IngredientUsageFromJSON)),
        'sum': ((json['sum'] as Array<any>).map(AmountFromJSON)),
        'ing': EntitySummaryFromJSON(json['ing']),
    };
}

export function UsageValueToJSON(value?: UsageValue | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'ings': ((value.ings as Array<any>).map(IngredientUsageToJSON)),
        'sum': ((value.sum as Array<any>).map(AmountToJSON)),
        'ing': EntitySummaryToJSON(value.ing),
    };
}


