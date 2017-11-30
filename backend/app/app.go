package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/nickysemenza/food/backend/app/handler"
	"github.com/nickysemenza/food/backend/app/model"
	"github.com/nickysemenza/food/backend/config"
	"os"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	db     *gorm.DB
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("%s:%s@/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal("Could not connect database")
	}

	a.db = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()

	pwd, _ := os.Getwd()
	pwd += "/recipes/"
	//model.LegacyImport(db, pwd)
	//model.Export(db, pwd)

}

// setRouters sets the all required routers
func (a *App) setRouters() {
	a.Router.HandleFunc("/api/recipes", a.GetAllRecipes).Methods("GET")
	a.Router.HandleFunc("/api/recipes/{slug}", a.GetRecipe).Methods("GET")

}

/*
** Recipe Handlers
 */
func (a *App) GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	handler.GetAllRecipes(a.db, w, r)
}
func (a *App) GetRecipe(w http.ResponseWriter, r *http.Request) {
	handler.GetRecipe(a.db, w, r)
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
