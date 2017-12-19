
import {
    REQUEST_RECIPES,
    RECEIVE_RECIPES,
    REQUEST_RECIPE_DETAIL,
    RECEIVE_RECIPE_DETAIL,
} from '../actions/recipe';


const INITIAL_STATE = {
    recipe_list: [],
    recipe_list_loading: false,
    
    recipe_detail: {},
    recipe_detail_loading: false,
};

export default function (state = INITIAL_STATE, action) {
    switch (action.type) {
        case REQUEST_RECIPES:
            return { ...state, recipe_list_loading: true };
        case RECEIVE_RECIPES:
            //todo: error checking
            return { ...state,
                recipe_list_loading: false,
                recipe_list: action.json
            };
        case REQUEST_RECIPE_DETAIL:
            return { ...state, recipe_detail_loading: true };
        case RECEIVE_RECIPE_DETAIL:
            //todo: error checking
            return { ...state,
                recipe_detail_loading: false,
                recipe_detail: {
                    ...state.recipe_detail,
                    [action.recipeSlug]: action.json
                }
            };
        default:
            return state;
    }
}