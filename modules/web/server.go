package web

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/episodes"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/featured"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/feeds"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/users"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	rootEmberIndexHtml = ""
)

type ServerOptions struct {
	RootEmberAppPath string
}

func Run(options ServerOptions) {
	index, err := ioutil.ReadFile(options.RootEmberAppPath + "/index.html")
	if err != nil {
		log.Printf("could not find the ember app's index.html in `%s`", options.RootEmberAppPath)
		log.Fatal(err)
	}

	rootEmberIndexHtml = string(index)

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.NotFoundHandler = http.HandlerFunc(HomeHandler)

	r.PathPrefix("/assets").Handler(http.FileServer(http.Dir(options.RootEmberAppPath)))

	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/episodes", episodes.Index).Methods("GET")
	apiRouter.HandleFunc("/episodes/{id}", episodes.Show).Methods("GET")
	apiRouter.HandleFunc("/featureds", featured.Index).Methods("GET")
	apiRouter.HandleFunc("/feeds", feeds.Index).Methods("GET")
	apiRouter.HandleFunc("/feeds/{id}", feeds.Show).Methods("GET")
	apiRouter.HandleFunc("/users", users.Create).Methods("POST")
	http.Handle("/", r)

	log.Printf("Starting server...\n")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, rootEmberIndexHtml)
}
