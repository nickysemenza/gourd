import apiFetch  from './index';

export const REQUEST_RECIPES = 'REQUEST_RECIPES';
export const RECEIVE_RECIPES = 'RECEIVE_RECIPES';

export function fetchRecipes () {
    return (dispatch) => {
        dispatch(requestRecipes());
        return apiFetch('recipes')
            .then((response) => response.json())
            .then((json) => dispatch(receiveRecipes(json)));
    };
}

function requestRecipes () {
    return {
        type: REQUEST_RECIPES
    };
}

function receiveRecipes (json) {
    return {
        type: RECEIVE_RECIPES,
        json,
        receivedAt: Date.now()
    };
}
export const REQUEST_RECIPE_DETAIL = 'REQUEST_RECIPE_DETAIL';
export const RECEIVE_RECIPE_DETAIL = 'RECEIVE_RECIPE_DETAIL';

export function fetchRecipeDetail (recipeSlug) {
    return (dispatch) => {
        dispatch(requestRecipeDetail(recipeSlug));
        return apiFetch('recipes/'+recipeSlug)
            .then((response) => response.json(recipeSlug))
            .then((json) => dispatch(receiveRecipeDetail(json,recipeSlug)));
    };
}

function requestRecipeDetail (recipeSlug) {
    return {
        type: REQUEST_RECIPE_DETAIL,
        recipeSlug
    };
}

function receiveRecipeDetail (json, recipeSlug) {
    return {
        type: RECEIVE_RECIPE_DETAIL,
        recipeSlug,
        json,
        receivedAt: Date.now()
    };
}

