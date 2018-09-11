import { LOGIN_FROM_JWT_SUCCESS, RECEIVE_ME } from '../actions/users';
const INITIAL_STATE = {
  authenticated: false,
  me: {
    first_name: '',
    last_name: ''
  },
  token: null
};

export default function(state = INITIAL_STATE, action) {
  switch (action.type) {
    case LOGIN_FROM_JWT_SUCCESS:
      return {
        ...state,
        authenticated: true,
        token: action.token
      };
    case RECEIVE_ME:
      //todo: error checking
      return {
        ...state,
        me: action.me
      };
    default:
      return state;
  }
}
