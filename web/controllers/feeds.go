package controllers

import (
	"net/http"
	"strconv"

	"github.com/nerdyworm/go-flowfeeds/web/serializers"
	"github.com/gorilla/mux"
)

type FeedsController struct {
	ApplicationController
}

func (c *FeedsController) Index() error {
	feeds, err := c.Store.Feeds.List()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, serializers.NewFeeds(feeds))
}

func (c *FeedsController) Show() error {
	id, err := strconv.Atoi(mux.Vars(c.Request)["id"])
	if err != nil {
		return err
	}

	feed, err := c.Store.Feeds.Get(int64(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, serializers.NewShowFeed(feed))
}
