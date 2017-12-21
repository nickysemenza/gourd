import apiFetch from './index';

export const REQUEST_RECIPES = 'REQUEST_RECIPES';
export const RECEIVE_RECIPES = 'RECEIVE_RECIPES';

export function fetchRecipes() {
  return dispatch => {
    dispatch(requestRecipes());
    return apiFetch('recipes')
      .then(response => response.json())
      .then(json => dispatch(receiveRecipes(json)));
  };
}

function requestRecipes() {
  return {
    type: REQUEST_RECIPES
  };
}

function receiveRecipes(json) {
  return {
    type: RECEIVE_RECIPES,
    json,
    receivedAt: Date.now()
  };
}
export const REQUEST_RECIPE_DETAIL = 'REQUEST_RECIPE_DETAIL';
export const RECEIVE_RECIPE_DETAIL = 'RECEIVE_RECIPE_DETAIL';

export function fetchRecipeDetail(recipeSlug) {
  return dispatch => {
    dispatch(requestRecipeDetail(recipeSlug));
    return apiFetch('recipes/' + recipeSlug)
      .then(response => response.json(recipeSlug))
      .then(json => dispatch(receiveRecipeDetail(json, recipeSlug)));
  };
}

function requestRecipeDetail(recipeSlug) {
  return {
    type: REQUEST_RECIPE_DETAIL,
    recipeSlug
  };
}

function receiveRecipeDetail(json, recipeSlug) {
  return {
    type: RECEIVE_RECIPE_DETAIL,
    recipeSlug,
    json,
    receivedAt: Date.now()
  };
}

export const EDIT_TOP_LEVEL_ITEM = 'EDIT_TOP_LEVEL_ITEM';
export function editTopLevelItem(slug, fieldName, value) {
  return {
    type: EDIT_TOP_LEVEL_ITEM,
    slug,
    fieldName,
    value
  };
}
export const DELETE_SECTION = 'DELETE_SECTION';
export function deleteSectionByIndex(slug, index) {
  return {
    type: DELETE_SECTION,
    slug,
    index
  };
}
export const ADD_SECTION = 'ADD_SECTION';
export function addSection(slug, index) {
  return {
    type: ADD_SECTION,
    slug,
    index
  };
}
export const DELETE_INSTRUCTION = 'DELETE_INSTRUCTION';
export function deleteInstruction(slug, sectionNum, instructionNum) {
  return {
    type: DELETE_INSTRUCTION,
    slug,
    sectionNum,
    instructionNum
  };
}
export const ADD_INSTRUCTION = 'ADD_INSTRUCTION';
export function addInstruction(slug, sectionNum, instructionNum) {
  return {
    type: ADD_INSTRUCTION,
    slug,
    sectionNum,
    instructionNum
  };
}
export const EDIT_INSTRUCTION = 'EDIT_INSTRUCTION';
export function editInstruction(slug, sectionNum, instructionNum, value) {
  return {
    type: EDIT_INSTRUCTION,
    slug,
    sectionNum,
    instructionNum,
    value
  };
}
export const DELETE_INGREDIENT = 'DELETE_INGREDIENT';
export function deleteIngredient(slug, sectionNum, ingredientNum) {
  return {
    type: DELETE_INGREDIENT,
    slug,
    sectionNum,
    ingredientNum
  };
}
export const ADD_INGREDIENT = 'ADD_INGREDIENT';
export function addIngredient(slug, sectionNum, ingredientNum) {
  return {
    type: ADD_INGREDIENT,
    slug,
    sectionNum,
    ingredientNum
  };
}
export const EDIT_INGREDIENT = 'EDIT_INGREDIENT';
export function editIngredient(slug, sectionNum, ingredientNum, field, value) {
  return {
    type: EDIT_INGREDIENT,
    slug,
    sectionNum,
    ingredientNum,
    field,
    value
  };
}
export function addRecipeNote(slug, note) {
  return dispatch => {
    return apiFetch(`recipes/${slug}/notes`, {
      method: 'POST',
      body: JSON.stringify({ note })
    })
      .then(response => response.json())
      .then(() => {
        dispatch(fetchRecipeDetail(slug));
      });
  };
}
