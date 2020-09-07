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
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/plugin/othttp"

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
	NewFloat64Measure("http.requests.duration", metric.WithKeys(key.New("method")))

func timing(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)
		d := time.Since(t).Seconds()
		httpRequestsDurationMetric.Record(r.Context(), d, key.String("method", r.Method))
	}
	return http.HandlerFunc(fn)
}

func (s *Server) GetResolver() generated.ResolverRoot {
	return &graph.Resolver{
		Manager: s.Manager,
		DB:      s.DB,
	}
}

func (s *Server) Run(_ context.Context) error {
	r := chi.NewRouter()

	sentryMiddleware := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})

	r.Use(middleware.RequestID)
	r.Use(NewStructuredLogger(log.New()))
	r.Use(cors.Handler(cors.Options{}))
	r.Use(timing)
	r.Use(middleware.Recoverer)
	r.Use(sentryMiddleware.Handle)

	_, hf, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		return nil
	}

	// gql server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: s.GetResolver()}))
	srv.Use(extension.FixedComplexityLimit(100))
	// srv.Use(graph.Observability{})

	// http routes
	r.Get("/metrics", hf)
	r.Mount("/debug", middleware.Profiler())
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", othttp.WithRouteTag("/query", srv))
	r.Get("/scrape", s.Scrape)
	r.Get("/h", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hi"))
	})

	r.Get("/recipes/{uuid}", s.GetRecipe)

	addr := fmt.Sprintf("%s:%d", s.HTTPHost, s.HTTPPort)
	log.Printf("running on: http://%s/", addr)
	return http.ListenAndServe(addr,
		othttp.NewHandler(http.TimeoutHandler(r, s.HTTPTimeout, "timeout"), "server",
			othttp.WithMessageEvents(othttp.ReadEvents, othttp.WriteEvents),
		),
	)
}

func (s *Server) GetRecipe(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id := chi.URLParam(req, "uuid")
	recipe, _ := s.Manager.GetRecipe(ctx, id)
	writeJSON(w, http.StatusOK, recipe)
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
