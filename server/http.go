package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	servertiming "github.com/mitchellh/go-server-timing"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/exporter/stackdriver/propagation"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/trace"

	"github.com/nickysemenza/gourd/api"
	"github.com/nickysemenza/gourd/auth"
	"github.com/nickysemenza/gourd/db"
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
	APIManager  *api.API
	BypassAuth  bool
}

func (s *Server) Run(_ context.Context) error {
	r := echo.New()

	r.Use(otelecho.Middleware("gourd-server"))
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.RequestID())
	r.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	r.Use(middleware.Recover())
	r.Use(echo.WrapMiddleware(func(h http.Handler) http.Handler { return servertiming.Middleware(h, nil) }))
	r.Use(TraceIDHeader)

	hf, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		return err
	}

	skipper := func(c echo.Context) bool {
		if s.BypassAuth {
			log.Debugf("bypassing auth for %s", c.Path())

			return true
		}
		switch c.Path() {
		case "/api/auth", "/spec":
			return true
		default:
			return false
		}
	}
	config := middleware.JWTConfig{
		Skipper:    skipper,
		Claims:     &auth.Claims{},
		SigningKey: s.Manager.Auth.Key,
	}
	jwtMiddleware := middleware.JWTWithConfig(config)

	// http routes
	r.GET("/metrics", echo.WrapHandler(hf))
	r.GET("/scrape", echo.WrapHandler(http.HandlerFunc(s.Scrape)))

	// r.Any("/api", echo.WrapHandler(otelhttp.NewHandler(api.Handler(apiManager), "/api")))
	g := r.Group("/api")

	g.Use(jwtMiddleware)
	api.RegisterHandlers(g, s.APIManager)

	r.GET("/routes", func(c echo.Context) error {
		return c.JSON(http.StatusOK, r.Routes())
	})
	r.GET("/debug", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{"headers": c.Request().Header})
	})
	r.GET("/spec", func(c echo.Context) error {
		spec, err := api.GetSwagger()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, api.Error{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, spec)
	})

	// r.GET("/auth/redirect", func(c echo.Context) error {
	// 	return c.Redirect(http.StatusTemporaryRedirect, s.APIManager.Google.GetURL())
	// })
	// r.GET("/auth/callback", func(c echo.Context) error {
	// 	code := c.Request().FormValue("code")
	// 	err := s.APIManager.Google.Finish(c.Request().Context(), code)
	// 	if err != nil {
	// 		return c.JSON(http.StatusInternalServerError, err.Error())
	// 	}
	// 	return c.JSON(http.StatusOK, "ok")
	// })

	// r.GET("/photos", func(c echo.Context) error {
	// 	pics, err := s.APIManager.Google.GetTest(c.Request().Context())
	// 	if err != nil {
	// 		spew.Dump(err)
	// 		return c.JSON(http.StatusInternalServerError, err.Error())
	// 	}
	// 	return c.JSON(http.StatusOK, pics)
	// })

	addr := fmt.Sprintf("%s:%d", s.HTTPHost, s.HTTPPort)
	log.Printf("running on: http://%s/", addr)
	return http.ListenAndServe(addr,
		wrapHandler(http.TimeoutHandler(r, s.HTTPTimeout, "timeout")),
		// otelhttp.NewHandler(http.TimeoutHandler(r, s.HTTPTimeout, "timeout"), "server",
		// 	otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		// ),
	)
}
func wrapHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// opts := []trace.SpanOption{}
		sc, ok := (&propagation.HTTPFormat{}).SpanContextFromRequest(r)
		if ok && sc.SpanID.String() != "" {
			sc2 := trace.SpanContext{
				SpanID:  trace.SpanID(sc.SpanID),
				TraceID: trace.TraceID(sc.TraceID),
			}
			sc2.TraceFlags |= trace.FlagsSampled
			spew.Dump("ctx-before", r.Context())

			ctx := trace.ContextWithRemoteSpanContext(r.Context(), sc2)
			r = r.Clone(ctx)
			spew.Dump("extracted span", sc2)
		}

		// h2 := otelhttp.NewHandler(h, "server",
		// 	otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		// 	// otelhttp.WithSpanOptions(opts...),
		// )
		h.ServeHTTP(w, r) // call original
	})
}

func (s *Server) Scrape(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	recipe, err := scraper.FetchAndTransform(ctx, "https://www.seriouseats.com/recipes/2013/12/roasted-kabocha-squash-soy-sauce-butter-shichimi-recipe.html", s.APIManager.IngredientIdByName)
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

func TraceIDHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		sc := trace.SpanContextFromContext(ctx)
		if sc.IsValid() {

			c.Response().Header().Set(echo.HeaderXRequestID, sc.TraceID.String())
		}
		return next(c)
	}
}
