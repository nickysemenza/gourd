
import {
    REQUEST_USERS,
    RECEIVE_USERS,
    REQUEST_APPLICATION_DETAIL,
    RECEIVE_APPLICATION_DETAIL,
} from '../actions/exec';


const INITIAL_STATE = {
    user_list: [],
    user_list_loading: false,
    application_list: [],
    application_list_loading: false,
    user_detail: {},
    user_detail_loading: false,
    application_detail: {},
    application_detail_loading: false,
};

export default function (state = INITIAL_STATE, action) {
    switch (action.type) {
        case REQUEST_USERS:
            return { ...state, user_list_loading: true };
        case RECEIVE_USERS:
            //todo: error checking
            return { ...state,
                user_list_loading: false,
                user_list: action.json.data
            };
        case REQUEST_APPLICATION_DETAIL:
            return { ...state, application_detail_loading: true };
        case RECEIVE_APPLICATION_DETAIL:
            //todo: error checking
            return { ...state,
                application_detail_loading: false,
                application_detail: {
                    ...state.application_detail,
                    [action.applicationId]: action.json.data
                }
            };
        default:
            return state;
    }
}