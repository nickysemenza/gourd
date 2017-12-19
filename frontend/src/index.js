import React from 'react';
import ReactDOM from 'react-dom';
import './assets/index.css';
import App from './App';
import registerServiceWorker from './registerServiceWorker';
import "bootstrap/dist/css/bootstrap.css";

ReactDOM.render(<App />, document.getElementById('root'));
registerServiceWorker();
