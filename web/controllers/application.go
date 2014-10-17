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
	if err != nil {
		log.Println("handlers.Default session.CurrentUser", err)
	} else {
		log.Printf("current user: (%d) %s\n", c.CurrentUser.Id, c.CurrentUser.Email)
	}

	return c.Base.Init(rw, r)
}

func (c *ApplicationController) JSON(status int, a interface{}) error {
	c.ResponseWriter.WriteHeader(status)
	c.ResponseWriter.Header().Set("content-type", "application/json; charset=utf-8")
	data, err := json.Marshal(a)
	if err != nil {
		c.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		log.Printf("ERROR JSON MarshalIndent %v\n", err)
		return err
	}

	_, err = c.ResponseWriter.Write(data)
	return err
}

func (c *ApplicationController) Decode(target interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(target)
}
