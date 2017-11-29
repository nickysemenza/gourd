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
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
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

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()

	i := model.Ingredient{}
	db.FirstOrCreate(&i, model.Ingredient{Name: "flour"})

	test := model.Recipe{
		Title: "test",
		Slug:  "derp",
		Sections: []model.Section{
			{
				Ingredients: []model.SectionIngredient{
					{Grams: 4, Unit: "cups", Item: i},
				},
				Instructions: []model.SectionInstruction{
					{Name: "thing1"},
					{Name: "thing2"},
				},
			},
			{
				Instructions: []model.SectionInstruction{
					{Name: "thing3"},
					{Name: "thing4"},
				},
			}},
	}
	db.Create(&test)
}

// setRouters sets the all required routers
func (a *App) setRouters() {
	// Routing for handling the projects
	a.Get("/projects", a.GetAllProjects)
	a.Post("/projects", a.CreateProject)
	a.Get("/projects/{title}", a.GetProject)
	a.Put("/projects/{title}", a.UpdateProject)
	a.Delete("/projects/{title}", a.DeleteProject)
	a.Put("/projects/{title}/archive", a.ArchiveProject)
	a.Delete("/projects/{title}/archive", a.RestoreProject)

	a.Get("/api/recipes", a.GetAllRecipes)
	a.Get("/api/recipes/{slug}", a.GetRecipe)

	// Routing for handling the tasks
	a.Get("/projects/{title}/tasks", a.GetAllTasks)
	a.Post("/projects/{title}/tasks", a.CreateTask)
	a.Get("/projects/{title}/tasks/{id:[0-9]+}", a.GetTask)
	a.Put("/projects/{title}/tasks/{id:[0-9]+}", a.UpdateTask)
	a.Delete("/projects/{title}/tasks/{id:[0-9]+}", a.DeleteTask)
	a.Put("/projects/{title}/tasks/{id:[0-9]+}/complete", a.CompleteTask)
	a.Delete("/projects/{title}/tasks/{id:[0-9]+}/complete", a.UndoTask)
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

/*
** Recipe Handlers
 */
func (a *App) GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	handler.GetAllRecipes(a.DB, w, r)
}
func (a *App) GetRecipe(w http.ResponseWriter, r *http.Request) {
	handler.GetRecipe(a.DB, w, r)
}

/*
** Projects Handlers
 */
func (a *App) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	handler.GetAllProjects(a.DB, w, r)
}

func (a *App) CreateProject(w http.ResponseWriter, r *http.Request) {
	handler.CreateProject(a.DB, w, r)
}

func (a *App) GetProject(w http.ResponseWriter, r *http.Request) {
	handler.GetProject(a.DB, w, r)
}

func (a *App) UpdateProject(w http.ResponseWriter, r *http.Request) {
	handler.UpdateProject(a.DB, w, r)
}

func (a *App) DeleteProject(w http.ResponseWriter, r *http.Request) {
	handler.DeleteProject(a.DB, w, r)
}

func (a *App) ArchiveProject(w http.ResponseWriter, r *http.Request) {
	handler.ArchiveProject(a.DB, w, r)
}

func (a *App) RestoreProject(w http.ResponseWriter, r *http.Request) {
	handler.RestoreProject(a.DB, w, r)
}

/*
** Tasks Handlers
 */
func (a *App) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	handler.GetAllTasks(a.DB, w, r)
}

func (a *App) CreateTask(w http.ResponseWriter, r *http.Request) {
	handler.CreateTask(a.DB, w, r)
}

func (a *App) GetTask(w http.ResponseWriter, r *http.Request) {
	handler.GetTask(a.DB, w, r)
}

func (a *App) UpdateTask(w http.ResponseWriter, r *http.Request) {
	handler.UpdateTask(a.DB, w, r)
}

func (a *App) DeleteTask(w http.ResponseWriter, r *http.Request) {
	handler.DeleteTask(a.DB, w, r)
}

func (a *App) CompleteTask(w http.ResponseWriter, r *http.Request) {
	handler.CompleteTask(a.DB, w, r)
}

func (a *App) UndoTask(w http.ResponseWriter, r *http.Request) {
	handler.UndoTask(a.DB, w, r)
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
