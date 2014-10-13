package sessions

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/sessions"

	"bitbucket.org/nerdyworm/go-flowfeeds/datastore"
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/ctx"
)

const (
	USER_SESSION_KEY = "user_id"
)

var (
	sessionStore = sessions.NewCookieStore([]byte("1234568900"))
)

type CreateSessionRequest struct {
	Session struct {
		Email    string
		Password string
	}
}

func Create(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	createSessionRequest := CreateSessionRequest{}
	err := json.NewDecoder(r.Body).Decode(&createSessionRequest)
	if err != nil {
		log.Println("sessions.Create json.Decode", err)
		return err
	}
	params := createSessionRequest.Session

	if params.Email == "" || params.Password == "" {
		w.WriteHeader(422)
		return errors.New("Invalid email or password")
	}

	user, err := models.FindUserForSignin(params.Email)
	if err == models.ErrNotFound {
		w.WriteHeader(422)
		return errors.New("Invalid email or password")
	}

	if err != nil {
		return err
	}

	err = user.CheckPassword(params.Password)
	if err != nil {
		w.WriteHeader(422)
		return errors.New("Invalid email or password")
	}

	err = signin(user, w, r)
	if err != nil {
		return err
	}

	serializer := serializers.ShowUser{
		serializers.NewUser(user),
	}

	w.WriteHeader(http.StatusCreated)
	serializers.JSON(w, serializer)
	return nil
}

func Destroy(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	err := signout(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusAccepted)
	return nil
}

func signin(user models.User, w http.ResponseWriter, r *http.Request) error {
	session, err := sessionStore.Get(r, "__flowfeeds_session")
	if err != nil {
		return err
	}

	session.Values[USER_SESSION_KEY] = user.Id
	return session.Save(r, w)
}

func signout(w http.ResponseWriter, r *http.Request) error {
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
