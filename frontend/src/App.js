import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';

import Home from './pages/static/Home';
import About from './pages/static/About';
import Nav from './components/Nav';
import RecipePage from './pages/RecipePage';
import EditorPage from './pages/RecipeEditorPage';
import { Container } from 'semantic-ui-react';
import Footer from './components/Footer';

const App = () => (
  <Router>
    <div>
      <Nav />
      <Container
        fluid
        style={{ marginTop: '7em', width: '80%', minHeight: '100vh' }}
      >
        <Switch>
          <Route exact path="/" component={Home} />
          <Route path="/about" component={About} />
          <Route path="/editor/:recipe_id" component={EditorPage} />
          <Route path="/:recipe_id" component={RecipePage} />
        </Switch>
      </Container>
      <Footer />
    </div>
  </Router>
);
export default App;
