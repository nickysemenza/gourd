import React from "react";

import ApolloClient from "apollo-boost";
import { ApolloProvider } from "@apollo/react-hooks";
import Test from "./Test";
import RecipeList from "./pages/RecipeList";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
  RouteProps,
} from "react-router-dom";
import RecipeDetail from "./pages/RecipeDetail";
import NavBar from "./components/NavBar";
import IngredientList from "./pages/IngredientList";
import CreateRecipe from "./pages/CreateRecipe";
import Food from "./pages/Food";
import Playground from "./pages/Playground";

import "./tailwind.output.css";
import { RestfulProvider } from "restful-react";
import Photos from "./pages/Photos";
import Meals from "./pages/Meals";
import {
  getAPIURL,
  getGQLURL,
  getJWT,
  isLoggedIn,
  onAPIRequest,
} from "./config";
import { CookiesProvider } from "react-cookie";
import { Docs } from "./pages/Misc";

const PrivateRoute = ({ children, ...rest }: RouteProps) => {
  return (
    <Route
      {...rest}
      render={({ location }) =>
        isLoggedIn() ? (
          children
        ) : (
          <Redirect
            to={{
              pathname: "/login",
              state: { from: location },
            }}
          />
        )
      }
    />
  );
};

function App() {
  const client = new ApolloClient({
    uri: getGQLURL(),
    headers: { authorization: "Bearer " + getJWT() },
  });
  return (
    <CookiesProvider>
      <RestfulProvider base={getAPIURL()} onRequest={onAPIRequest}>
        <ApolloProvider client={client}>
          <Router>
            <NavBar />
            <div className="lg:container lg:mx-auto">
              <Switch>
                <Route path="/recipe/:uuid">
                  <RecipeDetail />
                </Route>
                <Route path="/recipes">
                  <RecipeList />
                </Route>
                <Route path="/ingredients">
                  <IngredientList />
                </Route>
                <Route path="/create">
                  <CreateRecipe />
                </Route>
                <Route path="/food">
                  <Food />
                </Route>
                <Route path="/docs">
                  <Docs />
                </Route>
                <Route path="/playground">
                  <Playground />
                </Route>
                <PrivateRoute path="/photos">
                  <Photos />
                </PrivateRoute>
                <PrivateRoute path="/meals">
                  <Meals />
                </PrivateRoute>
                <Route path="/">
                  <Test />
                </Route>
              </Switch>
            </div>
          </Router>

          <hr />
        </ApolloProvider>
      </RestfulProvider>
    </CookiesProvider>
  );
}

export default App;
