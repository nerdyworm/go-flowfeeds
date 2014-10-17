package serializers

import (
	"fmt"

	"bitbucket.org/nerdyworm/go-flowfeeds/config"
	"bitbucket.org/nerdyworm/go-flowfeeds/models"
	"bitbucket.org/nerdyworm/go-flowfeeds/web/helpers"
)

type Episode struct {
	*models.Episode
	Thumb string
	Cover string
	Links EpisodeLinks `json:"links"`
}

type EpisodeLinks struct {
	Favorites string `json:"favorites"`
	Listens   string `json:"listens"`
	Related   string `json:"related"`
}

func NewEpisode(episode *models.Episode) Episode {
	return Episode{
		episode,
		fmt.Sprintf("%s/feeds/%d/thumb-x2.jpg", config.ASSETS, episode.FeedId),
		fmt.Sprintf("%s/feeds/%d/cover.jpg", config.ASSETS, episode.FeedId),
		EpisodeLinks{
			Favorites: fmt.Sprintf("/api/v1/episodes/%d/favorites", episode.Id),
			Listens:   fmt.Sprintf("/api/v1/episodes/%d/listens", episode.Id),
			Related:   fmt.Sprintf("/api/v1/episodes/%d/related", episode.Id),
		},
	}
}

type ShowEpisode struct {
	Episode Episode
	Feeds   []Feed
}

func NewShowEpisode(episode *models.Episode, feed *models.Feed) ShowEpisode {
	return ShowEpisode{
		Episode: NewEpisode(episode),
		Feeds:   []Feed{NewFeed(feed)},
	}
}

type Episodes struct {
	Episodes []Episode
	Feeds    []Feed
}

func NewEpisodes(episodes []*models.Episode, feeds []*models.Feed) Episodes {
	s := Episodes{}
	s.Episodes = make([]Episode, len(episodes))
	s.Feeds = make([]Feed, len(feeds))

	for i, episode := range episodes {
		s.Episodes[i] = NewEpisode(episode)
	}

	for i, feed := range feeds {
		s.Feeds[i] = NewFeed(feed)
	}

	return s
}

type Feed struct {
	*models.Feed
	Thumb string
	Cover string
}

func NewFeed(feed *models.Feed) Feed {
	return Feed{
		feed,
		fmt.Sprintf("%s/feeds/%d/thumb-x2.jpg", config.ASSETS, feed.Id),
		fmt.Sprintf("%s/feeds/%d/cover.jpg", config.ASSETS, feed.Id),
	}
}

type ShowFeed struct {
	Feed Feed
}

type Feeds struct {
	Feeds []Feed
}

func NewFeeds(feeds []*models.Feed) Feeds {
	s := Feeds{}
	s.Feeds = make([]Feed, len(feeds))
	for i, feed := range feeds {
		s.Feeds[i] = NewFeed(feed)
	}
	return s
}

func NewShowFeed(feed *models.Feed) ShowFeed {
	return ShowFeed{
		Feed: NewFeed(feed),
	}
}

type User struct {
	Id     int64
	Avatar string
}

func NewUser(user *models.User) User {
	return User{
		Id:     user.Id,
		Avatar: helpers.Gravatar(user.Email),
	}
}

type ShowUser struct {
	User User
}

func NewShowUser(user *models.User) ShowUser {
	return ShowUser{NewUser(user)}
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

func NewListen(listen *models.Listen) Listen {
	return Listen{listen.Id, listen.UserId, listen.EpisodeId}
}

func NewShowListen(listen *models.Listen, episode *models.Episode) ShowListen {
	s := ShowListen{}
	s.Listen = NewListen(listen)

	s.Episodes = make([]Episode, 1)
	s.Episodes[0] = NewEpisode(episode)

	return s
}

func NewListens(listens []*models.Listen, users []*models.User) Listens {
	s := Listens{}
	s.Listens = make([]Listen, len(listens))
	s.Users = make([]User, len(users))

	for i, listen := range listens {
		s.Listens[i] = NewListen(listen)
	}

	for i, user := range users {
		s.Users[i] = NewUser(user)
	}

	return s
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

func NewFavorite(favorite *models.Favorite) Favorite {
	return Favorite{favorite.Id, favorite.UserId, favorite.EpisodeId}
}

func NewFavorites(favorites []*models.Favorite, users []*models.User) Favorites {
	s := Favorites{}
	s.Favorites = make([]Favorite, len(favorites))
	s.Users = make([]User, len(favorites))

	for i, favorite := range favorites {
		s.Favorites[i] = NewFavorite(favorite)
	}

	for i, user := range users {
		s.Users[i] = NewUser(user)
	}

	return s
}
