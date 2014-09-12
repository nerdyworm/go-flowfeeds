package serializers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"bitbucket.org/nerdyworm/go-flowfeeds/config"
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/modules/web/helpers"
)

type Episodes struct {
	Episodes []Episode
}

type ShowEpisode struct {
	Episode Episode
}

type Episode struct {
	Id          int64
	Feed        int64
	Title       string
	Description string
	Url         string
	Thumb       string
	Cover       string
	Published   time.Time
	PlaysCount  int
	LovesCount  int
	Favorited   bool
	Links       EpisodeLinks `json:"links"`
}

type EpisodeLinks struct {
	Favorites string `json:"favorites"`
	Listens   string `json:"listens"`
	Related   string `json:"related"`
}

func NewShowEpisode(episode models.Episode) ShowEpisode {
	return ShowEpisode{Episode: NewEpisode(episode)}
}

func NewEpisode(episode models.Episode) Episode {
	return Episode{
		Id:          episode.Id,
		Feed:        episode.FeedId,
		Title:       episode.Title,
		Description: episode.Description,
		Url:         episode.Url,
		Thumb:       fmt.Sprintf("http://s3.amazonaws.com/%s/feeds/%d/thumb-x2.jpg", config.S3Bucket, episode.FeedId),
		Cover:       fmt.Sprintf("http://s3.amazonaws.com/%s/feeds/%d/cover.jpg", config.S3Bucket, episode.FeedId),
		Published:   episode.Published,
		Links: EpisodeLinks{
			Favorites: fmt.Sprintf("/api/v1/episodes/%d/favorites", episode.Id),
			Listens:   fmt.Sprintf("/api/v1/episodes/%d/listens", episode.Id),
			Related:   fmt.Sprintf("/api/v1/episodes/%d/related", episode.Id),
		},
	}
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
	Favorited   bool
}

type FeaturedsSerializer struct {
	Featureds []models.Featured
	Feeds     []Feed
	Teasers   []Teaser
}

type FeedsSerializer struct {
	Feeds []Feed
}

type Teasers struct {
	Teasers []Teaser
}

func NewTeasers(teasers []models.Teaser) Teasers {
	serializer := Teasers{}
	serializer.Teasers = make([]Teaser, len(teasers))
	for i, r := range teasers {
		serializer.Teasers[i] = NewTeaser(r)
	}

	return serializer
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
		Url:         teaser.Url,
		Thumb:       fmt.Sprintf("http://s3.amazonaws.com/%s/feeds/%d/thumb-x2.jpg", config.S3Bucket, teaser.FeedId),
		Cover:       fmt.Sprintf("http://s3.amazonaws.com/%s/feeds/%d/cover.jpg", config.S3Bucket, teaser.FeedId),
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

type Listen struct {
	Id      int64
	User    int64
	Episode int64
}

type Listens struct {
	Listens []Listen
}

type ShowListen struct {
	Listen Listen
}

func NewListen(listen models.Listen) Listen {
	return Listen{
		Id:      listen.Id,
		User:    listen.UserId,
		Episode: listen.EpisodeId,
	}
}

func NewShowListen(listen models.Listen) ShowListen {
	return ShowListen{NewListen(listen)}
}

func NewListens(listens []models.Listen) Listens {
	serializer := Listens{}
	serializer.Listens = make([]Listen, len(listens))

	for i, listen := range listens {
		serializer.Listens[i] = NewListen(listen)
	}

	return serializer
}

type Favorite struct {
	Id      int64
	User    int64
	Episode int64
}

type Favorites struct {
	Favorites []Favorite
}

type ShowFavorite struct {
	Favorite Favorite
}

func NewFavorite(favorite models.Favorite) Favorite {
	return Favorite{
		Id:      favorite.Id,
		User:    favorite.UserId,
		Episode: favorite.EpisodeId,
	}
}

func NewShowFavorite(favorite models.Favorite) ShowFavorite {
	return ShowFavorite{NewFavorite(favorite)}
}

func NewFavorites(favorites []models.Favorite) Favorites {
	serializer := Favorites{}
	serializer.Favorites = make([]Favorite, len(favorites))

	for i, favorite := range favorites {
		serializer.Favorites[i] = NewFavorite(favorite)
	}

	return serializer
}
