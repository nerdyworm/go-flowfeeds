package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/episodes"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/featured"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/feeds"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/listens"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/sessions"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/users"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/ctx"

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
	r.HandleFunc("/", Default(HomeHandler))
	r.NotFoundHandler = Default(HomeHandler)

	r.PathPrefix("/assets").Handler(http.FileServer(http.Dir(options.RootEmberAppPath)))

	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/episodes/{id}/favorites", Default(episodes.Favorites)).Methods("GET")
	apiRouter.HandleFunc("/episodes/{id}/listens", Default(episodes.Listens)).Methods("GET")
	apiRouter.HandleFunc("/episodes/{id}/related", Default(episodes.Related)).Methods("GET")
	apiRouter.HandleFunc("/episodes/{id}", Default(episodes.Show)).Methods("GET")
	apiRouter.HandleFunc("/episodes/{id}", Default(episodes.ToggleFavorite)).Methods("PUT")
	apiRouter.HandleFunc("/episodes", Default(episodes.Index)).Methods("GET")

	apiRouter.HandleFunc("/listens", Default(listens.Create)).Methods("POST")

	apiRouter.HandleFunc("/featureds", Default(featured.Index)).Methods("GET")
	apiRouter.HandleFunc("/feeds", Default(feeds.Index)).Methods("GET")
	apiRouter.HandleFunc("/feeds/{id}", Default(feeds.Show)).Methods("GET")
	apiRouter.HandleFunc("/users", Default(users.Create)).Methods("POST")
	apiRouter.HandleFunc("/users/{id}", Default(users.Show)).Methods("GET")
	apiRouter.HandleFunc("/sessions", Default(sessions.Create)).Methods("POST")
	apiRouter.HandleFunc("/sessions", Default(sessions.Destroy)).Methods("DELETE")
	http.Handle("/", r)

	log.Printf("Starting server on 3000\n")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}

type Manifest struct {
	CurrentUser int64
	Payload     struct {
		User serializers.User
	}
}

func HomeHandler(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "text/html")

	html := rootEmberIndexHtml

	if ctx.User.Id != 0 {
		m := Manifest{CurrentUser: ctx.User.Id}
		m.Payload.User = serializers.NewUser(ctx.User)

		b, err := json.Marshal(m)
		if err != nil {
			return err
		}

		scripts := fmt.Sprintf("<script>window.FlowfeedsManifest = %s;</script>\n</head>", b)
		html = strings.Replace(html, "</head>", scripts, 1)
	}

	fmt.Fprint(w, html)
	return nil
}
