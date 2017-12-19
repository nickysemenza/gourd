import { combineReducers } from 'redux'
import recipe from './reducer_recipe';
const rootReducer = combineReducers({
    recipe,
});

export default rootReducer