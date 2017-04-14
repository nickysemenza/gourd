import React from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";

import Home from "./Home";
import About from "./About";
import Nav from "./Nav";
import RecipePage from "./RecipePage";

const App = () => (
  <Router>
    <div>
      <Nav />
      <div className="container">
        <Switch>
          <Route exact path="/" component={Home} />
          <Route path="/about" component={About} />
          <Route path="/:recipe_id" component={RecipePage} />
        </Switch>
      </div>
    </div>
  </Router>
);
export default App;
