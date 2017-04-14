import React from "react";
import ReactDOM from "react-dom";
import Home from "./Home";
import { shallow, mount } from "enzyme";
import { MemoryRouter } from "react-router-dom";

beforeEach(function() {
  window.fetch = jest.fn().mockImplementation(() =>
    Promise.resolve({
      ok: true,
      Id: "123",
      json: () => ["thing1", "thing2", "thing3"]
    })
  );
});

it("renders without crashing", () => {
  const div = document.createElement("div");
  ReactDOM.render(
    <MemoryRouter>
      <Home />
    </MemoryRouter>,
    div
  );
});

// it('renders without crashing2', () => {
//     shallow(<Home />);
// });

it("renders welcome message", () => {
  const wrapper = shallow(<Home />);
  const welcome = <h2>Nicky's Recipe Stash</h2>;
  // console.log(wrapper.debug());
  expect(wrapper.contains(welcome)).toEqual(true);
});
