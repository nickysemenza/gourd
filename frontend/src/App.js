import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';

import Home from './pages/static/Home';
import About from './pages/static/About';
import Nav from './components/Nav';
import RecipePage from './pages/RecipePage';
import EditorPage from './pages/RecipeEditorPage';

const App = () => (
  <Router>
    <div>
      <Nav />
      <div className="container-fluid">
        <Switch>
          <Route exact path="/" component={Home} />
          <Route path="/about" component={About} />
          <Route path="/editor/:recipe_id" component={EditorPage} />
          <Route path="/:recipe_id" component={RecipePage} />
        </Switch>
      </div>
    </div>
  </Router>
);
export default App;
