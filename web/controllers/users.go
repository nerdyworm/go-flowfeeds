package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"bitbucket.org/nerdyworm/go-flowfeeds/datastore"
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/web/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/web/sessions"
	"github.com/gorilla/mux"
)

type UsersController struct {
	ApplicationController
}

func (c *UsersController) Create() error {
	createUserRequest := CreateUserRequest{}

	err := json.NewDecoder(c.Request.Body).Decode(&createUserRequest)
	if err != nil {
		log.Println("users.Create", err)
		return err
	}

	errors, err := createUserRequest.Validate(c.Store)
	if err != nil {
		return err
	}

	if errors.Any() {
		return c.JSON(422, errors)
	}

	params := createUserRequest.User

	user := models.NewUser(params.Email, params.Password)
	err = c.Store.Users.Insert(&user)
	if err != nil {
		log.Println("users.Create models.UserCreate", err)
		return err
	}

	err = sessions.Signin(user, c.ResponseWriter, c.Request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, serializers.NewShowUser(user))
}

func (c *UsersController) Show() error {
	id, err := strconv.Atoi(mux.Vars(c.Request)["id"])
	if err != nil {
		return err
	}

	user, err := c.Store.Users.Get(int64(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, serializers.NewShowUser(*user))
}

type CreateUserRequest struct {
	User struct {
		Email    string
		Password string
	}
}

func (r CreateUserRequest) Validate(store *datastore.Datastore) (models.ValidationErrors, error) {
	validationErrors := models.NewValidationErrors()
	params := r.User

	if params.Email == "" {
		validationErrors.Add("Email", "can't be blank")
	} else {
		exists, err := store.Users.Exists(params.Email)
		if err != nil {
			log.Println("users.Create models.UserExistsWithEmail", err)
			return validationErrors, err
		}

		if exists {
			validationErrors.Add("Email", "already taken")
		}
	}

	if params.Password == "" {
		validationErrors.Add("Password", "can't be blank")
	}

	return validationErrors, nil
}
