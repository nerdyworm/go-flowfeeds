package web

import (
	"log"
	"net/http"
	"os"

	"bitbucket.org/nerdyworm/go-flowfeeds/config"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/controllers"
	"github.com/codegangsta/controller"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func Run() {
	if _, err := os.Stat(config.EmberApp + "/index.html"); os.IsNotExist(err) {
		log.Printf("could not find the ember app's index.html in `%s`", config.EmberApp)
		log.Fatal(err)
	}

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
	)

	r := mux.NewRouter()
	r.Handle("/", controller.Action((*controllers.HomeController).Index))
	r.NotFoundHandler = controller.Action((*controllers.HomeController).Index)

	r.PathPrefix("/assets").Handler(http.FileServer(http.Dir(config.EmberApp)))

	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	apiRouter.Handle("/episodes/{id}/favorites", controller.Action((*controllers.EpisoidesController).Favorites)).Methods("GET")
	apiRouter.Handle("/episodes/{id}/listens", controller.Action((*controllers.EpisoidesController).Listens)).Methods("GET")
	apiRouter.Handle("/episodes/{id}/related", controller.Action((*controllers.EpisoidesController).Related)).Methods("GET")
	apiRouter.Handle("/episodes/{id}", controller.Action((*controllers.EpisoidesController).Show)).Methods("GET")
	apiRouter.Handle("/episodes/{id}", controller.Action((*controllers.EpisoidesController).Update)).Methods("PUT")
	apiRouter.Handle("/episodes", controller.Action((*controllers.EpisoidesController).Index)).Methods("GET")

	apiRouter.Handle("/listens", controller.Action((*controllers.ListensController).Create)).Methods("POST")

	apiRouter.Handle("/feeds", controller.Action((*controllers.FeedsController).Index)).Methods("GET")
	apiRouter.Handle("/feeds/{id}", controller.Action((*controllers.FeedsController).Show)).Methods("GET")

	apiRouter.Handle("/users", controller.Action((*controllers.UsersController).Create)).Methods("POST")
	apiRouter.Handle("/users/{id}", controller.Action((*controllers.UsersController).Show)).Methods("POST")

	apiRouter.Handle("/sessions", controller.Action((*controllers.SessionsController).Create)).Methods("POST")
	apiRouter.Handle("/sessions", controller.Action((*controllers.SessionsController).Delete)).Methods("DELETE")
	http.Handle("/", r)

	n.UseHandler(r)
	n.Run(":3000")
}
