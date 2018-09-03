package app

import (
	"github.com/gin-contrib/cors"
	"github.com/jinzhu/gorm"
	"github.com/nickysemenza/food/backend/app/config"
	h "github.com/nickysemenza/food/backend/app/handler"
	"github.com/nickysemenza/food/backend/app/model"
	log "github.com/sirupsen/logrus"
)
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type App struct {
}

func (a *App) Initialize(c *config.Config) *config.Env {
	db, err := gorm.Open(c.DB.Dialect, c.GetDBURI())
	db.LogMode(true)
	if err != nil {
		log.Fatal("Could not connect database: ", c.GetDBURI())
	}
	//set up the env
	env := &config.Env{
		DB:   model.DBMigrate(db),
		Port: c.Port,
		//Router: &a.R,
	}
	return env
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
func (a *App) RunServer(host string, db *gorm.DB) {
	log.Println("Running API server on", host)

	router := SetupRouter(db)

	router.Run()
}
