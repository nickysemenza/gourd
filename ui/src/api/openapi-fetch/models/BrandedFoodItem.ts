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
import type { BrandedFoodItemLabelNutrients } from './BrandedFoodItemLabelNutrients';
import {
    BrandedFoodItemLabelNutrientsFromJSON,
    BrandedFoodItemLabelNutrientsFromJSONTyped,
    BrandedFoodItemLabelNutrientsToJSON,
} from './BrandedFoodItemLabelNutrients';
import type { FoodNutrient } from './FoodNutrient';
import {
    FoodNutrientFromJSON,
    FoodNutrientFromJSONTyped,
    FoodNutrientToJSON,
} from './FoodNutrient';
import type { FoodUpdateLog } from './FoodUpdateLog';
import {
    FoodUpdateLogFromJSON,
    FoodUpdateLogFromJSONTyped,
    FoodUpdateLogToJSON,
} from './FoodUpdateLog';

/**
 * 
 * @export
 * @interface BrandedFoodItem
 */
export interface BrandedFoodItem {
    /**
     * 
     * @type {number}
     * @memberof BrandedFoodItem
     */
    fdc_id: number;
    /**
     * 
     * @type {string}
     * @memberof BrandedFoodItem
     */
    available_date?: string;
    /**
     * 
     * @type {string}
     * @memberof BrandedFoodItem
     */
    brand_owner?: string;
    /**
     * 
     * @type {string}
     * @memberof BrandedFoodItem
     */
    data_source?: string;
    /**
     * 
     * @type {string}
     * @memberof BrandedFoodItem
     */
    data_type: string;
    /**
     * 
     * @type {string}
     * @memberof BrandedFoodItem
     */
    description: string;
    /**
     * 
     * @type {string}
     * @memberof BrandedFoodItem
     */
    food_class?: string;
    /**
     * 
     * @type {string}
     * @memberof BrandedFoodItem
     */
    gtin_upc?: string;
    /**
     * 
     * @type {string}
     * @memberof BrandedFoodItem
     */
    household_serving_full_text?: string;
    /**
     * 
     * @type {string}
     * @memberof BrandedFoodItem
     */
    ingredients?: string;
    /**
     * 
     * @type {string}
     * @memberof BrandedFoodItem
     */
    modified_date?: string;
    /**
     * 
     * @type {string}
     * @memberof BrandedFoodItem
     */
    publication_date?: string;
    /**
     * 
     * @type {number}
     * @memberof BrandedFoodItem
     */
    serving_size?: number;
    /**
     * 
     * @type {string}
     * @memberof BrandedFoodItem
     */
    serving_size_unit?: string;
    /**
     * 
     * @type {string}
     * @memberof BrandedFoodItem
     */
    preparation_state_code?: string;
    /**
     * 
     * @type {string}
     * @memberof BrandedFoodItem
     */
    branded_food_category?: string;
    /**
     * 
     * @type {Array<string>}
     * @memberof BrandedFoodItem
     */
    trade_channel?: Array<string>;
    /**
     * 
     * @type {number}
     * @memberof BrandedFoodItem
     */
    gpc_class_code?: number;
    /**
     * 
     * @type {Array<FoodNutrient>}
     * @memberof BrandedFoodItem
     */
    food_nutrients?: Array<FoodNutrient>;
    /**
     * 
     * @type {Array<FoodUpdateLog>}
     * @memberof BrandedFoodItem
     */
    food_update_log?: Array<FoodUpdateLog>;
    /**
     * 
     * @type {BrandedFoodItemLabelNutrients}
     * @memberof BrandedFoodItem
     */
    label_nutrients?: BrandedFoodItemLabelNutrients;
}

/**
 * Check if a given object implements the BrandedFoodItem interface.
 */
export function instanceOfBrandedFoodItem(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "fdc_id" in value;
    isInstance = isInstance && "data_type" in value;
    isInstance = isInstance && "description" in value;

    return isInstance;
}

export function BrandedFoodItemFromJSON(json: any): BrandedFoodItem {
    return BrandedFoodItemFromJSONTyped(json, false);
}

export function BrandedFoodItemFromJSONTyped(json: any, ignoreDiscriminator: boolean): BrandedFoodItem {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'fdc_id': json['fdcId'],
        'available_date': !exists(json, 'availableDate') ? undefined : json['availableDate'],
        'brand_owner': !exists(json, 'brandOwner') ? undefined : json['brandOwner'],
        'data_source': !exists(json, 'dataSource') ? undefined : json['dataSource'],
        'data_type': json['dataType'],
        'description': json['description'],
        'food_class': !exists(json, 'foodClass') ? undefined : json['foodClass'],
        'gtin_upc': !exists(json, 'gtinUpc') ? undefined : json['gtinUpc'],
        'household_serving_full_text': !exists(json, 'householdServingFullText') ? undefined : json['householdServingFullText'],
        'ingredients': !exists(json, 'ingredients') ? undefined : json['ingredients'],
        'modified_date': !exists(json, 'modifiedDate') ? undefined : json['modifiedDate'],
        'publication_date': !exists(json, 'publicationDate') ? undefined : json['publicationDate'],
        'serving_size': !exists(json, 'servingSize') ? undefined : json['servingSize'],
        'serving_size_unit': !exists(json, 'servingSizeUnit') ? undefined : json['servingSizeUnit'],
        'preparation_state_code': !exists(json, 'preparationStateCode') ? undefined : json['preparationStateCode'],
        'branded_food_category': !exists(json, 'brandedFoodCategory') ? undefined : json['brandedFoodCategory'],
        'trade_channel': !exists(json, 'tradeChannel') ? undefined : json['tradeChannel'],
        'gpc_class_code': !exists(json, 'gpcClassCode') ? undefined : json['gpcClassCode'],
        'food_nutrients': !exists(json, 'foodNutrients') ? undefined : ((json['foodNutrients'] as Array<any>).map(FoodNutrientFromJSON)),
        'food_update_log': !exists(json, 'foodUpdateLog') ? undefined : ((json['foodUpdateLog'] as Array<any>).map(FoodUpdateLogFromJSON)),
        'label_nutrients': !exists(json, 'labelNutrients') ? undefined : BrandedFoodItemLabelNutrientsFromJSON(json['labelNutrients']),
    };
}

export function BrandedFoodItemToJSON(value?: BrandedFoodItem | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'fdcId': value.fdc_id,
        'availableDate': value.available_date,
        'brandOwner': value.brand_owner,
        'dataSource': value.data_source,
        'dataType': value.data_type,
        'description': value.description,
        'foodClass': value.food_class,
        'gtinUpc': value.gtin_upc,
        'householdServingFullText': value.household_serving_full_text,
        'ingredients': value.ingredients,
        'modifiedDate': value.modified_date,
        'publicationDate': value.publication_date,
        'servingSize': value.serving_size,
        'servingSizeUnit': value.serving_size_unit,
        'preparationStateCode': value.preparation_state_code,
        'brandedFoodCategory': value.branded_food_category,
        'tradeChannel': value.trade_channel,
        'gpcClassCode': value.gpc_class_code,
        'foodNutrients': value.food_nutrients === undefined ? undefined : ((value.food_nutrients as Array<any>).map(FoodNutrientToJSON)),
        'foodUpdateLog': value.food_update_log === undefined ? undefined : ((value.food_update_log as Array<any>).map(FoodUpdateLogToJSON)),
        'labelNutrients': BrandedFoodItemLabelNutrientsToJSON(value.label_nutrients),
    };
}

