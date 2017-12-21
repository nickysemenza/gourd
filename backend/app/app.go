package app

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
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
	HandlerFunc func(e *h.Env, w http.ResponseWriter, r *http.Request) error
}

type Routes []Route

func (a *App) Initialize(config *Config) *h.Env {
	db, err := gorm.Open(config.DB.Dialect, config.getDBURI())
	if err != nil {
		log.Fatal("Could not connect database")
	}
	//set up the env
	env := &h.Env{
		DB:     model.DBMigrate(db),
		Port:   config.Port,
		Router: &a.R,
	}
	a.buildRoutes(env)
	return env
}

func (a *App) buildRoutes(env *h.Env) {

	var routes = Routes{
		{"GET", "/api", h.ErrorTest},
		{"GET", "/api/recipes", h.GetAllRecipes},
		{"GET", "/api/recipes/{slug}", h.GetRecipe},
		{"PUT", "/api/recipes/{slug}", h.PutRecipe},
		{"POST", "/api/recipes/{slug}/notes", h.AddNote},
	}

	//add them all
	a.R = mux.NewRouter()
	for _, route := range routes {
		a.R.Handle(route.Pattern, h.Handler{env, route.HandlerFunc}).Methods(route.Method)
	}
	a.R.NotFoundHandler = h.Handler{env, h.NotFoundRoute}
}

func (a *App) RunServer(host string) {
	log.Println("Running API server on", host)
	log.Fatal(http.ListenAndServe(host, a.R))
}
