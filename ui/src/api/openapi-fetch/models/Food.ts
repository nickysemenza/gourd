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
    BrandedFood,
    BrandedFoodFromJSON,
    BrandedFoodFromJSONTyped,
    BrandedFoodToJSON,
    FoodCategory,
    FoodCategoryFromJSON,
    FoodCategoryFromJSONTyped,
    FoodCategoryToJSON,
    FoodDataType,
    FoodDataTypeFromJSON,
    FoodDataTypeFromJSONTyped,
    FoodDataTypeToJSON,
    FoodNutrient,
    FoodNutrientFromJSON,
    FoodNutrientFromJSONTyped,
    FoodNutrientToJSON,
    FoodPortion,
    FoodPortionFromJSON,
    FoodPortionFromJSONTyped,
    FoodPortionToJSON,
} from './';

/**
 * A top level food
 * @export
 * @interface Food
 */
export interface Food {
    /**
     * FDC Id
     * @type {number}
     * @memberof Food
     */
    fdc_id: number;
    /**
     * Food description
     * @type {string}
     * @memberof Food
     */
    description: string;
    /**
     * 
     * @type {FoodDataType}
     * @memberof Food
     */
    data_type: FoodDataType;
    /**
     * 
     * @type {FoodCategory}
     * @memberof Food
     */
    category?: FoodCategory;
    /**
     * todo
     * @type {Array<FoodNutrient>}
     * @memberof Food
     */
    nutrients: Array<FoodNutrient>;
    /**
     * portion datapoints
     * @type {Array<FoodPortion>}
     * @memberof Food
     */
    portions?: Array<FoodPortion>;
    /**
     * 
     * @type {BrandedFood}
     * @memberof Food
     */
    branded_info?: BrandedFood;
}

export function FoodFromJSON(json: any): Food {
    return FoodFromJSONTyped(json, false);
}

export function FoodFromJSONTyped(json: any, ignoreDiscriminator: boolean): Food {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'fdc_id': json['fdc_id'],
        'description': json['description'],
        'data_type': FoodDataTypeFromJSON(json['data_type']),
        'category': !exists(json, 'category') ? undefined : FoodCategoryFromJSON(json['category']),
        'nutrients': ((json['nutrients'] as Array<any>).map(FoodNutrientFromJSON)),
        'portions': !exists(json, 'portions') ? undefined : ((json['portions'] as Array<any>).map(FoodPortionFromJSON)),
        'branded_info': !exists(json, 'branded_info') ? undefined : BrandedFoodFromJSON(json['branded_info']),
    };
}

export function FoodToJSON(value?: Food | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'fdc_id': value.fdc_id,
        'description': value.description,
        'data_type': FoodDataTypeToJSON(value.data_type),
        'category': FoodCategoryToJSON(value.category),
        'nutrients': ((value.nutrients as Array<any>).map(FoodNutrientToJSON)),
        'portions': value.portions === undefined ? undefined : ((value.portions as Array<any>).map(FoodPortionToJSON)),
        'branded_info': BrandedFoodToJSON(value.branded_info),
    };
}


