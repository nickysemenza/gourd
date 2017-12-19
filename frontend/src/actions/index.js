import {API_BASE_URL } from '../config';

export default function apiFetch(endpoint, options = {}) {
    return fetch(`${API_BASE_URL}/${endpoint}`,options);
}