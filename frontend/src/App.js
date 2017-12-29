import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';

import Home from './pages/static/Home';
import About from './pages/static/About';
import Nav from './components/Nav';
import RecipePage from './pages/RecipePage';
import ImageList from './pages/ImageList';
import MealList from './pages/MealList';
import Settings from './pages/Settings';
import EditorPage from './pages/RecipeEditorPage';
import { Container } from 'semantic-ui-react';
import Footer from './components/Footer';
import ReduxToastr from 'react-redux-toastr';

const App = () => (
  <Router>
    <div>
      <ReduxToastr
        timeOut={4000}
        newestOnTop={false}
        preventDuplicates
        position="top-left"
        transitionIn="fadeIn"
        transitionOut="fadeOut"
        progressBar
      />
      <Nav />
      <Container
        fluid
        style={{ marginTop: '7em', width: '80%', minHeight: '100vh' }}
      >
        <Switch>
          <Route exact path="/" component={Home} />
          <Route path="/about" component={About} />
          <Route path="/images" component={ImageList} />
          <Route path="/meals" component={MealList} />
          <Route path="/settings" component={Settings} />
          <Route path="/editor/:recipe_id" component={EditorPage} />
          <Route path="/:recipe_id" component={RecipePage} />
        </Switch>
      </Container>
      <Footer />
    </div>
  </Router>
);
export default App;
