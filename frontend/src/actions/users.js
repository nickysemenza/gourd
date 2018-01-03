import cookie from 'react-cookie';
import apiFetch from './index';
export const LOGIN_FROM_JWT_SUCCESS = 'LOGIN_FROM_JWT_SUCCESS';
export function loginFromJWT(token) {
  cookie.save('token', token, { path: '/' });
  return dispatch => {
    dispatch(saveToken(token));
    setTimeout(() => {
      dispatch(fetchMe());
    }, 50);
  };
}

function saveToken(token) {
  return {
    type: LOGIN_FROM_JWT_SUCCESS,
    token: token
  };
}

export const REQUEST_ME = 'REQUEST_ME';
export const RECEIVE_ME = 'RECEIVE_ME';

export function fetchMe() {
  return dispatch => {
    dispatch(requestMe());
    return apiFetch('me')
      .then(response => response.json())
      .then(json => dispatch(receiveMe(json)));
  };
}

function requestMe() {
  return {
    type: REQUEST_ME
  };
}

function receiveMe(json) {
  return {
    type: RECEIVE_ME,
    me: json
  };
}
