import { Helmet } from "react-helmet";

import Test from "./Homepage";
import RecipeList from "./pages/recipe/RecipeList";
import {
  BrowserRouter as Router,
  Route,
  Navigate,
  Routes,
  useLocation,
} from "react-router-dom";
import RecipeDetail from "./pages/recipe/RecipeDetail";
import NavBar from "./components/nav/NavBar";
import IngredientList from "./pages/IngredientList";
import CreateRecipe from "./pages/recipe/CreateRecipe";
import Food from "./pages/Food";
import Playground from "./pages/Playground";
import RecipeDiff from "./pages/recipe/RecipeDiff";
import Search from "./pages/Search";

import Photos from "./pages/Photos";
import Meals from "./pages/Meals";
import { CookiesProvider } from "react-cookie";
import Albums from "./pages/Albums";
import "react-toastify/dist/ReactToastify.css";
import { ToastContainer } from "react-toastify";
import IngredientInfo from "./pages/IngredientInfo";
import { WasmContextProvider } from "./util/wasmContext";
import Graph from "./pages/Graph";
import ErrorBoundary from "./components/ui/ErrorBoundary";
import { isLoggedIn } from "./auth/auth";
import { registerTracing } from "./util/tracing";
import ErrorPage from "./components/ui/ErrorPage";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { getTracingURL } from "./util/urls";
import { Dialog } from "./components/SearchPopover";

registerTracing(getTracingURL(), true);

const RequireAuth = ({ children }: { children: JSX.Element }) => {
  const authed = isLoggedIn() || true;
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
const queryClient = new QueryClient();

function App() {
  return (
    <CookiesProvider>
      <QueryClientProvider client={queryClient}>
        <WasmContextProvider>
          {/* @ts-expect-error bug */}
          <Helmet>
            <title>gourd</title>
          </Helmet>
          <ReactQueryDevtools initialIsOpen={false} />

          <ToastContainer position="bottom-right" />
          <Router>
            <NavBar />
            <div className="lg:container lg:mx-auto">
              <ErrorBoundary>
                <Dialog />
                <Routes>
                  <Route index element={<Test />} />
                  <Route path="recipe/:id" element={<RecipeDetail />} />
                  <Route path="recipes" element={<RecipeList />} />
                  <Route path="ingredients/:id" element={<IngredientInfo />} />
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
      </QueryClientProvider>

      <hr />
    </CookiesProvider>
  );
}

export default App;

function NoMatch() {
  const location = useLocation();

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
