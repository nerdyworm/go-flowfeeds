package controllers

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/sessions"
)

type SessionsController struct {
	ApplicationController
}

type CreateSessionRequest struct {
	Session struct {
		Email    string
		Password string
	}
}

func (c *SessionsController) Create() error {
	createSessionRequest := CreateSessionRequest{}
	err := json.NewDecoder(c.Request.Body).Decode(&createSessionRequest)
	if err != nil {
		return err
	}
	params := createSessionRequest.Session

	if params.Email == "" || params.Password == "" {
		c.ResponseWriter.WriteHeader(422)
		return nil
	}

	user, err := c.Store.Users.FindForSignin(params.Email)
	if err == models.ErrNotFound {
		c.ResponseWriter.WriteHeader(422)
		return nil
	}

	if err != nil {
		return err
	}

	err = user.CheckPassword(params.Password)
	if err != nil {
		c.ResponseWriter.WriteHeader(422)
		return nil
	}

	err = sessions.Signin(*user, c.ResponseWriter, c.Request)
	if err != nil {
		return err
	}

	serializer := serializers.ShowUser{
		serializers.NewUser(*user),
	}

	return c.JSON(http.StatusCreated, serializer)
}

func (c *SessionsController) Delete() error {
	err := sessions.Signout(c.ResponseWriter, c.Request)
	if err != nil {
		http.Error(c.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return err
	}

	c.ResponseWriter.WriteHeader(http.StatusAccepted)
	return nil
}
