package app

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/nickysemenza/food/backend/app/handler"
	"github.com/nickysemenza/food/backend/app/model"
	"github.com/nickysemenza/food/backend/config"
	"log"
	"net/http"
	"os"
)

// App has router and db instances
type App struct {
	Router *mux.Router
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config) {
	db, err := gorm.Open(config.DB.Dialect, config.DB.URI)
	if err != nil {
		log.Fatal("Could not connect database")
	}
	db = model.DBMigrate(db)

	env := &handler.Env{}
	env.DB = db

	a.Router = mux.NewRouter()
	a.setRouters(env)

	pwd, _ := os.Getwd()
	pwd += "/recipes/"
	//utils.Utils{env}.Export(pwd)
}

// setRouters sets the all required routers
func (a *App) setRouters(env *handler.Env) {
	a.Router.Handle("/api", handler.Handler{env, handler.ErrorTest}).Methods("GET")
	a.Router.Handle("/api/recipes", handler.Handler{env, handler.GetAllRecipes}).Methods("GET")
	a.Router.Handle("/api/recipes/{slug}", handler.Handler{env, handler.GetRecipe}).Methods("GET")

}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
