/* tslint:disable */
/* eslint-disable */
/**
* @param {string} input
* @returns {string}
*/
export function parse(input: string): string;
/**
* @param {string} input
* @returns {Ingredient}
*/
export function parse2(input: string): Ingredient;
/**
* @param {string} input
* @returns {any}
*/
export function parse3(input: string): any;
/**
* @param {string} input
* @returns {any}
*/
export function parse4(input: string): any;
/**
* @param {Ingredient} val
* @returns {string}
*/
export function format_ingredient(val: Ingredient): string;
/**
* @param {any} recipe_detail
* @returns {any}
*/
export function sum_ingr(recipe_detail: any): any;
/**
* @param {any} conversion_request
* @returns {Amount}
*/
export function dolla(conversion_request: any): Amount;
/**
* @param {string} input
* @returns {Amount[]}
*/
export function parse_amount(input: string): Amount[];
/**
* @param {any} recipe_detail
* @returns {string}
*/
export function encode_recipe_text(recipe_detail: any): string;
/**
* @param {any} recipe_detail
* @returns {CompactR[][]}
*/
export function encode_recipe_to_compact_json(recipe_detail: any): CompactR[][];
/**
* @param {string} r
* @returns {any}
*/
export function decode_recipe_text(r: string): any;
/**
* @param {any} conversion_request
* @returns {string}
*/
export function make_dag(conversion_request: any): string;

interface Ingredient {
    amounts: Amount[];
    modifier?: string;
    name: string;
  }
  
  interface Amount {
    unit: string;
    value: number;
  }

  interface CompactR {
    Ing?: Ingredient;
    Ins?: string;
  }


