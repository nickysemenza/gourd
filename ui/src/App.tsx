import React from "react";
import "./App.css";

import ApolloClient from "apollo-boost";
import { ApolloProvider } from "@apollo/react-hooks";
import Test from "./Test";
import { ThemeProvider, Box } from "theme-ui";
import theme from "./theme";
import RecipeList from "./RecipeList";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
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
          <Box
            sx={{
              maxWidth: "80%",
              mx: "auto",
              px: 3,
            }}
          >
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
          </Box>
        </Router>

        <hr />
      </ApolloProvider>
    </ThemeProvider>
  );
}

export default App;
