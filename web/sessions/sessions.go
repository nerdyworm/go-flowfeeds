package sessions

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/nerdyworm/go-flowfeeds/config"
	"github.com/nerdyworm/go-flowfeeds/datastore"
	"github.com/nerdyworm/go-flowfeeds/models"
)

const (
	USER_SESSION_KEY = "user_id"
	SESSION_KEY      = "__flowfeeds_session"
)

var (
	sessionStore = sessions.NewCookieStore(
		[]byte(config.SESSION_AUTH_KEY),
		[]byte(config.SESSION_CRYPT_KEY),
	)

	NoCurrentUser = errors.New("No current user")
)

func Signin(user *models.User, w http.ResponseWriter, r *http.Request) error {
	session, err := sessionStore.Get(r, SESSION_KEY)
	if err != nil {
		return err
	}

	session.Values[USER_SESSION_KEY] = user.Id
	return session.Save(r, w)
}

func Signout(w http.ResponseWriter, r *http.Request) error {
	http.SetCookie(w, &http.Cookie{Name: SESSION_KEY, MaxAge: -1, Path: "/"})
	return nil
}

func CurrentUser(r *http.Request, store *datastore.Datastore) (models.User, error) {
	session, err := sessionStore.Get(r, SESSION_KEY)
	if err != nil {
		return models.User{}, err
	}

	if id, ok := session.Values[USER_SESSION_KEY]; ok {
		user, err := store.Users.Get(id.(int64))
		if err == models.ErrNotFound {
			return models.User{}, nil
		}

		if err != nil {
			return models.User{}, err
		}

		return *user, err

	}

	return models.User{}, NoCurrentUser
}
