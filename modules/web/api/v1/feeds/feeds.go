package feeds

import (
	"log"
	"net/http"
	"strconv"

	"bitbucket.org/nerdyworm/go-flowfeeds/datastore"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/sessions"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/ctx"
	"github.com/codegangsta/controller"
	"github.com/gorilla/mux"
)

type ApplicationController struct {
	controller.Base
	Context ctx.Context
	Store   *datastore.Datastore
}

func (c *ApplicationController) Init(rw http.ResponseWriter, r *http.Request) error {
	context := ctx.Context{}
	context.Store = datastore.NewDatastore()

	user, err := sessions.CurrentUser(r, context.Store)
	if err != nil {
		log.Println("handlers.Default session.CurrentUser", err)
	} else {
		context.User = user
		log.Printf("current user: %d\n", user.Id)
	}
	c.Context = context
	c.Store = context.Store

	return c.Base.Init(rw, r)
}

type FeedsController struct {
	ApplicationController
}

func (c *FeedsController) Index() error {
	serializer := serializers.FeedsSerializer{}
	serializer.Feeds = make([]serializers.Feed, 0)

	feeds, err := c.Store.Feeds.List()
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		serializer.Feeds = append(serializer.Feeds, serializers.NewFeed(*feed))
	}

	return serializers.JSON(c.ResponseWriter, serializer)
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

	return serializers.JSON(c.ResponseWriter, serializers.NewShowFeed(feed))
}
