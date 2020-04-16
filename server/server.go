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
	"github.com/nickysemenza/food/graph"
	"github.com/nickysemenza/food/graph/generated"
	"github.com/nickysemenza/food/manager"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// Server represents a server
type Server struct {
	Manager  *manager.Manager
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

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{Manager: s.Manager}}))
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	r.Get("/h", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hi"))
	})
	r.Get("/recipes/{uuid}", s.GetRecipe)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.HTTPPort),
		Handler: http.TimeoutHandler(r, time.Second*30, "timeout"),
	}
	log.Printf("connect to http://localhost:%d/ for GraphQL playground", s.HTTPPort)
	return server.ListenAndServe()
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
