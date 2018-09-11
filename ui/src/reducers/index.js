import { combineReducers } from 'redux';
import { reducer as toastrReducer } from 'react-redux-toastr';
import recipe from './reducer_recipe';
import user from './reducer_user';
const rootReducer = combineReducers({
  recipe,
  user,
  toastr: toastrReducer
});

export default rootReducer;
