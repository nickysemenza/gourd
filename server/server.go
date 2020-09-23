package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	sentryecho "github.com/getsentry/sentry-go/echo"
	log "github.com/sirupsen/logrus"
	echotrace "go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo"
	othttp "go.opentelemetry.io/contrib/instrumentation/net/http"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"

	"github.com/nickysemenza/gourd/api"
	"github.com/nickysemenza/gourd/db"
	"github.com/nickysemenza/gourd/graph"
	"github.com/nickysemenza/gourd/graph/generated"
	"github.com/nickysemenza/gourd/manager"
	"github.com/nickysemenza/gourd/scraper"
)

// Server represents a server
type Server struct {
	Manager     *manager.Manager
	DB          *db.Client
	HTTPHost    string
	HTTPPort    uint
	HTTPTimeout time.Duration
}

// nolint:gochecknoglobals
var httpRequestsDurationMetric = metric.Must(global.Meter("ex.com/basic")).
	NewFloat64Counter("http.requests.duration")

// func timing(next http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		t := time.Now()
// 		next.ServeHTTP(w, r)
// 		d := time.Since(t).Seconds()
// 		httpRequestsDurationMetric.Add(r.Context(), d, label.String("method", r.Method))
// 	}
// 	return http.HandlerFunc(fn)
// }

func (s *Server) GetResolver() generated.ResolverRoot {
	return &graph.Resolver{
		Manager: s.Manager,
		DB:      s.DB,
	}
}

func (s *Server) Run(_ context.Context) error {
	r := echo.New()

	r.Use(echotrace.Middleware("gourd-backend"))
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.RequestID())
	r.Use(middleware.Recover())
	r.Use(sentryecho.New(sentryecho.Options{}))

	hf, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		return nil
	}

	// gql server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: s.GetResolver()}))
	srv.Use(extension.FixedComplexityLimit(100))
	// srv.Use(graph.Observability{})

	// http routes
	r.GET("/metrics", echo.WrapHandler(hf))
	r.Any("/", echo.WrapHandler(http.HandlerFunc(playground.Handler("GraphQL playground", "/query"))))
	r.Any("/query", echo.WrapHandler(othttp.WithRouteTag("/query", srv)))
	r.GET("/scrape", echo.WrapHandler(http.HandlerFunc(s.Scrape)))

	apiManager := api.NewAPI(s.Manager)
	// r.Any("/api", echo.WrapHandler(othttp.NewHandler(api.Handler(apiManager), "/api")))
	g := r.Group("/api")
	api.RegisterHandlers(g, apiManager)

	r.GET("/recipes/{uuid}", s.GetRecipe)

	r.GET("/routes", func(c echo.Context) error {
		return c.JSON(http.StatusOK, r.Routes())
	})

	addr := fmt.Sprintf("%s:%d", s.HTTPHost, s.HTTPPort)
	log.Printf("running on: http://%s/", addr)
	return http.ListenAndServe(addr,
		othttp.NewHandler(http.TimeoutHandler(r, s.HTTPTimeout, "timeout"), "server",
			othttp.WithMessageEvents(othttp.ReadEvents, othttp.WriteEvents),
		),
	)
}

func (s *Server) GetRecipe(c echo.Context) error {
	recipe, _ := s.Manager.GetRecipe(c.Request().Context(), c.QueryParam("uuid"))
	return c.JSON(http.StatusOK, recipe)
	// writeJSON(w, http.StatusOK, recipe)
}

func (s *Server) Scrape(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	recipe, err := scraper.FetchAndTransform(ctx, "https://www.seriouseats.com/recipes/2013/12/roasted-kabocha-squash-soy-sauce-butter-shichimi-recipe.html", s.GetResolver().Mutation().UpsertIngredient)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, recipe)
}
func writeErr(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusInternalServerError, err.Error())
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(body); err != nil {
		panic(err)
	}
}
