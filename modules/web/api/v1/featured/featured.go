package featured

import (
	"fmt"
	"log"
	"net/http"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/api/v1/serializers"
)

func Index(w http.ResponseWriter, r *http.Request) {
	serializer := serializers.FeaturedsSerializer{}
	serializer.Featureds = make([]models.Featured, 0)
	serializer.Teasers = make([]serializers.Teaser, 0)
	serializer.Feeds = make([]serializers.Feed, 0)

	teasers, feeds, err := models.FeaturedEpisodeTeasers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	for i, teaser := range teasers {
		serializer.Featureds = append(serializer.Featureds, models.Featured{Rank: i, Teaser: teaser.Id})
		serializer.Teasers = append(serializer.Teasers, serializers.NewTeaser(teaser))
	}

	for _, feed := range feeds {
		serializer.Feeds = append(serializer.Feeds, serializers.Feed{
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
