package app

import (
	//"github.com/gorilla/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/nickysemenza/food/backend/app/config"
	h "github.com/nickysemenza/food/backend/app/handler"
	"github.com/nickysemenza/food/backend/app/model"
	log "github.com/sirupsen/logrus"
	"net/http"
)
import "github.com/gin-gonic/gin"

type App struct {
	R *mux.Router
}
type Route struct {
	Method      string
	Pattern     string
	Protected   bool
	HandlerFunc func(e *config.Env, w http.ResponseWriter, r *http.Request) error
}

func (a *App) Initialize(c *config.Config) *config.Env {
	db, err := gorm.Open(c.DB.Dialect, c.GetDBURI())
	if err != nil {
		log.Fatal("Could not connect database")
	}
	//set up the env
	env := &config.Env{
		DB:     model.DBMigrate(db),
		Port:   c.Port,
		Router: &a.R,
	}
	a.buildRoutes(env)
	return env
}

type Routes []Route

func (a *App) buildRoutes(env *config.Env) {

}
func DatabaseInjector(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}

func (a *App) RunServer(host string, db *gorm.DB) {
	log.Println("Running API server on", host)
	//headersOk := handlers.AllowedHeaders([]string{"x-jwt"})
	//originsOk := handlers.AllowedOrigins([]string{"*"})
	//methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	//
	//log.Fatal(http.ListenAndServe(host, handlers.CORS(originsOk, headersOk, methodsOk)(a.R)))
	router := gin.Default()
	router.Use(DatabaseInjector(db))

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-JWT"}
	router.Use(cors.New(config))

	router.GET("/recipes", h.GetAllRecipes)

	router.Handle("GET", "/me", h.GetMe)                     //todo: protect
	router.Handle("PUT", "/imageupload", h.PutImageUpload)   //todo: protect
	router.Handle("POST", "/recipes", h.CreateRecipe)        //todo: protect
	router.Handle("PUT", "/recipes/:slug", h.PutRecipe)      //todo: protect
	router.Handle("POST", "/recipes/:slug/notes", h.AddNote) //todo: protect
	router.Handle("PUT", "/meals/:id", h.UpdateMealByID)     //todo: protect
	router.Handle("GET", "/recipes/:slug", h.GetRecipe)
	router.Handle("GET", "/images", h.GetAllImages)
	router.Handle("GET", "/categories", h.GetAllCategories)
	router.Handle("GET", "/meals", h.GetAllMeals)
	router.Handle("GET", "/meals/:id", h.GetMealByID)
	//router.Handle("GET", "/auth/facebook/login", h.HandleFacebookLogin)
	//router.Handle("GET", "/auth/facebook/callback", h.HandleFacebookCallback)

	router.Run()
}
