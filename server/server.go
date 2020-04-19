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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/plugin/othttp"
)

// Server represents a server
type Server struct {
	Manager  *manager.Manager
	DB       *db.Client
	HTTPPort uint
}

// Run runs http
func (s *Server) Run() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{}))

	r.Get("/_metrics", promhttp.Handler().ServeHTTP)
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
		othttp.NewHandler(http.TimeoutHandler(r, time.Second*30, "timeout"), "server",
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
