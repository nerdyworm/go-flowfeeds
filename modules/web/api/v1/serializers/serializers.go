package serializers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
)

type EpisodesSerializer struct {
	Episodes []models.Episode
}

type EpisodeSerializer struct {
	Episode struct {
		models.Episode
		Comments []int64
	}
	Comments []models.Comment
}

type FeaturedsSerializer struct {
	Featureds []models.Featured
	Teasers   []models.Teaser
}

type FeedsSerializer struct {
	Feeds []FeedSerializer
}

type FeedSerializer struct {
	Id          int64
	Title       string
	Description string
	Url         string
	Thumb       string
	Cover       string
	Updated     time.Time
}

type FeedShowSerializer struct {
	Feed models.Feed
}

func JSON(w http.ResponseWriter, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ERROR JSON MarshalIndent %v\n", err)
		return err
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	_, err = w.Write(data)
	return err
}
