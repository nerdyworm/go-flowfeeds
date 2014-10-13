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

type ShowEpisode struct {
	Episode Episode
	Feeds   []Feed
}

type Episode struct {
	Id             int64
	Feed           int64
	Title          string
	Description    string
	Url            string
	Thumb          string
	Cover          string
	Published      time.Time
	FavoritesCount int
	ListensCount   int
	Favorited      bool
	Listened       bool
	Links          EpisodeLinks `json:"links"`
}

type EpisodeLinks struct {
	Favorites string `json:"favorites"`
	Listens   string `json:"listens"`
	Related   string `json:"related"`
}

func NewShowEpisode(episode *models.Episode, feed *models.Feed) ShowEpisode {
	return ShowEpisode{
		Episode: NewEpisode(*episode),
		Feeds:   []Feed{NewFeed(*feed)},
	}
}

func NewEpisode(episode models.Episode) Episode {
	return Episode{
		Id:             episode.Id,
		Feed:           episode.FeedId,
		Title:          episode.Title,
		Description:    episode.Description,
		Url:            episode.Url,
		Thumb:          fmt.Sprintf("http://s3.amazonaws.com/%s/feeds/%d/thumb-x2.jpg", config.S3Bucket, episode.FeedId),
		Cover:          fmt.Sprintf("http://s3.amazonaws.com/%s/feeds/%d/cover.jpg", config.S3Bucket, episode.FeedId),
		Published:      episode.Published,
		Favorited:      episode.Favorited,
		Listened:       episode.Listened,
		ListensCount:   episode.ListensCount,
		FavoritesCount: episode.FavoritesCount,
		Links: EpisodeLinks{
			Favorites: fmt.Sprintf("/api/v1/episodes/%d/favorites", episode.Id),
			Listens:   fmt.Sprintf("/api/v1/episodes/%d/listens", episode.Id),
			Related:   fmt.Sprintf("/api/v1/episodes/%d/related", episode.Id),
		},
	}
}

type FeaturedsSerializer struct {
	Featureds []models.Featured
	Feeds     []Feed
	Episodes  []Episode
}

type FeedsSerializer struct {
	Feeds []Feed
}

type Episodes struct {
	Episodes []Episode
	Feeds    []Feed
}

func NewEpisodes(episodes []*models.Episode, feeds []*models.Feed) Episodes {
	serializer := Episodes{}
	serializer.Episodes = make([]Episode, len(episodes))
	serializer.Feeds = make([]Feed, len(feeds))

	for i, episode := range episodes {
		serializer.Episodes[i] = NewEpisode(*episode)
	}

	for i, feed := range feeds {
		serializer.Feeds[i] = NewFeed(*feed)
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

func NewFeed(feed models.Feed) Feed {
	return Feed{
		Id:          feed.Id,
		Title:       feed.Title,
		Description: feed.Description,
		Url:         feed.Url,
		Thumb:       fmt.Sprintf("http://s3.amazonaws.com/%s/feeds/%d/thumb-x2.jpg", config.S3Bucket, feed.Id),
		Cover:       fmt.Sprintf("http://s3.amazonaws.com/%s/feeds/%d/cover.jpg", config.S3Bucket, feed.Id),
	}
}

type User struct {
	Id     int64
	Avatar string
}

type ShowUser struct {
	User User
}

func NewShowUser(user models.User) ShowUser {
	return ShowUser{NewUser(user)}
}

type FeedShowSerializer struct {
	Feed Feed
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

func NewUser(user models.User) User {
	return User{
		Id:     user.Id,
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
	Users   []User
}

type ShowListen struct {
	Listen   Listen
	Episodes []Episode
}

func NewListen(listen models.Listen) Listen {
	return Listen{
		Id:      listen.Id,
		User:    listen.UserId,
		Episode: listen.EpisodeId,
	}
}

func NewShowListen(listen models.Listen, episode models.Episode) ShowListen {
	serializer := ShowListen{}
	serializer.Listen = NewListen(listen)

	serializer.Episodes = make([]Episode, 1)
	serializer.Episodes[0] = NewEpisode(episode)

	return serializer
}

func NewListens(listens []*models.Listen, users []*models.User) Listens {
	serializer := Listens{}
	serializer.Listens = make([]Listen, len(listens))
	serializer.Users = make([]User, len(users))

	for i, listen := range listens {
		serializer.Listens[i] = NewListen(*listen)
	}

	for i, user := range users {
		serializer.Users[i] = NewUser(*user)
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
	Users     []User
}

type ShowFavorite struct {
	Favorite Favorite
	Episodes []Episode
}

type DeleteFavorite struct {
	Episodes []Episode
}

func NewFavorite(favorite models.Favorite) Favorite {
	return Favorite{
		Id:      favorite.Id,
		User:    favorite.UserId,
		Episode: favorite.EpisodeId,
	}
}

func NewShowFavorite(favorite models.Favorite, episode models.Episode) ShowFavorite {
	s := ShowFavorite{}
	s.Favorite = NewFavorite(favorite)

	s.Episodes = make([]Episode, 1)
	s.Episodes[0] = NewEpisode(episode)

	return s
}

func NewDeleteFavorite(favorite models.Favorite, episode models.Episode) DeleteFavorite {
	s := DeleteFavorite{}
	s.Episodes = make([]Episode, 1)
	s.Episodes[0] = NewEpisode(episode)

	return s
}

func NewFavorites(favorites []*models.Favorite, users []*models.User) Favorites {
	serializer := Favorites{}
	serializer.Favorites = make([]Favorite, len(favorites))
	serializer.Users = make([]User, len(favorites))

	for i, favorite := range favorites {
		serializer.Favorites[i] = NewFavorite(*favorite)
	}

	for i, user := range users {
		serializer.Users[i] = NewUser(*user)
	}

	return serializer
}

type ShowFeed struct {
	Feed Feed
}

func NewShowFeed(feed *models.Feed) ShowFeed {
	return ShowFeed{
		Feed: NewFeed(*feed),
	}
}
