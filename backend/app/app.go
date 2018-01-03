package app

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/nickysemenza/food/backend/app/config"
	h "github.com/nickysemenza/food/backend/app/handler"
	"github.com/nickysemenza/food/backend/app/model"
	"log"
	"net/http"
)

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
		DB:   model.DBMigrate(db),
		Port: c.Port,
		//Router: &a.R,
	}
	a.buildRoutes(env)
	return env
}

type Routes []Route

func (a *App) buildRoutes(env *config.Env) {

	var routes = Routes{
		{Method: "GET", Pattern: "/", HandlerFunc: h.ErrorTest, Protected: false},
		{Method: "GET", Pattern: "/me", HandlerFunc: h.GetMe, Protected: true},

		{Method: "GET", Pattern: "/recipes", HandlerFunc: h.GetAllRecipes, Protected: false},
		{Method: "POST", Pattern: "/recipes", HandlerFunc: h.CreateRecipe, Protected: true},
		{Method: "GET", Pattern: "/recipes/{slug}", HandlerFunc: h.GetRecipe, Protected: false},
		{Method: "PUT", Pattern: "/recipes/{slug}", HandlerFunc: h.PutRecipe, Protected: true},
		{Method: "POST", Pattern: "/recipes/{slug}/notes", HandlerFunc: h.AddNote, Protected: true},

		{Method: "GET", Pattern: "/images", HandlerFunc: h.GetAllImages, Protected: false},
		{Method: "PUT", Pattern: "/imageupload", HandlerFunc: h.PutImageUpload, Protected: true},

		{Method: "GET", Pattern: "/categories", HandlerFunc: h.GetAllCategories, Protected: false},
		{Method: "GET", Pattern: "/meals", HandlerFunc: h.GetAllMeals, Protected: false},

		{Method: "GET", Pattern: "/auth/facebook/login", HandlerFunc: h.HandleFacebookLogin, Protected: false},
		{Method: "GET", Pattern: "/auth/facebook/callback", HandlerFunc: h.HandleFacebookCallback, Protected: false},
	}

	//add them all
	a.R = mux.NewRouter()
	for _, route := range routes {
		a.R.Handle(route.Pattern, h.Handler{Env: env, H: route.HandlerFunc, P: route.Protected}).Methods(route.Method)
	}
	a.R.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	a.R.NotFoundHandler = h.Handler{Env: env, H: h.NotFoundRoute}
}

func (a *App) RunServer(host string) {
	log.Println("Running API server on", host)
	headersOk := handlers.AllowedHeaders([]string{"x-jwt"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Fatal(http.ListenAndServe(host, handlers.CORS(originsOk, headersOk, methodsOk)(a.R)))
}
