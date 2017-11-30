package app

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/nickysemenza/food/backend/app/handler"
	"github.com/nickysemenza/food/backend/app/model"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize(config *Config) *handler.Env {
	db, err := gorm.Open(config.DB.Dialect, config.DB.URI)
	if err != nil {
		log.Fatal("Could not connect database")
	}
	//set up the env
	env := &handler.Env{
		DB:   model.DBMigrate(db),
		Port: config.Port,
	}

	a.Router = mux.NewRouter()
	a.setRouters(env)
	return env
}

func (a *App) setRouters(env *handler.Env) {
	a.Router.Handle("/api", handler.Handler{env, handler.ErrorTest}).Methods("GET")
	a.Router.Handle("/api/recipes", handler.Handler{env, handler.GetAllRecipes}).Methods("GET")
	a.Router.Handle("/api/recipes/{slug}", handler.Handler{env, handler.GetRecipe}).Methods("GET")

}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
