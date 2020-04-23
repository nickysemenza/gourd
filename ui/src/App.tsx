import React from "react";
import "./App.css";

import ApolloClient from "apollo-boost";
import { ApolloProvider } from "@apollo/react-hooks";
import Test from "./Test";
import { ThemeProvider } from "theme-ui";
import theme from "./theme";
import RecipeList from "./RecipeList";
import { BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";
import RecipeDetail from "./RecipeDetail";

function App() {
  const client = new ApolloClient({
    uri: "http://localhost:4242/query",
  });
  return (
    <ThemeProvider theme={theme}>
      <ApolloProvider client={client}>
        <Router>
          {/* TODO: nav */}
          <Switch>
            <Route path="/recipe/:uuid">
              <RecipeDetail />
            </Route>
            <Route path="/recipes">
              <RecipeList />
            </Route>
            <Route path="/">
              <Test />
            </Route>
          </Switch>
        </Router>

        <hr />
      </ApolloProvider>
    </ThemeProvider>
  );
}

export default App;
