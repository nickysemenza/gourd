package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	servertiming "github.com/mitchellh/go-server-timing"

	mw2 "github.com/neko-neko/echo-logrus/v2"
	log2 "github.com/neko-neko/echo-logrus/v2/log"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/exporter/stackdriver/propagation"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/trace"

	"github.com/nickysemenza/gourd/api"
	"github.com/nickysemenza/gourd/auth"
	"github.com/nickysemenza/gourd/db"
)

// Server represents a server
type Server struct {
	DB          *db.Client
	HTTPHost    string
	HTTPPort    uint
	HTTPTimeout time.Duration
	APIManager  *api.API
	BypassAuth  bool
	Logger      *log.Logger
}

func (s *Server) Run(_ context.Context) error {
	r := echo.New()

	log2.Logger().Logger = s.Logger
	r.Logger = log2.Logger()
	r.Use(otelecho.Middleware("gourd-server"))
	r.Use(middleware.CORS())
	r.Use(mw2.Logger())
	r.Use(middleware.RequestID())
	r.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	r.Use(middleware.Recover())
	r.Use(echo.WrapMiddleware(func(h http.Handler) http.Handler { return servertiming.Middleware(h, nil) }))
	r.Use(TraceIDHeader)

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
		SigningKey: s.APIManager.Auth.Key,
	}
	jwtMiddleware := middleware.JWTWithConfig(config)

	// r.Add("/images", echo.WrapHandler(s.APIManager.ImageStore.Handler))
	r.Static("/images", s.APIManager.ImageStore.Dir())
	// http routes
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

	r.GET("/debug/notion", s.APIManager.NotionTest)
	r.GET("/misc", s.APIManager.Misc)

	r.GET("/*", func(c echo.Context) error {
		root := "ui/build"
		fs := http.FileServer(http.Dir(root))
		if _, err := os.Stat(root + c.Request().RequestURI); os.IsNotExist(err) {
			http.StripPrefix(c.Request().RequestURI, fs).ServeHTTP(c.Response().Writer, c.Request())
		} else {
			fs.ServeHTTP(c.Response().Writer, c.Request())
		}
		return nil
	})

	log.Printf("running on: %s", s.GetBaseURL())
	return http.ListenAndServe(s.getAddress(),
		// nolint: contextcheck
		wrapHandler(http.TimeoutHandler(r, s.HTTPTimeout, "timeout")),
	)
}
func (s *Server) getAddress() string {
	return fmt.Sprintf("%s:%d", s.HTTPHost, s.HTTPPort)
}
func (s *Server) GetBaseURL() string {
	return fmt.Sprintf("http://%s", s.getAddress())
}

func wrapHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// opts := []trace.SpanOption{}
		sc, ok := (&propagation.HTTPFormat{}).SpanContextFromRequest(r)
		if ok && sc.SpanID.String() != "" {
			sc2 := trace.NewSpanContext(trace.SpanContextConfig{
				SpanID:     trace.SpanID(sc.SpanID),
				TraceID:    trace.TraceID(sc.TraceID),
				TraceFlags: trace.FlagsSampled,
			})

			ctx := trace.ContextWithRemoteSpanContext(r.Context(), sc2)
			r = r.Clone(ctx)
		}

		h.ServeHTTP(w, r) // call original
	})
}

func (s *Server) Scrape(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	recipe, err := s.APIManager.FetchAndTransform(ctx, "https://www.seriouseats.com/recipes/2013/12/roasted-kabocha-squash-soy-sauce-butter-shichimi-recipe.html", s.APIManager.IngredientIdByName)
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

			c.Response().Header().Set(echo.HeaderXRequestID, sc.TraceID().String())
		}
		return next(c)
	}
}
