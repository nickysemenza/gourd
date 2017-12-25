import { combineReducers } from 'redux';
import { reducer as toastrReducer } from 'react-redux-toastr';
import recipe from './reducer_recipe';
const rootReducer = combineReducers({
  recipe,
  toastr: toastrReducer
});

export default rootReducer;
