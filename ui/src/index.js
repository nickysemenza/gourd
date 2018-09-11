import React from 'react';
import ReactDOM from 'react-dom';
import './assets/index.css';
import 'react-big-calendar/lib/css/react-big-calendar.css';
import App from './App';
import registerServiceWorker from './registerServiceWorker';
import configureStore from './store/configureStore';
import { Provider } from 'react-redux';
import { loginFromJWT } from './actions/users';
import cookie from 'react-cookie';

const store = configureStore();

const token = cookie.load('token');
if (token) store.dispatch(loginFromJWT(token)); //log us in from cookie

ReactDOM.render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById('root')
);
registerServiceWorker();
