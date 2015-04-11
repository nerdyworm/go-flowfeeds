package controllers

import (
	"net/http"
	"strconv"

	"github.com/nerdyworm/go-flowfeeds/datastore"
	"github.com/nerdyworm/go-flowfeeds/models"
	"github.com/nerdyworm/go-flowfeeds/web/serializers"
	"github.com/gorilla/mux"
)

type EpisoidesController struct {
	ApplicationController
}

func (c *EpisoidesController) Index() error {
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	feed, _ := strconv.Atoi(c.Request.URL.Query().Get("feed"))
	baseOptions := datastore.ListOptions{PerPage: 24, Page: page}
	options := datastore.EpisodeListOptions{baseOptions, int64(feed)}
	episodes, feeds, err := c.Store.Episodes.ListFor(&c.CurrentUser, options)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, serializers.NewEpisodes(episodes, feeds, options))
}

func (c *EpisoidesController) Show() error {
	episode, feed, err := c.getEpisode()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, serializers.NewShowEpisode(episode, feed))
}

func (c *EpisoidesController) Update() error {
	if c.CurrentUser.IsAnonymous() {
		c.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	id, err := c.getId()
	if err != nil {
		return err
	}

	err = c.Store.Episodes.ToggleFavoriteForUser(&c.CurrentUser, id)
	if err != nil {
		return err
	}

	episode, feed, err := c.getEpisode()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, serializers.NewShowEpisode(episode, feed))
}

func (c *EpisoidesController) Related() error {
	id, err := c.getId()
	if err != nil {
		return err
	}

	related, feeds, err := c.Store.Episodes.Related(id)
	if err != nil {
		return err
	}

	episodes := datastore.Episodes{Episodes: related}
	return c.JSON(http.StatusOK, serializers.NewEpisodes(episodes, feeds, datastore.EpisodeListOptions{}))
}

func (c *EpisoidesController) Listens() error {
	id, err := c.getId()
	if err != nil {
		return err
	}

	listens, users, err := c.Store.Episodes.Listens(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, serializers.NewListens(listens, users))
}

func (c *EpisoidesController) Favorites() error {
	id, err := c.getId()
	if err != nil {
		return err
	}

	favorites, users, err := c.Store.Episodes.Favorites(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, serializers.NewFavorites(favorites, users))
}

func (c *EpisoidesController) getId() (int64, error) {
	id, err := strconv.Atoi(mux.Vars(c.Request)["id"])
	if err != nil {
		return 0, err
	}

	return int64(id), nil
}

func (c *EpisoidesController) getEpisode() (*models.Episode, *models.Feed, error) {
	id, err := c.getId()
	if err != nil {
		return nil, nil, err
	}

	episode, err := c.Store.Episodes.GetForUser(&c.CurrentUser, id)
	if err != nil {
		return nil, nil, err
	}

	feed, err := c.Store.Feeds.Get(episode.Feed)
	if err != nil {
		return nil, nil, err
	}

	return episode, feed, nil
}
