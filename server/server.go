package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nickysemenza/food/manager"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	r.Get("/_metrics", promhttp.Handler().ServeHTTP)
	r.Mount("/debug", middleware.Profiler())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hi"))
	})
	r.Get("/recipes/{uuid}", s.GetRecipe)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.HTTPPort),
		Handler: http.TimeoutHandler(r, time.Second*30, "timeout"),
	}
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
