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
	"github.com/davecgh/go-spew/spew"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/exporter/stackdriver/propagation"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"

	"github.com/nickysemenza/gourd/api"
	"github.com/nickysemenza/gourd/auth"
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
	APIManager  *api.API
	BypassAuth  bool
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

	r.Use(otelecho.Middleware("gourd-server"))
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.RequestID())
	r.Use(middleware.Recover())
	r.Use(sentryecho.New(sentryecho.Options{}))

	hf, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		return nil
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

	// gql server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: s.GetResolver()}))
	srv.Use(extension.FixedComplexityLimit(100))
	// srv.Use(graph.Observability{})

	// http routes
	r.GET("/metrics", echo.WrapHandler(hf))
	r.Any("/", echo.WrapHandler(http.HandlerFunc(playground.Handler("GraphQL playground", "/query"))))
	r.Any("/query", echo.WrapHandler(otelhttp.WithRouteTag("/query", srv)), jwtMiddleware)
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
		log.Println("Before")

		// opts := []trace.SpanOption{}
		sc, ok := (&propagation.HTTPFormat{}).SpanContextFromRequest(r)
		if ok && sc.SpanID.String() != "" {
			sc2 := trace.SpanContext{
				SpanID:  trace.SpanID(sc.SpanID),
				TraceID: trace.ID(sc.TraceID),
			}
			sc2.TraceFlags |= trace.FlagsSampled
			spew.Dump("ctx-before", r.Context())

			ctx := trace.ContextWithRemoteSpanContext(r.Context(), sc2)
			r = r.Clone(ctx)
			spew.Dump("extracted span", sc2)
		}

		h2 := otelhttp.NewHandler(h, "server",
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
			// otelhttp.WithSpanOptions(opts...),
		)
		h2.ServeHTTP(w, r) // call original
		log.Println("After")
	})
}

func (s *Server) Scrape(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	recipe, err := scraper.FetchAndTransform(ctx, "https://www.seriouseats.com/recipes/2013/12/roasted-kabocha-squash-soy-sauce-butter-shichimi-recipe.html", s.APIManager.IngredientUUIDByName)
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
