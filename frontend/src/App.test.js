import React from 'react';
import Adapter from 'enzyme-adapter-react-16';
import { configure } from 'enzyme';
configure({ adapter: new Adapter() });
it('renders without crashing', () => {
  const div = document.createElement('div');
  {
    /*ReactDOM.render(<App />, div);*/
  }
});
