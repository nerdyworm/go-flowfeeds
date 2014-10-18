package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"bitbucket.org/nerdyworm/go-flowfeeds/datastore"
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/web/sessions"
	"github.com/codegangsta/controller"
)

type ApplicationController struct {
	controller.Base
	Store       *datastore.Datastore
	CurrentUser models.User
}

func (c *ApplicationController) Init(rw http.ResponseWriter, r *http.Request) error {
	c.Store = datastore.NewDatastore()

	var err error
	c.CurrentUser, err = sessions.CurrentUser(r, c.Store)
	if err != nil && err != sessions.NoCurrentUser {
		log.Printf("ApplicationController#Init sessions.GetCurrentUser %v\n", err)
		return err
	}

	return c.Base.Init(rw, r)
}

func (c *ApplicationController) JSON(status int, a interface{}) error {
	data, err := json.Marshal(a)
	if err != nil {
		log.Printf("ApplicationController#JSON json.Marshal %v\n", err)
		c.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return err
	}

	c.ResponseWriter.WriteHeader(status)
	c.ResponseWriter.Header().Set("content-type", "application/json; charset=utf-8")
	_, err = c.ResponseWriter.Write(data)
	return err
}

func (c *ApplicationController) Decode(target interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(target)
}
