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
 * Instructions in a single section
 * @export
 * @interface SectionInstructionInput
 */
export interface SectionInstructionInput {
    /**
     * instruction
     * @type {string}
     * @memberof SectionInstructionInput
     */
    instruction: string;
}

export function SectionInstructionInputFromJSON(json: any): SectionInstructionInput {
    return SectionInstructionInputFromJSONTyped(json, false);
}

export function SectionInstructionInputFromJSONTyped(json: any, ignoreDiscriminator: boolean): SectionInstructionInput {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'instruction': json['instruction'],
    };
}

export function SectionInstructionInputToJSON(value?: SectionInstructionInput | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'instruction': value.instruction,
    };
}

