import React from 'react';
import ReactDOM from 'react-dom';
import Home from './Home';
import { shallow, mount } from 'enzyme';
import { MemoryRouter } from 'react-router-dom';

import Adapter from 'enzyme-adapter-react-16';
import { configure } from 'enzyme';
import { createStore } from 'redux';
import configureStore from '../../store/configureStore';
import { Provider } from 'react-redux';
configure({ adapter: new Adapter() });

beforeEach(function() {
  window.fetch = jest.fn().mockImplementation(() =>
    Promise.resolve({
      ok: true,
      Id: '123',
      json: () => ['thing1', 'thing2', 'thing3']
    })
  );
});

it('renders without crashing', () => {
  const div = document.createElement('div');
  const store = configureStore();
  ReactDOM.render(
    <MemoryRouter>
      <Provider store={store}>
        <Home />
      </Provider>
    </MemoryRouter>,
    div
  );
});
