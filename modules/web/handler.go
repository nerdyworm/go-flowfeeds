package web

import (
	"log"
	"net/http"
	"time"

	"bitbucket.org/nerdyworm/go-flowfeeds/datastore"
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/sessions"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/ctx"
)

type ApplicationHandler func(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error

func Default(handler ApplicationHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		addr := r.Header.Get("X-Real-IP")
		if addr == "" {
			addr = r.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = r.RemoteAddr
			}
		}

		context := ctx.Context{}
		context.Store = datastore.NewDatastore()

		user, err := sessions.CurrentUser(r, context.Store)
		if err != nil {
			log.Println("handlers.Default session.CurrentUser", err)
		} else {
			context.User = user
			log.Printf("current user: %d\n", user.Id)
		}

		log.Printf("%s %s %s\n", r.Method, r.URL.Path, addr)

		err = handler(context, w, r)
		if err != nil {
			handleError(err, w, r)
			log.Printf("Errored %v", err)
		} else {
			log.Printf("Completed runtime %v\n", time.Since(start))
		}
	}
}

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	switch err := err.(type) {
	case models.ValidationErrors:
		w.WriteHeader(422)
		serializers.JSON(w, err)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
