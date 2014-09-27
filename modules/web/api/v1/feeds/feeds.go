package feeds

import (
	"net/http"
	"strconv"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/ctx"
	"github.com/gorilla/mux"
)

func Index(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	serializer := serializers.FeedsSerializer{}
	serializer.Feeds = make([]serializers.Feed, 0)

	feeds, err := models.Feeds()
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		serializer.Feeds = append(serializer.Feeds, serializers.NewFeed(feed))
	}

	serializers.JSON(w, serializer)
	return nil
}

func Show(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	feed, err := models.FindFeedById(int64(id))
	if err != nil {
		return err
	}

	serializer := serializers.FeedShowSerializer{Feed: feed}
	serializers.JSON(w, serializer)
	return nil
}
