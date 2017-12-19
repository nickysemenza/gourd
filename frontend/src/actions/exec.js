import apiFetch  from './index';

export const REQUEST_USERS = 'REQUEST_USERS';
export const RECEIVE_USERS = 'RECEIVE_USERS';

export function fetchUsers () {
    return (dispatch) => {
        dispatch(requestUsers());
        return apiFetch('exec/users')
            .then((response) => response.json())
            .then((json) => dispatch(receiveUsers(json)));
    };
}

function requestUsers () {
    return {
        type: REQUEST_USERS
    };
}

function receiveUsers (json) {
    return {
        type: RECEIVE_USERS,
        json,
        receivedAt: Date.now()
    };
}

export function searchUsers (data) {
    return (dispatch) => {
        dispatch(requestUsers());
        return apiFetch('exec/users/search',
            {
                method: 'POST',
                body: JSON.stringify(data)
            })
            .then((response) => response.json())
            .then((json) => {
                if(json.success) {
                    dispatch(receiveUsers(json));
                } else {
                }
            });
    };
}

export const REQUEST_APPLICATION_DETAIL = 'REQUEST_APPLICATION_DETAIL';
export const RECEIVE_APPLICATION_DETAIL = 'RECEIVE_APPLICATION_DETAIL';

export function fetchApplicationDetail (applicationId) {
    return (dispatch) => {
        dispatch(requestApplicationDetail(applicationId));
        return apiFetch('exec/applications/'+applicationId)
            .then((response) => response.json(applicationId))
            .then((json) => dispatch(receiveApplicationDetail(json,applicationId)));
    };
}

function requestApplicationDetail (applicationId) {
    return {
        type: REQUEST_APPLICATION_DETAIL,
        applicationId
    };
}

function receiveApplicationDetail (json, applicationId) {
    return {
        type: RECEIVE_APPLICATION_DETAIL,
        applicationId,
        json,
        receivedAt: Date.now()
    };
}

