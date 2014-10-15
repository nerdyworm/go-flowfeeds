package sessions

import (
	"errors"
	"net/http"
	"github.com/gorilla/sessions"

	"bitbucket.org/nerdyworm/go-flowfeeds/config"
	"bitbucket.org/nerdyworm/go-flowfeeds/datastore"
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
)

const (
	USER_SESSION_KEY = "user_id"
)

var (
	sessionStore = sessions.NewCookieStore(
		[]byte(config.SESSION_AUTH_KEY),
		[]byte(config.SESSION_CRYPT_KEY),
	)
)

func Signin(user *models.User, w http.ResponseWriter, r *http.Request) error {
	session, err := sessionStore.Get(r, "__flowfeeds_session")
	if err != nil {
		return err
	}

	session.Values[USER_SESSION_KEY] = user.Id
	return session.Save(r, w)
}

func Signout(w http.ResponseWriter, r *http.Request) error {
	session, err := sessionStore.Get(r, "__flowfeeds_session")
	if err != nil {
		return err
	}

	delete(session.Values, USER_SESSION_KEY)
	return session.Save(r, w)
}

func CurrentUser(r *http.Request, store *datastore.Datastore) (models.User, error) {
	session, err := sessionStore.Get(r, "__flowfeeds_session")
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

	return models.User{}, errors.New("No current user")
}
