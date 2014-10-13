package feeds

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"bitbucket.org/nerdyworm/go-flowfeeds/datastore"
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/sessions"
	"github.com/codegangsta/controller"
	"github.com/gorilla/mux"
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
	return serializers.JSON(c.ResponseWriter, a)
}

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

type ListensController struct {
	ApplicationController
}

type ListenParams struct {
	Listen struct {
		Episode string
	}
}

func (c *ListensController) Create() error {
	if c.CurrentUser.Id == 0 {
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

	return c.JSON(http.StatusCreated, serializers.NewShowListen(*listen, *episode))
}

func (c *ListensController) getEpisodeId() (int64, error) {
	params := ListenParams{}

	err := json.NewDecoder(c.Request.Body).Decode(&params)
	if err != nil {
		return 0, err
	}

	id, err := strconv.Atoi(params.Listen.Episode)
	if err != nil {
		return 0, err
	}

	return int64(id), nil
}

type EpisoidesController struct {
	ApplicationController
}

func (c *EpisoidesController) Index() error {
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	options := models.ListOptions{PerPage: 24, Page: page}
	episodes, feeds, err := c.Store.Episodes.ListFor(&c.CurrentUser, options)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, serializers.NewEpisodes(episodes, feeds))
}

func (c *EpisoidesController) Show() error {
	episode, feed, err := c.getEpisode()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, serializers.NewShowEpisode(episode, feed))
}

func (c *EpisoidesController) Update() error {
	if c.CurrentUser.Id == 0 {
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

	return c.JSON(http.StatusOK, serializers.NewEpisodes(related, feeds))
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

	feed, err := c.Store.Feeds.Get(episode.FeedId)
	if err != nil {
		return nil, nil, err
	}

	return episode, feed, nil
}
