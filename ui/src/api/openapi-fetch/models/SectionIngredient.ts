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
    Ingredient,
    IngredientFromJSON,
    IngredientFromJSONTyped,
    IngredientToJSON,
    IngredientKind,
    IngredientKindFromJSON,
    IngredientKindFromJSONTyped,
    IngredientKindToJSON,
    RecipeDetail,
    RecipeDetailFromJSON,
    RecipeDetailFromJSONTyped,
    RecipeDetailToJSON,
} from './';

/**
 * Ingredients in a single section
 * @export
 * @interface SectionIngredient
 */
export interface SectionIngredient {
    /**
     * id
     * @type {string}
     * @memberof SectionIngredient
     */
    id: string;
    /**
     * 
     * @type {IngredientKind}
     * @memberof SectionIngredient
     */
    kind: IngredientKind;
    /**
     * 
     * @type {RecipeDetail}
     * @memberof SectionIngredient
     */
    recipe?: RecipeDetail;
    /**
     * 
     * @type {Ingredient}
     * @memberof SectionIngredient
     */
    ingredient?: Ingredient;
    /**
     * weight in grams
     * @type {number}
     * @memberof SectionIngredient
     */
    grams: number;
    /**
     * amount
     * @type {number}
     * @memberof SectionIngredient
     */
    amount?: number;
    /**
     * unit
     * @type {string}
     * @memberof SectionIngredient
     */
    unit?: string;
    /**
     * adjective
     * @type {string}
     * @memberof SectionIngredient
     */
    adjective?: string;
    /**
     * optional
     * @type {boolean}
     * @memberof SectionIngredient
     */
    optional?: boolean;
    /**
     * raw line item (pre-import/scrape)
     * @type {string}
     * @memberof SectionIngredient
     */
    original?: string;
    /**
     * x
     * @type {Array<SectionIngredient>}
     * @memberof SectionIngredient
     */
    substitutes?: Array<SectionIngredient>;
}

export function SectionIngredientFromJSON(json: any): SectionIngredient {
    return SectionIngredientFromJSONTyped(json, false);
}

export function SectionIngredientFromJSONTyped(json: any, ignoreDiscriminator: boolean): SectionIngredient {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'kind': IngredientKindFromJSON(json['kind']),
        'recipe': !exists(json, 'recipe') ? undefined : RecipeDetailFromJSON(json['recipe']),
        'ingredient': !exists(json, 'ingredient') ? undefined : IngredientFromJSON(json['ingredient']),
        'grams': json['grams'],
        'amount': !exists(json, 'amount') ? undefined : json['amount'],
        'unit': !exists(json, 'unit') ? undefined : json['unit'],
        'adjective': !exists(json, 'adjective') ? undefined : json['adjective'],
        'optional': !exists(json, 'optional') ? undefined : json['optional'],
        'original': !exists(json, 'original') ? undefined : json['original'],
        'substitutes': !exists(json, 'substitutes') ? undefined : ((json['substitutes'] as Array<any>).map(SectionIngredientFromJSON)),
    };
}

export function SectionIngredientToJSON(value?: SectionIngredient | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'kind': IngredientKindToJSON(value.kind),
        'recipe': RecipeDetailToJSON(value.recipe),
        'ingredient': IngredientToJSON(value.ingredient),
        'grams': value.grams,
        'amount': value.amount,
        'unit': value.unit,
        'adjective': value.adjective,
        'optional': value.optional,
        'original': value.original,
        'substitutes': value.substitutes === undefined ? undefined : ((value.substitutes as Array<any>).map(SectionIngredientToJSON)),
    };
}


