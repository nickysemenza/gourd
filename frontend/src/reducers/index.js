import { combineReducers } from 'redux'
import exec from './reducer_exec';
const rootReducer = combineReducers({
    exec,
});

export default rootReducer