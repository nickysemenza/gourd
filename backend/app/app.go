package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nickysemenza/food/backend/app/config"
	h "github.com/nickysemenza/food/backend/app/handler"
	"github.com/nickysemenza/food/backend/app/model"
	log "github.com/sirupsen/logrus"
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

func NewApp(c *config.Config) *App {
	db, err := gorm.Open(c.DB.Dialect, c.DB.GetURI())
	if err != nil {
		log.Fatal("Could not connect database: ", err)
	}
	db.LogMode(true)
	//set up the env
	env := &Env{
		DB: model.DBMigrate(db),
	}
	return &App{
		Env:  env,
		Port: c.Port,
	}
}

func DatabaseInjector(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}

func Authorized() gin.HandlerFunc {
	return func(c *gin.Context) {

		//check if token is present in header
		if jwt := c.GetHeader("X-Jwt"); jwt == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "no token")
		} else {
			//retrieve from DB
			db := c.MustGet("DB").(*gorm.DB)
			u, _ := h.GetUserFromToken(db, jwt)

			//check if they have admin role
			if u != nil && u.Admin == true {
				c.Set("user", u)
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "not authorized!")
			}

		}
	}
}

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-JWT"}
	router.Use(cors.New(corsConfig))
	router.Use(DatabaseInjector(db))

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

	router.GET("/me", Authorized(), h.GetMe)

	router.GET("/recipes", h.GetAllRecipes)
	router.POST("/recipes", Authorized(), h.CreateRecipe)
	router.GET("/recipes/:slug", h.GetRecipe)
	router.PUT("/recipes/:slug", Authorized(), h.PutRecipe)
	router.POST("/recipes/:slug/notes", Authorized(), h.AddNote)

	router.GET("/images", h.GetAllImages)
	router.PUT("/imageupload", Authorized(), h.PutImageUpload)

	router.GET("/categories", h.GetAllCategories)

	router.GET("/meals", h.GetAllMeals)
	router.GET("/meals/:id", h.GetMealByID)
	router.PUT("/meals/:id", Authorized(), h.UpdateMealByID)

	router.GET("/auth/facebook/login", h.HandleFacebookLogin)
	router.GET("/auth/facebook/callback", h.HandleFacebookCallback)

	return router
}
func (a *App) RunServer() error {
	host := fmt.Sprintf("%s:%s", a.Host, a.Port)
	log.Println("Running API server on", host)

	router := SetupRouter(a.Env.DB)

	if err := router.Run(); err != nil {
		return err
	}
	return nil
}
