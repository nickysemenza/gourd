import React from "react";
import { Helmet } from "react-helmet";

import Test from "./Test";
import RecipeList from "./pages/RecipeList";
import {
  BrowserRouter as Router,
  Route,
  Navigate,
  Routes,
  useLocation,
} from "react-router-dom";
import RecipeDetail from "./pages/RecipeDetail";
import NavBar from "./components/NavBar";
import IngredientList from "./pages/IngredientList";
import CreateRecipe from "./pages/CreateRecipe";
import Food from "./pages/Food";
import Playground from "./pages/Playground";
import RecipeDiff from "./pages/RecipeDiff";
import Search from "./pages/Search";

import "./tailwind.output.css";
import { RestfulProvider } from "restful-react";
import Photos from "./pages/Photos";
import Meals from "./pages/Meals";
import { getAPIURL, onAPIRequest, onAPIError, getTracingURL } from "./config";
import { CookiesProvider } from "react-cookie";
import Albums from "./pages/Albums";
import "react-toastify/dist/ReactToastify.css";
import { ToastContainer } from "react-toastify";
import IngredientDetail from "./pages/IngredientDetail";
import { WasmContextProvider } from "./wasmContext";
import Graph from "./pages/Graph";
import ErrorBoundary from "./components/ErrorBoundary";
import { isLoggedIn } from "./auth";
import { registerTracing } from "./tracing";
import ErrorPage from "./components/ErrorPage";

registerTracing(getTracingURL(), true);

const RequireAuth = ({ children }: { children: JSX.Element }) => {
  let authed = isLoggedIn() || true;
  return authed ? (
    children
  ) : (
    <Navigate
      to={{
        pathname: "login",
      }}
    />
  );
};

function App() {
  return (
    <CookiesProvider>
      {/* @ts-ignore */}
      <RestfulProvider
        base={getAPIURL()}
        onRequest={onAPIRequest}
        onError={onAPIError}
      >
        <WasmContextProvider>
          {/* @ts-ignore */}
          <Helmet>
            <title>gourd</title>
          </Helmet>
          <ToastContainer position="bottom-right" />
          <Router>
            <NavBar />
            <div className="lg:container lg:mx-auto">
              <ErrorBoundary>
                <Routes>
                  <Route index element={<Test />} />
                  <Route path="recipe/:id" element={<RecipeDetail />} />
                  <Route path="recipes" element={<RecipeList />} />
                  <Route
                    path="ingredients/:id"
                    element={<IngredientDetail />}
                  />
                  <Route path="ingredients" element={<IngredientList />} />
                  <Route path="create" element={<CreateRecipe />} />
                  <Route path="food" element={<Food />} />
                  <Route path="playground" element={<Playground />} />
                  <Route path="diff" element={<RecipeDiff />} />
                  <Route path="graph" element={<Graph />} />
                  <Route path="search" element={<Search />} />
                  <Route
                    path="photos"
                    element={
                      <RequireAuth>
                        <Photos />
                      </RequireAuth>
                    }
                  />
                  <Route
                    path="meals"
                    element={
                      <RequireAuth>
                        <Meals />
                      </RequireAuth>
                    }
                  />
                  <Route
                    path="albums"
                    element={
                      <RequireAuth>
                        <Albums />
                      </RequireAuth>
                    }
                  />
                  <Route path="*" element={<NoMatch />} />
                </Routes>
              </ErrorBoundary>
            </div>
          </Router>
        </WasmContextProvider>

        <hr />
      </RestfulProvider>
    </CookiesProvider>
  );
}

export default App;

function NoMatch() {
  let location = useLocation();

  return (
    <ErrorPage
      title="not found"
      message={
        <>
          no match for<code>{location.pathname}</code>
        </>
      }
    />
  );
}
