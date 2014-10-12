package feeds

import (
	"net/http"
	"strconv"

	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/ctx"
	"github.com/gorilla/mux"
)

// TODO: remove serilizer stuff from here into serializers
func Index(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	serializer := serializers.FeedsSerializer{}
	serializer.Feeds = make([]serializers.Feed, 0)

	feeds, err := ctx.Store.Feeds.List()
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		serializer.Feeds = append(serializer.Feeds, serializers.NewFeed(*feed))
	}

	serializers.JSON(w, serializer)
	return nil
}

func Show(ctx ctx.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	feed, err := ctx.Store.Feeds.Get(int64(id))
	if err != nil {
		return err
	}

	serializer := serializers.FeedShowSerializer{Feed: *feed}
	serializers.JSON(w, serializer)
	return nil
}
