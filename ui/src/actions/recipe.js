import apiFetch from './index';
import { toastr } from 'react-redux-toastr';

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
export function deleteInstruction(slug, sectionIndex, instructionIndex) {
  return {
    type: DELETE_INSTRUCTION,
    slug,
    sectionIndex,
    instructionIndex
  };
}
export const ADD_INSTRUCTION = 'ADD_INSTRUCTION';
export function addInstruction(slug, sectionIndex, instructionIndex) {
  return {
    type: ADD_INSTRUCTION,
    slug,
    sectionIndex,
    instructionIndex
  };
}
export const EDIT_INSTRUCTION = 'EDIT_INSTRUCTION';
export function editInstruction(slug, sectionIndex, instructionIndex, value) {
  return {
    type: EDIT_INSTRUCTION,
    slug,
    sectionIndex,
    instructionIndex,
    value
  };
}
export const MOVE_INSTRUCTION = 'MOVE_INSTRUCTION';
export function moveInstruction(
  slug,
  sectionIndex,
  instructionIndex,
  targetIndex
) {
  return {
    type: MOVE_INSTRUCTION,
    slug,
    sectionIndex,
    instructionIndex,
    targetIndex
  };
}
export const DELETE_INGREDIENT = 'DELETE_INGREDIENT';
export function deleteIngredient(slug, sectionIndex, ingredientIndex) {
  return {
    type: DELETE_INGREDIENT,
    slug,
    sectionIndex,
    ingredientIndex
  };
}
export const ADD_INGREDIENT = 'ADD_INGREDIENT';
export function addIngredient(slug, sectionIndex, ingredientIndex) {
  return {
    type: ADD_INGREDIENT,
    slug,
    sectionIndex,
    ingredientIndex
  };
}
export const EDIT_INGREDIENT = 'EDIT_INGREDIENT';
export function editIngredient(
  slug,
  sectionIndex,
  ingredientIndex,
  field,
  value
) {
  return {
    type: EDIT_INGREDIENT,
    slug,
    sectionIndex,
    ingredientIndex,
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

export const RECEIVE_IMAGES = 'RECEIVE_IMAGES';

export function fetchImages() {
  return dispatch => {
    return apiFetch('images')
      .then(response => response.json())
      .then(json => dispatch(receiveImages(json)));
  };
}

function receiveImages(json) {
  return {
    type: RECEIVE_IMAGES,
    json
  };
}

export const RECEIVE_MEAL_LIST = 'RECEIVE_MEAL_LIST';

export function fetchMealList() {
  return dispatch => {
    return apiFetch('meals')
      .then(response => response.json())
      .then(json => dispatch(receiveMealList(json)));
  };
}

function receiveMealList(json) {
  return {
    type: RECEIVE_MEAL_LIST,
    json
  };
}

export const RECEIVE_MEAL_DETAIL = 'RECEIVE_MEAL_DETAIL';

export function fetchMealDetail(meal_id) {
  return dispatch => {
    return apiFetch('meals/' + meal_id)
      .then(response => response.json())
      .then(json => dispatch(receiveMealDetail(meal_id, json)));
  };
}

function receiveMealDetail(meal_id, json) {
  return {
    type: RECEIVE_MEAL_DETAIL,
    json,
    meal_id
  };
}
export const EDIT_MEAL_RECIPE_MULTIPLIER = 'EDIT_MEAL_RECIPE_MULTIPLIER';
export function editMealMultiplier(meal_id, recipeIndex, event) {
  return {
    type: EDIT_MEAL_RECIPE_MULTIPLIER,
    meal_id,
    recipeIndex,
    value: parseFloat(event.target.value)
  };
}

export const EDIT_MEAL_RECIPE = 'EDIT_MEAL_RECIPE';
export function editMealRecipe(meal_id, recipeIndex, value) {
  return {
    type: EDIT_MEAL_RECIPE,
    meal_id,
    recipeIndex,
    value
  };
}

export const ADD_MEAL_RECIPE = 'ADD_MEAL_RECIPE';
export function addMealRecipe(meal_id) {
  return {
    type: ADD_MEAL_RECIPE,
    meal_id
  };
}

export const DELETE_MEAL_RECIPE = 'DELETE_MEAL_RECIPE';
export function deleteMealRecipe(meal_id, index) {
  return {
    type: DELETE_MEAL_RECIPE,
    meal_id,
    index
  };
}
export const EDIT_MEAL_FIELD = 'EDIT_MEAL_FIELD';

export function editMealField(meal_id, field, value) {
  return {
    type: EDIT_MEAL_FIELD,
    meal_id,
    field,
    value
  };
}
export function saveMeal(meal_id) {
  return (dispatch, getState) => {
    let meal = getState().recipe.meal_detail[meal_id];
    return apiFetch(`meals/${meal_id}`, {
      method: 'PUT',
      body: JSON.stringify(meal)
    })
      .then(response => response.json())
      .then(json => {
        if (json.error) toastr.error('Oops!', json.error);
        else toastr.success('Success!', ` meal #${meal_id} updated!`);
        dispatch(fetchMealDetail(meal_id));
      });
  };
}

export function createRecipe(slug, title) {
  return dispatch => {
    return apiFetch('recipes', {
      method: 'POST',
      body: JSON.stringify({ slug, title })
    })
      .then(response => response.json())
      .then(json => {
        if (json.error) toastr.error('Oops!', `slug ${slug} already exists!`);
        else toastr.success('Success!', `${title} (${slug}) created!`);
        dispatch(fetchRecipeDetail(slug));
      });
  };
}

export function saveRecipe(slug) {
  return (dispatch, getState) => {
    let recipe = getState().recipe.recipe_detail[slug];
    return apiFetch(`recipes/${slug}`, {
      method: 'PUT',
      body: JSON.stringify(recipe)
    })
      .then(response => response.json())
      .then(json => {
        if (json.error) toastr.error('Oops!', json.error);
        else toastr.success('Success!', `${recipe.title} (${slug}) updated!`);
        dispatch(fetchRecipeDetail(slug));
      });
  };
}

export const RECEIVE_CATEGORIES = 'RECEIVE_CATEGORIES';

export function fetchCategories() {
  return dispatch => {
    return apiFetch('categories')
      .then(response => response.json())
      .then(json => dispatch(receiveCategories(json)));
  };
}

function receiveCategories(json) {
  return {
    type: RECEIVE_CATEGORIES,
    json
  };
}
export const REMOVE_CATEGORY_FROM_RECIPE = 'REMOVE_CATEGORY_FROM_RECIPE';
export function removeCategoryFromRecipe(slug, categoryId) {
  return {
    type: REMOVE_CATEGORY_FROM_RECIPE,
    slug,
    categoryId
  };
}

export const ADD_CATEGORY_TO_RECIPE = 'ADD_CATEGORY_TO_RECIPE';
export function addCategoryToRecipe(slug, categoryId) {
  return {
    type: ADD_CATEGORY_TO_RECIPE,
    slug,
    categoryId
  };
}
