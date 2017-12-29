import React from 'react';
import ReactDOM from 'react-dom';
import './assets/index.css';
import 'react-big-calendar/lib/css/react-big-calendar.css';
import App from './App';
import registerServiceWorker from './registerServiceWorker';
import configureStore from './store/configureStore';
import { Provider } from 'react-redux';

const store = configureStore();
ReactDOM.render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById('root')
);
registerServiceWorker();
