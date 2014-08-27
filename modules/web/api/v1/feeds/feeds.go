package feeds

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	serializer := serializers.FeedsSerializer{}
	serializer.Feeds = make([]serializers.FeedSerializer, 0)

	feeds, err := models.Feeds()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	for _, feed := range feeds {
		serializer.Feeds = append(serializer.Feeds, serializers.FeedSerializer{
			Id:          feed.Id,
			Title:       feed.Title,
			Description: feed.Description,
			Url:         feed.Url,
			Thumb:       fmt.Sprintf("http://s3.amazonaws.com/flowfeeds2/feeds/%d/thumb.jpg", feed.Id),
			Cover:       fmt.Sprintf("http://s3.amazonaws.com/flowfeeds2/feeds/%d/cover.jpg", feed.Id),
		})
	}

	serializers.JSON(w, serializer)
}

func Show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Fatal(err)
	}

	feed, err := models.FindFeedById(int64(id))
	if err != nil {
		log.Fatal(err)
	}

	serializer := serializers.FeedShowSerializer{Feed: feed}
	serializers.JSON(w, serializer)
}
