package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/nickysemenza/food/db"
	"github.com/nickysemenza/food/graph"
	"github.com/nickysemenza/food/graph/generated"
	"github.com/nickysemenza/food/manager"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/plugin/othttp"
)

// Server represents a server
type Server struct {
	Manager     *manager.Manager
	DB          *db.Client
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

func (s *Server) Run() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(NewStructuredLogger(log.New()))
	r.Use(cors.Handler(cors.Options{}))
	r.Use(timing)

	_, hf, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		return nil
	}

	r.Get("/metrics", hf)
	r.Mount("/debug", middleware.Profiler())

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{
			Manager: s.Manager,
			DB:      s.DB,
		}},
	))
	srv.Use(graph.Observability{})
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", othttp.WithRouteTag("/query", srv))

	r.Get("/h", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hi"))
	})
	r.Get("/recipes/{uuid}", s.GetRecipe)

	// server := &http.Server{
	// 	Addr:    fmt.Sprintf(":%d", s.HTTPPort),
	// 	Handler: http.TimeoutHandler(r, time.Second*30, "timeout"),
	// }
	log.Printf("connect to http://localhost:%d/ for GraphQL playground", s.HTTPPort)
	// return server.ListenAndServe()
	return http.ListenAndServe(fmt.Sprintf(":%d", s.HTTPPort),
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
