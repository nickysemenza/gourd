package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nickysemenza/food/backend/app/config"
	h "github.com/nickysemenza/food/backend/app/handler"
	"github.com/nickysemenza/food/backend/app/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	log "github.com/sirupsen/logrus"
	otgorm "github.com/smacker/opentracing-gorm"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"

	"github.com/zsais/go-gin-prometheus"
)

//App is the app server
type App struct {
	Env  *Env
	Port string
	Host string
}

//Env holds misc env stuff like the DB connection object.
type Env struct {
	DB          *gorm.DB
	CurrentUser *model.User
}

//NewApp creates a new app
func NewApp(c *config.Config) *App {
	db, err := gorm.Open(c.DB.Dialect, c.DB.GetURI())
	if err != nil {
		log.Fatal("Could not connect database: ", err)
	}
	db.LogMode(true)
	otgorm.AddGormCallbacks(db)
	//set up the env
	env := &Env{
		DB: model.DBMigrate(db),
	}
	return &App{
		Env:  env,
		Port: c.Port,
	}
}

func databaseInjector(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		ctx = context.WithValue(ctx, model.DBKey, db)
		c.Set("ctx", ctx)
		c.Next()
	}
}

func mustAuthorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		//check if token is present in header
		span, ctx := opentracing.StartSpanFromContext(ctx, "mustAuthorize")
		span.LogEvent("check auth")
		defer span.Finish()
		c.Set("tracing-context", span)

		if jwt := c.GetHeader("X-Jwt"); jwt == "" {
			span.SetTag("user-auth-result", "guest")
			c.AbortWithStatusJSON(http.StatusUnauthorized, "no token")
		} else {
			//retrieve from DB
			db := model.GetDBFromContext(ctx)
			u, _ := h.GetUserFromToken(db, jwt)

			span.SetTag("user-id", u.ID)
			//check if they have admin role
			if u != nil && u.Admin == true {
				c.Set("user", u)
				span.SetTag("user-auth-result", "ok")
				c.Next()

			} else {
				span.SetTag("user-auth-result", "not-authorized")
				c.AbortWithStatusJSON(http.StatusUnauthorized, "not authorized!")

			}

		}
	}
}

func tracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		span, ctx := opentracing.StartSpanFromContext(ctx, "request: "+c.Request.Method+" "+c.Request.URL.Path)
		span.SetTag(string(ext.HTTPMethod), c.Request.Method)
		span.SetTag(string(ext.HTTPUrl), c.Request.URL.Path)
		span.LogEvent("begin")
		defer span.SetTag(string(ext.HTTPStatusCode), c.Writer.Status())
		defer span.Finish()
		c.Set("tracing-context", span)
		c.Set("ctx", ctx)
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			c.Writer.Header().Set("x-trace-id", sc.TraceID().String())
		}

		c.Next()
	}
}

//SetupRouter builds the router
func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	//proper cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-JWT"}
	router.Use(cors.New(corsConfig))

	//inject DB into ctx
	router.Use(tracingMiddleware())
	router.Use(databaseInjector(db))

	// prometheus setup
	p := ginprometheus.NewPrometheus("gin")
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.String()
		for _, p := range c.Params {
			if p.Key == "slug" {
				url = strings.Replace(url, p.Value, ":slug", 1)
				break
			} else if p.Key == "id" {
				url = strings.Replace(url, p.Value, ":id", 1)
				break
			}
		}
		return url
	}
	p.Use(router)

	// routes
	router.Static("/public", "./public")

	router.GET("/me", mustAuthorize(), h.GetMe)

	router.GET("/recipes", h.GetAllRecipes)
	router.POST("/recipes", mustAuthorize(), h.CreateRecipe)
	router.GET("/recipes/:slug", h.GetRecipe)
	router.PUT("/recipes/:slug", mustAuthorize(), h.PutRecipe)
	router.POST("/recipes/:slug/notes", mustAuthorize(), h.AddNote)

	router.GET("/images", h.GetAllImages)
	router.PUT("/imageupload", mustAuthorize(), h.PutImageUpload)

	router.GET("/categories", h.GetAllCategories)

	router.GET("/meals", h.GetAllMeals)
	router.GET("/meals/:id", h.GetMealByID)
	router.PUT("/meals/:id", mustAuthorize(), h.UpdateMealByID)

	router.GET("/auth/facebook/login", h.HandleFacebookLogin)
	router.GET("/auth/facebook/callback", h.HandleFacebookCallback)

	return router
}

//RunServer starts up the server on host:port
func (a *App) RunServer() error {
	host := fmt.Sprintf("%s:%s", a.Host, a.Port)
	log.Println("Running API server on", host)

	initTracing()
	router := SetupRouter(a.Env.DB)

	if err := router.Run(); err != nil {
		return err
	}
	return nil
}
func initTracing() {
	// sender := transport.NewHTTPTransport(
	// 	"localhost:6831",
	// 	transport.HTTPBatchSize(1),
	// )

	sender, err := jaeger.NewUDPTransport("localhost:6831", 0)
	if err != nil {
		log.Fatal(err)
	}
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1.0, // sample all traces
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, _, _ := cfg.New("food-api-backend",
		jaegercfg.Reporter(jaeger.NewRemoteReporter(
			sender,
			jaeger.ReporterOptions.BufferFlushInterval(1*time.Second),
			jaeger.ReporterOptions.Logger(jaegerlog.StdLogger),
		)))

	// defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
}
