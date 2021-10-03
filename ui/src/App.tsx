import React from "react";
import { Helmet } from "react-helmet";

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
  isLoggedIn,
  onAPIRequest,
  onAPIError,
  getTracingURL,
} from "./config";
import { CookiesProvider } from "react-cookie";
import { Docs } from "./pages/Misc";
import Albums from "./pages/Albums";
import "react-toastify/dist/ReactToastify.css";
import { ToastContainer } from "react-toastify";
import IngredientDetail from "./pages/IngredientDetail";
import { WasmContextProvider } from "./wasm";
import Graph from "./pages/Graph";
import { FetchInstrumentation } from "@opentelemetry/instrumentation-fetch";

import { WebTracerProvider } from "@opentelemetry/web";
import { CollectorTraceExporter } from "@opentelemetry/exporter-collector";
import { SimpleSpanProcessor } from "@opentelemetry/tracing";
import { registerInstrumentations } from "@opentelemetry/instrumentation";
import { ZoneContextManager } from "@opentelemetry/context-zone";
import { JaegerPropagator } from "@opentelemetry/propagator-jaeger";
import { diag, DiagConsoleLogger, DiagLogLevel } from "@opentelemetry/api";

const registerTracing = (url: string) => {
  if (url === "") return;
  console.info("enabled tracing", url);
  diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.ERROR);

  const exporter = new CollectorTraceExporter({
    url: url,
    // https://github.com/open-telemetry/opentelemetry-js/issues/2321#issuecomment-889861080
    headers: {
      "Content-Type": "application/json",
    },
    // serviceName: "auto-instrumentations-web",
  });

  const provider = new WebTracerProvider();

  registerInstrumentations({
    instrumentations: [
      new FetchInstrumentation({
        propagateTraceHeaderCorsUrls: /.+/,
      }),
    ],
  });
  provider.register({
    contextManager: new ZoneContextManager(),
    // propagator: new B3Propagator(),
    propagator: new JaegerPropagator(),
  });

  provider.addSpanProcessor(new SimpleSpanProcessor(exporter));
  // provider.addSpanProcessor(new SimpleSpanProcessor(new ConsoleSpanExporter()));
};
registerTracing(getTracingURL());

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
  return (
    <CookiesProvider>
      <RestfulProvider
        base={getAPIURL()}
        onRequest={onAPIRequest}
        onError={onAPIError}
      >
        <WasmContextProvider>
          <Router>
            <Helmet>
              <title>gourd</title>
            </Helmet>
            <ToastContainer position="bottom-right" />
            <NavBar />
            <div className="lg:container lg:mx-auto">
              <Switch>
                <Route path="/recipe/:id">
                  <RecipeDetail />
                </Route>
                <Route path="/recipes">
                  <RecipeList />
                </Route>
                <Route path="/ingredients/:id">
                  <IngredientDetail />
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
                <Route path="/graph">
                  <Graph />
                </Route>
                <PrivateRoute path="/photos">
                  <Photos />
                </PrivateRoute>
                <PrivateRoute path="/meals">
                  <Meals />
                </PrivateRoute>
                <PrivateRoute path="/albums">
                  <Albums />
                </PrivateRoute>
                <Route path="/">
                  <Test />
                </Route>
              </Switch>
            </div>
          </Router>
        </WasmContextProvider>

        <hr />
      </RestfulProvider>
    </CookiesProvider>
  );
}

export default App;
