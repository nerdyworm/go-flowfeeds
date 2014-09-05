package serializers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/helpers"
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

type Teaser struct {
	Id          int64
	Feed        int64
	Episode     int64
	Title       string
	Description string
	Url         string
	Thumb       string
	Cover       string
	Published   time.Time
	PlaysCount  int
	LovesCount  int
}

type FeaturedsSerializer struct {
	Featureds []models.Featured
	Teasers   []Teaser
	Feeds     []Feed
}

type FeedsSerializer struct {
	Feeds []Feed
}

type Feed struct {
	Id          int64
	Title       string
	Description string
	Url         string
	Thumb       string
	Cover       string
	Updated     time.Time
}

type User struct {
	Id     int64
	Email  string
	Avatar string
}

type ShowUser struct {
	User User
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

func NewTeaser(teaser models.Teaser) Teaser {
	return Teaser{
		Id:          teaser.Id,
		Feed:        teaser.FeedId,
		Episode:     teaser.Episode,
		Title:       teaser.Title,
		Description: teaser.Description,
		Thumb:       fmt.Sprintf("http://s3.amazonaws.com/flowfeeds2/feeds/%d/thumb-x2.jpg", teaser.FeedId),
		Cover:       fmt.Sprintf("http://s3.amazonaws.com/flowfeeds2/feeds/%d/cover.jpg", teaser.FeedId),
		Published:   teaser.Published,
	}
}

func NewUser(user models.User) User {
	return User{
		Id:     user.Id,
		Email:  user.Email,
		Avatar: helpers.Gravatar(user.Email),
	}
}
