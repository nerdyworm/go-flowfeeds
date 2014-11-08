package controllers

import (
	"net/http"
	"strconv"

	"bitbucket.org/nerdyworm/go-flowfeeds/web/serializers"
)

type ListensController struct {
	ApplicationController
}

type ListenParams struct {
	Listen struct {
		Episode string
	}
}

func (c *ListensController) Create() error {
	if c.CurrentUser.IsAnonymous() {
		return nil
	}

	id, err := c.getEpisodeId()
	if err != nil {
		return err
	}

	listen, err := c.Store.Listens.Create(&c.CurrentUser, id)
	if err != nil {
		return err
	}

	episode, err := c.Store.Episodes.GetForUser(&c.CurrentUser, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, serializers.NewShowListen(listen, episode))
}

func (c *ListensController) getEpisodeId() (int64, error) {
	params := ListenParams{}

	err := c.Decode(&params)
	if err != nil {
		return 0, err
	}

	id, err := strconv.Atoi(params.Listen.Episode)
	if err != nil {
		return 0, err
	}

	return int64(id), nil
}
